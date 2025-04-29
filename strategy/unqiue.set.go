package strategy

import (
	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type UniqueSetSolver struct {
	Cells []sudoku.CellLocation
	Area  sudoku.Area
}

func (slv UniqueSetSolver) Name() string {
	return "UniqueSetSolver"
}

func (slv UniqueSetSolver) Solve(s sudoku.Sudoku) ([]sudoku.Strategy, error) {
	digits := make([]sudoku.Digits, 0, len(slv.Cells))
	for _, cell := range slv.Cells {
		digits = append(digits, s.Get(cell))
	}
	sets, err := FindSets(digits)
	if err != nil {
		return nil, err
	}
	solvers := make([]sudoku.Strategy, 0, len(sets))
	for _, set := range sets {
		cells := make([]sudoku.CellLocation, 0, len(set.Indices))
		for _, index := range set.Indices {
			cell := slv.Cells[index]
			cells = append(cells, cell)
			if err := s.Mask(cell, set.Mask); err != nil {
				return nil, err
			}
		}
		if len(cells) > 1 {
			solvers = append(solvers, UniqueSetSolver{
				Cells: cells,
				Area:  sudoku.NewArea(cells...),
			})
		}
	}
	return solvers, nil
}

func (slv UniqueSetSolver) AreaFilter() sudoku.Area {
	return slv.Area
}

func UniqueSetSolverFactory(restrictions []sudoku.Restriction) []sudoku.Strategy {
	solvers := []sudoku.Strategy{}
	for _, r := range restrictions {
		if unique, ok := r.(restriction.UniqueRestriction); ok {
			cells := []sudoku.CellLocation{}
			for _, cell := range unique.Area().Locations {
				cells = append(cells, cell)
			}
			solvers = append(solvers, UniqueSetSolver{
				Cells: cells,
				Area:  unique.Area(),
			})
		}
	}
	return solvers
}
