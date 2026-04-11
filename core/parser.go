package core

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

const (
	Esc     string = "\x1b"
	OpenCsi string = "\x1b["
)

func Parse_Ansi(data []byte) Grid {
	var buf bytes.Buffer

	m := make(map[Pos]Cell)
	row, col := 0, 0

	p := ansi.GetParser()

	p.SetHandler(ansi.Handler{
		Print: func(r rune) {
			pos := Pos{Row: row, Col: col}
			m[pos] = Cell{Value: r}
			col++
		},
		Execute: func(b byte) {
			if b == '\n' {
				row++
				col = 0
			}
		},
		HandleEsc: func(cmd ansi.Cmd) {
			buf.WriteString(Esc)

			if cmd.Intermediate() != 0 {
				buf.WriteByte(cmd.Intermediate())
			}
			buf.WriteByte(cmd.Final())
		},
		HandleCsi: func(cmd ansi.Cmd, params ansi.Params) {
			buf.WriteString(OpenCsi)

			if cmd.Prefix() != 0 {
				buf.WriteByte(cmd.Prefix())
			}
			params.ForEach(0, func(i, param int, hasMore bool) {
				fmt.Fprintf(&buf, "%d", param)

				if i < len(params)-1 {
					write_separator(&buf, hasMore)
				}
			})
			if cmd.Intermediate() != 0 {
				buf.WriteByte(cmd.Intermediate())
			}
			buf.WriteByte(cmd.Final())
		},
	})
	p.Parse(data)

	width, height := grid_dimensions(m)
	grid := make(Grid, height)

	for i := range grid {
		grid[i] = make([]Cell, width)
	}
	for pos, cell := range m {
		grid[pos.Row][pos.Col] = cell
	}
	return grid
}

func grid_dimensions(m map[Pos]Cell) (int, int) {
	width, height := 0, 0

	for pos := range m {
		w := pos.Col + 1
		h := pos.Row + 1

		if w > width {
			width = w
		}
		if h > height {
			height = h
		}
	}
	return width, height
}

func write_separator(buf *bytes.Buffer, hasMore bool) {
	if hasMore {
		buf.WriteByte(':')
		return
	}
	buf.WriteByte(';')
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
		values := []rune(string(line))
		size := width - len(values)
		pad(&values, size)

		cells := []Cell{}
		for _, v := range values {
			cells = append(cells, Cell{Value: v})
		}
		grid = append(grid, cells)
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
