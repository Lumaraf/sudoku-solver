package test

import (
	"context"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"testing"
)

func TestClassic(t *testing.T) {
	sudokuTests{
		"easy": {
			rows: []string{
				" 3       ",
				"   195   ",
				"  8    6 ",
				"8   6    ",
				"4  8    1",
				"    2    ",
				" 6    28 ",
				"   419  5",
				"       7 ",
			},
		},
		"loneliest number": {
			rows: []string{
				"  23 67  ",
				"   4 5   ",
				"3       8",
				"21     97",
				"    1    ",
				"45     13",
				"8       4",
				"   6 2   ",
				"  79 85  ",
			},
		},
		"impossible": {
			rows: []string{
				"8        ",
				"  36     ",
				" 7  9 2  ",
				" 5   7   ",
				"    457  ",
				"   1   3 ",
				"  1    68",
				"  85   1 ",
				" 9    4  ",
			},
		},
	}.Run(t)
}

func BenchmarkClassic(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		s := sudoku.NewSudoku()
		s.SetChainLimit(1)
		//s.EnableGuessing()
		sudokuSpec{
			rows: []string{
				" 3       ",
				"   195   ",
				"  8    6 ",
				"8   6    ",
				"4  8    1",
				"    2    ",
				" 6    28 ",
				"   419  5",
				"       7 ",
			},
		}.setup(s)

		b.StartTimer()
		s.Solve(context.Background())
		b.StopTimer()
	}
}
