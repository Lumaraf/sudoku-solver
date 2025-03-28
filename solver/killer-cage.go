package solver

import (
	"errors"
	"github.com/lumaraf/sudoku-solver/restriction"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

var killerCageCombinations = map[[2]int][]sudoku.Digits{}

func init() {
	for d := sudoku.Digits(1); d <= sudoku.AllDigits; d++ {
		sum := 0
		for v := range d.Values {
			sum += v
		}
		k := [2]int{d.Count(), sum}
		killerCageCombinations[k] = append(killerCageCombinations[k], d)
	}
}

type KillerCageSolver struct {
	masks []sudoku.Digits
	area  sudoku.Area
}

func (slv KillerCageSolver) Name() string {
	return "KillerCageSolver"
}

func (slv KillerCageSolver) Solve(s sudoku.Sudoku) ([]sudoku.Solver, error) {
	possibleMasks := make([]sudoku.Digits, 0, len(slv.masks))
	var errs []error
	for _, mask := range slv.masks {
		if err := s.Try(func(s sudoku.Sudoku) error {
			for _, cell := range slv.area.Locations {
				if err := s.Mask(cell, mask); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			errs = append(errs, err)
			continue
		}
		possibleMasks = append(possibleMasks, mask)
	}

	if len(possibleMasks) == 0 {
		return nil, errors.Join(errs...)
	}

	return []sudoku.Solver{
		KillerCageSolver{
			masks: possibleMasks,
			area:  slv.area,
		},
	}, nil
}

func (slv KillerCageSolver) AreaFilter() sudoku.Area {
	return slv.area
}

func KillerCageSolverFactory(restrictions []sudoku.Restriction) []sudoku.Solver {
	solvers := []sudoku.Solver{}
	for _, r := range restrictions {
		if cage, ok := r.(restriction.KillerCageRestriction); ok {
			solvers = append(solvers, KillerCageSolver{
				masks: killerCageCombinations[[2]int{cage.Area.Size(), cage.Sum}],
				area:  cage.Area,
			})
		}
	}
	return solvers
}
