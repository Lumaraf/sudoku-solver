package strategy

import (
	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type UniqueExclusionSolver struct {
	area sudoku.Area
	mask sudoku.Digits
}

func (slv UniqueExclusionSolver) Name() string {
	return "UniqueExclusionSolver"
}

func (slv UniqueExclusionSolver) Solve(s sudoku.Sudoku) ([]sudoku.Strategy, error) {
	for _, cell := range slv.area.And(s.SolvedArea()).Locations {
		slv.mask = slv.mask & ^s.Get(cell)
	}
	slv.area = slv.area.And(s.SolvedArea().Not())

	if slv.area.Empty() {
		return nil, nil
	}

	for v := range slv.mask.Values {
		area := s.SolvedArea().Not()
		for _, cell := range slv.area.Locations {
			d := s.Get(cell)
			if d.CanContain(v) {
				area = area.And(s.GetExclusionArea(cell))
			}
		}
		if area.Size() == 81 {
			continue
		}
		area = area.And(slv.area.Not())
		if !area.Empty() {
			for _, cell := range area.Locations {
				if err := s.RemoveOption(cell, v); err != nil {
					return nil, err
				}
			}
		}
	}
	return []sudoku.Strategy{slv}, nil
}

func (slv UniqueExclusionSolver) AreaFilter() sudoku.Area {
	return slv.area
}

func UniqueExclusionSolverFactory(restrictions []sudoku.Restriction) []sudoku.Strategy {
	solvers := []sudoku.Strategy{}
	for _, r := range restrictions {
		if unique, ok := r.(restriction.UniqueRestriction); ok {
			if unique.Area().Size() == 9 {
				solvers = append(solvers, UniqueExclusionSolver{
					area: unique.Area(),
					mask: sudoku.AllDigits,
				})
			}
		}
	}
	return solvers
}
