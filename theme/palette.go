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

var Kanagawa = Palette{
	Name: "Kanagawa",
	Foreground: []Color{
		// page 1
		{Ansi: "38;2;220;215;186", Display: "#DCDBB7"}, // fujiWhite
		{Ansi: "38;2;199;189;143", Display: "#C7BD8F"}, // oldWhite
		{Ansi: "38;2;84;84;109", Display: "#54546D"},   // sumiInk4
		{Ansi: "38;2;126;156;216", Display: "#7E9CD8"}, // crystalBlue
		{Ansi: "38;2;118;148;106", Display: "#76946A"}, // autumnGreen
		{Ansi: "38;2;196;166;96", Display: "#C4A660"},  // boatYellow2
		{Ansi: "38;2;195;64;67", Display: "#C34043"},   // autumnRed
		{Ansi: "38;2;210;126;153", Display: "#D27E99"}, // sakuraPink
		// page 2
		{Ansi: "38;2;149;127;184", Display: "#957FB8"}, // oniViolet
		{Ansi: "38;2;127;180;202", Display: "#7FB4CA"}, // springBlue
		{Ansi: "38;2;152;187;108", Display: "#98BB6C"}, // springGreen
		{Ansi: "38;2;230;195;132", Display: "#E6C384"}, // carpYellow
		{Ansi: "38;2;228;104;118", Display: "#E46876"}, // waveRed
		{Ansi: "38;2;255;93;98", Display: "#FF5D62"},   // peachRed
		{Ansi: "38;2;255;160;102", Display: "#FFA066"}, // surimiOrange
		{Ansi: "38;2;106;149;137", Display: "#6A9589"}, // waveAqua1
	},
	Background: []Color{
		// page 1
		{Ansi: "48;2;13;12;12", Display: "#0D0C0C"},    // sumiInk0
		{Ansi: "48;2;22;22;29", Display: "#16161D"},    // sumiInk1
		{Ansi: "48;2;22;22;40", Display: "#161628"},    // sumiInk2
		{Ansi: "48;2;31;31;40", Display: "#1F1F28"},    // sumiInk3
		{Ansi: "48;2;43;61;86", Display: "#2B3D56"},    // winterBlue
		{Ansi: "48;2;59;76;119", Display: "#3B4C77"},   // waveBlue1
		{Ansi: "48;2;151;124;178", Display: "#977CB2"}, // wisteria
		{Ansi: "48;2;215;152;166", Display: "#D798A6"}, // lotusPink
		// page 2
		{Ansi: "48;2;34;50;73", Display: "#223249"},    // waveBlue1 deep
		{Ansi: "48;2;45;79;103", Display: "#2D4F67"},   // waveBlue2
		{Ansi: "48;2;43;51;40", Display: "#2B3328"},    // winterGreen
		{Ansi: "48;2;73;68;60", Display: "#49443C"},    // winterYellow
		{Ansi: "48;2;67;36;43", Display: "#43242B"},    // winterRed
		{Ansi: "48;2;37;37;53", Display: "#252535"},    // winterBlue dark
		{Ansi: "48;2;54;54;70", Display: "#363646"},    // sumiInk3 light
		{Ansi: "48;2;101;133;89", Display: "#658594"},  // dragonBlue
	},
}

var Themes = []Palette{Kanagawa}
