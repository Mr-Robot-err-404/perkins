package main

import (
	"bytes"

	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
)

func color_ascii_bytes(grid core.Grid, start_col int, end_col int) []byte {
	cv := bytes.Buffer{}

	for row, line := range grid {
		for col, r := range line {
			if col == start_col {
				cv.WriteString(theme.BrownBG)
			}
			if col == end_col {
				cv.WriteString(theme.Reset)
			}
			cv.WriteRune(r)
		}
		if row < len(grid)-1 {
			cv.WriteByte('\n')
		}
	}
	cv.WriteString(theme.Reset)
	return cv.Bytes()
}
