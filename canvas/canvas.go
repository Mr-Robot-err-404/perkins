package canvas

import (
	"github.com/Mr-Robot-err-404/perkins/common"
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Grid [][]rune

type Model struct {
	width  int
	height int
	Cursor *common.Pos
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
		Cursor: &common.Pos{},
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
			m.Cursor.Col = max(0, m.Cursor.Col-1)
		case VIM_RIGHT:
			m.Cursor.Col = min(len(m.Grid[m.Cursor.Row])-1, m.Cursor.Col+1)
		case VIM_UP:
			m.Cursor.Row = max(0, m.Cursor.Row-1)
		case VIM_DOWN:
			m.Cursor.Row = min(len(m.Grid)-1, m.Cursor.Row+1)

		case JUMP_DOWN:
			m.Cursor.Row = len(m.Grid) - 1
		case JUMP_UP:
			m.Cursor.Row = 0
		case CENTER:
			*m.Cursor = find_center(m.Grid)

		case "G":
			m.Cursor.Row = len(m.Grid) - 1
		case "$":
			m.Cursor.Col = len(m.Grid[m.Cursor.Row]) - 1
		case "_":
			m.Cursor.Col = 0

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
		Render(grid_to_canvas(m.Grid, *m.Cursor))
}

func (m Model) Resize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}
