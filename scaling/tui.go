package scaling

import (
	"image"
	"strings"

	"github.com/Mr-Robot-err-404/perkins/canvas"
	"github.com/Mr-Robot-err-404/perkins/component"
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/debug"
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width    int
	height   int
	mode     int
	grid     core.Grid
	base     core.Dimensions
	bitmap   *core.ImageBitmap
	factor   float64
	img      image.Image
	cmd      []rune
	ch       chan<- core.Grid
	inverted bool
}

func New(img image.Image, size core.Dimensions, ch chan<- core.Grid) Model {
	grid, bitmap := core.Image_To_Ascii(img, size, false)
	return Model{
		grid:   grid,
		bitmap: &bitmap,
		base:   size,
		img:    img,
		factor: 1.0,
		ch:     ch,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.mode == canvas.COMMAND_MODE {
			switch msg.String() {
			case "esc":
				m.mode = canvas.NORMAL_MODE
				m.cmd = m.cmd[:0]

			case "backspace":
				if len(m.cmd) > 0 {
					m.cmd = m.cmd[:len(m.cmd)-1]
				}

			case " ":
				m.cmd = append(m.cmd, ' ')

			case "enter":
				switch string(m.cmd) {
				case "q", "quit":
					return m, tea.Quit
				}
				m.mode = canvas.NORMAL_MODE
				m.cmd = m.cmd[:0]

			default:
				m.cmd = append(m.cmd, msg.Runes...)
			}
			return m, nil
		}

		switch msg.String() {
		case "+":
			size, factor := zoom_in(m.base, m.factor, m.img.Bounds())
			m.grid, *m.bitmap = core.Image_To_Ascii(m.img, size, m.inverted)
			m.factor = factor

		case "-", "_":
			size, factor := zoom_out(m.base, m.factor)
			m.grid, *m.bitmap = core.Image_To_Ascii(m.img, size, m.inverted)
			m.factor = factor

		case "i":
			m.inverted = !m.inverted
			scale := core.Dimensions{
				Width:  amplify(m.base.Width, m.factor),
				Height: amplify(m.base.Height, m.factor),
			}
			m.bitmap.Invert()
			m.grid = core.Image_To_Grid(*m.bitmap, scale)

		case "ctrl+c":
			return m, tea.Quit
		case ":":
			m.mode = canvas.COMMAND_MODE
		case "enter":
			m.ch <- m.grid
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	indicator := canvas.Status_Bar(canvas.Status{
		Mode:  m.mode,
		Width: m.width,
		Cmd:   string(m.cmd),
	})
	grid := window(m.grid, core.Dimensions{Width: m.width, Height: m.height - 1})
	ascii := canvas.Grid_To_Canvas(grid, core.Selected{}, core.Pos{Row: -1, Col: -1}, false)

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 1).
		Background(theme.SumiInk1).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(ascii)
	screen := lipgloss.JoinVertical(lipgloss.Left, content, indicator)
	overlay, err := component.Overlay(screen, menu(), 1, 2, true)

	if err != nil {
		debug.Logf("overlay failed: %s", err.Error())
		return screen
	}
	return overlay
}

func menu() string {
	list := lipgloss.JoinVertical(lipgloss.Left,
		info("enter", "continue", 2, 0),
		info("+", "zoom in", 6, 1),
		info("-", "zoom out", 6, 0),
		info("i", "invert", 6, 2),
		info(":q", "quit", 5, 4),
	)
	return component.Notification(list, 20, 7, theme.WaveBlue, theme.SumiInk0)
}

func info(key string, value string, space int, pad int) string {
	primary := lipgloss.NewStyle().
		Background(theme.SumiInk0).
		Foreground(theme.Cursor).
		PaddingRight(1).
		Render(key)
	gap := lipgloss.NewStyle().
		Background(theme.SumiInk0).
		Render(strings.Repeat(" ", space))

	secondary := lipgloss.NewStyle().
		Background(theme.SumiInk0).
		Foreground(theme.FujiGray).
		PaddingRight(pad).
		Render(value)
	return lipgloss.JoinHorizontal(lipgloss.Bottom, primary, gap, secondary)
}

func Run(img image.Image, size core.Dimensions, ch chan<- core.Grid) error {
	p := tea.NewProgram(New(img, size, ch), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
