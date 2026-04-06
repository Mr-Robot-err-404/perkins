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
	m.Selected[pos] = true
}
func (m Model) update_cursor(pos core.Pos) {
	*m.prev_cursor = *m.Cursor
	*m.Cursor = pos
	m.expand_selection()
}
func (m Model) set_normal_mode() Model {
	m.mode = NORMAL_MODE
	clear(m.Selected)
	return m
}
