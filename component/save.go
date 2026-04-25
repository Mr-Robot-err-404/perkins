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
}

var containerBg = lipgloss.NewStyle().Background(theme.SumiInk0)

func (m Modal) IsActive() bool {
	return m.Active
}
func (m Modal) info(key string, action string) string {
	primary := containerBg.Foreground(theme.Cursor).PaddingRight(1).Render(key)
	secondary := containerBg.Foreground(theme.FujiGray).Render(action)
	return lipgloss.JoinHorizontal(lipgloss.Bottom, primary, secondary)
}

func (m Modal) infoList() string {
	s := m.info("enter", "save")
	t := m.info("esc", "cancel")
	bg := theme.SumiInk0
	return containerBg.
		Render(JustifyBetween(Justify{
			Left:  s,
			Right: t,
			Width: m.Config.Width - 2,
			Bg:    &bg,
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
	box := containerBg.
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.Cursor).
		BorderBackground(theme.SumiInk0).
		PaddingLeft(m.Config.XPadding).
		Width(m.Config.Width)

	inputBg := lipgloss.NewStyle().Background(theme.SumiInk3).Padding(1, 1, 1, 1)
	input := inputBg.Render(m.Input.View())

	title := containerBg.Width(m.Config.Width - 2).Foreground(theme.RoninYellow).PaddingBottom(1).Render(m.Config.Title)
	content := lipgloss.JoinVertical(lipgloss.Left, title, input, " ", m.infoList())
	return box.Render(content)
}
func NewModal(config ModalConfig, initial string) Modal {
	in := textinput.New()
	in.Width = config.Width - 5
	in.SetValue(initial)
	in.Prompt = config.Prompt
	in.TextStyle = lipgloss.NewStyle().Background(theme.SumiInk3).Foreground(theme.Cursor)

	if config.CharLimit > 0 {
		in.CharLimit = config.CharLimit
	}
	in.Cursor.SetMode(cursor.CursorStatic)
	in.Focus()
	return Modal{Config: config, Input: in}
}
