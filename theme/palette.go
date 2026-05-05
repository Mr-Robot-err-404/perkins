package theme

import "github.com/charmbracelet/lipgloss"

type Palette struct {
	Name       string
	Foreground []Color
	Background []Color
}
type Color struct {
	Ansi    string
	Display lipgloss.Color
}

const PageSize int = 8

func (p Palette) ForegroundPage(page int) []Color {
	return paged_colors(p.Foreground, page)
}

func (p Palette) BackgroundPage(page int) []Color {
	return paged_colors(p.Background, page)
}

func (p Palette) ForegroundPages() int {
	return page_total(len(p.Foreground))
}

func (p Palette) BackgroundPages() int {
	return page_total(len(p.Background))
}

func paged_colors(colors []Color, page int) []Color {
	start := page * PageSize

	if start >= len(colors) {
		return []Color{}
	}
	end := min(start+PageSize, len(colors))
	return colors[start:end]
}

func page_total(n int) int {
	return max(1, (n+PageSize-1)/PageSize)
}

const (
	FujuWhite    = "38;2;220;215;186"
	OldWhite     = "38;2;199;189;143"
	SumiInk4     = "38;2;84;84;109"
	WaveBlue2    = "38;2;126;156;216"
	SpringGreen  = "38;2;118;148;106"
	AutumnYellow = "38;2;196;166;96"
	SamuraiRedFG = "38;2;195;64;67"
	SakuraPink   = "38;2;210;126;153"

	SumiInk0BG = "48;2;13;12;12"
	SumiInk1BG = "48;2;22;22;29"
	SumiInk2BG = "48;2;22;22;40"
	SumiInk3BG = "48;2;31;31;40"
	WinterBlue = "48;2;43;61;86"
	WaveBlue1  = "48;2;59;76;119"
	WisteriaBG = "48;2;151;124;178"
	LotusPink  = "48;2;215;152;166"
)

var Themes = []Palette{Kanagawa, Gruvbox, TokyoNight, Nordic, RosePine}
