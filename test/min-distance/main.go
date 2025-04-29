package main

import (
	"context"
	"fmt"
	restriction2 "github.com/lumaraf/sudoku-solver/extra/restriction"
	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func main() {
	s := sudoku.NewSudoku()
	restriction.AddClassicRestrictions(s)

	//masks := map[int]sudoku.Digits{}
	//for n := 1; n <= 9; n++ {
	//	mask := sudoku.AllDigits
	//	for i := 1; i <= 9; i++ {
	//		if i+n == 5 || i+n == 10 {
	//			continue
	//		}
	//		mask.RemoveOption(i)
	//	}
	//	masks[n] = mask
	//}
	//offsets := sudoku.Offsets{
	//	{-1, 0},
	//	{1, 0},
	//	{0, -1},
	//	{0, 1},
	//}
	//offsets = append(offsets, restriction.KnightsMove...)
	//s.AddRestriction(restriction.AntiRelationRestriction{
	//	Offsets: offsets,
	//	Masks:   masks,
	//})
	//restriction.AddAntiKingRestriction(s)
	for n := 2; n <= 6; n++ {
		restriction2.AddKillerCageRestriction(
			s,
			sudoku.NewArea(
				sudoku.CellLocation{n, 2},
				sudoku.CellLocation{n, 3},
				sudoku.CellLocation{n, 4},
				sudoku.CellLocation{n, 5},
				sudoku.CellLocation{n, 6},
			),
			25,
		)
		restriction2.AddKillerCageRestriction(
			s,
			sudoku.NewArea(
				sudoku.CellLocation{2, n},
				sudoku.CellLocation{3, n},
				sudoku.CellLocation{4, n},
				sudoku.CellLocation{5, n},
				sudoku.CellLocation{6, n},
			),
			25,
		)
	}

	err := s.Solve(context.Background())
	fmt.Println(err)
	s.Print()
}
