package canvas

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/common"
)

const (
	canvasBG   = "\x1b[48;2;22;22;29m"
	cursorAnsi = "\x1b[48;2;200;192;147;38;2;13;12;12m"
	reset      = "\x1b[0m"
)

func grid_to_canvas(grid Grid, cursor common.Pos) string {
	cv := strings.Builder{}
	cv.WriteString(canvasBG)

	for row, line := range grid {
		for col, r := range line {
			if row == cursor.Row && col == cursor.Col {
				cv.WriteString(reset)
				cv.WriteString(cursorAnsi)
				cv.WriteRune(r)
				cv.WriteString(reset)
				cv.WriteString(canvasBG)
				continue
			}
			cv.WriteRune(r)
		}
		if row < len(grid)-1 {
			cv.WriteString("\n")
		}
	}
	cv.WriteString(reset)
	return cv.String()
}

func find_center(grid Grid) common.Pos {
	return common.Pos{
		Row: len(grid) / 2,
		Col: (len(grid[0]) / 2) - 1,
	}
}
