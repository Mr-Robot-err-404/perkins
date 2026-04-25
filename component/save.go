package component

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/theme"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModalConfig struct {
	Width       int
	XPadding    int
	Prompt      string
	Title       string
	Placeholder string
	CharLimit   int
	Id          int
}

type ModalSubmit struct {
	Value  string
	Cancel bool
	Id     int
}

type Modal struct {
	Active bool
	Config ModalConfig
	Input  textinput.Model
	Get    theme.Get
}

func (m Modal) IsActive() bool {
	return m.Active
}

func (m Modal) info(key string, action string) string {
	theme := m.Get.Theme()
	bg := m.Get.Bg()

	primary := bg.Container.Foreground(theme.Fg.Title).PaddingRight(1).Render(key)
	secondary := bg.Container.Foreground(theme.Fg.Subtle).PaddingRight(2).Render(action)
	return lipgloss.JoinHorizontal(lipgloss.Bottom, primary, secondary)
}

func (m Modal) infoList() string {
	theme := m.Get.Theme()
	bg := m.Get.Bg()

	s := m.info("enter", "save")
	t := m.info("esc", "cancel")
	return bg.Container.
		PaddingTop(1).
		Render(JustifyBetween(Justify{
			Left:  s,
			Right: t,
			Width: m.Config.Width,
			Bg:    &theme.Bg.Container,
		}))
}

func (m Modal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Active = false
			return m, func() tea.Msg { return ModalSubmit{Cancel: true} }
		case "enter":
			v := m.Input.Value()
			m.Input.SetValue("")
			m.Active = false

			return m, func() tea.Msg {
				return ModalSubmit{Value: strings.TrimSpace(v), Id: m.Config.Id}
			}
		}
	}
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m Modal) View() string {
	theme := m.Get.Theme()
	bg := m.Get.Bg()

	box := bg.Container.
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.Fg.Main).
		BorderBackground(theme.Bg.Container).
		Padding(0, m.Config.XPadding).
		Width(m.Config.Width)

	title := bg.Container.Foreground(theme.Fg.Title).PaddingBottom(1).Render(m.Config.Title)
	content := lipgloss.JoinVertical(lipgloss.Left, title, m.Input.View(), m.infoList())

	return box.Render(content)
}

func NewModal(config ModalConfig, get theme.Get) Modal {
	theme := get.Theme()
	bg := get.Bg()

	in := textinput.New()
	in.Placeholder = config.Placeholder
	in.Width = config.Width - 5
	in.Prompt = config.Prompt
	in.TextStyle = bg.Container.Foreground(theme.Fg.Main)
	in.PlaceholderStyle = bg.Container.Foreground(theme.Fg.Subtle)

	if config.CharLimit > 0 {
		in.CharLimit = config.CharLimit
	}
	in.Cursor.SetMode(cursor.CursorStatic)
	in.Focus()

	return Modal{Config: config, Input: in, Get: get}
}
