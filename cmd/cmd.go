package main

import (
	"context"
	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func main() {
	s, err := sudoku.NewSudoku9x9(
		restriction.ClassicRestrictions[sudoku.Digits9, sudoku.Area9x9],
		restriction.GivenDigits[sudoku.Digits9, sudoku.Area9x9](
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
	if err != nil {
		panic(err)
	}

	s.SetChainLimit(0)
	s.SolveWith(context.Background())
	s.Print()

	//for r := range sudoku.GetRestrictions[sudoku.Digits9, sudoku.Area9x9, restriction.UniqueRestriction[sudoku.Digits9, sudoku.Area9x9]](s) {
	//	fmt.Printf("%+v\n", r)
	//}
}
