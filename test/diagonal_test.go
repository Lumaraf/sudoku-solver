package test

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"testing"
)

func TestDiagonal(t *testing.T) {
	SudokuTests[sudoku.Digits9, sudoku.Area9x9]{
		"diagonal": {
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
				"37  98",
				"1",
				"  4",
				" 23 864",
				"",
				"  845 96",
				"      6",
				"        8",
				"   36  97",
			),
			rule.DiagonalRule[sudoku.Digits9, sudoku.Area9x9]{},
		},
	}.Run(t, sudoku.NewSudokuBuilder9x9)
}
