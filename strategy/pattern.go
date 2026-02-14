package strategy

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

// finds all possible placement patterns for every digit and overlays them to eliminate impossible options
func PatternOverlayStrategyFactory[D sudoku.Digits[D], A sudoku.Area](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	requiredAreas := make([]A, 0, s.Size()*3)
	for r := range sudoku.GetRestrictions[D, A, rule.UniqueRestriction[D, A]](s) {
		a := r.Area()
		if a.Size() == s.Size() {
			requiredAreas = append(requiredAreas, a)
		}
	}

	return []sudoku.Strategy[D, A]{PatternOverlayStrategy[D, A]{
		area:          s.InvertArea(s.NewArea()),
		requiredAreas: requiredAreas,
	}}
}

type PatternOverlayStrategy[D sudoku.Digits[D], A sudoku.Area] struct {
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

func (st PatternOverlayStrategy[D, A]) Solve(s sudoku.Sudoku[D, A]) ([]sudoku.Strategy[D, A], error) {
	valueAreas := make([]A, s.Size())
	for row := 0; row < s.Size(); row++ {
		for col := 0; col < s.Size(); col++ {
			l := sudoku.CellLocation{
				Row: row,
				Col: col,
			}
			for v := range s.Get(l).Values {
				s.AreaWith(&valueAreas[v-1], l)
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
				return sudoku.Strategies[D, A]{st}, nil
			}
		}
	}

	// overlay patterns and eliminate impossible options
	patternUnions := make([]A, s.Size())
	for vp := range st.findValidPatterns(s, s.NewArea(), valuePatterns) {
		patternUnions[vp.value-1] = s.UnionAreas(patternUnions[vp.value-1], vp.area)
	}

	changes := make(map[sudoku.CellLocation]D, s.Size())
	for v, area := range patternUnions {
		eliminationArea := s.IntersectAreas(valueAreas[v], s.InvertArea(area))
		for _, l := range eliminationArea.Locations {
			changes[l] = s.UnionDigits(changes[l], s.NewDigits(v+1))
		}
	}

	for l, d := range changes {
		if err := s.RemoveMask(l, d); err != nil {
			return nil, err
		}
	}

	return sudoku.Strategies[D, A]{st}, nil
}

type valuePattern[A sudoku.Area] struct {
	value int
	area  A
}

func (st PatternOverlayStrategy[D, A]) findPlacementPatterns(s sudoku.Sudoku[D, A], valueArea A, requiredAreas []A) func(func(A) bool) {
	return func(yield func(A) bool) {
		if len(requiredAreas) == 0 {
			yield(valueArea)
			return
		}

		intersection := s.IntersectAreas(valueArea, requiredAreas[0])
		for _, l := range intersection.Locations {
			candidateArea := s.IntersectAreas(valueArea, s.InvertArea(s.GetExclusionArea(l)))
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

func (st PatternOverlayStrategy[D, A]) findValidPatterns(s sudoku.Sudoku[D, A], pattern A, valuePatterns [][]valuePattern[A]) func(func(valuePattern[A]) bool) {
	return func(yield func(valuePattern[A]) bool) {
		for _, otherPattern := range valuePatterns[0] {
			intersection := s.IntersectAreas(pattern, otherPattern.area)
			if !intersection.Empty() {
				continue
			}

			matched := false
			if len(valuePatterns) == 1 {
				matched = true
			} else {
				for vp := range st.findValidPatterns(s, s.UnionAreas(pattern, otherPattern.area), valuePatterns[1:]) {
					matched = true
					if !yield(vp) {
						return
					}
				}
			}
			if matched {
				if !yield(otherPattern) {
					return
				}
			}
		}
		return
	}
}
