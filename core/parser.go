package core

import (
	"bytes"
	"fmt"
	"slices"

	"github.com/charmbracelet/x/ansi"
)

const (
	Esc     string = "\x1b"
	OpenCsi string = "\x1b["
)
const (
	FG int = iota
	BG
)

var (
	FG_sequence = []int{30, 31, 32, 33, 34, 35, 36, 37, 39, 90, 91, 92, 93, 94, 95, 96, 97}
	BG_sequence = []int{40, 41, 42, 43, 44, 45, 46, 47, 49, 100, 101, 102, 103, 104, 105, 106, 107}
)

func Parse_Ansi(data []byte) Grid {
	m := make(map[Pos]Cell)
	row, col := 0, 0

	var fg bytes.Buffer
	var bg bytes.Buffer
	p := ansi.GetParser()

	p.SetHandler(ansi.Handler{
		Print: func(r rune) {
			cell, pos := get_cell(m, row, col)
			cell.Value = r
			cell.Bg = bg.String()
			cell.Fg = fg.String()

			m[pos] = cell
			col++
		},
		Execute: func(b byte) {
			if b == '\n' {
				row++
				col = 0
			}
		},
		HandleEsc: func(cmd ansi.Cmd) {
			reset_buffers(&fg, &bg)
		},
		HandleCsi: func(cmd ansi.Cmd, params ansi.Params) {
			reset_buffers(&fg, &bg)

			if cmd.Final() == 'm' {
				if len(params) == 0 || params[0] == 0 {
					return
				}
			}
			seq := FG
			params.ForEach(0, func(i, param int, hasMore bool) {
				check_sequence(param, &seq)
				buf := derive_buffer(seq, &fg, &bg)

				fmt.Fprintf(buf, "%d", param)
				if i < len(params)-1 {
					write_separator(buf, hasMore)
				}
			})
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

func Construct(fg string, bg string) string {
	if len(fg) == 0 {
		return OpenCsi + bg + "m"
	}
	if len(bg) == 0 {
		return OpenCsi + fg + "m"
	}
	return OpenCsi + fg + ";" + bg + "m"
}

func check_sequence(param int, seq *int) {
	if slices.Contains(FG_sequence, param) {
		*seq = FG
		return
	}
	if slices.Contains(BG_sequence, param) {
		*seq = BG
	}
}

func reset_buffers(fg *bytes.Buffer, bg *bytes.Buffer) {
	fg.Reset()
	bg.Reset()
}

func derive_buffer(seq int, fg *bytes.Buffer, bg *bytes.Buffer) *bytes.Buffer {
	if seq == BG {
		return bg
	}
	return fg
}

func get_cell(m map[Pos]Cell, row int, col int) (Cell, Pos) {
	pos := Pos{Row: row, Col: col}
	cell := m[pos]
	return cell, pos
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
