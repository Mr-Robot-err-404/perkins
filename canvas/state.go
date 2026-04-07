package canvas

import "github.com/Mr-Robot-err-404/perkins/core"

func (m Model) toggle_mode(mode int) int {
	if mode == m.mode {
		return NORMAL_MODE
	}
	return mode
}
func (m Model) expand_selection() {
	if m.mode != VISUAL_BLOCK {
		return
	}
	pos := *m.Cursor
	min_row := min(m.harpoon.min.Row, pos.Row)
	min_col := min(m.harpoon.min.Col, pos.Col)
	max_row := max(m.harpoon.max.Row, pos.Row)
	max_col := max(m.harpoon.max.Col, pos.Col)

	clear(m.Selected)
	for row := min_row; row <= max_row; row++ {
		for col := min_col; col <= max_col; col++ {
			pos := core.Pos{Row: row, Col: col}
			m.Selected[pos] = true
		}
	}
}
func (m Model) update_cursor(pos core.Pos) {
	*m.prev_cursor = *m.Cursor
	*m.Cursor = pos
	m.expand_selection()
}
func (m Model) set_normal_mode() Model {
	m.mode = NORMAL_MODE
	*m.harpoon = Harpoon{}
	clear(m.Selected)
	return m
}
