package test

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"testing"
)

func TestChess(t *testing.T) {
	SudokuTests[sudoku.Digits9, sudoku.Area9x9]{
		"anti knight": {
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
				" 5   9   ",
				"8        ",
				"     3 4 ",
				"7 8   1 9",
				"         ",
				"    3    ",
				"         ",
				"  3 1   8",
				"   9   2 ",
			),
			rule.AntiKnightRule[sudoku.Digits9, sudoku.Area9x9]{},
		},
		"miracle": {
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
				"         ",
				"         ",
				"         ",
				"         ",
				"  1      ",
				"      2  ",
				"         ",
				"         ",
				"         ",
			),
			rule.AntiKingRule[sudoku.Digits9, sudoku.Area9x9]{},
			rule.AntiKnightRule[sudoku.Digits9, sudoku.Area9x9]{},
			rule.NonConsecutiveRule[sudoku.Digits9, sudoku.Area9x9]{},
		},
		//"159": {
		//	Rows: []string{
		//		"         ",
		//		"         ",
		//		"    E    ",
		//		"E   E    ",
		//		"E   E    ",
		//		"E   E    ",
		//		"E        ",
		//		"         ",
		//		"         ",
		//	},
		//	AntiKnight: true,
		//	Rule159:    true,
		//},
	}.Run(t, sudoku.NewSudokuBuilder9x9)
}
