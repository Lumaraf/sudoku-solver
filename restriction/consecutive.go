package restriction

import (
	"errors"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

type ConsecutiveRestriction struct {
	Offsets    sudoku.Offsets
	Difference int
}

func (r ConsecutiveRestriction) Name() string {
	return "Consecutive"
}

func (r ConsecutiveRestriction) Validate(s sudoku.Sudoku) error {
outer:
	for _, cell := range s.SolvedArea().Locations {
		d := s.Get(cell)
		mask := (d<<r.Difference | d>>r.Difference) & sudoku.AllDigits
		for cell := range r.Offsets.Locations(cell) {
			if s.Get(cell)&mask != 0 {
				continue outer
			}
		}
		return errors.New("cell has no consecutive neighbour")
	}
	return nil
}
