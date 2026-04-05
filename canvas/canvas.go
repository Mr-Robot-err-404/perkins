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

const (
	VIM_LEFT  string = "h"
	VIM_RIGHT string = "l"
	VIM_DOWN  string = "j"
	VIM_UP    string = "k"
	JUMP_DOWN string = "ctrl+d"
	JUMP_UP   string = "ctrl+u"
	CENTER    string = "c"
)

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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case VIM_LEFT:
			m.cursor.col = max(0, m.cursor.col-1)

		case VIM_RIGHT:
			m.cursor.col = min(len(m.Grid[m.cursor.row])-1, m.cursor.col+1)
		case VIM_UP:
			m.cursor.row = max(0, m.cursor.row-1)

		case VIM_DOWN:
			m.cursor.row = min(len(m.Grid)-1, m.cursor.row+1)
		case JUMP_DOWN:
			m.cursor.row = len(m.Grid) - 1
		case JUMP_UP:
			m.cursor.row = 0
		case "$":
			m.cursor.col = len(m.Grid[m.cursor.row]) - 1
		case "_":
			m.cursor.col = 0

		}
	}
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
