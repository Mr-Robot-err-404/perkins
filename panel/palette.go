package panel

import (
	"fmt"
	"strings"

	"github.com/Mr-Robot-err-404/perkins/component"
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

func get_color(pos core.Pos, colors []theme.Color) theme.Color {
	idx := pos_to_idx[pos]
	if idx >= len(colors) {
		return theme.Color{}
	}
	return colors[idx]
}

func (p *Palette) get_color_palette() []theme.Color {
	t := theme.Themes[p.theme_idx]
	if p.Layer == FOREGROUND_LAYER {
		return t.ForegroundPage(p.page)
	}
	return t.BackgroundPage(p.page)
}

func (p *Palette) next_page() {
	t := theme.Themes[p.theme_idx]
	pages := t.ForegroundPages()

	if p.Layer == BACKGROUND_LAYER {
		pages = t.BackgroundPages()
	}
	p.page = min(p.page+1, pages-1)
}

func (p *Palette) prev_page() {
	p.page = max(p.page-1, 0)
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

		if offset == 4 && i == offset && idx < offset {
			append_ascii(&s, ' ')
		}
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

		// weird edge case that I won't bother trying to polish
		if offset == 4 && i == offset && idx < offset {
			append_ascii(&s, ' ')
		}
		if offset == 0 && i == offset && idx > 3 {
			append_ascii(&s, ' ')
		}
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

func (p *Palette) get_pages(t theme.Palette) int {
	if p.Layer == BACKGROUND_LAYER {
		return t.BackgroundPages()
	}
	return t.ForegroundPages()
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
		return component.Notification("Foreground", PALETTE_WIDTH, 3, theme.RoninYellow, theme.SumiInk1)
	}
	return component.Notification("Background", PALETTE_WIDTH, 3, theme.WaveBlue, theme.SumiInk1)
}

func (p *Palette) column(offset int, color []theme.Color, selected int, column int) []string {
	items := []string{}

	for i := offset; i < 4+offset; i++ {
		items = append(items, p.x_gap(i, selected, column))

		if i >= len(color) {
			items = append(items, square().Render())
			continue
		}
		items = append(items, square().Background(color[i].Display).Render())
	}
	items = append(items, p.x_gap(offset+4, selected, column))
	return items
}

func (p *Palette) page_indicator() string {
	t := theme.Themes[p.theme_idx]
	pages := p.get_pages(t)

	fg := p.get_highlight()
	base := lipgloss.NewStyle().Background(theme.SumiInk3).Foreground(theme.FujiGray)
	active := lipgloss.NewStyle().Background(theme.SumiInk3).Foreground(fg).Bold(true)

	prev := active.Render("<")
	if p.page == 0 {
		prev = base.Render("<")
	}
	next := active.Render(">")
	if p.page >= pages-1 {
		next = base.Render(">")
	}
	num := active.PaddingLeft(1).PaddingRight(1).Render(fmt.Sprintf("%d", p.page+1))

	return lipgloss.NewStyle().
		Width(PALETTE_WIDTH).
		Background(theme.SumiInk3).
		AlignHorizontal(lipgloss.Center).
		Render(lipgloss.JoinHorizontal(lipgloss.Bottom, prev, num, next))
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
	return lipgloss.JoinVertical(lipgloss.Left, layer_state(p.Layer), content, p.page_indicator())
}
