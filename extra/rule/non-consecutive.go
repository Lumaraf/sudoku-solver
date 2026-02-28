package rule

import (
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type NonConsecutiveRule[D sudoku.Digits[D], A sudoku.Area[A]] struct{}

func (r NonConsecutiveRule[D, A]) Name() string {
	return "non-consecutive"
}

func (r NonConsecutiveRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	masks := make(map[int]D)
	for v := 1; v <= sb.Size(); v++ {
		digits := make([]int, 0, 2)
		if v > 1 {
			digits = append(digits, v-1)
		}
		if v < sb.Size() {
			digits = append(digits, v+1)
		}
		masks[v] = sb.NewDigits(digits...).Not()
	}

	for v := 1; v <= sb.Size(); v++ {
		sb.AddOffsetMask(v, sudoku.Offset{-1, 0}, masks[v])
		sb.AddOffsetMask(v, sudoku.Offset{1, 0}, masks[v])
		sb.AddOffsetMask(v, sudoku.Offset{0, -1}, masks[v])
		sb.AddOffsetMask(v, sudoku.Offset{0, 1}, masks[v])
	}
	return nil
}
