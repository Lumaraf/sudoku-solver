package main

import (
	"context"
	"fmt"
	"github.com/lumaraf/sudoku-solver/restriction"
	_ "github.com/lumaraf/sudoku-solver/strategy"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"math"
)

func main() {
	s := sudoku.NewSudoku()
	s.SetChainLimit(0)
	restriction.AddClassicRestrictions(s)
	//// anti xv restriction
	//s.AddRestriction(restriction.AntiRelationRestriction{
	//	Offsets: sudoku.Offsets{
	//		{Row: 0, Col: -1},
	//		{Row: -1, Col: 0},
	//		{Row: 0, Col: 1},
	//		{Row: 1, Col: 0},
	//	},
	//	Masks: map[int]sudoku.Digits{
	//		1: sudoku.NewDigits(4, 9),
	//		2: sudoku.NewDigits(3, 8),
	//		3: sudoku.NewDigits(2, 7),
	//		4: sudoku.NewDigits(1, 6),
	//		6: sudoku.NewDigits(4),
	//		7: sudoku.NewDigits(3),
	//		8: sudoku.NewDigits(2),
	//		9: sudoku.NewDigits(1),
	//	},
	//})
	//restriction.AddNonConsecutiveRestriction(s)

	offsets := sudoku.Offsets{
		{Row: 0, Col: -1},
		{Row: -1, Col: 0},
		{Row: 0, Col: 1},
		{Row: 1, Col: 0},
	}
	avoidMap := map[sudoku.Digits]sudoku.Digits{
		//// kropki
		//sudoku.NewDigits(1): sudoku.NewDigits(2),
		//sudoku.NewDigits(2): sudoku.NewDigits(1, 3),
		//sudoku.NewDigits(3): sudoku.NewDigits(2, 4, 6),
		//sudoku.NewDigits(4): sudoku.NewDigits(2, 3, 5, 8),
		//sudoku.NewDigits(5): sudoku.NewDigits(4, 6),
		//sudoku.NewDigits(6): sudoku.NewDigits(3, 5, 7),
		//sudoku.NewDigits(7): sudoku.NewDigits(6, 8),
		//sudoku.NewDigits(8): sudoku.NewDigits(4, 7, 9),
		//sudoku.NewDigits(9): sudoku.NewDigits(8),

		// prime sums
		sudoku.NewDigits(1): sudoku.NewDigits(2, 4, 6),
		sudoku.NewDigits(2): sudoku.NewDigits(1, 3, 5, 9),
		sudoku.NewDigits(3): sudoku.NewDigits(2, 4, 8),
		sudoku.NewDigits(4): sudoku.NewDigits(1, 3, 7, 9),
		sudoku.NewDigits(5): sudoku.NewDigits(2, 6, 8),
		sudoku.NewDigits(6): sudoku.NewDigits(1, 5, 7),
		sudoku.NewDigits(7): sudoku.NewDigits(4, 6),
		sudoku.NewDigits(8): sudoku.NewDigits(3, 5, 9),
		sudoku.NewDigits(9): sudoku.NewDigits(2, 4, 8),

		//// xv
		//sudoku.NewDigits(1): sudoku.NewDigits(4, 9),
		//sudoku.NewDigits(2): sudoku.NewDigits(3, 8),
		//sudoku.NewDigits(3): sudoku.NewDigits(2, 7),
		//sudoku.NewDigits(4): sudoku.NewDigits(1, 6),
		//sudoku.NewDigits(6): sudoku.NewDigits(4),
		//sudoku.NewDigits(7): sudoku.NewDigits(3),
		//sudoku.NewDigits(8): sudoku.NewDigits(2),
		//sudoku.NewDigits(9): sudoku.NewDigits(1),
	}

	bestCount := math.MaxInt
	for s := range s.GuessSolutions(context.Background(), func(s sudoku.Sudoku) (sudoku.CellLocation, sudoku.Values) {
		cell := s.SolvedArea().Not().RandomLocation()

		avoidMask := sudoku.NewDigits()
		for neighbor := range offsets.Locations(cell) {
			if d := s.Get(neighbor); d.Count() == 1 {
				avoidMask = avoidMask | avoidMap[d]
			}
		}

		d := s.Get(cell)
		return cell, func(yield func(int) bool) {
			for v := range d.Values {
				if !avoidMask.CanContain(v) {
					if !yield(v) {
						return
					}
				}
			}
			for v := range d.Values {
				if avoidMask.CanContain(v) {
					if !yield(v) {
						return
					}
				}
			}
		}
	}) {
		count := countRelations(s, avoidMap)
		if count < bestCount {
			bestCount = count
			s.Print()
			fmt.Println("new best count", count)
		}
	}
}

func countRelations(s sudoku.Sudoku, avoidMap map[sudoku.Digits]sudoku.Digits) int {
	c := 0
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			d := s.Get(sudoku.CellLocation{Row: row, Col: col})
			if avoidMap[d]&s.Get(sudoku.CellLocation{Row: row, Col: col + 1}) != 0 {
				c++
			}
			if avoidMap[d]&s.Get(sudoku.CellLocation{Row: row + 1, Col: col}) != 0 {
				c++
			}
		}
	}
	return c
}
