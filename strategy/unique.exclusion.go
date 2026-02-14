package strategy

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

// removes options in the puzzle that are excluded by all possible placements of a digit in a unique area
func UniqueExclusionStrategyFactory[D sudoku.Digits[D], A sudoku.Area](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	strategies := make([]sudoku.Strategy[D, A], 0)
	for r := range sudoku.GetRestrictions[D, A, rule.UniqueRestriction[D, A]](s) {
		a := r.Area()
		if a.Size() < s.Size() {
			continue
		}

		strategies = append(strategies, UniqueExclusionStrategy[D, A]{
			area: a,
		})
	}
	return strategies
}

type UniqueExclusionStrategy[D sudoku.Digits[D], A sudoku.Area] struct {
	area A
}

func (st UniqueExclusionStrategy[D, A]) Name() string {
	return "UniqueExclusionStrategy"
}

func (st UniqueExclusionStrategy[D, A]) Difficulty() sudoku.Difficulty {
	return sudoku.DIFFICULTY_HARD
}

func (st UniqueExclusionStrategy[D, A]) AreaFilter() A {
	return st.area
}

func (st UniqueExclusionStrategy[D, A]) Solve(s sudoku.Sudoku[D, A]) ([]sudoku.Strategy[D, A], error) {
	area := s.IntersectAreas(st.area, s.InvertArea(s.SolvedArea()))
	if area.Empty() {
		return nil, nil
	}

	// by value
	candidates := make([]A, s.Size())
	for _, l := range area.Locations {
		for v := range s.Get(l).Values {
			s.AreaWith(&candidates[v-1], l)
		}
	}

	for v, a := range candidates {
		if s.IntersectAreas(a, s.ChangedArea()).Empty() {
			continue
		}

		v += 1
		var changed A
		clones := make([]sudoku.Sudoku[D, A], 0, a.Size())
		for _, l := range a.Locations {
			_ = s.Try(func(s sudoku.Sudoku[D, A]) error {
				if err := s.Set(l, v); err != nil {
					return err
				}
				changed = s.UnionAreas(changed, s.NextChangedArea())
				clones = append(clones, s)
				return nil
			})
		}

		if err := st.maskChangedCells(s, changed, clones); err != nil {
			return nil, err
		}
	}

	// by cell
	for _, l := range s.IntersectAreas(area, s.ChangedArea()).Locations {
		d := s.Get(l)
		var changed A
		clones := make([]sudoku.Sudoku[D, A], 0, d.Count())
		for v := range d.Values {
			_ = s.Try(func(s sudoku.Sudoku[D, A]) error {
				if err := s.Set(l, v); err != nil {
					return err
				}
				changed = s.UnionAreas(changed, s.NextChangedArea())
				clones = append(clones, s)
				return nil
			})
		}

		if err := st.maskChangedCells(s, changed, clones); err != nil {
			return nil, err
		}
	}

	return []sudoku.Strategy[D, A]{st}, nil
}

func (st UniqueExclusionStrategy[D, A]) maskChangedCells(s sudoku.Sudoku[D, A], changed A, clones []sudoku.Sudoku[D, A]) error {
	for _, l := range s.IntersectAreas(changed, s.InvertArea(st.area)).Locations {
		var mask D
		for _, clone := range clones {
			mask = s.UnionDigits(mask, clone.Get(l))
		}
		if err := s.Mask(l, mask); err != nil {
			return err
		}
	}
	return nil
}
