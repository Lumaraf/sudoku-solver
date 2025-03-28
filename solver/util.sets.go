package solver

import (
	"errors"
	"math/bits"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

var sets = make(map[int][][]int)

func init() {
	for c := 2; c <= 9; c++ {
		buf := make([][][]int, c)
		for mask := 1; mask < (1 << c); mask++ {
			ones := bits.OnesCount16(uint16(mask))
			indices := make([]int, 0, ones)
			for b := 0; b < c; b++ {
				if mask&(1<<b) != 0 {
					indices = append(indices, b)
				}
			}
			buf[ones-1] = append(buf[ones-1], indices)
		}
		for _, b := range buf {
			sets[c] = append(sets[c], b...)
		}
	}
}

type Set struct {
	Indices []int
	Mask    sudoku.Digits
}

func FindSets(cells []sudoku.Digits) ([]Set, error) {
	result := make([]Set, 0, 2)
	used := map[int]bool{}
outer:
	for _, set := range sets[len(cells)] {
		mask := sudoku.Digits(0)
		for _, index := range set {
			if used[index] {
				continue outer
			}
			mask = mask | cells[index]
		}
		if bits.OnesCount16(uint16(mask)) == len(set) {
			lastIndex := 0
			for _, index := range set {
				for n := lastIndex; n < index; n++ {
					cells[n] = cells[n] & ^mask
					if cells[n] == 0 {
						return nil, errors.New("invalid sets")
					}
				}
				lastIndex = index + 1
			}
			for n := lastIndex; n < len(cells); n++ {
				cells[n] = cells[n] & ^mask
				if cells[n] == 0 {
					return nil, errors.New("invalid sets")
				}
			}

			for _, index := range set {
				used[index] = true
			}
			result = append(result, Set{
				Indices: set,
				Mask:    mask,
			})
		}
	}
	if len(used) < len(cells) {
		rest := make([]int, 0, len(cells))
		mask := sudoku.Digits(0)
		for index := 0; index < len(cells); index++ {
			if used[index] {
				continue
			}
			rest = append(rest, index)
			mask = mask | cells[index]
		}
		result = append(result, Set{
			Indices: rest,
			Mask:    mask,
		})
	}
	return result, nil
}

func CheckAreaMask(s sudoku.Sudoku, area sudoku.Area, mask sudoku.Digits) []Set {
	if area.Size() != mask.Count() {
		return nil
	}

	cells := make([]sudoku.Digits, 0, area.Size())
	for _, cell := range area.Locations {
		d := s.Get(cell)
		if d&mask == 0 {
			return nil
		}
		cells = append(cells, d&mask)
	}
	if sets, err := FindSets(cells); err == nil {
		return sets
	}
	return nil
}
