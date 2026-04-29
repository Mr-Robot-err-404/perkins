package core

import (
	"image"

	"github.com/Mr-Robot-err-404/perkins/debug"
)

func Image_To_Grid(buf [][]float64, width int, height int) Grid {
	imgH, imgW := len(buf), len(buf[0])
	dotW, dotH := width*2, height*4
	debug.Logf("image=%dx%d canvas=%dx%d (dots=%dx%d) ratio=%.2fx%.2f",
		imgW, imgH, width, height, dotW, dotH,
		float64(imgW)/float64(dotW), float64(imgH)/float64(dotH),
	)
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
