package main

import (
	"github.com/Mr-Robot-err-404/perkins/canvas"
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/panel"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const panelWidth int = 42

type model struct {
	width    int
	height   int
	panel    panel.Model
	canvas   canvas.Model
	grid     canvas.Grid
	selected core.Selected
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.canvas.Init(), m.panel.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		canvasWidth := max(0, m.width-panelWidth)
		m.canvas = m.canvas.Resize(canvasWidth, m.height)
		m.panel = m.panel.Resize(panelWidth, m.height)
		return m, nil
	case panel.FlipMsg:
		pos := m.canvas.Cursor
		r := m.grid[pos.Row][pos.Col]

		if !core.Is_Braille(r) {
			return m, nil
		}
		b := core.Bitmap(r)
		b = core.Flip(b, msg.Bit)
		m.grid[pos.Row][pos.Col] = core.Bitmap_To_Braille(b)

	case panel.FillMsg:
		pos := m.canvas.Cursor
		r := m.grid[pos.Row][pos.Col]

		if !core.Is_Braille(r) {
			return m, nil
		}
		m.grid[pos.Row][pos.Col] = core.Full

	case panel.ClearMsg:
		pos := m.canvas.Cursor
		r := m.grid[pos.Row][pos.Col]

		if !core.Is_Braille(r) {
			return m, nil
		}
		m.grid[pos.Row][pos.Col] = core.Base
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

func (m model) GetCell() rune {
	return m.grid[m.canvas.Cursor.Row][m.canvas.Cursor.Col]
}

func newModel(grid canvas.Grid) model {
	selected := make(core.Selected)
	m := model{
		canvas:   canvas.New(0, 0, grid, selected),
		grid:     grid,
		selected: selected,
	}
	m.panel = panel.New(panelWidth, 0, m.GetCell)
	return m
}
