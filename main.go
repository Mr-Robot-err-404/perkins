package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/Mr-Robot-err-404/perkins/canvas"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	b, err := os.ReadFile("ascii")

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

	for line := range bytes.SplitSeq(b, []byte("\n")) {
		grid = append(grid, []rune(string(line)))
	}
	return grid
}
