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

func New(width, height int, get func() rune) Model {
	return Model{width: width, height: height, Cell: get}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
