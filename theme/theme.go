package theme

import (
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/charmbracelet/lipgloss"
)

const (
	SumiInk0    = lipgloss.Color("#0D0C0C")
	SumiInk1    = lipgloss.Color("#16161D")
	SumiInk2    = lipgloss.Color("#2A2A37")
	SumiInk3    = lipgloss.Color("#1F1F28")
	Cursor      = lipgloss.Color("#C8C093")
	SamuraiRed  = lipgloss.Color("#C34043")
	Brown       = lipgloss.Color("#49443C")
	WaveBlue    = lipgloss.Color("#7E9CD8")
	RoninYellow = lipgloss.Color("#DCA561")
	Wisteria    = lipgloss.Color("#957FB8")
)

const (
	CanvasBG     = "\x1b[48;2;22;22;29m"
	CellBG       = "\x1b[48;2;13;12;12m"
	MagnifierGap = "\x1b[48;2;42;42;55m"
	PanelBG      = "\x1b[48;2;31;31;40m"
	CursorAnsi   = "\x1b[48;2;200;192;147;38;2;13;12;12m"
	SelectionBG  = "\x1b[48;2;90;90;110m"
	CropBG       = "\x1b[48;2;140;50;53m"
	ActiveCellFG = "\x1b[38;2;126;156;216m"
	BrownBG      = "\x1b[48;2;62;57;50m"
	WaveBlueBG   = "\x1b[48;2;72;118;214m"
	Reset        = "\x1b[0m"
)

type Background struct {
	Main      lipgloss.Color
	Container lipgloss.Color
	Selection lipgloss.Color
	Border    lipgloss.Color
	Notify    lipgloss.Color
}
type Foreground struct {
	Command lipgloss.Color
	Main    lipgloss.Color
	Subtle  lipgloss.Color
	Title   lipgloss.Color
	Border  lipgloss.Color
}

type Theme struct {
	Label string
	Bg    Background
	Fg    Foreground
}

func BgStyle(theme Theme) *Bg {
	return &Bg{
		Container: lipgloss.NewStyle().Background(theme.Bg.Container),
		Base:      lipgloss.NewStyle().Background(theme.Bg.Main),
		Border:    lipgloss.NewStyle().Background(theme.Bg.Border),
		Selection: lipgloss.NewStyle().Background(theme.Bg.Selection),
		Notify:    lipgloss.NewStyle().Background(theme.Bg.Notify),
	}
}

type Bg struct {
	Container lipgloss.Style
	Base      lipgloss.Style
	Border    lipgloss.Style
	Selection lipgloss.Style
	Notify    lipgloss.Style
}

type Get struct {
	Theme func() *Theme
	Bg    func() *Bg
}

func Selection_Ansi(slt int) string {
	if slt == core.Crop {
		return CropBG
	}
	return SelectionBG
}
