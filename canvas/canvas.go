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
	Window      *core.Window
	cmd         *[]rune
	message     *string
	mirror      *Mirror
	v_mirror    *Mirror
	n           *int
	save_modal  component.Modal
	theme_modal component.SltModal
	help_open   bool
	help_page   int
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
type ThemeMsg struct{ Idx int }
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
var ThemeConfig = component.SltModalConfig{
	Width:    60,
	Title:    "Select palette theme",
	XPadding: 1,
}

func New(width, height int, grid core.Grid, selected core.Selected, file_path string) Model {
	midpoint := core.Find_Center(grid)
	window := core.Get_Window(core.Dimensions{Width: width, Height: height - 1}, grid, midpoint)
	return Model{
		width:       width,
		height:      height,
		Grid:        grid,
		Selected:    selected,
		Draw:        make(core.Selected),
		prev_cursor: new(core.Pos),
		harpoon:     new(Harpoon),
		cmd:         new([]rune),
		message:     new(string),
		mirror:      new(Mirror),
		n:           new(int),
		Window:      &window,
		Cursor:      &midpoint,
		save_modal:  component.NewModal(SaveConfig, file_path),
		theme_modal: component.NewSltModal(ThemeConfig),
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
	if m.theme_modal.Active {
		var cmd tea.Cmd
		m.theme_modal, cmd = m.theme_modal.Update(msg)
		return m, cmd
	}
	if m.help_open {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "1", "2", "3", "4":
				d := int(msg.Runes[0] - '0')
				m.help_page = d - 1
				return m, nil
			case "?", "esc":
				m.help_open = false
				return m, nil
			}
		}
	}
	switch msg := msg.(type) {
	case component.SltModalSubmit:
		if !msg.Cancel {
			return m, emit(ThemeMsg{Idx: msg.Idx})
		}

	case component.ModalSubmit:
		if msg.Cancel {
			return m, nil
		}
		ascii := Grid_To_Canvas(m.Grid, core.Selected{}, core.Pos{Row: -1, Col: -1}, false)
		return m, emit(SaveMsg{Path: strings.TrimSpace(msg.Value), Ascii: []byte(ascii)})

	case Flush:
		*m.message = ""

	case StatusMsg:
		*m.message = msg.Status
		return m, flush()

	case tea.MouseMsg:
		if msg.X > m.width || msg.Y >= m.height {
			return m, nil
		}
		pos, ok := m.mouse_to_grid(msg.X, msg.Y, *m.Window)
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
				case "t", "th", "theme":
					names := make([]string, len(theme.Themes))
					for i, t := range theme.Themes {
						names[i] = t.Name
					}
					m.theme_modal.Slt.List = names
					m.theme_modal.Active = true
				case "h", "help":
					m.help_open = !m.help_open
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

		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			d := int(msg.Runes[0] - '0')
			*m.n *= 10
			*m.n += d
			return m, nil

		case VIM_LEFT:
			n := m.consume()
			m.update_cursor(core.Pos{Row: m.Cursor.Row, Col: max(0, m.Cursor.Col-n)})
		case VIM_RIGHT:
			n := m.consume()
			m.update_cursor(core.Pos{
				Row: m.Cursor.Row,
				Col: min(len(m.Grid[m.Cursor.Row])-1, m.Cursor.Col+n),
			})
		case VIM_UP:
			n := m.consume()
			m.update_cursor(core.Pos{
				Col: m.Cursor.Col,
				Row: max(0, m.Cursor.Row-n),
			})
		case VIM_DOWN:
			n := m.consume()
			m.update_cursor(core.Pos{
				Col: m.Cursor.Col,
				Row: min(len(m.Grid)-1, m.Cursor.Row+n),
			})

		case JUMP_DOWN:
			m.update_cursor(core.Pos{
				Col: m.Cursor.Col,
				Row: min(len(m.Grid)-1, m.Cursor.Row+10),
			})
		case JUMP_UP:
			m.update_cursor(core.Pos{
				Col: m.Cursor.Col,
				Row: max(0, m.Cursor.Row-10),
			})

		case "w":
			n := m.consume() * 10
			m.update_cursor(core.Pos{
				Row: m.Cursor.Row,
				Col: min(len(m.Grid[0])-1, m.Cursor.Col+n),
			})
		case "b":
			n := m.consume() * 10
			m.update_cursor(core.Pos{
				Col: max(0, m.Cursor.Col-n),
				Row: m.Cursor.Row,
			})

		case CENTER:
			midpoint := core.Find_Center(m.Grid)
			*m.Window = core.Get_Window(core.Dimensions{Width: m.width, Height: m.height - 1}, m.Grid, midpoint)
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
			if !m.can_crop(Y_AXIS) {
				m.Mode = NORMAL_MODE
				m.Reset_to_normal()
				return m, Notify("canvas too small to crop")
			}
			m.mirror.enabled = true
			m.set_mirror_axis(Y_AXIS)
			m.init_cropping_block()
			m.expand_selection()
			update_window(m.Window, *m.Cursor)

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
			next_axis := m.toggle_mirror_axis()

			if m.Mode == CROP_MODE && !m.can_crop(next_axis) {
				return m, Notify("canvas too small to crop on this axis")
			}
			m.mirror.axis = next_axis

			switch m.Mode {
			case VISUAL_BLOCK, DRAW_MODE:
				m.expand_selection()
			case CROP_MODE:
				m.init_cropping_block()
				m.expand_selection()
			}
		case "?":
			m.help_open = !m.help_open
		case "esc":
			m.Mode = NORMAL_MODE
			m.Reset_to_normal()
			m.help_open = false
		case ":":
			m.Mode = COMMAND_MODE
		}
	}
	return m, nil
}

func (m Model) View() string {
	ascii := Canvas_Window(m.Grid, m.Selected, *m.Cursor, m.Window)

	indicator := Status_Bar(Status{
		Mode:    m.Mode,
		Width:   m.width,
		Cmd:     string(*m.cmd),
		Message: *m.message,
		Mirror:  m.mirror.enabled,
		Axis:    m.mirror.axis,
	})
	centered := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 1).
		Background(theme.SumiInk1).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(ascii)

	modals := []ComponentModal{m.save_modal, m.theme_modal}

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

	screen := lipgloss.JoinVertical(lipgloss.Left, centered, indicator)

	if m.help_open {
		col := m.width - component.HelpWidth - 1
		overlay, err := component.Overlay(screen, component.Help(m.help_page), 1, col, true)
		if err != nil {
			log.Printf("Failed to overlay help: %s", err.Error())
			return screen
		}
		return overlay
	}
	return screen
}

func (m Model) Resize(width, height int) Model {
	m.width = width
	m.height = height
	*m.Window = core.Get_Window(core.Dimensions{Width: width, Height: height - 1}, m.Grid, core.Find_Center(m.Grid))
	return m
}

type Status struct {
	Mode    int
	Width   int
	Cmd     string
	Message string
	Mirror  bool
	Axis    int
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

func Status_Bar(status Status) string {
	var label string
	var msg string
	var color lipgloss.Color

	switch status.Mode {
	case COMMAND_MODE:
		label := status_label(theme.RoninYellow, "COMMAND ")
		cmd := lipgloss.NewStyle().
			Background(theme.SumiInk2).
			Bold(true).
			Render(":" + status.Cmd + "█")
		return lipgloss.NewStyle().
			Background(theme.SumiInk2).
			PaddingLeft(1).
			Width(status.Width).
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
		msg = status.Message
		color = theme.WaveBlue
	}
	var indicators []string
	indicators = append(indicators, status_label(color, label))

	if status.Mirror {
		axis_label := "Y_AXIS"
		if status.Axis == X_AXIS {
			axis_label = "X_AXIS"
		}
		indicators = append(indicators, subtitle("MIRROR::"+axis_label))
	}
	indicators = append(indicators, subtitle(msg))

	return lipgloss.NewStyle().
		Background(theme.SumiInk2).
		PaddingLeft(1).
		Width(status.Width).
		AlignHorizontal(lipgloss.Left).
		Render(lipgloss.JoinHorizontal(lipgloss.Bottom, indicators...))
}
