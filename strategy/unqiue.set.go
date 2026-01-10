package strategy

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

// checks unique areas for complete sets of digits and removes those from other cells in the area
func UniqueSetStrategyFactory[D sudoku.Digits, A sudoku.Area](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	strategies := make([]sudoku.Strategy[D, A], 0)
	for r := range sudoku.GetRestrictions[D, A, rule.UniqueRestriction[D, A]](s) {
		strategies = append(strategies, UniqueSetStrategy[D, A]{
			Area: r.Area(),
		})
	}
	return strategies
}

type UniqueSetStrategy[D sudoku.Digits, A sudoku.Area] struct {
	Area A
}

func (st UniqueSetStrategy[D, A]) Name() string {
	return "UniqueSetStrategy"
}

func (st UniqueSetStrategy[D, A]) Difficulty() sudoku.Difficulty {
	return sudoku.DIFFICULTY_EASY
}

func (st UniqueSetStrategy[D, A]) AreaFilter() A {
	return st.Area
}

func (st UniqueSetStrategy[D, A]) Solve(s sudoku.Sudoku[D, A]) ([]sudoku.Strategy[D, A], error) {
	cells := make([]D, 0, st.Area.Size())
	for _, cell := range st.Area.Locations {
		d := s.Get(cell)
		if d.Count() == 1 {
			s.AreaWithout(&st.Area, cell)
			continue
		}
		cells = append(cells, d)
	}
	if len(cells) <= 1 {
		return nil, nil
	}

	bestSet := s.AllDigits()
	for set := range st.findSets(s, cells, 0, *new(D)) {
		if set.Count() < bestSet.Count() {
			bestSet = set
		}
	}
	if bestSet == s.AllDigits() || bestSet.Count() == len(cells) {
		return []sudoku.Strategy[D, A]{
			st,
		}, nil
	}

	mask := s.InvertDigits(bestSet)
	inSet := s.NewArea()
	notInSet := s.NewArea()
	for _, cell := range st.Area.Locations {
		d := s.Get(cell)
		if s.IntersectDigits(d, mask).Empty() {
			s.AreaWith(&inSet, cell)
		} else {
			s.AreaWith(&notInSet, cell)
			if err := s.RemoveMask(cell, bestSet); err != nil {
				return nil, err
			}
		}
	}
	return []sudoku.Strategy[D, A]{
		UniqueSetStrategy[D, A]{
			Area: inSet,
		},
		UniqueSetStrategy[D, A]{
			Area: notInSet,
		},
	}, nil
}

func (st UniqueSetStrategy[D, A]) findSets(s sudoku.Sudoku[D, A], contents []D, count int, mask D) func(yield func(D) bool) {
	return func(yield func(D) bool) {
		count++
		for i, d := range contents {
			combined := s.UnionDigits(mask, d)
			if combined.Count() == count {
				if !yield(combined) {
					return
				}
				continue
			}
			for set := range st.findSets(s, contents[i+1:], count, combined) {
				if !yield(set) {
					return
				}
			}
		}
	}
}
