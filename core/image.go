package core

import (
	"image"
)

func Image_To_Grid(buf [][]float64, width int) Grid {
	dm := Dimensions{
		Width:  width * 2,
		Height: width * 4,
	}
	height := width / 2

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

type Cluster struct {
	On  int
	Off int
}

func Image_To_Grid_2(buf [][]float64, canvas Dimensions) Grid {
	points := map[Coords]*Cluster{}
	braille := map[Coords]byte{}

	grid := make(Grid, canvas.Height)
	dm := Dimensions{Width: canvas.Width * 2, Height: canvas.Height * 4}

	ratio_x := len(buf[0]) / dm.Width
	ratio_y := len(buf) / dm.Height

	y := 0
	for row := range buf {
		x := 0
		for col := range buf[row] {
			pixel := uint8(buf[row][col])
			p := Coords{X: x, Y: y}
			cluster, ok := points[p]

			if !ok {
				cluster = &Cluster{}
				points[p] = cluster
			}
			update_cluster(pixel, cluster)

			if (col+1)%ratio_x == 0 {
				x++
			}
		}
		if (row+1)%ratio_y == 0 {
			y++
		}
	}
	for coords, cluster := range points {
		cell := Coords{
			X: coords.X / 2,
			Y: coords.Y / 4,
		}
		bit := Coords_map[Coords{
			X: coords.X % 2,
			Y: coords.Y % 4,
		}]
		if cluster.On >= cluster.Off {
			braille[cell] |= 1 << bit
		}
	}
	for row := range canvas.Height {
		grid[row] = make([]Cell, canvas.Width)

		for col := range canvas.Width {
			b := braille[Coords{X: col, Y: row}]
			grid[row][col] = Cell{Value: Bitmap_To_Braille(b)}
		}
	}
	return grid
}

func update_cluster(pixel uint8, cluster *Cluster) {
	if pixel == 0 {
		cluster.On++
		return
	}
	cluster.Off++
}

func image_to_buffer(img image.Image) [][]float64 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	buf := make([][]float64, height)

	for y := range height {
		current := make([]float64, width)

		for x := range width {
			r, g, b, _ := img.At(x, y).RGBA()
			current[x] = luminance(r, g, b)
		}
		buf[y] = current
	}
	return buf
}

func luminance(r, g, b uint32) float64 {
	n := (r >> 8) + (g >> 8) + (b >> 8)
	return float64(n) / 3.0
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
