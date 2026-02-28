package strategy

import (
	"github.com/lumaraf/sudoku-solver/sudoku"
)

// checks if values in a cell would break a rule
func LogicChainStrategyFactory[D sudoku.Digits[D], A sudoku.Area[A]](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	return []sudoku.Strategy[D, A]{LogicChainStrategy[D, A]{
		s.NewArea().Not(),
	}}
}

type LogicChainStrategy[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	area A
}

func (slv LogicChainStrategy[D, A]) Name() string {
	return "LogicChainStrategy"
}

func (slv LogicChainStrategy[D, A]) Difficulty() sudoku.Difficulty {
	return sudoku.DIFFICULTY_NORMAL
}

func (slv LogicChainStrategy[D, A]) AreaFilter() A {
	return slv.area
}

func (slv LogicChainStrategy[D, A]) Solve(s sudoku.Sudoku[D, A], push func(sudoku.Strategy[D, A])) error {
	for _, cell := range s.SolvedArea().Not().Locations {
		d := s.Get(cell)
		results := make([]sudoku.Sudoku[D, A], 0, d.Count())
		for v := range d.Values {
			err := s.Try(func(s sudoku.Sudoku[D, A]) error {
				err := s.Set(cell, v)
				if err == nil {
					err = s.Validate()
				}
				if err == nil {
					results = append(results, s)
				}
				return err
			})
			if err != nil {
				err = s.RemoveOption(cell, v)
				if err != nil {
					return err
				}
			}
		}

		if len(results) == 0 {
			continue
		}

		affectedArea := s.NewArea().All()
		for _, r := range results {
			affectedArea = affectedArea.And(r.NextChangedArea())
		}
		for _, l := range affectedArea.Locations {
			mask := s.AllDigits()
			for _, r := range results {
				mask = mask.And(r.Get(l).Not())
			}
			if err := s.RemoveMask(l, mask); err != nil {
				return err
			}
		}
	}
	push(slv)
	return nil
}
