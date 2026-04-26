package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Mr-Robot-err-404/perkins/canvas"
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/debug"
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
	history  *core.History
	meta     meta
}
type meta struct {
	home      string
	file_path string
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
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		canvas_width := max(0, m.width-PANEL_WIDTH)
		m.canvas = m.canvas.Resize(canvas_width, m.height)
		m.panel = m.panel.Resize(
			panel.Dimensions{Width: PANEL_WIDTH, Height: m.height},
			panel.Dimensions{Width: m.width, Height: m.height},
		)

	case canvas.CropMsg:
		from := m.grid
		m.grid = msg.Grid
		to := m.grid
		m.history.Branch(core.State{From_Grid: from, To_Grid: to})

		m.canvas.Grid = msg.Grid
		m.canvas.Mode = canvas.NORMAL_MODE
		m.canvas.Reset_to_normal()
		*m.canvas.Cursor = core.Pos{}
		m.panel.Cell = m.get_cell

	case panel.FlipMsg:
		if m.canvas.Mode == canvas.CROP_MODE {
			return m, nil
		}
		pos := *m.canvas.Cursor
		cell := m.grid[pos.Row][pos.Col]
		filter := core.Filter_Cells(m.grid, cell, m.selected)

		from := core.MakeSnapshot(m.grid, m.selected)
		for p := range filter {
			core.Flip_Cell(m.grid, p, msg.Bit)
		}
		to := core.MakeSnapshot(m.grid, m.selected)
		m.history.Branch(core.State{From: from, To: to})

	case panel.ActionMsg:
		from := core.MakeSnapshot(m.grid, m.selected)
		m.apply_action(msg.Action)
		to := core.MakeSnapshot(m.grid, m.selected)
		m.history.Branch(core.State{From: from, To: to})

	case panel.ColorMsg:
		from := core.MakeSnapshot(m.grid, m.selected)
		m.apply_colours(msg)
		to := core.MakeSnapshot(m.grid, m.selected)
		m.history.Branch(core.State{From: from, To: to})

		m.canvas.Mode = canvas.NORMAL_MODE
		m.canvas.Reset_to_normal()

	case canvas.UndoMsg:
		is_snapshot := m.history.Undo(&m.grid)
		if !is_snapshot {
			m.canvas.Grid = m.grid
			m.panel.Cell = m.get_cell
			clear(m.selected)
			*m.canvas.Cursor = core.Pos{}
		}
		return m, nil

	case canvas.RedoMsg:
		is_snapshot := m.history.Redo(&m.grid)
		if !is_snapshot {
			m.canvas.Grid = m.grid
			m.panel.Cell = m.get_cell
			clear(m.selected)
			*m.canvas.Cursor = core.Pos{}
		}
		return m, nil

	case canvas.SaveMsg:
		path := msg.Path

		if strings.HasPrefix(msg.Path, "~/") {
			path = filepath.Join(m.meta.home, path[2:])
		}
		if strings.HasPrefix(msg.Path, "$HOME/") {
			path = filepath.Join(m.meta.home, path[6:])
		}
		err := os.WriteFile(path, msg.Ascii, 0644)

		if err != nil {
			msg := fmt.Sprintf("Failed to save file: %s", err.Error())
			debug.Log(msg)
			return m, canvas.Notify(msg)
		}
		return m, canvas.Notify(fmt.Sprintf("Saved file to %s!", path))
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

func newModel(grid core.Grid, meta meta) model {
	selected := make(core.Selected)
	selected[core.Pos{Row: 0, Col: 0}] = core.Highlight

	m := model{
		canvas:   canvas.New(0, 0, grid, selected, meta.file_path),
		grid:     grid,
		selected: selected,
		history:  core.MakeHistory(),
		meta:     meta,
	}
	m.panel = panel.New(panel.Dimensions{Width: PANEL_WIDTH}, m.get_cell)
	return m
}
