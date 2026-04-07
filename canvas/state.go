package canvas

import (
	"github.com/Mr-Robot-err-404/perkins/core"
)

func (m Model) toggle_mode(mode int) int {
	if mode == m.mode {
		return NORMAL_MODE
	}
	return mode
}
func mirror_pos(pos core.Pos, axis int, w, h int) core.Pos {
	switch axis {
	case HORIZONTAL:
		col := w - 1 - pos.Col
		return core.Pos{Row: pos.Row, Col: col}
	case VERTICAL:
		row := h - 1 - pos.Row
		return core.Pos{Row: row, Col: pos.Col}
	default:
		return core.Pos{}
	}
}

func (m Model) expand_selection() {
	if m.mode == NORMAL_MODE {
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
	if m.selector.mirror == MIRROR_DISABLE {
		return
	}
	mirror := make(core.Selected, len(m.Selected))
	w, h := len(m.Grid[0]), len(m.Grid)

	for pos := range m.Selected {
		p := mirror_pos(pos, m.selector.mirror, w, h)
		mirror[p] = true
	}
	for pos := range mirror {
		m.Selected[pos] = true
	}
}
func (m Model) update_cursor(pos core.Pos) {
	*m.prev_cursor = *m.Cursor
	*m.Cursor = pos
	m.expand_selection()
}

func (m Model) set_mirror_axis(axis int) {
	if m.mode != VISUAL_BLOCK {
		return
	}
	if m.selector.mirror == axis {
		m.selector.mirror = MIRROR_DISABLE
		return
	}
	m.selector.mirror = axis
}

func (m Model) set_normal_mode() Model {
	m.mode = NORMAL_MODE
	*m.selector = Selector{}
	*m.harpoon = Harpoon{}
	clear(m.Selected)
	return m
}
