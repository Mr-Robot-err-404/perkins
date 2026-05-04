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
	window   *core.Window
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
		window: new(core.Window),
	}
}

func (m Model) reset_window(grid core.Grid) {
	midpoint := core.Find_Center(grid)
	*m.window = core.Get_Window(core.Dimensions{Width: m.width, Height: m.height - 1}, grid, midpoint)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.reset_window(m.grid)

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
			m.reset_window(m.grid)

		case "-", "_":
			size, factor := zoom_out(m.base, m.factor)
			m.grid, *m.bitmap = core.Image_To_Ascii(m.img, size, m.inverted)
			m.factor = factor
			m.reset_window(m.grid)

		case "i":
			m.inverted = !m.inverted
			scale := core.Dimensions{
				Width:  amplify(m.base.Width, m.factor),
				Height: amplify(m.base.Height, m.factor),
			}
			m.bitmap.Invert()
			m.grid = core.Image_To_Grid(*m.bitmap, scale)

		case "k", "up":
			if m.window.Start.Row == 0 {
				return m, nil
			}
			m.window.Start.Row = max(0, m.window.Start.Row-1)
			m.window.End.Row--
		case "j", "down":
			if m.window.End.Row == len(m.grid) {
				return m, nil
			}
			m.window.Start.Row++
			m.window.End.Row = min(len(m.grid), m.window.End.Row+1)

		case "ctrl+u":
			if m.window.Start.Row == 0 {
				return m, nil
			}
			dy := min(10, m.window.Start.Row)
			m.window.Start.Row -= dy
			m.window.End.Row -= dy
		case "ctrl+d":
			if m.window.End.Row == len(m.grid) {
				return m, nil
			}
			dy := min(10, len(m.grid)-m.window.End.Row)
			m.window.Start.Row += dy
			m.window.End.Row += dy

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
	grid := window(m.grid, core.Dimensions{Width: m.width, Height: m.height - 1}, *m.window)
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
		info("j/k", "up/down", 4, 1),
		info(":q", "quit", 5, 4),
	)
	return component.Notification(list, 20, 8, theme.WaveBlue, theme.SumiInk0)
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
