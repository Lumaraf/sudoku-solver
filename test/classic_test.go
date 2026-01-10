package test

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClassic(t *testing.T) {
	SudokuTests[sudoku.Digits9, sudoku.Area9x9]{
		"easy": {
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
				" 3       ",
				"   195   ",
				"  8    6 ",
				"8   6    ",
				"4  8    1",
				"    2    ",
				" 6    28 ",
				"   419  5",
				"       7 ",
			),
		},
		"loneliest number": {
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
				"  23 67  ",
				"   4 5   ",
				"3       8",
				"21     97",
				"    1    ",
				"45     13",
				"8       4",
				"   6 2   ",
				"  79 85  ",
			),
		},
		"unsolvable #680": {
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
				"  3",
				"4    5 1",
				"   7  96",
				"     253",
				"6       9",
				" 524  7",
				" 17  8",
				"28 6    5",
				"      8",
			),
		},
		//"impossible": {
		//	rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
		//	rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
		//		"8        ",
		//		"  36     ",
		//		" 7  9 2  ",
		//		" 5   7   ",
		//		"    457  ",
		//		"   1   3 ",
		//		"  1    68",
		//		"  85   1 ",
		//		" 9    4  ",
		//	),
		//},
	}.Run(t, sudoku.NewSudokuBuilder9x9)
}

func BenchmarkClassic(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		s, err := sudoku.NewSudoku9x9(
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
				" 3       ",
				"   195   ",
				"  8    6 ",
				"8   6    ",
				"4  8    1",
				"    2    ",
				" 6    28 ",
				"   419  5",
				"       7 ",
			),
		)
		assert.NoError(b, err)

		b.StartTimer()
		assert.NoError(b, s.NewSolver().Solve(b.Context()))
		b.StopTimer()
	}
}
