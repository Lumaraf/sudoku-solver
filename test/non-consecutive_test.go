package test

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"testing"
)

func TestNonConsecutive(t *testing.T) {
	SudokuTests[sudoku.Digits9, sudoku.Area9x9]{
		"one": {
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
				"        5",
				" 1    7  ",
				"7        ",
				"    7  59",
				"         ",
				"42  9    ",
				"        8",
				"  1    7 ",
				"8        ",
			),
			rule.NonConsecutiveRule[sudoku.Digits9, sudoku.Area9x9]{},
		},
	}.Run(t, sudoku.NewSudokuBuilder9x9)
}
