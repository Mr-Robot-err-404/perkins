package theme

import "github.com/charmbracelet/lipgloss"

type Palette struct {
	Foreground [8]Color
	Background [8]Color
}
type Color struct {
	Ansi    string
	Display lipgloss.Color
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
	Wisteria   = "48;2;151;124;178"
	LotusPink  = "48;2;215;152;166"
)

var Kanagawa = Palette{
	Foreground: [8]Color{
		{Ansi: "38;2;220;215;186", Display: "#DCDBB7"},
		{Ansi: "38;2;199;189;143", Display: "#C7BD8F"},
		{Ansi: "38;2;84;84;109", Display: "#54546D"},
		{Ansi: "38;2;126;156;216", Display: "#7E9CD8"},
		{Ansi: "38;2;118;148;106", Display: "#76946A"},
		{Ansi: "38;2;196;166;96", Display: "#C4A660"},
		{Ansi: "38;2;195;64;67", Display: "#C34043"},
		{Ansi: "38;2;210;126;153", Display: "#D27E99"},
	},
	Background: [8]Color{
		{Ansi: "48;2;13;12;12", Display: "#0D0C0C"},
		{Ansi: "48;2;22;22;29", Display: "#16161D"},
		{Ansi: "48;2;22;22;40", Display: "#161628"},
		{Ansi: "48;2;31;31;40", Display: "#1F1F28"},
		{Ansi: "48;2;43;61;86", Display: "#2B3D56"},
		{Ansi: "48;2;59;76;119", Display: "#3B4C77"},
		{Ansi: "48;2;151;124;178", Display: "#977CB2"},
		{Ansi: "48;2;215;152;166", Display: "#D798A6"},
	},
}
