package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/debug"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	dev := flag.Bool("dev", false, "dev mode")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--dev] <file>\n", os.Args[0])
		os.Exit(1)
	}
	filePath := flag.Arg(0)
	b, err := os.ReadFile(filePath)

	if err != nil {
		panic(err.Error())
	}
	err = debug.Init()

	if err != nil {
		panic(err.Error())
	}
	grid := core.Parse_Grid(b)

	if *dev {
		ascii := color_ascii_bytes(grid, 18, 56)
		os.WriteFile("quick_brown_fox", ascii, 0644)
		return
	}
	p := tea.NewProgram(newModel(grid), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
