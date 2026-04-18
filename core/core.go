package core

type Pos struct {
	Row int
	Col int
}
type Cell struct {
	Value rune
	Bg    string
	Fg    string
}

type Selected map[Pos]int
type Grid [][]Cell

const (
	Highlight int = iota + 1
	Crop
)

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

func Flip_Cell(grid Grid, pos Pos, bit byte) {
	cell := grid[pos.Row][pos.Col]

	if !Is_Braille(cell.Value) {
		return
	}
	b := Bitmap(cell.Value)
	b = Flip(b, bit)
	grid[pos.Row][pos.Col].Value = Bitmap_To_Braille(b)
}

func Filter_Cells(grid Grid, cell Cell, selected Selected) Selected {
	filter := make(Selected)

	for pos, value := range selected {
		current := grid[pos.Row][pos.Col]

		if current.Value != cell.Value {
			continue
		}
		filter[pos] = value
	}
	return filter
}

func Out_Of_Bounds(pos Pos, grid Grid) bool {
	if pos.Row < 0 || pos.Col < 0 {
		return true
	}
	if len(grid) == 0 {
		return true
	}
	if pos.Row >= len(grid) || pos.Col >= len(grid[0]) {
		return true
	}
	return false
}

func Is_Braille(r rune) bool {
	return r >= Base && r <= Full
}
func Find_Center(grid Grid) Pos {
	return Pos{
		Row: len(grid) / 2,
		Col: (len(grid[0]) / 2) - 1,
	}
}
