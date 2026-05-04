package core

import (
	"image"
)

func Image_To_Ascii(img image.Image, size Dimensions) Grid {
	resized := Scale_Down(img, size)
	buf := Floyd_Steinberg(resized)
	return Image_To_Grid(buf, size)
}

func Scale_Down(src image.Image, canvas Dimensions) image.Image {
	dm := Dimensions{Width: canvas.Width * 2, Height: canvas.Height * 4}
	bounds := image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(dm.Width, dm.Height)}
	img := image.NewRGBA(bounds)

	dx := src.Bounds().Max.X / bounds.Max.X
	dy := src.Bounds().Max.Y / bounds.Max.Y

	for x := range dm.Width {
		for y := range dm.Height {
			color := src.At(x*dx, y*dy)
			img.Set(x, y, color)
		}
	}
	return img
}

func Image_To_Grid(bitmap ImageBitmap, canvas Dimensions) Grid {
	m := map[Coords]byte{}
	grid := make(Grid, canvas.Height)

	for y := range bitmap.height {
		for x := range bitmap.width {
			coords := Coords{
				X: x / 2,
				Y: y / 4,
			}
			bit := Coords_map[Coords{
				X: x % 2,
				Y: y % 4,
			}]
			idx := bitmap.idx(x, y)
			n := bitmap.bit(x, y)
			b := bitmap.buf[idx] >> n & 1

			if b == 1 {
				m[coords] |= 1 << bit
			}
		}
	}
	for row := range canvas.Height {
		grid[row] = make([]Cell, canvas.Width)

		for col := range canvas.Width {
			bitmap := m[Coords{X: col, Y: row}]
			grid[row][col] = Cell{Value: Bitmap_To_Braille(bitmap)}
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
