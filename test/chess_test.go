package test

import (
	"testing"

	extraRule "github.com/lumaraf/sudoku-solver/extra/rule"
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
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
			extraRule.AntiKnightRule[sudoku.Digits9, sudoku.Area9x9]{},
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
			extraRule.AntiKingRule[sudoku.Digits9, sudoku.Area9x9]{},
			extraRule.AntiKnightRule[sudoku.Digits9, sudoku.Area9x9]{},
			extraRule.NonConsecutiveRule[sudoku.Digits9, sudoku.Area9x9]{},
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
