package scaling

import (
	"image"

	"github.com/Mr-Robot-err-404/perkins/canvas"
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width  int
	height int
	mode   int
	grid   core.Grid
	base   core.Dimensions
	factor float64
	img    image.Image
	cmd    []rune
}

func New(img image.Image, size core.Dimensions) Model {
	return Model{
		grid:   core.Image_To_Ascii(img, size),
		base:   size,
		img:    img,
		factor: 1.0,
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
			m.grid = core.Image_To_Ascii(m.img, size)
			m.factor = factor

		case "-", "_":
			size, factor := zoom_out(m.base, m.factor)
			m.grid = core.Image_To_Ascii(m.img, size)
			m.factor = factor

		case "ctrl+c":
			return m, tea.Quit
		case ":":
			m.mode = canvas.COMMAND_MODE
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

	return lipgloss.JoinVertical(lipgloss.Left, content, indicator)
}

func Run(img image.Image, size core.Dimensions) error {
	p := tea.NewProgram(New(img, size), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
