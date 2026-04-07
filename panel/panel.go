package panel

import (
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width  int
	height int
	Cell   func() rune
}
type FlipMsg struct{ Bit byte }
type ActionMsg struct{ Action int }

const (
	FILL_ACTION int = iota + 1
	CLEAR_ACTION
)

func New(width, height int, get func() rune) Model {
	return Model{width: width, height: height, Cell: get}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func FlipBit(b byte) tea.Cmd {
	return func() tea.Msg {
		return FlipMsg{Bit: b}
	}
}

func Clear() tea.Cmd {
	return func() tea.Msg {
		return ActionMsg{Action: CLEAR_ACTION}
	}
}
func Fill() tea.Cmd {
	return func() tea.Msg {
		return ActionMsg{Action: FILL_ACTION}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1", "q":
			return m, FlipBit(0)
		case "2", "w":
			return m, FlipBit(1)
		case "3", "e":
			return m, FlipBit(2)
		case "4", "r":
			return m, FlipBit(6)
		case "5", "u":
			return m, FlipBit(3)
		case "6", "i":
			return m, FlipBit(4)
		case "7", "o":
			return m, FlipBit(5)
		case "8", "p":
			return m, FlipBit(7)
		case "x":
			return m, Clear()
		case "f":
			return m, Fill()
		}
	}
	return m, nil
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(theme.SumiInk3).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(m.magnify())
}

func (m Model) magnify() string {
	r := m.Cell()
	bits := magnifier(r)
	return render_magnifier(bits)
}

func (m Model) Resize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}
