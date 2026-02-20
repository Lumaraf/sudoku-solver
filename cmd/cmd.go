package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/strategy"
	"github.com/lumaraf/sudoku-solver/sudoku"
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

	{
		s, err := sudoku.NewSudoku12x12(
			rule.ClassicRules[sudoku.Digits12, sudoku.Area12x12]{},
			rule.GivenDigitsFromString[sudoku.Digits12, sudoku.Area12x12](
				"    A5 2    ",
				"    4   AB3 ",
				"2 93       1",
				"  3 1 8 9 7 ",
				" 4 B3 7     ",
				"8 5   B 34  ",
				"  1C 7   A 9",
				"     B A6 4 ",
				" 5 6 4 3 8  ",
				"3       B9 6",
				" 761   B    ",
				"    2 AC    ",
			),
		)
		if err != nil {
			panic(err)
		}
		//s.SetLogger(sudoku.NewLogger[sudoku.Digits12]())
		slv := s.NewSolver()
		slv.Use(
			strategy.AllStrategies[sudoku.Digits12, sudoku.Area12x12](),
		)
		start := time.Now()
		slv.Solve(context.Background())
		s.Print()
		fmt.Printf("time: %v\nerr: %s\n", time.Since(start), err)
	}
}
