package main

import (
	"github.com/Mr-Robot-err-404/perkins/canvas"
	"github.com/Mr-Robot-err-404/perkins/panel"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const panelWidth int = 42

type model struct {
	width  int
	height int
	panel  panel.Model
	canvas canvas.Model
	grid   canvas.Grid
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.canvas.Init(), m.panel.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		canvasWidth := max(0, m.width-panelWidth)
		m.canvas = m.canvas.Resize(canvasWidth, m.height)
		m.panel = m.panel.Resize(panelWidth, m.height)
		return m, nil
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.canvas, cmd = m.canvas.Update(msg)
	cmds = append(cmds, cmd)

	m.panel, cmd = m.panel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, m.canvas.View(), m.panel.View())
}

func newModel(grid canvas.Grid) model {
	return model{
		canvas: canvas.New(0, 0, grid),
		panel:  panel.New(panelWidth, 0),
		grid:   grid,
	}
}
