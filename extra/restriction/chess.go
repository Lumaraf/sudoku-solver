package restriction

import "github.com/lumaraf/sudoku-solver/sudoku"

var KingsMove = sudoku.Offsets{
	{-1, -1},
	{-1, 1},
	{1, -1},
	{1, 1},
}
var KnightsMove = sudoku.Offsets{
	{-1, -2}, {-2, -1},
	{-1, 2}, {-2, 1},
	{1, -2}, {2, -1},
	{1, 2}, {2, 1},
}
var QueensMove = make(sudoku.Offsets, 0, 32)

func init() {
	for offset := 1; offset <= 8; offset++ {
		QueensMove = append(
			QueensMove,
			sudoku.Offset{-offset, -offset},
			sudoku.Offset{-offset, offset},
			sudoku.Offset{offset, -offset},
			sudoku.Offset{offset, offset},
		)
	}
}

func AddAntiKingRestriction(s sudoku.Sudoku) {
	s.AddRestriction(AntiRelationRestriction{
		Offsets: KingsMove,
		Masks: map[int]sudoku.Digits{
			1: sudoku.NewDigits(1),
			2: sudoku.NewDigits(2),
			3: sudoku.NewDigits(3),
			4: sudoku.NewDigits(4),
			5: sudoku.NewDigits(5),
			6: sudoku.NewDigits(6),
			7: sudoku.NewDigits(7),
			8: sudoku.NewDigits(8),
			9: sudoku.NewDigits(9),
		},
	})
}

func AddAntiKnighRestriction(s sudoku.Sudoku) {
	s.AddRestriction(AntiRelationRestriction{
		Offsets: KnightsMove,
		Masks: map[int]sudoku.Digits{
			1: sudoku.NewDigits(1),
			2: sudoku.NewDigits(2),
			3: sudoku.NewDigits(3),
			4: sudoku.NewDigits(4),
			5: sudoku.NewDigits(5),
			6: sudoku.NewDigits(6),
			7: sudoku.NewDigits(7),
			8: sudoku.NewDigits(8),
			9: sudoku.NewDigits(9),
		},
	})
}
