package main

import (
	"context"
	"fmt"
	"github.com/lumaraf/sudoku-solver/restriction"
	_ "github.com/lumaraf/sudoku-solver/solver"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

//╔═══════╤═══════╤═══════╦═══════╤═══════╤═══════╦═══════╤═══════╤═══════╗
//║     3 │       │       ║   2   │       │       ║ 1     │       │       ║
//║       │       │   5   ║       │       │       X       │     6 X 4     ║
//║       │   8   │       ║       │ 7     │     9 ║       │       │       ║
//╟───────┼───X───┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢
//║       │   2   │       ║       │ 1     │       ║       │     3 │       ║
//║     6 │       │       ║ 4     V       │   5   ║       │       │       ║
//║       │       │ 7     ║       │       │       ║   8   │       │     9 ║
//╟───────┼───────┼───────╫───X───┼───────┼───────╫───────┼───X───┼───────╢
//║       │       │ 1     ║       │     3 │       ║       │       │   2   ║
//║       │ 4     V       ║     6 │       │       ║   5   │       │       ║
//║     9 │       │       ║       │       │   8   ║       │ 7     │       ║
//╠═══════╪═══════╪═══════╬═══════╪═══════╪═══════╬═══════╪═══════╪═══X═══╣
//║   2   │       │     3 ║ 1     │       │       ║       │       │       ║
//║       │       │       ║       │     6 X 4     ║       │   5   │       ║
//║       │     9 │       ║       │       │       ║ 7     │       │   8   ║
//╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢
//║       │       │       ║     3 │       │   2   ║       │ 1     │       ║
//║       │   5   │       ║       │       │       ║ 4     V       │     6 ║
//║ 7     │       │   8   ║       │     9 │       ║       │       │       ║
//╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───X───┼───────╢
//║       │ 1     │       ║       │       │       ║   2   │       │     3 ║
//║ 4     │       │     6 ║       │   5   │       ║       │       │       ║
//║       │       │       ║   8   │       │ 7     ║       │     9 │       ║
//╠═══V═══╪═══════╪═══════╬═══════╪═══════╪═══X═══╬═══════╪═══════╪═══X═══╣
//║ 1     │       │   2   ║       │       │     3 ║       │       │       ║
//║       │     6 │       ║   5   │       │       ║       │ 4     │       ║
//║       │       │       ║       │   8   │       ║     9 │       │ 7     ║
//╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢
//║       │     3 │       ║       │       │ 1     ║       │   2   │       ║
//║       │       │       ║       │ 4     │       ║     6 │       │   5   ║
//║   8   │       │     9 ║ 7     │       │       ║       │       │       ║
//╟───────┼───X───┼───────╫───────┼───────┼───────╫───────┼───X───┼───────╢
//║       │       │       ║       │   2   │       ║     3 │       │ 1     ║
//║   5   │       │ 4     ║       │       │     6 ║       │       │       ║
//║       │ 7     │       ║     9 │       │       ║       │   8   │       ║
//╚═══════╧═══════╧═══════╩═══════╧═══════╧═══════╩═══════╧═══════╧═══════╝

func main() {
	s := sudoku.NewSudoku()
	s.SetChainLimit(3)
	restriction.AddClassicRestrictions(s)
	restriction.AddNonConsecutiveRestriction(s)
	s.AddRestriction(restriction.AntiRelationRestriction{
		Offsets: sudoku.Offsets{
			{Row: 0, Col: -1},
			{Row: -1, Col: 0},
			{Row: 0, Col: 1},
			{Row: 1, Col: 0},
		},
		Masks: map[int]sudoku.Digits{
			1: sudoku.NewDigits(4, 9),
			2: sudoku.NewDigits(3, 8),
			3: sudoku.NewDigits(2, 7),
			4: sudoku.NewDigits(1, 6),
			6: sudoku.NewDigits(4),
			7: sudoku.NewDigits(3),
			8: sudoku.NewDigits(2),
			9: sudoku.NewDigits(1),
		},
		Exceptions: map[sudoku.Area]bool{
			// rows
			sudoku.NewArea(sudoku.CellLocation{Row: 0, Col: 5}, sudoku.CellLocation{Row: 0, Col: 6}): true,
			sudoku.NewArea(sudoku.CellLocation{Row: 0, Col: 7}, sudoku.CellLocation{Row: 0, Col: 8}): true,

			sudoku.NewArea(sudoku.CellLocation{Row: 1, Col: 3}, sudoku.CellLocation{Row: 1, Col: 4}): true,

			sudoku.NewArea(sudoku.CellLocation{Row: 2, Col: 1}, sudoku.CellLocation{Row: 2, Col: 2}): true,

			sudoku.NewArea(sudoku.CellLocation{Row: 3, Col: 4}, sudoku.CellLocation{Row: 3, Col: 5}): true,

			sudoku.NewArea(sudoku.CellLocation{Row: 4, Col: 6}, sudoku.CellLocation{Row: 4, Col: 7}): true,

			// cols
			sudoku.NewArea(sudoku.CellLocation{Row: 5, Col: 0}, sudoku.CellLocation{Row: 6, Col: 0}): true,

			sudoku.NewArea(sudoku.CellLocation{Row: 0, Col: 1}, sudoku.CellLocation{Row: 1, Col: 1}): true,
			sudoku.NewArea(sudoku.CellLocation{Row: 7, Col: 1}, sudoku.CellLocation{Row: 8, Col: 1}): true,

			sudoku.NewArea(sudoku.CellLocation{Row: 1, Col: 3}, sudoku.CellLocation{Row: 2, Col: 3}): true,

			sudoku.NewArea(sudoku.CellLocation{Row: 5, Col: 5}, sudoku.CellLocation{Row: 6, Col: 5}): true,

			sudoku.NewArea(sudoku.CellLocation{Row: 1, Col: 7}, sudoku.CellLocation{Row: 2, Col: 7}): true,
			sudoku.NewArea(sudoku.CellLocation{Row: 4, Col: 7}, sudoku.CellLocation{Row: 5, Col: 7}): true,
			sudoku.NewArea(sudoku.CellLocation{Row: 7, Col: 7}, sudoku.CellLocation{Row: 8, Col: 7}): true,

			sudoku.NewArea(sudoku.CellLocation{Row: 2, Col: 8}, sudoku.CellLocation{Row: 3, Col: 8}): true,
			sudoku.NewArea(sudoku.CellLocation{Row: 5, Col: 8}, sudoku.CellLocation{Row: 6, Col: 8}): true,
		},
	})

	AddXVRestrictions(
		s,
		[][2]sudoku.CellLocation{
			{{Row: 0, Col: 5}, {Row: 0, Col: 6}}, //x
			{{Row: 0, Col: 7}, {Row: 0, Col: 8}}, //x
			{{Row: 3, Col: 4}, {Row: 3, Col: 5}}, //x

			{{Row: 0, Col: 1}, {Row: 1, Col: 1}}, //x
			{{Row: 7, Col: 1}, {Row: 8, Col: 1}}, //x
			{{Row: 1, Col: 3}, {Row: 2, Col: 3}}, //x
			{{Row: 5, Col: 5}, {Row: 6, Col: 5}}, //x
			{{Row: 1, Col: 7}, {Row: 2, Col: 7}}, //x
			{{Row: 4, Col: 7}, {Row: 5, Col: 7}}, //x
			{{Row: 7, Col: 7}, {Row: 8, Col: 7}}, //x
			{{Row: 2, Col: 8}, {Row: 3, Col: 8}}, //x
			{{Row: 5, Col: 8}, {Row: 6, Col: 8}}, //x
		},
		[][2]sudoku.CellLocation{
			{{Row: 1, Col: 3}, {Row: 1, Col: 4}}, //v
			{{Row: 2, Col: 1}, {Row: 2, Col: 2}}, //v
			{{Row: 4, Col: 6}, {Row: 4, Col: 7}}, //v

			{{Row: 5, Col: 0}, {Row: 6, Col: 0}}, //v
		},
	)

	//for s := range s.GuessSolutions(context.Background(), nil) {
	//	s.Print()
	//	time.Sleep(1 * time.Second)
	//}

	s.SetLogger(sudoku.NewLogger())
	fmt.Println(s.Solve(context.Background()))
	s.Print()
}

func AddXVRestrictions(s sudoku.Sudoku, x [][2]sudoku.CellLocation, v [][2]sudoku.CellLocation) {
	for _, cells := range x {
		restriction.AddKillerCageRestriction(s, sudoku.NewArea(cells[0], cells[1]), 10)
	}
	for _, cells := range v {
		restriction.AddKillerCageRestriction(s, sudoku.NewArea(cells[0], cells[1]), 5)
	}
}
