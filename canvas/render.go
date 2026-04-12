package canvas

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
)

func Grid_To_Canvas(grid core.Grid, selected core.Selected, pos core.Pos) string {
	cv := strings.Builder{}
	cv.WriteString(theme.CanvasBG)

	prev := core.Cell{}
	width := len(grid[0])

	for row, line := range grid {
		for col, cell := range line {
			transform_cell(&cv, row, col, pos, &prev, cell, selected, width)
		}
		if row < len(grid)-1 {
			cv.WriteString("\n")
		}
	}
	cv.WriteString(theme.Reset)
	return cv.String()
}

func transform_cell(
	cv *strings.Builder,
	row int,
	col int,
	pos core.Pos,
	prev *core.Cell,
	cell core.Cell,
	selected core.Selected,
	width int,
) {
	defer func() {
		*prev = cell
	}()
	if row == pos.Row && col == pos.Col {
		set_canvas_cell(cv, theme.CursorAnsi, cell.Value)
		return
	}
	slt, ok := selected[core.Pos{Row: row, Col: col}]

	if ok {
		ansi := theme.Selection_Ansi(slt)
		set_canvas_cell(cv, ansi, cell.Value)
		return
	}
	prev_pos := prev_position(row, col, width)

	if len(cell.Fg) == 0 && len(cell.Bg) == 0 {
		if has_ansi_changed(*prev, cell) || is_prev_highlighted(selected, pos, prev_pos) {
			set_ansi(cv, theme.CanvasBG)
		}
		cv.WriteRune(cell.Value)
		return
	}
	if has_ansi_changed(*prev, cell) || is_prev_highlighted(selected, pos, prev_pos) {
		set_ansi(cv, core.Construct(cell.Fg, cell.Bg))
	}
	cv.WriteRune(cell.Value)
}

func has_ansi_changed(prev core.Cell, current core.Cell) bool {
	if prev.Fg != current.Fg {
		return true
	}
	return prev.Bg != current.Bg
}

func is_prev_highlighted(selected core.Selected, cursor core.Pos, prev core.Pos) bool {
	if cursor.Col == prev.Col && cursor.Row == prev.Row {
		return true
	}
	_, ok := selected[prev]
	return ok
}

func prev_position(row int, col int, width int) core.Pos {
	if col > 0 {
		return core.Pos{Row: row, Col: col - 1}
	}
	if row == 0 {
		return core.Pos{Row: 0, Col: 0}
	}
	return core.Pos{Row: row - 1, Col: width - 1}
}

func set_ansi(cv *strings.Builder, ansi string) {
	cv.WriteString(theme.Reset)
	cv.WriteString(ansi)
}

func set_canvas_cell(cv *strings.Builder, ansi string, r rune) {
	cv.WriteString(theme.Reset)
	cv.WriteString(ansi)
	cv.WriteRune(r)
	cv.WriteString(theme.Reset)
	cv.WriteString(theme.CanvasBG)
}
