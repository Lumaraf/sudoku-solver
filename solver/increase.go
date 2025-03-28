package solver

import (
	"errors"

	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type IncreaseSolver struct {
	first  sudoku.CellLocation
	second sudoku.CellLocation
	area   sudoku.Area
}

func (slv IncreaseSolver) Name() string {
	return "IncreaseSolver"
}

func (slv IncreaseSolver) Solve(s sudoku.Sudoku) ([]sudoku.Solver, error) {
	for {
		first := s.Get(slv.first)
		second := s.Get(slv.second)
		if second.Min() > first.Min() {
			break
		}
		if second.Count() == 1 {
			return nil, errors.New("invalid increase")
		}
		if err := s.RemoveOption(slv.second, second.Min()); err != nil {
			return nil, err
		}
	}
	for {
		first := s.Get(slv.first)
		second := s.Get(slv.second)
		if first.Max() < second.Max() {
			break
		}
		if first.Count() == 1 {
			return nil, errors.New("invalid increase")
		}
		if err := s.RemoveOption(slv.first, first.Max()); err != nil {
			return nil, err
		}
	}

	first := s.Get(slv.first)
	second := s.Get(slv.second)
	if first.Count() == 1 || second.Count() == 1 {
		return nil, nil
	}
	return []sudoku.Solver{slv}, nil
}

func (slv IncreaseSolver) AreaFilter() sudoku.Area {
	return slv.area
}

func IncreaseSolverFactory(restrictions []sudoku.Restriction) []sudoku.Solver {
	solvers := []sudoku.Solver{}
	for _, r := range restrictions {
		if inc, ok := r.(restriction.IncreaseRestriction); ok {
			solvers = append(solvers, IncreaseSolver{
				first:  inc.First,
				second: inc.Second,
				area:   sudoku.NewArea(inc.First, inc.Second),
			})
		}
	}
	return solvers
}
