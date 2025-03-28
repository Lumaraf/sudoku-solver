package solver

import (
	"errors"
	"github.com/lumaraf/sudoku-solver/restriction"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

type AreaSumSolver struct {
	masks       []sudoku.Digits
	area        sudoku.Area
	outsideArea sudoku.Area
}

func (slv AreaSumSolver) Name() string {
	return "AreaSumSolver"
}

func (slv AreaSumSolver) Solve(s sudoku.Sudoku) ([]sudoku.Solver, error) {
	count := slv.area.Size()
	allMasks := sudoku.Digits(0)
	cellMasks := make([]sudoku.Digits, count)
	requiredDigits := sudoku.AllDigits
	possibleMasks := make([]sudoku.Digits, 0, len(slv.masks))
	for _, mask := range slv.masks {
		if mask.Count() == count {
			sets := CheckAreaMask(s, slv.area, mask)
			if sets == nil {
				continue
			}

			possibleMasks = append(possibleMasks, mask)

			for _, set := range sets {
				for _, index := range set.Indices {
					cellMasks[index] = cellMasks[index] | set.Mask
				}
			}

			allMasks = allMasks | mask
			requiredDigits = requiredDigits & mask
		}
	}

	if len(possibleMasks) == 0 {
		return nil, errors.New("no area sum combinations")
	}

	for index, cell := range slv.area.Locations {
		if err := s.Mask(cell, cellMasks[index]); err != nil {
			return nil, err
		}
	}

	if requiredDigits != 0 {
		for _, cell := range slv.outsideArea.Locations {
			if err := s.RemoveMask(cell, requiredDigits); err != nil {
				return nil, err
			}
		}
	}

	if requiredDigits.Count() == count {
		return nil, nil
	}

	return []sudoku.Solver{
		AreaSumSolver{
			masks:       possibleMasks,
			area:        slv.area,
			outsideArea: slv.outsideArea,
		},
	}, nil
}

func (slv AreaSumSolver) AreaFilter() sudoku.Area {
	return slv.area
}

func AreaSumSolverFactory(restrictions []sudoku.Restriction) []sudoku.Solver {
	solvers := []sudoku.Solver{}
	for _, r := range restrictions {
		if cage, ok := r.(restriction.KillerCageRestriction); ok {
			solvers = append(solvers, AreaSumSolver{
				masks: killerCageCombinations[[2]int{cage.Area.Size(), cage.Sum}],
				area:  cage.Area,
			})
		}
	}
	return solvers
}
