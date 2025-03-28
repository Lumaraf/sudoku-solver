package solver

import (
	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type RelationSolver struct {
	area    sudoku.Area
	offsets sudoku.Offsets
	masks   map[int]sudoku.Digits
}

func (slv RelationSolver) Name() string {
	return "RelationSolver"
}

func (slv RelationSolver) AreaFilter() sudoku.Area {
	return sudoku.Area{}.Not()
}

func (slv RelationSolver) Solve(s sudoku.Sudoku) ([]sudoku.Solver, error) {
	targetArea := sudoku.NewArea()
	for _, cell := range s.SolvedArea().Locations {
		for offsetCell := range slv.offsets.Locations(cell) {
			targetArea.Set(offsetCell, true)
		}
	}

	for _, cell := range targetArea.Locations {
		if !slv.area.Get(cell) {
			continue
		}

		mask := sudoku.NewDigits()
		for relationCell := range slv.offsets.Locations(cell) {
			mask = mask | s.Get(relationCell)
		}

		if mask.Count() == 9 {
			continue
		}

		//fmt.Println(mask.Count())

		// TODO apply slv.masks
		if err := s.Mask(cell, mask); err != nil {
			return nil, err
		}
	}

	return []sudoku.Solver{slv}, nil
}

func RelationSolverFactory(restrictions []sudoku.Restriction) []sudoku.Solver {
	solvers := []sudoku.Solver{}
	for _, r := range restrictions {
		if r, ok := r.(restriction.RelationRestriction); ok {
			solvers = append(solvers, RelationSolver{
				area:    r.Area,
				offsets: r.Offsets,
				masks:   r.Masks,
			})
		}
	}
	return solvers
}
