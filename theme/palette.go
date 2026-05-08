package theme

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

type Palette struct {
	Name       string
	Foreground []Color
	Background []Color
}

type Color struct {
	Hex     string
	Display lipgloss.Color
}

func NewColor(hex string) Color {
	return Color{Hex: hex, Display: lipgloss.Color(hex)}
}

func (c Color) FG() string {
	r, g, b := hexToRGB(c.Hex)
	return fmt.Sprintf("38;2;%d;%d;%d", r, g, b)
}

func (c Color) BG() string {
	r, g, b := hexToRGB(c.Hex)
	return fmt.Sprintf("48;2;%d;%d;%d", r, g, b)
}

func hexToRGB(hex string) (uint8, uint8, uint8) {
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	v, _ := strconv.ParseUint(hex, 16, 32)
	return uint8(v >> 16), uint8(v >> 8), uint8(v)
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

type paletteData struct {
	fg []string
	bg []string
}

func buildPalette(name string, data paletteData) Palette {
	fg := make([]Color, len(data.fg))
	bg := make([]Color, len(data.bg))
	for i, hex := range data.fg {
		fg[i] = NewColor(hex)
	}
	for i, hex := range data.bg {
		bg[i] = NewColor(hex)
	}
	return Palette{Name: name, Foreground: fg, Background: bg}
}

var Themes = []Palette{Kanagawa, Gruvbox, TokyoNight, Nordic, RosePine}
