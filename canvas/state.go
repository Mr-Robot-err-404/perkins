package canvas

import (
	"github.com/Mr-Robot-err-404/perkins/core"
)

func (m Model) toggle_mode(mode int) int {
	if mode == m.Mode {
		return NORMAL_MODE
	}
	return mode
}
func (m Model) toggle_mirror_axis() int {
	switch m.selector.mirror_axis {
	case X_AXIS:
		return Y_AXIS
	case Y_AXIS:
		return X_AXIS
	default:
		return MIRROR_DISABLE
	}
}

func (m Model) init_cropping_block() {
	switch m.selector.mirror_axis {
	case Y_AXIS:
		m.Cursor.Col = 0
		*m.harpoon = Harpoon{
			min:   core.Pos{Col: 0, Row: 0},
			max:   core.Pos{Col: 0, Row: len(m.Grid)},
			start: *m.Cursor,
		}
	case X_AXIS:
		m.Cursor.Row = 0
		*m.harpoon = Harpoon{
			min:   core.Pos{Row: 0, Col: 0},
			max:   core.Pos{Row: 0, Col: len(m.Grid[0])},
			start: *m.Cursor,
		}
	}
}

func mirror_pos(pos core.Pos, axis int, w, h int) core.Pos {
	switch axis {
	case Y_AXIS:
		col := w - 1 - pos.Col
		return core.Pos{Row: pos.Row, Col: col}
	case X_AXIS:
		row := h - 1 - pos.Row
		return core.Pos{Row: row, Col: pos.Col}
	default:
		return core.Pos{}
	}
}
func mirror_harpoon(harpoon *Harpoon, axis int, w, h int) Harpoon {
	return Harpoon{
		min: mirror_pos(harpoon.max, axis, w, h),
		max: mirror_pos(harpoon.min, axis, w, h),
	}
}

func (m Model) crop_canvas() core.Grid {
	grid := make(core.Grid, len(m.Grid))
	w, h := len(m.Grid[0]), len(m.Grid)
	mirror := mirror_harpoon(m.harpoon, m.selector.mirror_axis, w, h)

	switch m.selector.mirror_axis {
	case Y_AXIS:
		start := m.harpoon.max.Col + 1
		end := mirror.min.Col
		size := end - start

		for row := range m.Grid {
			current := make([]core.Cell, size)

			for col := start; col < end; col++ {
				idx := col - start
				current[idx] = m.Grid[row][col]
			}
			grid[row] = current
		}
	case X_AXIS:
		start := m.harpoon.max.Row + 1
		end := mirror.min.Row

		for row := start; row < end; row++ {
			idx := row - start
			grid[idx] = m.Grid[row]
		}
	}
	return grid
}

func (m Model) selection_type() int {
	if m.Mode == CROP_MODE {
		return core.Crop
	}
	return core.Highlight
}

func (h *Harpoon) selection(pos core.Pos) {
	h.min.Row = min(h.start.Row, pos.Row)
	h.min.Col = min(h.start.Col, pos.Col)
	h.max.Row = max(h.start.Row, pos.Row)
	h.max.Col = max(h.start.Col, pos.Col)
}
func (h *Harpoon) crop(pos core.Pos, axis int) {
	switch axis {
	case X_AXIS:
		h.max.Row = pos.Row
	case Y_AXIS:
		h.max.Col = pos.Col
	}
}

func (m Model) expand_selection() {
	if m.Mode == NORMAL_MODE {
		return
	}
	slt := m.selection_type()
	pos := *m.Cursor

	switch m.Mode {
	case VISUAL_BLOCK:
		m.harpoon.selection(pos)
	case CROP_MODE:
		m.harpoon.crop(pos, m.selector.mirror_axis)
	}
	clear(m.Selected)

	for row := m.harpoon.min.Row; row <= m.harpoon.max.Row; row++ {
		for col := m.harpoon.min.Col; col <= m.harpoon.max.Col; col++ {
			pos := core.Pos{Row: row, Col: col}
			m.Selected[pos] = slt
		}
	}
	if m.selector.mirror_axis == MIRROR_DISABLE {
		return
	}
	mirror := make(core.Selected, len(m.Selected))
	w, h := len(m.Grid[0]), len(m.Grid)

	for pos := range m.Selected {
		p := mirror_pos(pos, m.selector.mirror_axis, w, h)
		mirror[p] = slt
	}
	for pos := range mirror {
		m.Selected[pos] = slt
	}
}
func (m Model) update_cursor(pos core.Pos) {
	*m.prev_cursor = *m.Cursor
	*m.Cursor = pos
	m.expand_selection()
}

func (m Model) set_mirror_axis(axis int) {
	if m.selector.mirror_axis == axis {
		m.selector.mirror_axis = MIRROR_DISABLE
		return
	}
	m.selector.mirror_axis = axis
}

func (m Model) Reset_to_normal() {
	*m.selector = Selector{}
	*m.harpoon = Harpoon{}
	clear(m.Selected)
}
