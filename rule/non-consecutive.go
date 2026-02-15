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

	sb.AddChangeProcessor(NonConsecutiveChangeProcessor[D, A]{
		masks: masks,
		offsets: sudoku.Offsets{
			{-1, 0},
			{1, 0},
			{0, -1},
			{0, 1},
		},
	})
	return nil
}

type NonConsecutiveChangeProcessor[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	masks   map[int]D
	offsets sudoku.Offsets
}

func (p NonConsecutiveChangeProcessor[D, A]) Name() string {
	return "Non-Consecutive Change Processor"
}

func (p NonConsecutiveChangeProcessor[D, A]) ProcessChange(s sudoku.Sudoku[D, A], cell sudoku.CellLocation, mask D) error {
	combinedMask := s.NewDigits()
	for v := range mask.Values {
		combinedMask = combinedMask.Or(p.masks[v])
	}
	for _, l := range s.NewAreaFromOffsets(cell, p.offsets).Locations {
		if err := s.Mask(l, combinedMask); err != nil {
			return err
		}
	}
	return nil
}
