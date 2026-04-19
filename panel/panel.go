package panel

import (
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	palette  Palette
	terminal Dimensions
	panel    Dimensions
	coords   Coords
	toggle   map[core.Pos]bool
	Cell     func() rune
	Offset   func() int
}
type Palette struct {
	Layer  int
	fg_pos *core.Pos
	bg_pos *core.Pos
}
type Dimensions struct {
	Width  int
	Height int
}

type FlipMsg struct{ Bit byte }
type ActionMsg struct{ Action int }
type ColorMsg = struct {
	Layer int
	Color theme.Color
}

type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

const (
	PALETTE_RIGHT string = "L"
	PALETTE_LEFT  string = "H"
	PALETTE_TOP   string = "K"
	PALETTE_DOWN  string = "J"
	APPLY_COLOR   string = " "
	SWITCH_LAYER  string = "N"
)

const (
	FILL_ACTION int = iota + 1
	CLEAR_ACTION
)
const (
	FOREGROUND_LAYER int = iota
	BACKGROUND_LAYER int = iota
)

func New(dm Dimensions, get func() rune) Model {
	return Model{
		panel:   dm,
		Cell:    get,
		palette: Palette{fg_pos: &core.Pos{}, bg_pos: &core.Pos{}},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func FlipBit(b byte) tea.Cmd {
	return func() tea.Msg {
		return FlipMsg{Bit: b}
	}
}

func ApplyColor(layer int, color theme.Color) tea.Cmd {
	return func() tea.Msg {
		return ColorMsg{Layer: layer, Color: color}
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
	case tea.MouseMsg:
		// debug.Logf("%d:%d\n", msg.X, msg.Y)
		mouse := core.Pos{Row: msg.Y, Col: msg.X}
		pos := m.palette.get_palette_pos()

		switch msg.Action {
		case tea.MouseActionPress:
			switch msg.Button {
			case tea.MouseButtonLeft:
				square, ok := m.coords[mouse]
				if ok {
					*pos = square
				}
				_, ok = m.toggle[mouse]
				if ok {
					m.palette.Layer = m.toggle_layer()
				}
			}
		case tea.MouseActionMotion:
			if msg.Button == tea.MouseButtonLeft {
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case PALETTE_TOP:
			pos := m.palette.get_palette_pos()
			pos.Row = dec(pos.Row, 0)
		case PALETTE_RIGHT:
			pos := m.palette.get_palette_pos()
			pos.Col = inc(pos.Col, 1)
		case PALETTE_DOWN:
			pos := m.palette.get_palette_pos()
			pos.Row = inc(pos.Row, 3)
		case PALETTE_LEFT:
			pos := m.palette.get_palette_pos()
			pos.Col = dec(pos.Col, 0)

		case SWITCH_LAYER:
			m.palette.Layer = m.toggle_layer()

		case APPLY_COLOR:
			pos := m.palette.get_palette_pos()
			colors := m.palette.get_color_palette()
			return m, ApplyColor(m.palette.Layer, get_color(*pos, colors))

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

		case "x":
			return m, Clear()
		case "f":
			return m, Fill()
		}
	}
	return m, nil
}

func (m Model) View() string {
	divider := lipgloss.NewStyle().Height(3).Render()
	content := lipgloss.JoinVertical(lipgloss.Left,
		m.magnify(),
		divider,
		m.palette.render_palette(),
	)
	return lipgloss.NewStyle().
		Width(m.panel.Width).
		Height(m.panel.Height).
		Background(theme.SumiInk3).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)
}

func title(s string, padding Padding) string {
	return lipgloss.NewStyle().
		Background(theme.SumiInk3).
		Foreground(theme.WaveBlue).
		PaddingBottom(padding.Bottom).
		PaddingRight(padding.Right).
		Render(s)
}
func (m Model) toggle_layer() int {
	if m.palette.Layer == FOREGROUND_LAYER {
		return BACKGROUND_LAYER
	}
	return FOREGROUND_LAYER
}

func (m Model) magnify() string {
	r := m.Cell()
	bits := magnifier(r)
	return lipgloss.JoinVertical(lipgloss.Left, title("Magnifier", Padding{Right: 7, Bottom: 1}), render_magnifier(bits))
}

func (m Model) Resize(panel Dimensions, terminal Dimensions) Model {
	m.panel = panel
	m.terminal = terminal
	offset := m.palette_start()
	m.coords = palette_coords(offset.Col, offset.Row)

	offset = m.toggle_start()
	m.toggle = toggle_coords(offset.Col, offset.Row)
	return m
}

func (p *Palette) get_palette_pos() *core.Pos {
	if p.Layer == FOREGROUND_LAYER {
		return p.fg_pos
	}
	return p.bg_pos
}

func inc(n int, cap int) int {
	if n >= cap {
		return cap
	}
	return n + 1
}
func dec(n int, floor int) int {
	if n <= floor {
		return floor
	}
	return n - 1
}
