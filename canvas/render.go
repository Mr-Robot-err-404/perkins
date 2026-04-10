package canvas

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
)

func grid_to_canvas(grid core.Grid, selected core.Selected, cursor core.Pos) string {
	cv := strings.Builder{}
	cv.WriteString(theme.CanvasBG)

	for row, line := range grid {
		for col, cell := range line {
			if row == cursor.Row && col == cursor.Col {
				set_canvas_cell(&cv, theme.CursorAnsi, cell.Value)
				continue
			}
			slt, ok := selected[core.Pos{Row: row, Col: col}]
			if ok {
				ansi := theme.Selection_Ansi(slt)
				set_canvas_cell(&cv, ansi, cell.Value)
				continue
			}
			cv.WriteRune(cell.Value)
		}
		if row < len(grid)-1 {
			cv.WriteString("\n")
		}
	}
	cv.WriteString(theme.Reset)
	return cv.String()
}

func set_canvas_cell(cv *strings.Builder, ansi string, r rune) {
	cv.WriteString(theme.Reset)
	cv.WriteString(ansi)
	cv.WriteRune(r)
	cv.WriteString(theme.Reset)
	cv.WriteString(theme.CanvasBG)
}
