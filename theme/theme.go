package theme

import "github.com/charmbracelet/lipgloss"

const (
	SumiInk0 = lipgloss.Color("#0D0C0C")
	SumiInk1 = lipgloss.Color("#16161D")
	SumiInk3 = lipgloss.Color("#1F1F28")
	Cursor   = lipgloss.Color("#C8C093")
)

const (
	CanvasBG     = "\x1b[48;2;22;22;29m"
	CellBG       = "\x1b[48;2;13;12;12m"
	PanelBG      = "\x1b[48;2;31;31;40m"
	CursorAnsi   = "\x1b[48;2;200;192;147;38;2;13;12;12m"
	SelectionBG  = "\x1b[48;2;90;90;110m"
	ActiveCellFG = "\x1b[38;2;126;156;216m"
	Reset        = "\x1b[0m"
)
