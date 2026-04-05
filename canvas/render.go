package canvas

import (
	"strings"
)

type Cursor struct {
	row int
	col int
}

const (
	canvasBG   = "\x1b[48;2;22;22;29m"
	cursorAnsi = "\x1b[48;2;200;192;147;38;2;13;12;12m"
	reset      = "\x1b[0m"
)

func grid_to_canvas(grid Grid, cursor Cursor) string {
	cv := strings.Builder{}
	cv.WriteString(canvasBG)

	for row, line := range grid {
		for col, r := range line {
			if row == cursor.row && col == cursor.col {
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
