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

func (slv LogicChainStrategy[D, A]) Solve(s sudoku.Sudoku[D, A]) ([]sudoku.Strategy[D, A], error) {
	for _, cell := range s.ChangedArea().Locations {
		d := s.Get(cell)
		if d.Count() >= 2 {
			for v := range d.Values {
				err := s.Try(func(s sudoku.Sudoku[D, A]) error {
					err := s.Set(cell, v)
					if err == nil {
						err = s.Validate()
					}
					return err
				})
				if err != nil {
					err = s.RemoveOption(cell, v)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return []sudoku.Strategy[D, A]{slv}, nil
}
