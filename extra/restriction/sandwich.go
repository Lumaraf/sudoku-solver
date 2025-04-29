package restriction

import (
	"errors"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

type SandwichRestriction struct {
	Sum   int
	Cells []sudoku.CellLocation
}

func (r SandwichRestriction) Name() string {
	return "Sandwich"
}

func NewRowSandwichRestriction(row, sum int) SandwichRestriction {
	cells := make([]sudoku.CellLocation, 0, 9)
	for col := 0; col < 9; col++ {
		cells = append(cells, sudoku.CellLocation{row, col})
	}
	return SandwichRestriction{
		Sum:   sum,
		Cells: cells,
	}
}

func NewColSandwichRestriction(col, sum int) SandwichRestriction {
	cells := make([]sudoku.CellLocation, 0, 9)
	for row := 0; row < 9; row++ {
		cells = append(cells, sudoku.CellLocation{row, col})
	}
	return SandwichRestriction{
		Sum:   sum,
		Cells: cells,
	}
}

func (r SandwichRestriction) Validate(s sudoku.Sudoku) error {
	inSandwich := false
	sum := 0
	for _, cell := range r.Cells {
		d := s.Get(cell)
		v, isSingle := d.Single()
		switch v {
		case 1, 9:
			inSandwich = !inSandwich
			continue
		}
		if inSandwich {
			if !isSingle {
				return nil
			}
			sum = sum + v
		}
	}
	if !inSandwich && sum != 0 && sum != r.Sum {
		return errors.New("invalid sandwich sum")
	}
	return nil
}
