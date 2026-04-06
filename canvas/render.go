package canvas

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/common"
	"github.com/Mr-Robot-err-404/perkins/theme"
)

func grid_to_canvas(grid Grid, cursor common.Pos) string {
	cv := strings.Builder{}
	cv.WriteString(theme.CanvasBG)

	for row, line := range grid {
		for col, r := range line {
			if row == cursor.Row && col == cursor.Col {
				cv.WriteString(theme.Reset)
				cv.WriteString(theme.CursorAnsi)
				cv.WriteRune(r)
				cv.WriteString(theme.Reset)
				cv.WriteString(theme.CanvasBG)
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

func find_center(grid Grid) common.Pos {
	return common.Pos{
		Row: len(grid) / 2,
		Col: (len(grid[0]) / 2) - 1,
	}
}
