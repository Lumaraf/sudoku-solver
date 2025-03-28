package restriction

import (
	"errors"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

type EqualRestriction struct {
	Cells []sudoku.CellLocation
}

func (r EqualRestriction) Name() string {
	return "Equal"
}

func (r EqualRestriction) Validate(s sudoku.Sudoku) error {
	mask := sudoku.AllDigits
	for _, cell := range r.Cells {
		mask = mask & s.Get(cell)
	}
	if mask == 0 {
		return errors.New("cells must be equal")
	}
	return nil
}
