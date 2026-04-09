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
	Mode        int
	Grid        Grid
	Selected    core.Selected
	harpoon     *Harpoon
	prev_cursor *core.Pos
	Cursor      *core.Pos
	selector    *Selector
}
type Harpoon struct {
	min   core.Pos
	max   core.Pos
	start core.Pos
}
type Selector struct {
	mirror int
}
type CropMsg struct{ Grid Grid }

const (
	NORMAL_MODE int = iota
	VISUAL_BLOCK
	CROP_MODE
)
const (
	MIRROR_DISABLE int = iota
	X_AXIS
	Y_AXIS
)

const (
	VIM_LEFT  string = "h"
	VIM_RIGHT string = "l"
	VIM_DOWN  string = "j"
	VIM_UP    string = "k"
	JUMP_DOWN string = "ctrl+d"
	JUMP_UP   string = "ctrl+u"
	CENTER    string = "t"
	CONFIRM   string = "ctrl+y"
	CROP      string = "c"
)

func New(width, height int, grid Grid, selected core.Selected) Model {
	return Model{
		width:       width,
		height:      height,
		Grid:        grid,
		Selected:    selected,
		Cursor:      &core.Pos{},
		prev_cursor: &core.Pos{},
		harpoon:     &Harpoon{},
		selector:    &Selector{},
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
		case CONFIRM:
			if m.Mode != CROP_MODE {
				return m, nil
			}
			grid := m.crop_canvas()

			return m, func() tea.Msg {
				return CropMsg{Grid: grid}
			}
		case CROP:
			m.Mode = m.toggle_mode(CROP_MODE)

			if m.Mode == NORMAL_MODE {
				m.Reset_to_normal()
				return m, nil
			}
			m.set_mirror_axis(Y_AXIS)
			m.init_cropping_block()
			m.expand_selection()

		case "v", "ctrl+v":
			m.Mode = m.toggle_mode(VISUAL_BLOCK)

			if m.Mode == NORMAL_MODE {
				m.Reset_to_normal()
				return m, nil
			}
			pos := *m.Cursor
			m.Selected[pos] = core.Highlight
			*m.harpoon = Harpoon{min: pos, max: pos, start: pos}

		case "tab":
			m.selector.mirror = m.toggle_mirror_axis()

			switch m.Mode {
			case VISUAL_BLOCK:
				m.expand_selection()
			case CROP_MODE:
				m.init_cropping_block()
				m.expand_selection()
			}

		case "=":
			switch m.Mode {
			case VISUAL_BLOCK:
				m.set_mirror_axis(X_AXIS)
				m.expand_selection()
			case CROP_MODE:
			}
		case "|":
			switch m.Mode {
			case VISUAL_BLOCK:
				m.set_mirror_axis(Y_AXIS)
				m.expand_selection()
			case CROP_MODE:
			}

		case "esc":
			m.Mode = NORMAL_MODE
			m.Reset_to_normal()
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
