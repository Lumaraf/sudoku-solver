package strategy

import (
	"errors"
	"github.com/lumaraf/sudoku-solver/extra/restriction"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

type EqualSolver struct {
	area sudoku.Area
}

func (slv EqualSolver) Name() string {
	return "EqualSolver"
}

func (slv EqualSolver) Solve(s sudoku.Sudoku) ([]sudoku.Strategy, error) {
	mask := sudoku.AllDigits
	for _, cell := range slv.area.Locations {
		mask = mask & s.Get(cell)
	}
	if mask == 0 {
		return nil, errors.New("no values in equal cells")
	}
	for _, cell := range slv.area.Locations {
		if err := s.Mask(cell, mask); err != nil {
			return nil, err
		}
	}
	if mask.Count() == 1 {
		return nil, nil
	}
	return []sudoku.Strategy{slv}, nil
}

func (slv EqualSolver) AreaFilter() sudoku.Area {
	return slv.area
}

func EqualSolverFactory(restrictions []sudoku.Restriction) []sudoku.Strategy {
	solvers := []sudoku.Strategy{}
	for _, r := range restrictions {
		if eq, ok := r.(restriction.EqualRestriction); ok {
			solvers = append(solvers, EqualSolver{
				area: sudoku.NewArea(eq.Cells...),
			})
		}
	}
	return solvers
}
