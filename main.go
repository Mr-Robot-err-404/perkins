package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/debug"
	"github.com/Mr-Robot-err-404/perkins/scaling"
	tea "github.com/charmbracelet/bubbletea"
	_ "golang.org/x/image/webp"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: perkins <convert|edit> <file>\n")
		os.Exit(1)
	}
	cmd := os.Args[1]
	file_path := os.Args[2]

	err := debug.Init()
	if err != nil {
		panic(err.Error())
	}
	abs, err := filepath.Abs(file_path)

	if err != nil {
		panic(err)
	}
	if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(abs, home) {
		abs = "~/" + abs[len(home)+1:]
	}
	meta := meta{file_path: abs}

	switch cmd {
	case "convert":
		f, err := os.Open(file_path)
		if err != nil {
			panic(err.Error())
		}
		defer f.Close()
		img, _, err := image.Decode(f)

		if err != nil {
			panic(err.Error())
		}
		size := core.Dimensions{Width: 102, Height: 51}
		ch := make(chan core.Grid, 1)

		if err := scaling.Run(img, size, ch); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		grid := <-ch
		p := tea.NewProgram(newModel(grid, meta), tea.WithAltScreen(), tea.WithMouseCellMotion())

		if _, err := p.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case "edit":
		b, err := os.ReadFile(file_path)

		if err != nil {
			panic(err.Error())
		}
		grid := core.Parse_Ansi(b)
		p := tea.NewProgram(newModel(grid, meta), tea.WithAltScreen(), tea.WithMouseCellMotion())

		if _, err := p.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q. Usage: %s <convert|edit> <file>\n", cmd, os.Args[0])
		os.Exit(1)
	}
}
