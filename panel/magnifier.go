package panel

import "github.com/Mr-Robot-err-404/perkins/common"

const (
	Cell string = "██"
	Base rune   = 0x2800
	Full rune   = 0x28FF
)

type Cells = [4][2]bool

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
