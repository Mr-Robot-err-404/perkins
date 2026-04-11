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

	for row, line := range grid {
		for col, cell := range line {
			transform_cell(&cv, row, col, pos, &prev, cell, selected)
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
	switch len(cell.Ansi) {
	case 0:
		if len(prev.Ansi) > 0 || is_prev_highlighted(selected, row, col, pos) {
			set_ansi(cv, theme.CanvasBG)
		}
		cv.WriteRune(cell.Value)

	default:
		if prev.Ansi != cell.Ansi || is_prev_highlighted(selected, row, col, pos) {
			set_ansi(cv, cell.Ansi)
		}
		cv.WriteRune(cell.Value)
	}
}

func is_prev_highlighted(selected core.Selected, row int, col int, cursor core.Pos) bool {
	col = max(0, col-1)
	pos := core.Pos{Row: row, Col: col}
	_, ok := selected[pos]
	return ok || (cursor.Col == col && cursor.Row == row)
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
