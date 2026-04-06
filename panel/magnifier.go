package panel

import (
	"strings"

	"github.com/Mr-Robot-err-404/perkins/common"
	"github.com/Mr-Robot-err-404/perkins/theme"
)

const (
	Cell string = "██"
	Base rune   = 0x2800
	Full rune   = 0x28FF
)

type Cells = [4][2]bool

func render_magnifier(cells Cells) string {
	s := strings.Builder{}

	for i, current := range cells {
		s.WriteString(derive_cell(current[0]))
		s.WriteString(" ")
		s.WriteString(derive_cell(current[1]))
		if i < len(cells)-1 {
			s.WriteString("\n")
		}
	}
	return s.String()
}

func derive_cell(on bool) string {
	if on {
		return set_cell(theme.ActiveCellFG, Cell)
	}
	return set_cell(theme.CellBG, "  ")
}
func set_cell(color string, s string) string {
	return theme.Reset + color + s + theme.Reset + theme.PanelBG
}

func magnifier(r rune) Cells {
	cells := Cells{}

	if r < Base || r > Full {
		return cells
	}
	b := bitmap(r)

	var n byte
	for ; n < 8; n++ {
		pos, ok := common.Ascii_map[n]
		if !ok {
			continue
		}
		bit := b & (1 << n)
		if bit == 0 {
			continue
		}
		cells[pos.Row][pos.Col] = true
	}
	return cells
}

func bitmap(r rune) byte {
	return byte(r - Base)
}
