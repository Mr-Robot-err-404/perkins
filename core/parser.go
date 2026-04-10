package core

import "bytes"

func parse_ansi() {
}

func Parse_Grid(b []byte) Grid {
	grid := Grid{}
	width := 0

	for line := range bytes.SplitSeq(b, []byte("\n")) {
		trim(&line)
		if len(line) == 0 {
			continue
		}
		r := []rune(string(line))
		width = max(width, len(r))
	}
	for line := range bytes.SplitSeq(b, []byte("\n")) {
		trim(&line)
		if len(line) == 0 {
			continue
		}
		r := []rune(string(line))
		size := width - len(r)
		pad(&r, size)
		grid = append(grid, r)
	}
	return grid
}

func trim(b *[]byte) {
	*b = bytes.TrimSuffix(*b, []byte(" "))
}

func pad(r *[]rune, size int) {
	for range size {
		*r = append(*r, Base)
	}
}
