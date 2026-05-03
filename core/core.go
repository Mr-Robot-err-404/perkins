package core

type Pos struct {
	Row int
	Col int
}
type Coords struct {
	X int
	Y int
}
type Cell struct {
	Value rune
	Bg    string
	Fg    string
}
type Dimensions struct {
	Width  int
	Height int
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

var Coords_map = map[Coords]byte{
	{X: 0, Y: 0}: 0,
	{X: 0, Y: 1}: 1,
	{X: 0, Y: 2}: 2,
	{X: 0, Y: 3}: 6,
	{X: 1, Y: 0}: 3,
	{X: 1, Y: 1}: 4,
	{X: 1, Y: 2}: 5,
	{X: 1, Y: 3}: 7,
}

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

func Window(dm Dimensions, grid Grid, midpoint Pos) (Pos, Pos) {
	dx := dm.Width / 2
	dy := dm.Height / 2

	start := Pos{
		Row: max(0, midpoint.Row-dy),
		Col: max(0, midpoint.Col-dx),
	}
	end := Pos{
		Row: max(dm.Height-1, midpoint.Row+dy),
		Col: max(dm.Width-1, midpoint.Col+dx),
	}
	end.Row = min(len(grid), end.Row)
	end.Col = min(len(grid[0]), end.Col)
	return start, end
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
