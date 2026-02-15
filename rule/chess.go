package rule

import "github.com/lumaraf/sudoku-solver/sudoku"

type AntiKingRule[D sudoku.Digits[D], A sudoku.Area[A]] struct{}

func (r AntiKingRule[D, A]) Name() string {
	return "anti-king"
}

func (r AntiKingRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	return sb.Use(RelativeExclusionRule[D, A]{
		offsets: sudoku.Offsets{
			{Row: -1, Col: -1},
			{Row: -1, Col: 0},
			{Row: -1, Col: 1},
			{Row: 0, Col: -1},
			{Row: 0, Col: 1},
			{Row: 1, Col: -1},
			{Row: 1, Col: 0},
			{Row: 1, Col: 1},
		},
	})
}

type AntiKnightRule[D sudoku.Digits[D], A sudoku.Area[A]] struct{}

func (r AntiKnightRule[D, A]) Name() string {
	return "anti-knight"
}

func (r AntiKnightRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	return sb.Use(RelativeExclusionRule[D, A]{
		offsets: sudoku.Offsets{
			{Row: 2, Col: 1},
			{Row: 2, Col: -1},
			{Row: -2, Col: 1},
			{Row: -2, Col: -1},
			{Row: 1, Col: 2},
			{Row: 1, Col: -2},
			{Row: -1, Col: 2},
			{Row: -1, Col: -2},
		},
	})
}

type RelativeExclusionRule[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	offsets sudoku.Offsets
}

func (r RelativeExclusionRule[D, A]) Name() string {
	return "relative exclusion"
}

func (r RelativeExclusionRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	for row := 0; row < sb.Size(); row++ {
		for col := 0; col < sb.Size(); col++ {
			cell := sudoku.CellLocation{row, col}
			sb.AddExclusionArea(cell, sb.NewAreaFromOffsets(cell, r.offsets))
		}
	}
	return nil
}
