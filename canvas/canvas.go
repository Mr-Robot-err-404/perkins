package canvas

import (
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Grid [][]rune

type Model struct {
	width       int
	height      int
	mode        int
	Cursor      *core.Pos
	prev_cursor *core.Pos
	Grid        Grid
	Selected    core.Selected
}

const (
	NORMAL_MODE int = iota
	VISUAL_BLOCK
)

const (
	VIM_LEFT  string = "h"
	VIM_RIGHT string = "l"
	VIM_DOWN  string = "j"
	VIM_UP    string = "k"
	JUMP_DOWN string = "ctrl+d"
	JUMP_UP   string = "ctrl+u"
	CENTER    string = "c"
)

func New(width, height int, grid Grid, selected core.Selected) Model {
	return Model{
		width:       width,
		height:      height,
		Grid:        grid,
		Selected:    selected,
		Cursor:      &core.Pos{},
		prev_cursor: &core.Pos{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) toggle_mode(mode int) int {
	if mode == m.mode {
		return NORMAL_MODE
	}
	return mode
}
func (m Model) expand_selection() {
	if m.mode != VISUAL_BLOCK {
		return
	}
	pos := *m.Cursor
	m.Selected[pos] = true
}
func (m Model) update_cursor(pos core.Pos) {
	*m.prev_cursor = *m.Cursor
	*m.Cursor = pos
	m.expand_selection()
}
func (m Model) set_normal_mode() Model {
	m.mode = NORMAL_MODE
	clear(m.Selected)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case VIM_LEFT:
			m.update_cursor(core.Pos{Row: m.Cursor.Row, Col: max(0, m.Cursor.Col-1)})
		case VIM_RIGHT:
			m.update_cursor(core.Pos{
				Row: m.Cursor.Row,
				Col: min(len(m.Grid[m.Cursor.Row])-1, m.Cursor.Col+1),
			})
		case VIM_UP:
			m.update_cursor(core.Pos{
				Col: m.Cursor.Col,
				Row: max(0, m.Cursor.Row-1),
			})
		case VIM_DOWN:
			m.update_cursor(core.Pos{
				Col: m.Cursor.Col,
				Row: min(len(m.Grid)-1, m.Cursor.Row+1),
			})

		case JUMP_DOWN:
			m.update_cursor(core.Pos{
				Col: m.Cursor.Col,
				Row: len(m.Grid) - 1,
			})
		case JUMP_UP:
			m.update_cursor(core.Pos{
				Col: m.Cursor.Col,
				Row: 0,
			})
		case CENTER:
			m.update_cursor(find_center(m.Grid))

		case "G":
			m.update_cursor(core.Pos{
				Col: m.Cursor.Col,
				Row: len(m.Grid) - 1,
			})
		case "$":
			m.update_cursor(core.Pos{
				Row: m.Cursor.Row,
				Col: len(m.Grid[m.Cursor.Row]) - 1,
			})
		case "_":
			m.update_cursor(core.Pos{
				Row: m.Cursor.Row,
				Col: 0,
			})

		case "v", "ctrl+v":
			mode := m.toggle_mode(VISUAL_BLOCK)

			if mode == NORMAL_MODE {
				m = m.set_normal_mode()
				return m, nil
			}
			m.mode = mode
			pos := *m.Cursor
			m.Selected[pos] = true
		case "esc":
			m = m.set_normal_mode()
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
		Render(grid_to_canvas(m.Grid, m.Selected, *m.Cursor))
}

func (m Model) Resize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}
