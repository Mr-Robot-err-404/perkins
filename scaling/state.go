package scaling

import (
	"image"

	"github.com/Mr-Robot-err-404/perkins/core"
)

const (
	MAX_FACTOR float64 = 10.0
	MIN_FACTOR float64 = 0.1
)

func zoom_out(base core.Dimensions, factor float64) (core.Dimensions, float64) {
	factor = max(factor-0.1, MIN_FACTOR)
	return core.Dimensions{
		Width:  max(amplify(base.Width, factor), 1),
		Height: max(amplify(base.Height, factor), 1),
	}, factor
}

func zoom_in(base core.Dimensions, factor float64, bounds image.Rectangle) (core.Dimensions, float64) {
	factor = min(factor+0.1, MAX_FACTOR)
	return core.Dimensions{
		Width:  min(amplify(base.Width, factor), bounds.Max.X-1),
		Height: min(amplify(base.Height, factor), bounds.Max.Y-1),
	}, factor
}

func amplify(n int, factor float64) int {
	return int(float64(n) * factor)
}

func window(grid core.Grid, dm core.Dimensions) core.Grid {
	midpoint := core.Find_Center(grid)
	start, end := core.Window(dm, grid, midpoint)
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
