package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/Mr-Robot-err-404/perkins/canvas"
	"github.com/Mr-Robot-err-404/perkins/debug"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	b, err := os.ReadFile("ascii")

	if err != nil {
		panic(err.Error())
	}
	err = debug.Init()
	if err != nil {
		panic(err.Error())
	}
	grid := parse_grid(b)

	p := tea.NewProgram(newModel(grid), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parse_grid(b []byte) canvas.Grid {
	grid := canvas.Grid{}
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
		*r = append(*r, ' ')
	}
}
