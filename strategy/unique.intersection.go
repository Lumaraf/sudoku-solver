package strategy

import (
	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type UniqueIntersectionSolver struct {
	area sudoku.Area
}

func (slv UniqueIntersectionSolver) Name() string {
	return "UniqueIntersectionSolver"
}

func (slv UniqueIntersectionSolver) AreaFilter() sudoku.Area {
	return slv.area
}

func (slv UniqueIntersectionSolver) Solve(s sudoku.Sudoku) ([]sudoku.Strategy, error) {
	valueCounts := [9]int{}
	for _, cell := range slv.area.And(s.SolvedArea().Not()).Locations {
		for v := range s.Get(cell).Values {
			valueCounts[v-1]++
		}
	}
	for v, count := range valueCounts {
		if count >= 2 && count <= 2 {
			if err := slv.checkValue(s, v+1); err != nil {
				return nil, err
			}
		}
	}
	return []sudoku.Strategy{slv}, nil
}

func (slv UniqueIntersectionSolver) checkValue(s sudoku.Sudoku, v int) error {
	masks := [9][9]sudoku.Digits{}
	for _, cell := range slv.area.Locations {
		if !s.Get(cell).CanContain(v) {
			continue
		}

		err := s.Try(func(s sudoku.Sudoku) error {
			if err := s.Set(cell, v); err != nil {
				return err
			}
			if err := s.Validate(); err != nil {
				return err
			}
			for row := 0; row < 9; row++ {
				for col := 0; col < 9; col++ {
					masks[row][col] |= s.Get(sudoku.CellLocation{Row: row, Col: col})
				}
			}
			return nil
		})
		if err != nil {
			err = s.RemoveOption(cell, v)
			if err != nil {
				return err
			}
		}
	}
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			err := s.Mask(sudoku.CellLocation{Row: row, Col: col}, masks[row][col])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UniqueIntersectionSolverFactory(restrictions []sudoku.Restriction) []sudoku.Strategy {
	solvers := []sudoku.Strategy{}
	for _, r := range restrictions {
		if unique, ok := r.(restriction.UniqueRestriction); ok {
			if unique.Area().Size() == 9 {
				solvers = append(solvers, UniqueIntersectionSolver{
					area: unique.Area(),
				})
			}
		}
	}
	return solvers
}
