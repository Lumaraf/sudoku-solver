package main

import (
	"context"
	"fmt"
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/strategy"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"time"
)

func main() {
	{
		s, err := sudoku.NewSudoku6x6(
			rule.ClassicRules[sudoku.Digits6, sudoku.Area6x6]{},
			rule.GivenDigitsFromString[sudoku.Digits6, sudoku.Area6x6](
				"  5  2",
				"6     ",
				"4    5",
				"5   4 ",
				"  12  ",
				"     1",
			),
		)
		if err != nil {
			panic(err)
		}
		//s.SetLogger(sudoku.NewLogger[sudoku.Digits6]())
		slv := s.NewSolver()
		slv.SetChainLimit(0)
		slv.Use(
			strategy.AllStrategies[sudoku.Digits6, sudoku.Area6x6](),
		)
		start := time.Now()
		err = slv.Solve(context.Background())
		s.Print()
		fmt.Printf("time: %v\nerr: %s\n", time.Since(start), err)
	}

	{
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
		if err != nil {
			panic(err)
		}
		//s.SetLogger(sudoku.NewLogger[sudoku.Digits9]())
		slv := s.NewSolver()
		slv.Use(
			strategy.AllStrategies[sudoku.Digits9, sudoku.Area9x9](),
		)
		start := time.Now()
		slv.Solve(context.Background())
		s.Print()
		fmt.Printf("time: %v\nerr: %s\n", time.Since(start), err)
	}
}
