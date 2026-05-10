package component

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/theme"
	"github.com/charmbracelet/lipgloss"
)

const HelpWidth int = 49

const (
	HelpModes int = iota
	HelpActions
	HelpCommands
	HelpNavigate
)

type Keypair = [2]string

var helpTabs = []string{"[1] MODE", "[2] ACTION", "[3] COMMAND", "[4] NAVIGATE"}

var pairs = [][]Keypair{
	{
		{"v / ctrl+v", "visual block"},
		{"V", "visual line"},
		{"esc", "normal mode"},
		{"d", "draw mode"},
		{"c", "crop mode"},
		{"enter", "confirm crop"},
	},
	{
		{"x / f", "clear / fill"},
		{"u / ctrl+r", "undo / redo"},
		{"Q/W/E/R", "flip left bits"},
		{"U/I/O/P", "flip right bits"},
		{"space", "apply color"},
		{"m", "toggle mirror"},
		{"|-> tab", "toggle axis"},
		{"N", "switch fg/bg"},
		{"< / >", "palette page"},
		{"?", "toggle help"},
	},
	{
		{":w / :write", "save file"},
		{":t / :theme", "palette theme"},
		{":q / :quit", "quit"},
		{":h / :help", "help"},
	},
	{
		{"h/j/k/l", "move cursor"},
		{"H/J/K/L", "navigate palette"},
		{"ctrl+d/u", "jump 10 rows"},
		{"w/b", "jump 10 cols"},
		{"_", "start of row"},
		{"$", "end of row"},
		{"G", "last row"},
		{"t", "center canvas"},
	},
}

var helpPages = makePages(pairs)

func makePages(pairs [][]Keypair) [][]string {
	pages := make([][]string, len(pairs))

	for i, list := range pairs {
		current := make([]string, len(list))
		_, m2 := column_widths(list)

		for j, keys := range list {
			k1, k2 := keys[0], keys[1]
			space, pad := calc_space(m2, len(k1), len(k2))
			current[j] = info(keys[0], keys[1], space, pad)
		}
		pages[i] = current
	}
	return pages
}

func Help(page int) string {
	tabs := helpTabBar(page)
	rows := helpPages[page]
	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{tabs}, rows...)...)
	h := len(rows) + 4
	return Notification(content, HelpWidth, h, theme.WaveBlue, theme.SumiInk0)
}

func helpTabBar(active int) string {
	bg := lipgloss.NewStyle().Background(theme.SumiInk0)
	var tabs []string
	for i, label := range helpTabs {
		pad := 2
		if i == len(helpTabs)-1 {
			pad = 0
		}
		if i == active {
			tabs = append(tabs, bg.Foreground(theme.RoninYellow).Bold(true).PaddingRight(pad).Render(label))
		} else {
			tabs = append(tabs, bg.Foreground(theme.FujiGray).PaddingRight(pad).Render(label))
		}
	}
	return lipgloss.NewStyle().
		PaddingBottom(1).
		Render(lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...))
}

func calc_space(m2, col1, col2 int) (int, int) {
	w := HelpWidth - 3
	space := w - col1 - m2
	pad := m2 - col2
	return space, pad
}

func column_widths(list [][2]string) (int, int) {
	m1, m2 := 0, 0
	for _, s := range list {
		if len(s[0]) > m1 {
			m1 = len(s[0])
		}
		if len(s[1]) > m2 {
			m2 = len(s[1])
		}
	}
	return m1, m2
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
