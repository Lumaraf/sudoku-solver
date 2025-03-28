package solver

import (
	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type ConsecutiveSolver struct {
	Offsets    sudoku.Offsets
	Difference int
}

func (slv ConsecutiveSolver) Name() string {
	return "ConsecutiveSolver"
}

func (slv ConsecutiveSolver) Solve(s sudoku.Sudoku) ([]sudoku.Solver, error) {
	for _, cell := range s.SolvedArea().Not().Locations {
		row := cell.Row
		col := cell.Col
		mask := sudoku.Digits(0)
		for cell := range slv.Offsets.Locations(cell) {
			d := s.Get(cell)
			mask = mask | d<<slv.Difference | d>>slv.Difference
		}
		if err := s.Mask(sudoku.CellLocation{row, col}, mask&sudoku.AllDigits); err != nil {
			return nil, err
		}
	}
	return []sudoku.Solver{slv}, nil
}

func (slv ConsecutiveSolver) AreaFilter() sudoku.Area {
	return sudoku.Area{}.Not()
}

func ConsecutiveSolverFactory(restrictions []sudoku.Restriction) []sudoku.Solver {
	solvers := []sudoku.Solver{}
	for _, r := range restrictions {
		if consecutive, ok := r.(restriction.ConsecutiveRestriction); ok {
			solvers = append(solvers, ConsecutiveSolver{
				Offsets:    consecutive.Offsets,
				Difference: consecutive.Difference,
			})
		}
	}
	return solvers
}
