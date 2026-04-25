package component

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Justify struct {
	Left   string
	Right  string
	Width  int
	Bg     *lipgloss.Color
	Offset int
}

func gapStyle(bg *lipgloss.Color) string {
	if bg == nil {
		return " "
	}
	return lipgloss.NewStyle().Background(*bg).Render(" ")
}

func JustifyBetween(j Justify) string {
	lw := lipgloss.Width(j.Left)
	rw := lipgloss.Width(j.Right)
	n := j.Width - lw - rw
	gap := max(n, 0) + j.Offset
	return lipgloss.JoinHorizontal(lipgloss.Bottom, j.Left, strings.Repeat(gapStyle(j.Bg), gap), j.Right)
}
