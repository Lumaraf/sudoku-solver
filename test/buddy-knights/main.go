package main

import (
	"context"
	"fmt"
	restriction2 "github.com/lumaraf/sudoku-solver/extra/restriction"
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func main() {
	s := sudoku.NewSudoku()
	rule.AddClassicRestrictions(s)

	s.AddRestriction(restriction2.AntiRelationRestriction{
		Offsets: restriction2.KnightsMove,
		Masks: map[int]sudoku.Digits{
			//1: sudoku.NewDigits(2),
			2: sudoku.NewDigits(1, 3),
			3: sudoku.NewDigits(2, 4),
			4: sudoku.NewDigits(3, 5),
			5: sudoku.NewDigits(4, 6),
			6: sudoku.NewDigits(5, 7),
			7: sudoku.NewDigits(6, 8),
			8: sudoku.NewDigits(7, 9),
			//9: sudoku.NewDigits(8),
		},
	})

	//s.AddRestriction(rule.RelationRestriction{
	//	Area:    sudoku.NewArea().Not(),
	//	Offsets: rule.KnightsMove,
	//	Masks: map[int]sudoku.Digits{
	//		1: sudoku.NewDigits(1),
	//		//2: sudoku.NewDigits(2),
	//		3: sudoku.NewDigits(3),
	//		//4: sudoku.NewDigits(4),
	//		5: sudoku.NewDigits(5),
	//		//6: sudoku.NewDigits(6),
	//		7: sudoku.NewDigits(7),
	//		//8: sudoku.NewDigits(8),
	//		9: sudoku.NewDigits(9),
	//	},
	//})

	//for i := 0; i < 9; i++ {
	//	s.Set(sudoku.CellLocation{0, i}, i+1)
	//}

	err := s.Solve(context.Background())
	fmt.Println(err)
	s.Print()
}
