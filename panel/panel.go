package panel

import (
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width   int
	height  int
	palette Palette
	Cell    func() rune
}
type Palette struct {
	Pos   core.Pos
	Layer int
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
)

const (
	FILL_ACTION int = iota + 1
	CLEAR_ACTION
)
const (
	FOREGROUND_LAYER int = iota
	BACKGROUND_LAYER int = iota
)

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
	case tea.KeyMsg:
		switch msg.String() {
		case PALETTE_TOP:
			m.palette.Pos.Row = dec(m.palette.Pos.Row, 0)
		case PALETTE_RIGHT:
			m.palette.Pos.Col = inc(m.palette.Pos.Col, 1)
		case PALETTE_DOWN:
			m.palette.Pos.Row = inc(m.palette.Pos.Row, 3)
		case PALETTE_LEFT:
			m.palette.Pos.Col = dec(m.palette.Pos.Col, 0)

		case APPLY_COLOR:
			return m, ApplyColor(m.palette.Layer, get_color(m.palette.Pos, theme.Kanagawa.Foreground))

		case "1", "q":
			return m, FlipBit(0)
		case "2", "w":
			return m, FlipBit(1)
		case "3", "e":
			return m, FlipBit(2)
		case "4", "r":
			return m, FlipBit(6)
		case "5", "u":
			return m, FlipBit(3)
		case "6", "i":
			return m, FlipBit(4)
		case "7", "o":
			return m, FlipBit(5)
		case "8", "p":
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
		render_palette(theme.Kanagawa.Foreground, m.palette.Pos),
	)
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
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

func (m Model) magnify() string {
	r := m.Cell()
	bits := magnifier(r)
	return lipgloss.JoinVertical(lipgloss.Left, title("Magnifier", Padding{Right: 5, Bottom: 1}), render_magnifier(bits))
}

func (m Model) Resize(width, height int) Model {
	m.width = width
	m.height = height
	return m
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
