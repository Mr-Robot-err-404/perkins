package canvas

import (
	"log"
	"strings"
	"time"

	"github.com/Mr-Robot-err-404/perkins/component"
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width       int
	height      int
	Mode        int
	Grid        core.Grid
	Selected    core.Selected
	Draw        core.Selected
	harpoon     *Harpoon
	prev_cursor *core.Pos
	Cursor      *core.Pos
	cmd         *[]rune
	message     *string
	mirror      *Mirror
	v_mirror    *Mirror
	save_modal  component.Modal
}
type Harpoon struct {
	min   core.Pos
	max   core.Pos
	start core.Pos
}
type Mirror struct {
	enabled bool
	axis    int
}
type CropMsg struct{ Grid core.Grid }
type UndoMsg struct{}
type RedoMsg struct{}
type SaveMsg struct {
	Path  string
	Ascii []byte
}
type StatusMsg struct{ Status string }
type Flush struct{}

type ComponentModal interface {
	IsActive() bool
	View() string
}

const (
	NORMAL_MODE int = iota
	VISUAL_BLOCK
	CROP_MODE
	COMMAND_MODE
	DRAW_MODE
)
const (
	Y_AXIS int = iota
	X_AXIS
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
	UNDO      string = "u"
	REDO      string = "ctrl+r"
)

var SaveConfig = component.ModalConfig{
	Width:       60,
	Title:       "Write to file",
	Placeholder: "~/pixel.art",
	XPadding:    1,
}

func New(width, height int, grid core.Grid, selected core.Selected, file_path string) Model {
	return Model{
		width:       width,
		height:      height,
		Grid:        grid,
		Selected:    selected,
		Draw:        make(core.Selected),
		Cursor:      new(core.Pos),
		prev_cursor: new(core.Pos),
		harpoon:     new(Harpoon),
		cmd:         new([]rune),
		message:     new(string),
		mirror:      new(Mirror),
		save_modal:  component.NewModal(SaveConfig, file_path),
	}
}

func emit[T any](msg T) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func flush() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second * 5)
		return Flush{}
	}
}

func Notify(message string) tea.Cmd {
	return func() tea.Msg {
		return StatusMsg{Status: message}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.save_modal.Active {
		var cmd tea.Cmd
		m.save_modal, cmd = m.save_modal.Update(msg)
		return m, cmd
	}
	switch msg := msg.(type) {
	case component.ModalSubmit:
		if msg.Cancel {
			return m, nil
		}
		ascii := Grid_To_Canvas(m.Grid, core.Selected{}, core.Pos{Row: -1, Col: -1})
		return m, emit(SaveMsg{Path: strings.TrimSpace(msg.Value), Ascii: []byte(ascii)})

	case Flush:
		*m.message = ""

	case StatusMsg:
		*m.message = msg.Status
		return m, flush()

	case tea.MouseMsg:
		pos, ok := m.mouse_to_grid(msg.X, msg.Y)
		if !ok {
			return m, nil
		}
		switch msg.Action {
		case tea.MouseActionPress:
			switch msg.Button {
			case tea.MouseButtonLeft:
				if m.Mode == VISUAL_BLOCK {
					m.Mode = NORMAL_MODE
					m.Reset_to_normal()
				}
				m.update_cursor(pos)
			case tea.MouseButtonRight:
			}
		case tea.MouseActionMotion:
			if msg.Button == tea.MouseButtonLeft {
				if m.Mode == NORMAL_MODE {
					m.Mode = VISUAL_BLOCK
					anchor := *m.Cursor
					*m.harpoon = Harpoon{min: anchor, max: anchor, start: anchor}
				}
				m.update_cursor(pos)
			}
		}

	case tea.KeyMsg:
		if m.Mode == COMMAND_MODE {
			switch msg.String() {
			case "esc":
				m.Mode = NORMAL_MODE
				m.Reset_to_normal()

			case "backspace":
				r := *m.cmd
				end := max(0, len(r)-1)
				*m.cmd = r[0:end]
			case " ":
				*m.cmd = append(*m.cmd, ' ')

			case "enter":
				cmd := strings.TrimSpace(string(*m.cmd))
				switch cmd {
				case "w", "write":
					m.save_modal.Active = true
				case "theme":
				case "h", "help":
				case "q", "quit":
					return m, tea.Quit
				default:
				}
				m.Mode = NORMAL_MODE
				m.Reset_to_normal()

			default:
				*m.cmd = append(*m.cmd, msg.Runes...)
			}
			return m, nil
		}
		switch msg.String() {
		case UNDO:
			return m, emit(UndoMsg{})
		case REDO:
			return m, emit(RedoMsg{})

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
			m.update_cursor(core.Find_Center(m.Grid))

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
			return m, emit(CropMsg{Grid: grid})

		case CROP:
			m.Mode = m.toggle_mode(CROP_MODE)

			if m.Mode == NORMAL_MODE {
				m.Reset_to_normal()
				return m, nil
			}
			m.mirror.enabled = true
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
			*m.harpoon = Harpoon{min: pos, max: pos, start: pos}
			m.expand_selection()

		case "V":
			m.Mode = m.toggle_mode(VISUAL_BLOCK)

			if m.Mode == NORMAL_MODE {
				m.Reset_to_normal()
				return m, nil
			}
			pos := *m.Cursor
			*m.harpoon = Harpoon{
				start: core.Pos{Row: pos.Row, Col: 0},
				min:   core.Pos{Row: pos.Row, Col: 0},
				max:   core.Pos{Row: pos.Row, Col: len(m.Grid[pos.Row]) - 1},
			}
			*m.Cursor = m.harpoon.max
			m.expand_selection()

		case "d":
			m.Mode = m.toggle_mode(DRAW_MODE)

			if m.Mode == NORMAL_MODE {
				m.Reset_to_normal()
				return m, nil
			}
			m.Draw[*m.Cursor] = core.Highlight

		case "m":
			m.toggle_mirror()

		case "tab":
			m.mirror.axis = m.toggle_mirror_axis()

			switch m.Mode {
			case VISUAL_BLOCK, DRAW_MODE:
				m.expand_selection()
			case CROP_MODE:
				m.init_cropping_block()
				m.expand_selection()
			}
		case "esc":
			m.Mode = NORMAL_MODE
			m.Reset_to_normal()
		case ":":
			m.Mode = COMMAND_MODE
		}
	}
	return m, nil
}

func (m Model) View() string {
	grid_str := Grid_To_Canvas(m.Grid, m.Selected, *m.Cursor)
	indicator := status_bar(Status{
		mode:    m.Mode,
		width:   m.width,
		cmd:     string(*m.cmd),
		message: *m.message,
		mirror:  m.mirror.enabled,
		axis:    m.mirror.axis,
	})

	grid_h := len(m.Grid)
	top_pad := (m.height - 1 - grid_h) / 2

	centered := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 1).
		Background(theme.SumiInk1).
		PaddingTop(top_pad).
		AlignHorizontal(lipgloss.Center).
		Render(grid_str)

	modals := []ComponentModal{m.save_modal}

	for _, current := range modals {
		if !current.IsActive() {
			continue
		}
		overlay, err := component.OverlayCenter(centered, current.View(), true)
		if err != nil {
			log.Printf("Failed to overlay: %s", err.Error())
			return centered
		}
		return lipgloss.JoinVertical(lipgloss.Left, overlay, indicator)
	}
	return lipgloss.JoinVertical(lipgloss.Left, centered, indicator)
}

func (m Model) Resize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

type Status struct {
	mode    int
	width   int
	cmd     string
	message string
	mirror  bool
	axis    int
}

func status_label(color lipgloss.Color, label string) string {
	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(label)
}

func subtitle(s string) string {
	return lipgloss.NewStyle().
		PaddingLeft(1).
		Foreground(theme.FujiGray).
		Background(theme.SumiInk2).
		Render(s)
}

func status_bar(status Status) string {
	var label string
	var msg string
	var color lipgloss.Color

	switch status.mode {
	case COMMAND_MODE:
		label := status_label(theme.RoninYellow, "COMMAND ")
		cmd := lipgloss.NewStyle().
			Background(theme.SumiInk2).
			Bold(true).
			Render(":" + status.cmd + "█")
		return lipgloss.NewStyle().
			Background(theme.SumiInk2).
			PaddingLeft(1).
			Width(status.width).
			AlignHorizontal(lipgloss.Left).
			Render(label + cmd)

	case VISUAL_BLOCK:
		label = "VISUAL"
		color = theme.Wisteria
	case CROP_MODE:
		label = "CROP"
		color = theme.SamuraiRed
	case DRAW_MODE:
		label = "DRAW"
		color = theme.Wisteria
	default:
		label = "NORMAL"
		msg = status.message
		color = theme.WaveBlue
	}
	var indicators []string
	indicators = append(indicators, status_label(color, label))

	if status.mirror {
		axis_label := "Y_AXIS"
		if status.axis == X_AXIS {
			axis_label = "X_AXIS"
		}
		indicators = append(indicators, subtitle("MIRROR::"+axis_label))
	}
	indicators = append(indicators, subtitle(msg))

	return lipgloss.NewStyle().
		Background(theme.SumiInk2).
		PaddingLeft(1).
		Width(status.width).
		AlignHorizontal(lipgloss.Left).
		Render(lipgloss.JoinHorizontal(lipgloss.Bottom, indicators...))
}
