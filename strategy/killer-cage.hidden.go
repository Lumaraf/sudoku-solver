package strategy

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func HiddenKillerCageStrategyFactory[D sudoku.Digits[D], A sudoku.Area](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	allMasks := generateAreaSumMasks(s)
	strategies := sudoku.Strategies[D, A]{}

	areaSumRestrictions := make([]rule.AreaSumRestriction[D, A], 0)
	for r := range sudoku.GetRestrictions[D, A, rule.AreaSumRestriction[D, A]](s) {
		areaSumRestrictions = append(areaSumRestrictions, r)
	}

	if len(areaSumRestrictions) == 0 {
		return strategies
	}

	// find hidden cages
	knownCages := map[A]bool{}
	for baseArea := range hiddenCageBaseAreas[D, A](s) {
		baseSum := (s.Size() * (s.Size() + 1) / 2) * (baseArea.Size() / s.Size())

		hits := 0
		for _, r2 := range areaSumRestrictions {
			if s.UnionAreas(baseArea, r2.Area()).Size() == baseArea.Size() {
				hits++
				baseArea = s.IntersectAreas(baseArea, s.InvertArea(r2.Area()))
				baseSum -= r2.Sum()
			}
		}

		if hits > 1 && baseArea.Size() < s.Size() && s.IsUniqueArea(baseArea) {
			if knownCages[baseArea] {
				continue
			}
			knownCages[baseArea] = true
			//fmt.Printf("Hidden cage found: sum=%d site=%d\n", baseSum, baseArea.Size())
			masks := make([]D, 0)
			for _, m := range allMasks[baseSum] {
				if m.Count() == baseArea.Size() {
					masks = append(masks, m)
				}
			}
			strategies = append(strategies, KillerCageStrategy[D, A]{
				area:  baseArea,
				masks: masks,
			})

			// check for inverted cage
			for _, r2 := range areaSumRestrictions {
				if s.IntersectAreas(baseArea, s.InvertArea(r2.Area())).Empty() {
					area := s.IntersectAreas(r2.Area(), s.InvertArea(baseArea))
					if !area.Empty() {
						masks := make([]D, 0)
						for _, m := range allMasks[r2.Sum()-baseSum] {
							if m.Count() == area.Size() {
								masks = append(masks, m)
							}
						}

						strategies = append(strategies, KillerCageStrategy[D, A]{
							area:  area,
							masks: masks,
						})
					}
				}
			}
		}
	}
	return strategies
}

func hiddenCageBaseAreas[D sudoku.Digits[D], A sudoku.Area](s sudoku.Sudoku[D, A]) func(func(A) bool) {
	var combineAreas func(current A, uniqueAreas []A) func(func(A) bool)
	combineAreas = func(current A, ua []A) func(func(A) bool) {
		return func(yield func(A) bool) {
			for _, a := range ua {
				if s.UnionAreas(current, a).Size() != current.Size()+a.Size() {
					continue
				}

				if !current.Empty() && !areasTouch(s, current, a) {
					continue
				}

				union := s.UnionAreas(current, a)
				if !yield(union) {
					return
				}
				if union.Size()/s.Size() < 4 {
					for combinedArea := range combineAreas(union, ua[1:]) {
						if !yield(combinedArea) {
							return
						}
					}
				}
			}
		}
	}

	uniqueAreas := make([]A, 0, s.Size()*3)
	for r := range sudoku.GetRestrictions[D, A, rule.UniqueRestriction[D, A]](s) {
		if r.Area().Size() == s.Size() {
			uniqueAreas = append(uniqueAreas, r.Area())
		}
	}
	return combineAreas(s.NewArea(), uniqueAreas)
}

func areasTouch[D sudoku.Digits[D], A sudoku.Area](s sudoku.Sudoku[D, A], a1, a2 A) bool {
	offsets := sudoku.Offsets{
		sudoku.Offset{Row: -1, Col: 0},
		sudoku.Offset{Row: 0, Col: -1},
		sudoku.Offset{Row: 0, Col: 1},
		sudoku.Offset{Row: 1, Col: 0},
	}
	for _, l1 := range a1.Locations {
		for _, l2 := range s.NewAreaFromOffsets(l1, offsets).Locations {
			if a2.Get(l2) {
				return true
			}
		}
	}
	return false
}
