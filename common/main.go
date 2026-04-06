package common

type Pos struct {
	Row int
	Col int
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
