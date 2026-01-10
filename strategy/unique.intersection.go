package strategy

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

// removes options from cells in unique areas if they are forced into an intersection with another unique area
func UniqueIntersectionStrategyFactory[D sudoku.Digits, A sudoku.Area](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	strategies := make([]sudoku.Strategy[D, A], 0)
	for r := range sudoku.GetRestrictions[D, A, rule.UniqueRestriction[D, A]](s) {
		a := r.Area()
		if a.Size() < s.Size() {
			continue
		}

		for r2 := range sudoku.GetRestrictions[D, A, rule.UniqueRestriction[D, A]](s) {
			a2 := r2.Area()
			if s.IntersectAreas(a, a2).Empty() || a == a2 {
				continue
			}

			strategies = append(strategies, UniqueIntersectionStrategy[D, A]{
				area:         s.UnionAreas(a, a2),
				source:       s.IntersectAreas(a, s.InvertArea(a2)),
				intersection: s.IntersectAreas(a, a2),
				target:       s.IntersectAreas(a2, s.InvertArea(a)),
			})
		}
	}
	return strategies
}

type UniqueIntersectionStrategy[D sudoku.Digits, A sudoku.Area] struct {
	area         A
	source       A
	intersection A
	target       A
}

func (st UniqueIntersectionStrategy[D, A]) Name() string {
	return "UniqueIntersectionStrategy"
}

func (st UniqueIntersectionStrategy[D, A]) Difficulty() sudoku.Difficulty {
	return sudoku.DIFFICULTY_EASY
}

func (st UniqueIntersectionStrategy[D, A]) AreaFilter() A {
	return st.area
}

func (st UniqueIntersectionStrategy[D, A]) Solve(s sudoku.Sudoku[D, A]) ([]sudoku.Strategy[D, A], error) {
	st.intersection = s.UnionAreas(st.intersection, s.InvertArea(s.SolvedArea()))
	if st.intersection.Empty() {
		return nil, nil
	}

	var d D
	for _, l := range st.intersection.Locations {
		d = s.UnionDigits(d, s.Get(l))
	}
	for _, l := range st.source.Locations {
		d = s.IntersectDigits(d, s.InvertDigits(s.Get(l)))
	}

	if d.Empty() {
		return []sudoku.Strategy[D, A]{st}, nil
	}

	for _, l := range st.target.Locations {
		if err := s.RemoveMask(l, d); err != nil {
			return []sudoku.Strategy[D, A]{st}, err
		}
	}

	return []sudoku.Strategy[D, A]{st}, nil
}
