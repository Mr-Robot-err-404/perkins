package panel

import (
	"strings"

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

func divider_cell(s *strings.Builder, r rune) {
	for range 3 {
		append_ascii(s, r)
	}
}

func append_ascii(s *strings.Builder, r rune) {
	s.WriteRune(r)
	s.WriteString("\n")
}

func left_divider(idx int, offset int) string {
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
		Foreground(theme.RoninYellow).
		Render(strings.TrimSuffix(s.String(), "\n"))
}

func right_divider(idx int, offset int) string {
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
		Foreground(theme.RoninYellow).
		Render(strings.TrimSuffix(s.String(), "\n"))
}

func x_gap(idx int, selected int) string {
	s := strings.Builder{}
	style := lipgloss.NewStyle().Background(theme.SumiInk3).Foreground(theme.RoninYellow)

	r := ' '
	if idx == selected || idx == selected+1 {
		r = HorizontalBar
	}
	for range 6 {
		s.WriteRune(r)
	}
	return style.Render(s.String())
}

func square() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.SumiInk3).
		Width(6).
		Height(3)
}

func column(offset int, color [8]theme.Color, selected int) []string {
	items := []string{}

	for i := offset; i < 4+offset; i++ {
		items = append(items,
			x_gap(i, selected),
			square().Background(color[i].Display).Render(),
		)
	}
	items = append(items, x_gap(offset+4, selected))
	return items
}

func render_palette(color [8]theme.Color) string {
	selected := 2

	left := lipgloss.JoinVertical(lipgloss.Left, column(0, color, selected)...)
	right := lipgloss.JoinVertical(lipgloss.Left, column(4, color, selected)...)

	palette := lipgloss.JoinHorizontal(lipgloss.Bottom,
		left_divider(selected, 0),
		left,
		right_divider(selected, 0),
		left_divider(selected, 4),
		right,
		right_divider(selected, 4),
	)
	return lipgloss.JoinVertical(lipgloss.Left, title(" Colour palette", Padding{}), palette)
}
