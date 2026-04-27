package core

func Image_To_Grid(buf [][]float64, width int, height int) Grid {
	dm := Dimensions{
		Width:  width * 2,
		Height: height * 4,
	}
	grid := make(Grid, height)

	for row := range height {
		grid[row] = make([]Cell, width)

		for col := range width {
			pos := Pos{Row: row, Col: col}
			b := bitmap(buf, pos, dm)
			grid[row][col] = Cell{Value: Bitmap_To_Braille(b)}
		}
	}
	return grid
}

func luminance(r, g, b uint32) float64 {
	return float64(299*r+587*g+114*b) / 1000.0 / 256.0
}

func is_pixel_on(buf [][]float64, src Pos) bool {
	return buf[src.Row][src.Col] < Threshold
}

func source_position(pos Pos, dy, dx int, bw, bh int, dm Dimensions) Pos {
	return Pos{
		Col: (pos.Col*2 + dx) * bw / dm.Width,
		Row: (pos.Row*4 + dy) * bh / dm.Height,
	}
}

func bitmap(buf [][]float64, pos Pos, dm Dimensions) byte {
	var bitmap byte
	bh := len(buf)
	bw := len(buf[0])

	for dy := range 4 {
		for dx := range 2 {
			src := source_position(pos, dy, dx, bw, bh, dm)
			if is_pixel_on(buf, src) {
				bitmap |= 1 << Pos_map[Pos{Row: dy, Col: dx}]
			}
		}
	}
	return bitmap
}
