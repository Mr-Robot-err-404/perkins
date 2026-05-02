package scaling

import (
	"github.com/Mr-Robot-err-404/perkins/canvas"
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width  int
	height int
	ascii  string
	mode   int
	cmd    []rune
}

func New(ascii string) Model {
	return Model{
		ascii: ascii,
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
	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-1).
		Background(theme.SumiInk1).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Padding(1, 2).
		Render(m.ascii)

	return lipgloss.JoinVertical(lipgloss.Left, content, indicator)
}

func Run(grid core.Grid) error {
	ascii := canvas.Grid_To_Canvas(grid, core.Selected{}, core.Pos{Row: -1, Col: -1}, false)
	p := tea.NewProgram(New(ascii), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
