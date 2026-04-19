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
	MAGNIFY_HEIGHT int = 6
	DIVIDER_HEIGHT int = 3
	TOGGLE_HEIGHT  int = 3
	PALETTE_WIDTH  int = 16
	BORDER_SIZE    int = 1
	PALETTE_HEIGHT int = 17
	CONTENT_HEIGHT int = MAGNIFY_HEIGHT + DIVIDER_HEIGHT + PALETTE_HEIGHT + TOGGLE_HEIGHT
)

type Coords map[core.Pos]core.Pos

func palette_coords(x_offset int, y_offset int) Coords {
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

func toggle_coords(x_offset int, y_offset int) map[core.Pos]bool {
	m := make(map[core.Pos]bool)

	for row := range TOGGLE_HEIGHT {
		for col := range PALETTE_WIDTH {
			pos := core.Pos{Row: row + y_offset, Col: col + x_offset}
			m[pos] = true
		}
	}
	return m
}

func (m Model) content_offset() (int, int) {
	panel_x := m.terminal.Width - m.panel.Width
	x := panel_x + (m.panel.Width-PALETTE_WIDTH)/2
	y := (m.panel.Height - CONTENT_HEIGHT) / 2
	return x, y
}

func (m Model) palette_start() core.Pos {
	x, y := m.content_offset()
	x += BORDER_SIZE
	y += MAGNIFY_HEIGHT + DIVIDER_HEIGHT + TOGGLE_HEIGHT + BORDER_SIZE
	return core.Pos{Row: y, Col: x}
}

func (m Model) toggle_start() core.Pos {
	x, y := m.content_offset()
	y += MAGNIFY_HEIGHT + DIVIDER_HEIGHT
	return core.Pos{Row: y, Col: x}
}
