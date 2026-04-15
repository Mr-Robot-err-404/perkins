package panel

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
	"github.com/charmbracelet/lipgloss"
)

const (
	RightBar        rune = 0x2595
	LeftBar         rune = 0x258F
	CenterBar       rune = 0x2502
	HorizontalBar   rune = 0x2500
	TopJoinRight    rune = 0x2514
	BottomJoinRight rune = 0x250C
	BottomJoinLeft  rune = 0x2510
	TopJoinLeft     rune = 0x2518
)

var pos_to_idx = map[core.Pos]int{
	{Row: 0, Col: 0}: 0,
	{Row: 1, Col: 0}: 1,
	{Row: 2, Col: 0}: 2,
	{Row: 3, Col: 0}: 3,
	{Row: 0, Col: 1}: 4,
	{Row: 1, Col: 1}: 5,
	{Row: 2, Col: 1}: 6,
	{Row: 3, Col: 1}: 7,
}

func get_color(pos core.Pos, colors [8]theme.Color) theme.Color {
	idx := pos_to_idx[pos]
	return colors[idx]
}

func (p *Palette) get_color_palette() [8]theme.Color {
	if p.Layer == FOREGROUND_LAYER {
		return theme.Kanagawa.Foreground
	}
	return theme.Kanagawa.Background
}

func divider_cell(s *strings.Builder, r rune) {
	for range 3 {
		append_ascii(s, r)
	}
}

func append_ascii(s *strings.Builder, r rune) {
	s.WriteRune(r)
	s.WriteString("\n")
}

func (p *Palette) left_divider(idx int, offset int) string {
	fg := p.get_highlight()
	s := strings.Builder{}

	for i := offset; i < 4+offset; i++ {
		if idx == i {
			append_ascii(&s, BottomJoinRight)
			divider_cell(&s, CenterBar)
			append_ascii(&s, TopJoinRight)
			continue
		}
		divider_cell(&s, ' ')
		s.WriteString(" \n")
	}
	return lipgloss.NewStyle().
		Background(theme.SumiInk3).
		Foreground(fg).
		Render(strings.TrimSuffix(s.String(), "\n"))
}

func (p *Palette) right_divider(idx int, offset int) string {
	fg := p.get_highlight()
	s := strings.Builder{}

	for i := offset; i < 4+offset; i++ {
		if idx == i {
			append_ascii(&s, BottomJoinLeft)
			divider_cell(&s, CenterBar)
			append_ascii(&s, TopJoinLeft)
			continue
		}
		divider_cell(&s, ' ')
		s.WriteString(" \n")
	}
	return lipgloss.NewStyle().
		Background(theme.SumiInk3).
		Foreground(fg).
		Render(strings.TrimSuffix(s.String(), "\n"))
}
func (p *Palette) get_highlight() lipgloss.Color {
	if p.Layer == FOREGROUND_LAYER {
		return theme.RoninYellow
	}
	return theme.WaveBlue
}

func (p *Palette) x_gap(idx int, selected int, column int) string {
	fg := p.get_highlight()
	s := strings.Builder{}
	style := lipgloss.NewStyle().Background(theme.SumiInk3).Foreground(fg)

	r := x_gap_rune(idx, selected, column)
	for range 6 {
		s.WriteRune(r)
	}
	return style.Render(s.String())
}

func x_gap_rune(idx int, selected int, column int) rune {
	r := ' '
	if column == 0 && selected == 4 {
		return r
	}
	if column == 1 && selected == 3 {
		return r
	}
	if idx == selected || idx == selected+1 {
		r = HorizontalBar
	}
	return r
}

func square() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.SumiInk3).
		Width(6).
		Height(3)
}

func layer_state(layer int) string {
	if layer == FOREGROUND_LAYER {
		return notification("Foreground", 16, theme.RoninYellow, theme.SumiInk1)
	}
	return notification("Background", 16, theme.WaveBlue, theme.SumiInk1)
}

func (p *Palette) column(offset int, color [8]theme.Color, selected int, column int) []string {
	items := []string{}

	for i := offset; i < 4+offset; i++ {
		items = append(items,
			p.x_gap(i, selected, column),
			square().Background(color[i].Display).Render(),
		)
	}
	items = append(items, p.x_gap(offset+4, selected, column))
	return items
}

func (p *Palette) render_palette() string {
	color := p.get_color_palette()
	pos := p.get_palette_pos()
	selected := pos_to_idx[*pos]

	left := lipgloss.JoinVertical(lipgloss.Left, p.column(0, color, selected, 0)...)
	right := lipgloss.JoinVertical(lipgloss.Left, p.column(4, color, selected, 1)...)

	content := lipgloss.JoinHorizontal(lipgloss.Bottom,
		p.left_divider(selected, 0),
		left,
		p.right_divider(selected, 0),
		p.left_divider(selected, 4),
		right,
		p.right_divider(selected, 4),
	)
	return lipgloss.JoinVertical(lipgloss.Left, layer_state(p.Layer), content)
}

func notification(s string, w int, fg lipgloss.Color, bg lipgloss.Color) string {
	bar := lipgloss.NewStyle().Foreground(fg).Background(bg)
	l := strings.Repeat(string(LeftBar)+"\n", 3)
	r := strings.Repeat(string(RightBar)+"\n", 3)

	leftBar := bar.Height(3).Render(strings.TrimSuffix(l, "\n"))
	rightBar := bar.Height(3).Render(strings.TrimSuffix(r, "\n"))

	empty := lipgloss.NewStyle().
		Background(bg).
		Width(w - 2).
		Height(1).
		Render("")
	msg := lipgloss.NewStyle().
		Background(bg).
		Foreground(fg).
		Width(w - 2).
		AlignHorizontal(lipgloss.Center).
		Render(s)
	center := lipgloss.JoinVertical(lipgloss.Left, empty, msg, empty)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftBar, center, rightBar)
}
