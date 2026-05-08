package component

import (
	"github.com/Mr-Robot-err-404/perkins/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Selection struct {
	List []string
	Idx  int
}

type SltModalConfig struct {
	Width    int
	XPadding int
	Prompt   string
	Title    string
	Id       int
	IsDelete bool
}

type SltModalSubmit struct {
	Value  string
	Cancel bool
	Id     int
	Idx    int
}
type SltModalDelete struct {
	Value string
	Id    int
	Idx   int
}

type SltModalToggle struct {
	Id  int
	Idx int
}

type SltModal struct {
	Config  SltModalConfig
	Active  bool
	Current int
	Slt     Selection
	prev    *string
}

const Circle string = "\u29BF"

func (m SltModal) IsActive() bool {
	return m.Active
}

var container = lipgloss.NewStyle().Background(theme.SumiInk0)

func (m SltModal) info(key string, action string) string {
	primary := container.Foreground(theme.Cursor).PaddingRight(1).Render(key)
	secondary := container.Foreground(theme.FujiGray).Render(action)
	return lipgloss.JoinHorizontal(lipgloss.Bottom, primary, secondary)
}

func (m SltModal) circle() string {
	return lipgloss.NewStyle().Foreground(theme.Cursor).Render(Circle)
}

func (m SltModal) infoList() string {
	s := m.info("enter", "select")
	t := m.info("esc", "cancel")
	bg := theme.SumiInk0
	return container.
		PaddingTop(1).
		Render(JustifyBetween(Justify{
			Left:  s,
			Right: t,
			Width: m.Config.Width - 2,
			Bg:    &bg,
		}))
}

func (m SltModal) listItemStyle(s string, idx int, width int) string {
	style := container.Foreground(theme.FujiGray).Width(width)

	if idx == m.Slt.Idx {
		style = style.Background(theme.SumiInk2).Foreground(theme.Cursor)

		if m.Config.IsDelete && *m.prev == "x" {
			style = style.Background(theme.SamuraiRed)
			s = "Press x again to delete"
		}
	}
	if idx == m.Current {
		str := JustifyBetween(Justify{
			Left:  s,
			Right: m.circle(),
			Width: width - 1,
			Bg:    nil,
		})
		return style.Bold(true).Render(str)
	}
	return style.Render(s)
}

func (m SltModal) list() string {
	var items []string
	for i, current := range m.Slt.List {
		width := m.Config.Width
		str := Truncate(current, width)

		item := m.listItemStyle(str, i, width-1)
		items = append(items, item)
	}
	return container.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

func (m SltModal) title() string {
	style := container.Foreground(theme.FujiGray)

	if !m.Config.IsDelete {
		return style.Width(m.Config.Width).Render(m.Config.Title)
	}
	bg := theme.SumiInk0
	return container.
		PaddingBottom(1).
		Render(JustifyBetween(Justify{
			Left:  style.Render(m.Config.Title),
			Right: m.info("x", "delete"),
			Width: m.Config.Width,
			Bg:    &bg,
		}))
}

func (m SltModal) Update(msg tea.Msg) (SltModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()

		if s == "x" && m.Config.IsDelete && *m.prev == "x" {
			*m.prev = ""
			v := m.Slt.List[m.Slt.Idx]
			return m, func() tea.Msg { return SltModalDelete{Value: v, Idx: m.Slt.Idx, Id: m.Config.Id} }
		}
		defer func() {
			*m.prev = s
		}()
		switch s {
		case "esc":
			if *m.prev == "x" {
				return m, nil
			}
			m.Active = false
			return m, func() tea.Msg { return SltModalSubmit{Cancel: true, Idx: m.Current, Id: m.Config.Id} }
		case "enter", "ctrl+y":
			m.Active = false
			v := m.Slt.List[m.Slt.Idx]
			m.Current = m.Slt.Idx

			return m, func() tea.Msg {
				return SltModalSubmit{Value: v, Id: m.Config.Id, Idx: m.Slt.Idx}
			}
		case "ctrl+p", "up", "k":
			idx := max(0, m.Slt.Idx-1)
			m.Slt.Idx = idx

			return m, func() tea.Msg {
				return SltModalToggle{Id: m.Config.Id, Idx: idx}
			}

		case "ctrl+n", "down", "j":
			idx := min(len(m.Slt.List)-1, m.Slt.Idx+1)
			m.Slt.Idx = idx

			return m, func() tea.Msg {
				return SltModalToggle{Id: m.Config.Id, Idx: idx}
			}
		}
	}
	return m, nil
}

func (m SltModal) View() string {
	box := container.
		BorderBackground(theme.SumiInk0).
		BorderForeground(theme.Cursor).
		Border(lipgloss.NormalBorder()).
		PaddingLeft(m.Config.XPadding).
		Width(m.Config.Width)

	content := lipgloss.JoinVertical(lipgloss.Left, m.title(), m.list(), m.infoList())
	return box.Render(content)
}

func NewSltModal(config SltModalConfig) SltModal {
	return SltModal{Config: config, prev: new(string)}
}
