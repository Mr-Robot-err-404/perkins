package panel

import (
	"github.com/Mr-Robot-err-404/perkins/core"
)

const (
	SQUARE_WIDTH  int = 6
	SQUARE_HEIGHT int = 3
	X_GAP         int = 2
	Y_GAP         int = 1
)
const (
	MAGNIFY_HEIGHT      int = 6
	DIVIDER_HEIGHT      int = 3
	LAYER_TOGGLE_HEIGHT int = 3
	BORDER_SIZE         int = 1
	PALETTE_HEIGHT      int = 20
	CONTENT_HEIGHT      int = MAGNIFY_HEIGHT + DIVIDER_HEIGHT + PALETTE_HEIGHT
)

type Coords map[core.Pos]core.Pos

func coordinates_to_idx(x_offset int, y_offset int) Coords {
	m := make(Coords)
	x := x_offset

	for col := range 2 {
		y := y_offset

		for row := range 4 {
			square_pos := core.Pos{Row: row, Col: col}
			map_square_coords(x, y, m, square_pos)
			y += SQUARE_HEIGHT + Y_GAP
		}
		x += SQUARE_WIDTH + X_GAP
	}
	return m
}

func map_square_coords(x int, y int, m Coords, pos core.Pos) {
	for row := y; row < y+SQUARE_HEIGHT; row++ {
		for col := x; col < x+SQUARE_WIDTH; col++ {
			p := core.Pos{Row: row, Col: col}
			m[p] = pos
		}
	}
}

func (m Model) palette_start() core.Pos {
	panel_x := m.terminal.Width - m.panel.Width
	content_x := panel_x + (m.panel.Width-16)/2
	x_offset := content_x + BORDER_SIZE

	content_y := (m.panel.Height - CONTENT_HEIGHT) / 2
	y_offset := content_y + MAGNIFY_HEIGHT + DIVIDER_HEIGHT + LAYER_TOGGLE_HEIGHT + BORDER_SIZE

	return core.Pos{Col: x_offset, Row: y_offset}
}
