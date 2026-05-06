package core

import (
	"image"
	"time"

	"github.com/Mr-Robot-err-404/perkins/debug"
)

type AsciiParams struct {
	Img       image.Image
	Size      Dimensions
	Invert    bool
	Algorithm int
}

func Image_To_Ascii(params AsciiParams) (Grid, ImageBitmap) {
	b := params.Img.Bounds()
	debug.Logf("[debug] image size: %dx%d", b.Max.X, b.Max.Y)

	t0 := time.Now()
	resized := Scale_Down(params.Img, params.Size)
	debug.Logf("[perf] resize: %s", time.Since(t0))

	t1 := time.Now()
	bitmap := Dithering(resized, params.Algorithm)
	debug.Logf("[perf] dithering: %s", time.Since(t1))

	if params.Invert {
		bitmap.Invert()
	}
	return Image_To_Grid(bitmap, params.Size), bitmap
}

func Scale_Down(src image.Image, canvas Dimensions) image.Image {
	dm := Dimensions{Width: canvas.Width * 2, Height: canvas.Height * 4}
	bounds := image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(dm.Width, dm.Height)}
	img := image.NewRGBA(bounds)

	dx := float64(src.Bounds().Max.X) / float64(bounds.Max.X)
	dy := float64(src.Bounds().Max.Y) / float64(bounds.Max.Y)

	for x := range dm.Width {
		for y := range dm.Height {
			color := src.At(int(float64(x)*dx), int(float64(y)*dy))
			img.Set(x, y, color)
		}
	}
	return img
}

func Image_To_Grid(bitmap ImageBitmap, canvas Dimensions) Grid {
	m := map[Coords]byte{}
	grid := make(Grid, canvas.Height)

	for y := range bitmap.Height {
		for x := range bitmap.Width {
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
			b := bitmap.Buf[idx] >> n & 1

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
	r >>= 8
	g >>= 8
	b >>= 8
	return float64((r*2126 + g*7152 + b*722) / 10000)
}
