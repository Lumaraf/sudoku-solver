package strategy

import (
	"errors"
	"fmt"
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func KillerCageStrategyFactory[D sudoku.Digits[D], A sudoku.Area](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	allMasks := generateAreaSumMasks(s)

	strategies := sudoku.Strategies[D, A]{}
	for r := range sudoku.GetRestrictions[D, A, rule.AreaSumRestriction[D, A]](s) {
		if !s.IsUniqueArea(r.Area()) {
			continue
		}

		masks := make([]D, 0)
		for _, m := range allMasks[r.Sum()] {
			if m.Count() == r.Area().Size() {
				masks = append(masks, m)
			}
		}
		strategies = append(strategies, KillerCageStrategy[D, A]{
			area:  r.Area(),
			masks: masks,
		})
	}

	if len(strategies) == 0 {
		return strategies
	}
	return strategies
}

type KillerCageStrategy[D sudoku.Digits[D], A sudoku.Area] struct {
	area  A
	masks []D
}

func (st KillerCageStrategy[D, A]) Name() string {
	return "KillerCageStrategy"
}

func (st KillerCageStrategy[D, A]) Difficulty() sudoku.Difficulty {
	return sudoku.DIFFICULTY_NORMAL
}

func (st KillerCageStrategy[D, A]) AreaFilter() A {
	return st.area
}

func (st KillerCageStrategy[D, A]) Solve(s sudoku.Sudoku[D, A]) ([]sudoku.Strategy[D, A], error) {
	area := s.IntersectAreas(st.area, s.InvertArea(s.SolvedArea()))
	if area.Empty() {
		return nil, nil
	}

	// filter masks and find forced digits
	forcedDigits := s.NewDigits()
	masks := make([]D, 0, len(st.masks))
	for _, m := range st.masks {
		if st.isMaskPlaceable(s, st.area, m) {
			masks = append(masks, m)
			forcedDigits = s.IntersectDigits(forcedDigits, m)
		}
	}
	if len(masks) == 0 {
		return nil, errors.New("no valid masks for killer cage")
	}
	//st.masks = masks

	// eliminate forced digits from other areas
	for v := range forcedDigits.Values {
		exclusionArea := s.InvertArea(st.area)
		for _, l := range area.Locations {
			d := s.Get(l)
			if d.CanContain(v) {
				exclusionArea = s.IntersectAreas(exclusionArea, s.GetExclusionArea(l))
			}
		}
		for _, l := range exclusionArea.Locations {
			fmt.Printf("Removing forced digit %d from %v due to killer cage\n", v, l)
			if err := s.RemoveOption(l, v); err != nil {
				return nil, err
			}
		}
	}

	// eliminate impossible digits
	for _, l := range area.Locations {
		d := s.Get(l)
		for v := range d.Values {
			if !st.isValuePlaceable(s, l, v) {
				if err := s.RemoveOption(l, v); err != nil {
					return nil, err
				}
			}
		}
	}

	return sudoku.Strategies[D, A]{st}, nil
}

func (st KillerCageStrategy[D, A]) isValuePlaceable(s sudoku.Sudoku[D, A], l sudoku.CellLocation, v int) bool {
	area := st.area
	s.AreaWithout(&area, l)
	for _, m := range st.masks {
		if !m.CanContain(v) {
			continue
		}
		if area.Empty() {
			return true
		}
		m = s.IntersectDigits(m, s.InvertDigits(s.NewDigits(v)))
		if st.isMaskPlaceable(s, area, m) {
			return true
		}
	}
	return false
}

func (st KillerCageStrategy[D, A]) isMaskPlaceable(s sudoku.Sudoku[D, A], area A, mask D) bool {
	for _, l := range area.Locations {
		d := s.Get(l)
		for v := range s.IntersectDigits(d, mask).Values {
			nextArea := area
			s.AreaWithout(&nextArea, l)
			if nextArea.Empty() {
				return true
			}
			nextMask := s.IntersectDigits(mask, s.InvertDigits(s.NewDigits(v)))
			if st.isMaskPlaceable(s, nextArea, nextMask) {
				return true
			}
		}
		break
	}
	return false
}

var areaSumMasksCache = map[int]any{}

func generateAreaSumMasks[D sudoku.Digits[D], A sudoku.Area](s sudoku.Sudoku[D, A]) map[int][]D {
	if cache, ok := areaSumMasksCache[s.Size()].(map[int][]D); ok {
		return cache
	}

	masks := make(map[int][]D)
	var buildMasks func(values []int, sum int, start int)
	buildMasks = func(values []int, sum int, start int) {
		values = append(values, 0)
		lastIndex := len(values) - 1
		for v := start; v <= s.Size(); v++ {
			values[lastIndex] = v
			masks[sum+v] = append(masks[sum+v], s.NewDigits(values...))
			buildMasks(values, sum+v, v+1)
		}
	}
	buildMasks(make([]int, 0, s.Size()), 0, 1)
	areaSumMasksCache[s.Size()] = masks
	return masks
}
