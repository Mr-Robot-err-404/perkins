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

func Parse_Ansi(data []byte) []byte {
	var buf bytes.Buffer

	p := ansi.GetParser()

	p.SetHandler(ansi.Handler{
		Print: func(r rune) {
			buf.WriteRune(r)
		},
		Execute: func(b byte) {
			buf.WriteByte(b)
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
	return buf.Bytes()
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
