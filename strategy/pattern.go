package strategy

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

// finds all possible placement patterns for every digit and overlays them to eliminate impossible options
func PatternOverlayStrategyFactory[D sudoku.Digits[D], A sudoku.Area[A]](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	requiredAreas := make([]A, 0, s.Size()*3)
	for r := range sudoku.GetRestrictions[D, A, rule.UniqueRestriction[D, A]](s) {
		a := r.Area()
		if a.Size() == s.Size() {
			requiredAreas = append(requiredAreas, a)
		}
	}

	return []sudoku.Strategy[D, A]{PatternOverlayStrategy[D, A]{
		area:          s.NewArea().Not(),
		requiredAreas: requiredAreas,
	}}
}

type PatternOverlayStrategy[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	area          A
	requiredAreas []A
	extraAreas    []A
}

func (st PatternOverlayStrategy[D, A]) Name() string {
	return "PatternOverlayStrategy"
}

func (st PatternOverlayStrategy[D, A]) Difficulty() sudoku.Difficulty {
	return sudoku.DIFFICULTY_IMPOSSIBLE
}

func (st PatternOverlayStrategy[D, A]) AreaFilter() A {
	return st.area
}

func (st PatternOverlayStrategy[D, A]) Solve(s sudoku.Sudoku[D, A], push func(sudoku.Strategy[D, A])) error {
	valueAreas := make([]A, s.Size())
	for row := 0; row < s.Size(); row++ {
		for col := 0; col < s.Size(); col++ {
			l := sudoku.CellLocation{
				Row: row,
				Col: col,
			}
			for v := range s.Get(l).Values {
				valueAreas[v-1] = valueAreas[v-1].With(l)
			}
		}
	}

	count := 0
	valuePatterns := make([][]valuePattern[A], s.Size())
	for v, valueArea := range valueAreas {
		valuePatterns[v] = make([]valuePattern[A], 0, s.Size())
		for area := range st.findPlacementPatterns(s, valueArea, st.requiredAreas) {
			valuePatterns[v] = append(valuePatterns[v], valuePattern[A]{
				area:  area,
				value: v + 1,
			})
			count++
			if count > 1000 {
				push(st)
				return nil
			}
		}
	}

	// create bitmasks of compatible patterns for higher values
	for v, patterns := range valuePatterns {
		for idx, p := range patterns {
			masks := make([][]uint64, 0, len(valuePatterns)-v-1)
			for _, otherPatterns := range valuePatterns[v+1:] {
				mask := make([]uint64, len(otherPatterns)/64+1)
				for i, op := range otherPatterns {
					intersection := p.area.And(op.area)
					if intersection.Empty() {
						mask[i/64] |= 1 << (i % 64)
					}
				}
				masks = append(masks, mask)
			}
			valuePatterns[v][idx].masks = masks
		}
	}

	patternUnions := make([]A, s.Size())
	for vp := range st.findValidPatterns(valuePatterns) {
		patternUnions[vp.value-1] = patternUnions[vp.value-1].Or(vp.area)
	}

	changes := make(map[sudoku.CellLocation]D, s.Size())
	for v, area := range patternUnions {
		eliminationArea := valueAreas[v].And(area.Not())
		for _, l := range eliminationArea.Locations {
			changes[l] = changes[l].Or(s.NewDigits(v + 1))
		}
	}

	for l, d := range changes {
		if err := s.RemoveMask(l, d); err != nil {
			return err
		}
	}

	push(st)
	return nil
}

type valuePattern[A sudoku.Area[A]] struct {
	value int
	area  A
	masks [][]uint64
}

func (st PatternOverlayStrategy[D, A]) findPlacementPatterns(s sudoku.Sudoku[D, A], valueArea A, requiredAreas []A) func(func(A) bool) {
	return func(yield func(A) bool) {
		if len(requiredAreas) == 0 {
			yield(valueArea)
			return
		}

		intersection := valueArea.And(requiredAreas[0])
		for _, l := range intersection.Locations {
			candidateArea := valueArea.And(s.GetExclusionArea(l).Not())
			if candidateArea.Size() < s.Size() {
				continue
			}
			for pattern := range st.findPlacementPatterns(s, candidateArea, requiredAreas[1:]) {
				if !yield(pattern) {
					return
				}
			}
		}
	}
}

func (st PatternOverlayStrategy[D, A]) findValidPatterns(valuePatterns [][]valuePattern[A]) func(func(valuePattern[A]) bool) {
	return func(yield func(valuePattern[A]) bool) {
		for _, pattern := range valuePatterns[0] {
			matched := false
			for vp := range st.combinePatterns(pattern.masks, valuePatterns[1:]) {
				matched = true
				if !yield(vp) {
					return
				}
			}
			if matched {
				if !yield(pattern) {
					return
				}
			}
		}
	}
}

func (st PatternOverlayStrategy[D, A]) combinePatterns(masks [][]uint64, valuePatterns [][]valuePattern[A]) func(func(valuePattern[A]) bool) {
	return func(yield func(valuePattern[A]) bool) {
		mask := masks[0]
		if len(masks) == 1 {
			for i, pattern := range valuePatterns[0] {
				if mask[i/64]&(1<<(i%64)) != 0 {
					if !yield(pattern) {
						return
					}
				}
			}
			return
		}

		masks = masks[1:]

	patternLoop:
		for i, pattern := range valuePatterns[0] {
			if mask[i/64]&(1<<(i%64)) == 0 {
				continue
			}

			// combine masks
			combinedMasks := make([][]uint64, len(masks))
			for idx, inputMask := range masks {
				combinedMask := make([]uint64, len(inputMask))
				matched := false
				for j, m := range inputMask {
					combinedMask[j] = m & pattern.masks[idx][j]
					if combinedMask[j] != 0 {
						matched = true
					}
				}
				if !matched {
					continue patternLoop
				}
				combinedMasks[idx] = combinedMask
			}

			matched := false
			for vp := range st.combinePatterns(combinedMasks, valuePatterns[1:]) {
				matched = true
				if !yield(vp) {
					return
				}
			}
			if matched {
				if !yield(pattern) {
					return
				}
			}
		}
	}
}
