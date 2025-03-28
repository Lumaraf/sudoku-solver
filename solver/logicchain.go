package solver

import (
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type LogicChainSolver struct {
}

func (slv LogicChainSolver) Name() string {
	return "LogicChainSolver"
}

func (slv LogicChainSolver) AreaFilter() sudoku.Area {
	return sudoku.Area{}.Not()
}

func (slv LogicChainSolver) Solve(s sudoku.Sudoku) ([]sudoku.Solver, error) {
	for _, cell := range s.ChangedArea().Locations {
		d := s.Get(cell)
		if d.Count() == 2 {
			for v := range d.Values {
				err := s.Try(func(s sudoku.Sudoku) error {
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
	return []sudoku.Solver{slv}, nil
}

func LogicChainSolverFactory(restrictions []sudoku.Restriction) []sudoku.Solver {
	return []sudoku.Solver{LogicChainSolver{}}
}
