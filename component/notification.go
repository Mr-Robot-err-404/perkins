package component

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	RightBar rune = 0x2595
	LeftBar  rune = 0x258F
)

func Notification(content string, w int, fg lipgloss.Color, bg lipgloss.Color) string {
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
		Render(content)
	center := lipgloss.JoinVertical(lipgloss.Left, empty, msg, empty)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftBar, center, rightBar)
}
