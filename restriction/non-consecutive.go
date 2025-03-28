package restriction

import (
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func AddNeighbourOffsetRestriction(s sudoku.Sudoku, offset int) {
	masks := map[int]sudoku.Digits{}
	for n := 1; n <= 9; n++ {
		d := sudoku.NewDigits(n)
		masks[n] = (d<<offset | d>>offset) & sudoku.AllDigits
	}
	s.AddRestriction(AntiRelationRestriction{
		Offsets: sudoku.Offsets{
			{-1, 0},
			{1, 0},
			{0, -1},
			{0, 1},
		},
		Masks: masks,
	})
}

func AddNonConsecutiveRestriction(s sudoku.Sudoku) {
	AddNeighbourOffsetRestriction(s, 1)
}
