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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			return m, FlipBit(0)
		case "2":
			return m, FlipBit(1)
		case "3":
			return m, FlipBit(2)
		case "4":
			return m, FlipBit(6)
		case "5":
			return m, FlipBit(3)
		case "6":
			return m, FlipBit(4)
		case "7":
			return m, FlipBit(5)
		case "8":
			return m, FlipBit(7)
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
	cells := magnifier(r)
	return render_magnifier(cells)
}

func (m Model) Resize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}
