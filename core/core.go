package core

type Pos struct {
	Row int
	Col int
}
type Selected map[Pos]bool

const (
	Base rune = 0x2800
	Full rune = 0x28FF
)

var Pos_map = map[Pos]byte{
	{Col: 0, Row: 0}: 0,
	{Col: 0, Row: 1}: 1,
	{Col: 0, Row: 2}: 2,
	{Col: 0, Row: 3}: 6,
	{Col: 1, Row: 0}: 3,
	{Col: 1, Row: 1}: 4,
	{Col: 1, Row: 2}: 5,
	{Col: 1, Row: 3}: 7,
}

var Ascii_map = map[byte]Pos{
	0: {Col: 0, Row: 0},
	1: {Col: 0, Row: 1},
	2: {Col: 0, Row: 2},
	6: {Col: 0, Row: 3},
	3: {Col: 1, Row: 0},
	4: {Col: 1, Row: 1},
	5: {Col: 1, Row: 2},
	7: {Col: 1, Row: 3},
}

func Bitmap(r rune) byte {
	return byte(r - Base)
}

func Bitmap_To_Braille(bitmap byte) rune {
	return rune(bitmap) + Base
}

func Flip(b byte, bit byte) byte {
	return b ^ (1 << bit)
}

func Is_Braille(r rune) bool {
	return r >= Base && r <= Full
}
