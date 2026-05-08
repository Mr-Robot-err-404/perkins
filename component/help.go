package component

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/theme"
	"github.com/charmbracelet/lipgloss"
)

const HelpWidth int = 46

const (
	HelpModes int = iota
	HelpActions
	HelpCommands
	HelpNavigate
)

var helpTabs = []string{"1 MODES", "2 ACTIONS", "3 COMMANDS", "4 NAVIGATE"}

var helpPages = [][]string{
	{
		info("v / ctrl+v", "visual block", 21, 0),
		info("V", "visual line", 30, 1),
		info("d", "draw mode", 30, 3),
		info("c", "crop mode", 30, 3),
		info("ctrl+y", "confirm crop", 25, 0),
	},
	{
		info("u / ctrl+r", "undo / redo", 2, 0),
		info("space", "apply color", 2, 0),
		info("N", "switch fg/bg layer", 2, 0),
		info("> / <", "palette page", 2, 0),
		info("x / f", "clear / fill", 2, 0),
		info("m", "toggle mirror", 2, 0),
		info("tab", "toggle mirror axis", 2, 0),
	},
	{
		info(":w / :write", "save file", 2, 0),
		info(":theme", "select palette theme", 2, 0),
		info(":q / :quit", "quit", 2, 0),
	},
	{
		info("h/j/k/l", "move cursor", 2, 0),
		info("ctrl+d/u", "jump 10 rows", 2, 0),
		info("w/b", "jump 10 cols", 2, 0),
		info("G", "last row", 2, 0),
		info("$/_", "end/start of row", 2, 0),
		info("t", "center canvas", 2, 0),
	},
}

func Help(page int) string {
	tabs := helpTabBar(page)
	rows := helpPages[page]
	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{tabs}, rows...)...)
	h := len(rows) + 2
	return Notification(content, HelpWidth, h, theme.WaveBlue, theme.SumiInk0)
}

func helpTabBar(active int) string {
	bg := lipgloss.NewStyle().Background(theme.SumiInk0)
	var tabs []string
	for i, label := range helpTabs {
		if i == active {
			tabs = append(tabs, bg.Foreground(theme.WaveBlue).Bold(true).PaddingRight(2).Render(label))
		} else {
			tabs = append(tabs, bg.Foreground(theme.FujiGray).PaddingRight(2).Render(label))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)
}

func info(key string, value string, space int, pad int) string {
	primary := lipgloss.NewStyle().
		Background(theme.SumiInk0).
		Foreground(theme.Cursor).
		PaddingRight(1).
		Render(key)
	gap := lipgloss.NewStyle().
		Background(theme.SumiInk0).
		Render(strings.Repeat(" ", space))
	secondary := lipgloss.NewStyle().
		Background(theme.SumiInk0).
		Foreground(theme.FujiGray).
		PaddingRight(pad).
		Render(value)
	return lipgloss.JoinHorizontal(lipgloss.Bottom, primary, gap, secondary)
}
