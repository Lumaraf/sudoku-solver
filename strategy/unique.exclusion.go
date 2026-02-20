package strategy

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

// removes options in the puzzle that are excluded by all possible placements of a digit in a unique area
func UniqueExclusionStrategyFactory[D sudoku.Digits[D], A sudoku.Area[A]](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
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

type UniqueExclusionStrategy[D sudoku.Digits[D], A sudoku.Area[A]] struct {
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

func (st UniqueExclusionStrategy[D, A]) Solve(s sudoku.Sudoku[D, A], push func(sudoku.Strategy[D, A])) error {
	area := st.area.And(s.SolvedArea().Not())
	if area.Empty() {
		return nil
	}

	// by value
	candidates := make([]A, s.Size())
	for _, l := range area.Locations {
		for v := range s.Get(l).Values {
			candidates[v-1] = candidates[v-1].With(l)
		}
	}

	for v, a := range candidates {
		if a.And(s.ChangedArea()).Empty() {
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
				changed = changed.Or(s.NextChangedArea())
				clones = append(clones, s)
				return nil
			})
		}

		if err := st.maskChangedCells(s, changed, clones); err != nil {
			return err
		}
	}

	// by cell
	for _, l := range area.And(s.ChangedArea()).Locations {
		d := s.Get(l)
		var changed A
		clones := make([]sudoku.Sudoku[D, A], 0, d.Count())
		for v := range d.Values {
			_ = s.Try(func(s sudoku.Sudoku[D, A]) error {
				if err := s.Set(l, v); err != nil {
					return err
				}
				changed = changed.Or(s.NextChangedArea())
				clones = append(clones, s)
				return nil
			})
		}

		if err := st.maskChangedCells(s, changed, clones); err != nil {
			return err
		}
	}

	push(st)
	return nil
}

func (st UniqueExclusionStrategy[D, A]) maskChangedCells(s sudoku.Sudoku[D, A], changed A, clones []sudoku.Sudoku[D, A]) error {
	for _, l := range changed.And(st.area.Not()).Locations {
		var mask D
		for _, clone := range clones {
			mask = mask.Or(clone.Get(l))
		}
		if err := s.Mask(l, mask); err != nil {
			return err
		}
	}
	return nil
}
