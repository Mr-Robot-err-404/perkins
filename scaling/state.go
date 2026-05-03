package scaling

import (
	"github.com/Mr-Robot-err-404/perkins/core"
)

func window(grid core.Grid, dm core.Dimensions) core.Grid {
	start, end := core.Window(dm, grid)
	m := make(core.Grid, end.Row-start.Row)

	for row := start.Row; row < end.Row; row++ {
		i := row - start.Row
		m[i] = make([]core.Cell, end.Col-start.Col)

		for col := start.Col; col < end.Col; col++ {
			j := col - start.Col
			m[i][j] = grid[row][col]
		}
	}
	return m
}
