package core

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

const Threshold float64 = 128

var Neighbors = map[Coords]float64{
	{X: 1, Y: 0}:  7.0 / 16.0,
	{X: 1, Y: 1}:  1.0 / 16.0,
	{X: 0, Y: 1}:  5.0 / 16.0,
	{X: -1, Y: 1}: 3.0 / 16.0,
}

type ImageBitmap struct {
	Buf    []uint64
	Width  int
	Height int
}

func (bitmap *ImageBitmap) idx(x, y int) int {
	return ((y * bitmap.Width) + x) / 64
}
func (bitmap *ImageBitmap) bit(x, y int) int {
	return ((y * bitmap.Width) + x) % 64
}

func (bitmap *ImageBitmap) Invert() {
	for idx := range bitmap.Buf {
		bitmap.Buf[idx] = ^bitmap.Buf[idx]
	}
}

func Floyd_Steinberg(img image.Image) ImageBitmap {
	buf := image_to_buffer(img)

	size := (len(buf)*len(buf[0]) + 63) / 64
	bitmap := ImageBitmap{
		Buf:    make([]uint64, size),
		Width:  len(buf[0]),
		Height: len(buf),
	}
	visited := map[Coords]bool{}
	bounds := img.Bounds()

	for y := range bounds.Max.Y {
		for x := range bounds.Max.X {
			coords := Coords{X: x, Y: y}
			visited[coords] = true

			q, diff := quantize(buf, coords)
			buf[y][x] = q
			diffuse(coords, buf, visited, diff, bounds)

			if buf[y][x] == 0 {
				idx := bitmap.idx(x, y)
				bitmap.Buf[idx] |= 1 << bitmap.bit(x, y)
			}
		}
	}
	return bitmap
}

func Buffer_to_image(buf [][]float64) image.Image {
	bounds := image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(len(buf[0]), len(buf))}
	img := image.NewGray(bounds)

	for y := range len(buf) {
		for x := range len(buf[y]) {
			n := buf[y][x]
			img.Set(x, y, color.Gray{Y: uint8(n)})
		}
	}
	return img
}

func SaveJPG(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
}

func Dither_Ye_NOT(img image.Image) [][]float64 {
	return image_to_buffer(img)
}

func diffuse(current Coords, buf [][]float64, visited map[Coords]bool, diff float64, bounds image.Rectangle) {
	for next, multiplier := range Neighbors {
		neighbor := Coords{
			X: current.X + next.X,
			Y: current.Y + next.Y,
		}
		if out_of_image_bounds(bounds, neighbor.X, neighbor.Y) {
			continue
		}
		_, seen := visited[neighbor]
		if seen {
			continue
		}
		n := diff * multiplier
		buf[neighbor.Y][neighbor.X] += n
	}
}

func quantize(m [][]float64, coords Coords) (float64, float64) {
	lum := m[coords.Y][coords.X]

	if lum < Threshold {
		return 0, lum
	}
	return 255, lum - 255
}

func out_of_image_bounds(bounds image.Rectangle, x, y int) bool {
	if x < 0 || y < 0 {
		return true
	}
	if x >= bounds.Max.X || y >= bounds.Max.Y {
		return true
	}
	return false
}
