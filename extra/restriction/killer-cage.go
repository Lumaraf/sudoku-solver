package restriction

import (
	"errors"
	restriction2 "github.com/lumaraf/sudoku-solver/rule"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

type KillerCageRestriction struct {
	Area sudoku.Area
	Sum  int
}

func (r KillerCageRestriction) Name() string {
	return "KillerCage"
}

func (r KillerCageRestriction) Validate(s sudoku.Sudoku) error {
	sum := 0
	for _, cell := range r.Area.Locations {
		d := s.Get(cell)
		v, isSingle := d.Single()
		if !isSingle {
			return nil
		}
		sum = sum + v
	}
	if sum != r.Sum {
		return errors.New("invalid killer cage sum")
	}
	return nil
}

func AddKillerCageRestriction(s sudoku.Sudoku, area sudoku.Area, sum int) {
	cells := []sudoku.CellLocation{}
	for _, cell := range area.Locations {
		cells = append(cells, cell)
	}
	s.AddRestriction(restriction2.NewUniqueRestriction("cage", cells...))

	s.AddRestriction(KillerCageRestriction{
		Area: area,
		Sum:  sum,
	})
}
