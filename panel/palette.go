package panel

import (
	"github.com/Mr-Robot-err-404/perkins/theme"
	"github.com/charmbracelet/lipgloss"
)

func render_palette(color [8]theme.Color) string {
	square := lipgloss.NewStyle().Width(6).Height(3)

	left := lipgloss.JoinVertical(lipgloss.Left,
		square.Background(color[0].Display).Render(),
		square.Background(color[1].Display).Render(),
		square.Background(color[2].Display).Render(),
		square.Background(color[3].Display).Render(),
	)
	right := lipgloss.JoinVertical(lipgloss.Left,
		square.Background(color[4].Display).Render(),
		square.Background(color[5].Display).Render(),
		square.Background(color[6].Display).Render(),
		square.Background(color[7].Display).Render(),
	)
	gap := lipgloss.NewStyle().Background(theme.SumiInk3).Width(2).Height(12).Render()
	palette := lipgloss.JoinHorizontal(lipgloss.Bottom, left, gap, right)

	return lipgloss.JoinVertical(lipgloss.Left, title("Colour palette", 0), palette)
}
