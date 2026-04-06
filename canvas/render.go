package canvas

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
)

func grid_to_canvas(grid Grid, cursor core.Pos) string {
	cv := strings.Builder{}
	cv.WriteString(theme.CanvasBG)

	for row, line := range grid {
		for col, r := range line {
			if row == cursor.Row && col == cursor.Col {
				set_canvas_cell(&cv, theme.CursorAnsi, r)
				continue
			}
			cv.WriteRune(r)
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

func find_center(grid Grid) core.Pos {
	return core.Pos{
		Row: len(grid) / 2,
		Col: (len(grid[0]) / 2) - 1,
	}
}
