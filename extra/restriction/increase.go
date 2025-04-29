package restriction

import (
	"errors"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

type IncreaseRestriction struct {
	First  sudoku.CellLocation
	Second sudoku.CellLocation
}

func (r IncreaseRestriction) Name() string {
	return "Increase"
}

func (r IncreaseRestriction) Validate(s sudoku.Sudoku) error {
	first, _ := s.Get(r.First).Single()
	second, _ := s.Get(r.Second).Single()
	if first == 0 || second == 0 {
		return nil
	}
	if first >= second {
		return errors.New("invalid increase")
	}
	return nil
}

func AddThermometerRestriction(s sudoku.Sudoku, cells []sudoku.CellLocation) {
	prev := cells[0]
	for _, cell := range cells[1:] {
		s.AddRestriction(IncreaseRestriction{
			First:  prev,
			Second: cell,
		})
		prev = cell
	}
}
