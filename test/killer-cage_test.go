package test

import (
	"testing"

	extraRule "github.com/lumaraf/sudoku-solver/extra/rule"
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func TestKillerCage(t *testing.T) {
	SudokuTests[sudoku.Digits9, sudoku.Area9x9]{
		"daily killer #6502": {
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			extraRule.KillerCageRulesFromString[sudoku.Digits9, sudoku.Area9x9](
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
		//"besties 2": { // not yet solvable
		//	rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
		//	extraRule.AntiKnightRule[sudoku.Digits9, sudoku.Area9x9]{},
		//	extraRule.KillerCageRulesFromString[sudoku.Digits9, sudoku.Area9x9](
		//		[]string{
		//			"         ",
		//			"         ",
		//			"  AB CC  ",
		//			"  ABHIJ K",
		//			"    HIJ K",
		//			"  DD FG  ",
		//			"  EE FG  ",
		//			"         ",
		//			"   LL    ",
		//		},
		//		map[rune]int{
		//			'A': 7,
		//			'B': 7,
		//			'C': 7,
		//			'D': 7,
		//			'E': 7,
		//			'F': 7,
		//			'G': 7,
		//			'H': 13,
		//			'I': 13,
		//			'J': 13,
		//			'K': 13,
		//			'L': 13,
		//		},
		//	),
		//	extraRule.AreaSumRule[sudoku.Digits9, sudoku.Area9x9]{
		//		Area: func() []sudoku.CellLocation {
		//			area := make([]sudoku.CellLocation, 0, 9)
		//			for n := 0; n < 9; n++ {
		//				area = append(area, sudoku.CellLocation{n, n})
		//			}
		//			return area
		//		}(),
		//		Sum: 39,
		//	},
		//},
	}.Run(t, sudoku.NewSudokuBuilder9x9)
}
