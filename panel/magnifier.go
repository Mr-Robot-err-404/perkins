package panel

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/core"
	"github.com/Mr-Robot-err-404/perkins/theme"
)

const (
	Cell string = "██"
)

type Bits = [4][2]bool

func render_magnifier(bits Bits) string {
	s := strings.Builder{}

	for i, current := range bits {
		s.WriteString(derive_bit(current[0]))
		s.WriteString(" ")
		s.WriteString(derive_bit(current[1]))
		if i < len(bits)-1 {
			s.WriteString("\n")
		}
	}
	return s.String()
}

func derive_bit(on bool) string {
	if on {
		return set_cell(theme.ActiveCellFG, Cell)
	}
	return set_cell(theme.CellBG, "  ")
}
func set_cell(color string, s string) string {
	return theme.Reset + color + s + theme.Reset + theme.PanelBG
}

func magnifier(r rune) Bits {
	bits := Bits{}

	if !core.Is_Braille(r) {
		return bits
	}
	b := core.Bitmap(r)

	var n byte
	for ; n < 8; n++ {
		pos, ok := core.Ascii_map[n]
		if !ok {
			continue
		}
		bit := b & (1 << n)
		if bit == 0 {
			continue
		}
		bits[pos.Row][pos.Col] = true
	}
	return bits
}
