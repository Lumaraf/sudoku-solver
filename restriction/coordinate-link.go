package restriction

import (
	"errors"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

type CoordinateLinkRestriction struct{}

func (r CoordinateLinkRestriction) Validate(s sudoku.Sudoku) error {
	for _, cell := range s.SolvedArea().Locations {
		v, _ := s.Get(cell).Single()

		if !s.Get(sudoku.CellLocation{cell.Row, v - 1}).CanContain(cell.Col + 1) {
			return errors.New("coordinate link broken")
		}
		if !s.Get(sudoku.CellLocation{v - 1, cell.Col}).CanContain(cell.Row + 1) {
			return errors.New("coordinate link broken")
		}
	}
	return nil
}

func (r CoordinateLinkRestriction) ProcessSolve(s sudoku.Sudoku, cell sudoku.CellLocation) error {
	v, _ := s.Get(cell).Single()
	if err := s.Set(sudoku.CellLocation{cell.Row, v - 1}, cell.Col+1); err != nil {
		return err
	}
	if err := s.Set(sudoku.CellLocation{v - 1, cell.Col}, cell.Row+1); err != nil {
		return err
	}
	return nil
}
