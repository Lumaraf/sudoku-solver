package main

import (
	"context"
	"fmt"
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/strategy"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func main() {
	sb1 := sudoku.NewSudokuBuilder9x9()
	sb1.Use(
		rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
		rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
			"8 5",
			"  9     2",
			" 7    456",
			"72   8  4",
			"",
			"4   2  8",
			" 63  7",
			"9  1 6",
		),
	)

	sb2 := sudoku.NewSudokuBuilder9x9()
	sb2.Use(
		rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
		rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
			" 154",
			"     3 6",
			" 3    5",
			"  6 9   2",
			"1   8",
			"    4  59",
			" 29 6",
			"   2  7 8",
			"7  5   9",
		),
	)

	mb := sudoku.MultiSudokuBuilder[sudoku.Digits9, sudoku.Area9x9]{}
	mb.Overlap(sb1, sudoku.CellLocation{8, 8}, 3, 3, sb2)

	m, err := mb.Build()
	fmt.Println(err)
	m.Solve(context.Background(), strategy.AllStrategies[sudoku.Digits9, sudoku.Area9x9]())
}
