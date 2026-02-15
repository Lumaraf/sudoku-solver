package strategy

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

// removes options from cells in unique areas if they are forced into an intersection with another unique area
func UniqueIntersectionStrategyFactory[D sudoku.Digits[D], A sudoku.Area[A]](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	strategies := make([]sudoku.Strategy[D, A], 0)
	for r := range sudoku.GetRestrictions[D, A, rule.UniqueRestriction[D, A]](s) {
		a := r.Area()
		if a.Size() < s.Size() {
			continue
		}

		for r2 := range sudoku.GetRestrictions[D, A, rule.UniqueRestriction[D, A]](s) {
			a2 := r2.Area()
			if a.And(a2).Empty() || a == a2 {
				continue
			}

			strategies = append(strategies, UniqueIntersectionStrategy[D, A]{
				area:         a.Or(a2),
				source:       a.And(a2.Not()),
				intersection: a.And(a2),
				target:       a2.And(a.Not()),
			})
		}
	}
	return strategies
}

type UniqueIntersectionStrategy[D sudoku.Digits[D], A sudoku.Area[A]] struct {
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
	st.intersection = st.intersection.Or(s.SolvedArea().Not())
	if st.intersection.Empty() {
		return nil, nil
	}

	var d D
	for _, l := range st.intersection.Locations {
		d = d.Or(s.Get(l))
	}
	for _, l := range st.source.Locations {
		d = d.And(s.Get(l).Not())
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
