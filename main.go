package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/Mr-Robot-err-404/perkins/canvas"
	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/debug"
	tea "github.com/charmbracelet/bubbletea"
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
		buf := core.Floyd_Steinberg(img)
		grid := core.Image_To_Grid_2(buf, core.Dimensions{Width: 102, Height: 51})
		ansi := canvas.Grid_To_Canvas(grid, core.Selected{}, core.Pos{Row: -1, Col: -1}, false)
		os.WriteFile("converted", []byte(ansi), 0644)

	case "edit":
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err.Error())
		}
		abs, err := filepath.Abs(file_path)
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(abs, home) {
			offset := len(home) + 1
			abs = "~/" + abs[offset:]
		}
		b, err := os.ReadFile(file_path)
		if err != nil {
			panic(err.Error())
		}
		grid := core.Parse_Ansi(b)
		meta := meta{home: home, file_path: abs}
		p := tea.NewProgram(newModel(grid, meta), tea.WithAltScreen(), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case "dev":
		f, err := os.Open(file_path)

		if err != nil {
			panic(err.Error())
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			panic(err.Error())
		}
		buf := core.Floyd_Steinberg(img)
		gray := core.Buffer_to_image(buf)
		err = core.SaveJPG(gray, "dithered.jpg")

		if err != nil {
			panic(err)
		}
		fmt.Println("done")

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q. Usage: %s <convert|edit> <file>\n", cmd, os.Args[0])
		os.Exit(1)
	}
}
