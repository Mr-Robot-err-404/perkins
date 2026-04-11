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
	m := make(map[Pos]Cell)
	row, col := 0, 0

	var buf bytes.Buffer
	p := ansi.GetParser()

	p.SetHandler(ansi.Handler{
		Print: func(r rune) {
			set_cell_value(m, row, col, r)
			set_cell_ansi(m, row, col, buf.String())
			col++
		},
		Execute: func(b byte) {
			if b == '\n' {
				row++
				col = 0
			}
		},
		HandleEsc: func(cmd ansi.Cmd) {
			buf.Reset()
		},
		HandleCsi: func(cmd ansi.Cmd, params ansi.Params) {
			buf.Reset()

			if cmd.Final() == 'm' {
				if len(params) == 0 || params[0] == 0 {
					return
				}
			}
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
	for row := range height {
		for col := range width {
			pos := Pos{Row: row, Col: col}
			cell, ok := m[pos]

			if !ok {
				grid[row][col] = Cell{Value: Base}
				continue
			}
			grid[row][col] = cell
		}
	}
	for pos, cell := range m {
		grid[pos.Row][pos.Col] = cell
	}
	return grid
}

func set_cell_ansi(m map[Pos]Cell, row int, col int, ansi string) {
	pos := Pos{Row: row, Col: col}
	cell := m[pos]
	m[pos] = Cell{Value: cell.Value, Ansi: ansi}
}

func set_cell_value(m map[Pos]Cell, row int, col int, value rune) {
	pos := Pos{Row: row, Col: col}
	cell := m[pos]
	m[pos] = Cell{Value: value, Ansi: cell.Ansi}
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
