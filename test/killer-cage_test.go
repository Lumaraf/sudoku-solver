package test

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"testing"
)

func TestKillerCage(t *testing.T) {
	SudokuTests[sudoku.Digits9, sudoku.Area9x9]{
		"daily killer #6502": {
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			rule.KillerCageRulesFromString[sudoku.Digits9, sudoku.Area9x9](
				[]string{
					"AAAAGOOVV",
					"BBHHGPOVV",
					"CBIIPPWWW",
					"CDIIQQQQX",
					"DDJJRRRRX",
					"EKKKSTTYX",
					"ELLLSTTYX",
					"FFMMSUUYZ",
					"FFNNNNUZZ",
				},
				map[rune]int{
					'A': 11,
					'B': 21,
					'C': 12,
					'D': 13,
					'E': 7,
					'F': 19,
					'G': 12,
					'H': 9,
					'I': 25,
					'J': 9,
					'K': 9,
					'L': 19,
					'M': 12,
					'N': 24,
					'O': 16,
					'P': 12,
					'Q': 21,
					'R': 15,
					'S': 19,
					'T': 25,
					'U': 12,
					'V': 19,
					'W': 19,
					'X': 19,
					'Y': 14,
					'Z': 12,
				},
			),
		},
	}.Run(t, sudoku.NewSudokuBuilder9x9)
}
