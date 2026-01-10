package strategy

import (
	"github.com/lumaraf/sudoku-solver/sudoku"
	"testing"
)

func TestUniqueSetStrategy_Solve(t *testing.T) {
	s, err := sudoku.NewSudoku9x9()
	if err != nil {
		t.Fatalf("failed to create sudoku: %v", err)
	}
	s.Mask(sudoku.CellLocation{0, 0}, s.NewDigits(1, 2, 3))
	s.Mask(sudoku.CellLocation{0, 1}, s.NewDigits(1, 2, 4))
	for n := 2; n < 9; n++ {
		s.RemoveMask(sudoku.CellLocation{0, n}, s.NewDigits(1, 2))
	}
	strategy := UniqueSetStrategy[sudoku.Digits9, sudoku.Area9x9]{
		Area: s.Row(0),
	}
	_, err = strategy.Solve(s)
	if err != nil {
		t.Errorf("UniqueSetStrategy.Solve failed: %v", err)
	}
	s.Print()
}
