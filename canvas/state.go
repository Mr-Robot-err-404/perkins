package canvas

import (
	"github.com/Mr-Robot-err-404/perkins/core"
)

func (m Model) mouse_to_grid(x, y int, window core.Window) (core.Pos, bool) {
	x_offset := offset(m.width, len(m.Grid[0]))
	y_offset := offset(m.height-1, len(m.Grid))
	pos := core.Pos{
		Col: x - x_offset + window.Start.Col,
		Row: y - y_offset + window.Start.Row,
	}
	if core.Out_Of_Bounds(pos, m.Grid) {
		return core.Pos{}, false
	}
	return pos, true
}
func offset(viewport int, size int) int {
	diff := (viewport - size)
	if diff < 0 {
		return 0
	}
	n := diff / 2

	if diff%2 != 0 {
		n++
	}
	return n
}

func (m Model) consume() int {
	n := max(1, *m.n)
	*m.n = 0
	return n
}

func (m Model) toggle_mode(mode int) int {
	if mode == m.Mode {
		return NORMAL_MODE
	}
	return mode
}

func (m Model) toggle_mirror() {
	if m.Mode != DRAW_MODE && m.Mode != VISUAL_BLOCK {
		return
	}
	m.mirror.enabled = !m.mirror.enabled
	m.expand_selection()
}

func (m Model) toggle_mirror_axis() int {
	if !m.mirror.enabled {
		return m.mirror.axis
	}
	if m.mirror.axis == X_AXIS {
		return Y_AXIS
	}
	return X_AXIS
}

func (m Model) init_cropping_block() {
	switch m.mirror.axis {
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
	mirror := mirror_harpoon(m.harpoon, m.mirror.axis, w, h)

	switch m.mirror.axis {
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
		size := end - start
		grid = make(core.Grid, size)

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
	w, h := len(m.Grid[0]), len(m.Grid)

	switch m.Mode {
	case NORMAL_MODE:
		clear(m.Selected)
		m.Selected[*m.Cursor] = core.Highlight
		return

	case DRAW_MODE:
		clear(m.Selected)
		m.Draw[*m.Cursor] = core.Highlight

		for pos := range m.Draw {
			m.Selected[pos] = core.Highlight

			if !m.mirror.enabled {
				continue
			}
			mirror := mirror_pos(pos, m.mirror.axis, w, h)
			m.Selected[mirror] = core.Highlight
		}
		return
	}
	slt := m.selection_type()
	pos := *m.Cursor

	switch m.Mode {
	case VISUAL_BLOCK:
		m.harpoon.selection(pos)
	case CROP_MODE:
		m.harpoon.crop(pos, m.mirror.axis)
	}
	clear(m.Selected)

	for row := m.harpoon.min.Row; row <= m.harpoon.max.Row; row++ {
		for col := m.harpoon.min.Col; col <= m.harpoon.max.Col; col++ {
			pos := core.Pos{Row: row, Col: col}
			m.Selected[pos] = slt
		}
	}
	if !m.mirror.enabled {
		return
	}
	mirror := make(core.Selected, len(m.Selected))

	for pos := range m.Selected {
		p := mirror_pos(pos, m.mirror.axis, w, h)
		mirror[p] = slt
	}
	for pos := range mirror {
		m.Selected[pos] = slt
	}
}
func (m Model) update_cursor(pos core.Pos) {
	*m.prev_cursor = *m.Cursor
	*m.Cursor = pos
	update_window(m.Window, pos)
	m.expand_selection()
}

func update_window(window *core.Window, pos core.Pos) {
	if pos.Row >= window.End.Row {
		diff := pos.Row - window.End.Row + 1
		window.End.Row += diff
		window.Start.Row += diff
	}
	if pos.Col >= window.End.Col {
		diff := pos.Col - window.End.Col + 1
		window.End.Col += diff
		window.Start.Col += diff
	}
	if pos.Row < window.Start.Row {
		diff := window.Start.Row - pos.Row
		window.Start.Row = pos.Row
		window.End.Row -= diff
	}
	if pos.Col < window.Start.Col {
		diff := window.Start.Col - pos.Col
		window.Start.Col = pos.Col
		window.End.Col -= diff
	}
}

func (m Model) set_mirror_axis(axis int) {
	m.mirror.axis = axis
}

func (m Model) Reset_to_normal() {
	*m.mirror = Mirror{}
	*m.harpoon = Harpoon{}
	*m.cmd = []rune{}
	clear(m.Selected)
	clear(m.Draw)
}

func (m Model) Reset_Window(grid core.Grid) {
	midpoint := core.Find_Center(grid)
	*m.Window = core.Get_Window(core.Dimensions{Width: m.width, Height: m.height - 1}, grid, midpoint)
	*m.Cursor = midpoint
}
