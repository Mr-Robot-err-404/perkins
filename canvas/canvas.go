package canvas

import (
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Grid [][]rune

type Model struct {
	width  int
	height int
	cursor Cursor
	Grid   Grid
}

func New(width, height int, grid Grid) Model {
	return Model{
		width:  width,
		height: height,
		Grid:   grid,
		cursor: Cursor{},
	}
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
		Background(theme.SumiInk1).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(grid_to_canvas(m.Grid, m.cursor))
}

func (m Model) Resize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}
