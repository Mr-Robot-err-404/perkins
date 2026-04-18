package core

import (
	"slices"
)

type History struct {
	States []State
	Idx    int
}
type State struct {
	From      Snapshot
	To        Snapshot
	From_Grid Grid
	To_Grid   Grid
}
type Snapshot map[Pos]Cell

func MakeHistory() *History {
	empty := State{From: make(Snapshot), To: make(Snapshot)}
	return &History{States: []State{empty}, Idx: 0}
}

func (hist *History) Redo(grid *Grid) bool {
	hist.Idx = min(len(hist.States)-1, hist.Idx+1)
	state := hist.Current()

	if is_snapshot(state) {
		hist.RestoreSnaphot(*grid, state.To)
		return true
	}
	*grid = state.To_Grid
	return false
}

func (hist *History) Undo(grid *Grid) bool {
	state := hist.Current()

	if is_snapshot(state) {
		hist.RestoreSnaphot(*grid, hist.Current().From)
	} else {
		*grid = state.From_Grid
	}
	hist.Idx = max(0, hist.Idx-1)
	return is_snapshot(state)
}

func (hist *History) Branch(state State) {
	if hist.Idx < len(hist.States)-1 {
		hist.States = slices.Delete(hist.States, hist.Idx+1, len(hist.States))
	}
	hist.States = append(hist.States, state)
	hist.Idx++
}

func (hist *History) RestoreSnaphot(grid Grid, snapshot Snapshot) {
	for pos, cell := range snapshot {
		grid[pos.Row][pos.Col] = cell
	}
}

func (hist *History) Current() State {
	return hist.States[hist.Idx]
}

func is_snapshot(state State) bool {
	return state.From != nil && state.To != nil
}

func MakeSnapshot(grid Grid, selected Selected) Snapshot {
	snapshot := make(Snapshot)

	for pos := range selected {
		cell := grid[pos.Row][pos.Col]
		snapshot[pos] = cell
	}
	return snapshot
}
