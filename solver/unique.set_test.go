package solver

import (
	"github.com/lumaraf/sudoku-solver/sudoku"
	"testing"
)

func TestUniqueSetSolver_Solve(t *testing.T) {
	s := sudoku.NewSudoku()
	s.Mask(sudoku.CellLocation{0, 0}, sudoku.NewDigits(1, 2, 3))
	s.Mask(sudoku.CellLocation{0, 1}, sudoku.NewDigits(1, 2, 4))
	for n := 2; n < 9; n++ {
		s.RemoveMask(sudoku.CellLocation{0, n}, sudoku.NewDigits(1, 2))
	}
	slv := UniqueSetSolver{
		Cells: []sudoku.CellLocation{
			{0, 0},
			{0, 1},
			{0, 2},
			{0, 3},
			{0, 4},
			{0, 5},
			{0, 6},
			{0, 7},
			{0, 8},
		},
		Area: sudoku.Area{},
	}
	slv.Solve(s)
	s.Print()
}
