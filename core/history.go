package core

import (
	"slices"
)

type History struct {
	States []State
	Idx    int
}
type State struct {
	From Snapshot
	To   Snapshot
}
type Snapshot map[Pos]Cell

func MakeHistory() *History {
	empty := State{From: make(Snapshot), To: make(Snapshot)}
	return &History{States: []State{empty}, Idx: 0}
}

func (hist *History) Redo(grid Grid) {
	hist.Idx = min(len(hist.States)-1, hist.Idx+1)
	hist.Restore(grid, hist.Current().To)
}

func (hist *History) Undo(grid Grid) {
	hist.Restore(grid, hist.Current().From)
	hist.Idx = max(0, hist.Idx-1)
}

func (hist *History) Branch(from Snapshot, to Snapshot) {
	if hist.Idx < len(hist.States)-1 {
		hist.States = slices.Delete(hist.States, hist.Idx+1, len(hist.States))
	}
	hist.States = append(hist.States, State{From: from, To: to})
	hist.Idx++
}

func (hist *History) Restore(grid Grid, snapshot Snapshot) {
	for pos, cell := range snapshot {
		grid[pos.Row][pos.Col] = cell
	}
}

func (hist *History) Current() State {
	return hist.States[hist.Idx]
}

func MakeSnapshot(grid Grid, selected Selected) Snapshot {
	snapshot := make(Snapshot)

	for pos := range selected {
		cell := grid[pos.Row][pos.Col]
		snapshot[pos] = cell
	}
	return snapshot
}
