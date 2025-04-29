package restriction

import (
	"errors"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type IndexValueLinkRestriction struct {
	Cells []sudoku.CellLocation
	Index int
}

func (r IndexValueLinkRestriction) Name() string {
	return "IndexValueLink"
}

func (r IndexValueLinkRestriction) Validate(s sudoku.Sudoku) error {
	d := s.Get(r.Cells[r.Index])
	for v := range d.Values {
		if v > len(r.Cells) {
			continue
		}
		if s.Get(r.Cells[v-1]).CanContain(r.Index + 1) {
			return nil
		}
	}
	return errors.New("wrong value in linked cell")
}

func (r IndexValueLinkRestriction) ProcessSolve(s sudoku.Sudoku, cell sudoku.CellLocation) error {
	d := s.Get(cell)
	v, _ := d.Single()
	if cell == r.Cells[r.Index] {
		return s.Set(r.Cells[v-1], r.Index+1)
	}
	if v == r.Index+1 {
		for i, c := range r.Cells {
			if c == cell {
				return s.Set(r.Cells[r.Index], i+1)
			}
		}
	}
	return nil
}

func Add159Restriction(s sudoku.Sudoku) {
	for row := 0; row < 9; row++ {
		cells := make([]sudoku.CellLocation, 9)
		for col := 0; col < 9; col++ {
			cells[col] = sudoku.CellLocation{Row: row, Col: col}
		}
		s.AddRestriction(IndexValueLinkRestriction{Cells: cells, Index: 0})
		s.AddRestriction(IndexValueLinkRestriction{Cells: cells, Index: 4})
		s.AddRestriction(IndexValueLinkRestriction{Cells: cells, Index: 8})
	}
}
