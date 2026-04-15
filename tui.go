package main

import (
	"github.com/Mr-Robot-err-404/perkins/canvas"
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/panel"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width    int
	height   int
	panel    panel.Model
	canvas   canvas.Model
	grid     core.Grid
	selected core.Selected
}

const PANEL_WIDTH int = 42

func (m model) apply_action(action int) {
	switch action {
	case panel.FILL_ACTION:
		for pos := range m.selected {
			set_grid_cell(pos, m.grid, core.Full)
		}
	case panel.CLEAR_ACTION:
		for pos := range m.selected {
			set_grid_cell(pos, m.grid, core.Base)
		}
	}
}
func (m model) apply_colours(msg panel.ColorMsg) {
	for pos := range m.selected {
		switch msg.Layer {
		case panel.FOREGROUND_LAYER:
			m.grid[pos.Row][pos.Col].Fg = msg.Color.Ansi
		case panel.BACKGROUND_LAYER:
			m.grid[pos.Row][pos.Col].Bg = msg.Color.Ansi
		}
	}
}

func set_grid_cell(pos core.Pos, grid core.Grid, value rune) {
	cell := grid[pos.Row][pos.Col]
	if !core.Is_Braille(cell.Value) {
		return
	}
	grid[pos.Row][pos.Col].Value = value
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

		canvas_width := max(0, m.width-PANEL_WIDTH)
		m.canvas = m.canvas.Resize(canvas_width, m.height)
		m.panel = m.panel.Resize(PANEL_WIDTH, m.height)
		return m, nil

	case canvas.CropMsg:
		m.grid = msg.Grid
		m.canvas.Grid = msg.Grid
		m.canvas.Mode = canvas.NORMAL_MODE
		m.canvas.Reset_to_normal()

	case panel.FlipMsg:
		pos := m.canvas.Cursor
		cell := m.grid[pos.Row][pos.Col]

		if !core.Is_Braille(cell.Value) {
			return m, nil
		}
		b := core.Bitmap(cell.Value)
		b = core.Flip(b, msg.Bit)
		m.grid[pos.Row][pos.Col].Value = core.Bitmap_To_Braille(b)

	case panel.ActionMsg:
		m.apply_action(msg.Action)

	case panel.ColorMsg:
		m.apply_colours(msg)
		m.canvas.Mode = canvas.NORMAL_MODE
		m.canvas.Reset_to_normal()
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

func (m model) get_cell() rune {
	return m.grid[m.canvas.Cursor.Row][m.canvas.Cursor.Col].Value
}

func newModel(grid core.Grid) model {
	selected := make(core.Selected)
	m := model{
		canvas:   canvas.New(0, 0, grid, selected),
		grid:     grid,
		selected: selected,
	}
	m.panel = panel.New(PANEL_WIDTH, 0, m.get_cell)
	return m
}
