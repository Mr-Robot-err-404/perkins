package core

import (
	"image"
)

const Threshold float64 = 128

var Neighbors = map[Coords]float64{
	{X: 1, Y: 0}:  7.0 / 16.0,
	{X: 1, Y: 1}:  1.0 / 16.0,
	{X: 0, Y: 1}:  5.0 / 16.0,
	{X: -1, Y: 1}: 3.0 / 16.0,
}

func Floyd_Steinberg(img image.Image) [][]float64 {
	buf := image_to_buffer(img)
	visited := map[Coords]bool{}
	bounds := img.Bounds()

	for y := range bounds.Max.Y {
		for x := range bounds.Max.X {
			coords := Coords{X: x, Y: y}
			visited[coords] = true

			q, diff := quantize(buf, coords)
			buf[y][x] = q
			diffuse(coords, buf, visited, diff, bounds)
		}
	}
	return buf
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
