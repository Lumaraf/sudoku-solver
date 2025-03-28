package solver

import (
	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type sandwichCombination struct {
	minCount int
	maxCount int
	masks    []sudoku.Digits
}

var sandwichCombinations = map[int]sandwichCombination{}
var bread = sudoku.Digits(0)

func init() {
	for d := sudoku.Digits(0); d <= sudoku.AllDigits; d++ {
		if d.CanContain(1) || d.CanContain(9) {
			continue
		}
		sum := 0
		for v := 2; v <= 8; v++ {
			if d.CanContain(v) {
				sum += v
			}
		}
		combos, found := sandwichCombinations[sum]
		if !found {
			combos = sandwichCombination{
				minCount: 9,
				maxCount: 0,
				masks:    []sudoku.Digits{},
			}
		}
		count := d.Count()
		if count < combos.minCount {
			combos.minCount = count
		}
		if count > combos.maxCount {
			combos.maxCount = count
		}
		combos.masks = append(combos.masks, d)
		sandwichCombinations[sum] = combos
	}

	bread.AddOption(1)
	bread.AddOption(9)
}

type SandwichSolver struct {
	combos sandwichCombination
	cells  []sudoku.CellLocation
	area   sudoku.Area
}

func (slv SandwichSolver) Name() string {
	return "SandwichSolver"
}

func (slv SandwichSolver) Solve(s sudoku.Sudoku) ([]sudoku.Solver, error) {
	start := -1
	end := -1
	for index, cell := range slv.cells {
		d := s.Get(cell)
		if d&^bread == 0 {
			if start >= 0 {
				end = index
				break
			} else {
				start = index
			}
		}
	}

	if end == -1 {
		possibleEnds := map[int]bool{}
		for start, cell := range slv.cells {
			d := s.Get(cell) & bread
			if d == 0 {
				continue
			}
			breadMask := bread
			if d.Count() == 1 {
				breadMask = breadMask & ^d
			}
			hit := false
			for end := start + slv.combos.minCount + 1; end < len(slv.cells); end++ {
				d := s.Get(slv.cells[end])
				endBread := d & breadMask
				if endBread == 0 {
					continue
				}
				area := sudoku.NewArea(slv.cells[start+1 : end]...)
				for _, mask := range slv.combos.masks {
					if CheckAreaMask(s, area, mask) != nil {
						hit = true
						possibleEnds[end] = true
					}
				}
			}
			if !hit && !possibleEnds[start] {
				if err := s.RemoveMask(cell, bread); err != nil {
					return nil, err
				}
			}
		}
		return []sudoku.Solver{slv}, nil
	}

	return SandwichContentSolver{
		start:  start,
		end:    end,
		combos: slv.combos,
		cells:  slv.cells,
		area:   slv.area,
	}.Solve(s)
}

func (slv SandwichSolver) AreaFilter() sudoku.Area {
	return slv.area
}

type SandwichContentSolver struct {
	start  int
	end    int
	combos sandwichCombination
	cells  []sudoku.CellLocation
	area   sudoku.Area
}

func (slv SandwichContentSolver) Solve(s sudoku.Sudoku) ([]sudoku.Solver, error) {
	count := slv.end - slv.start - 1
	if count == 0 {
		return nil, nil
	}

	masks := make([]sudoku.Digits, 0, 10)
	for _, mask := range slv.combos.masks {
		if mask.Count() == count {
			masks = append(masks, mask)
		}
	}

	insideArea := sudoku.Area{}
	outsideArea := sudoku.Area{}

	for index := slv.start + 1; index < slv.end; index++ {
		cell := slv.cells[index]
		insideArea.Set(cell, true)
	}

	for index := 0; index < slv.start; index++ {
		cell := slv.cells[index]
		outsideArea.Set(cell, true)
	}
	for index := slv.end + 1; index < len(slv.cells); index++ {
		cell := slv.cells[index]
		outsideArea.Set(cell, true)
	}

	return AreaSumSolver{
		masks:       masks,
		area:        insideArea,
		outsideArea: outsideArea,
	}.Solve(s)
}

func (slv SandwichContentSolver) AreaFilter() sudoku.Area {
	return slv.area
}

func SandwichSolverFactory(restrictions []sudoku.Restriction) []sudoku.Solver {
	solvers := []sudoku.Solver{}
	for _, r := range restrictions {
		if sandwich, ok := r.(restriction.SandwichRestriction); ok {
			solvers = append(solvers, SandwichSolver{
				combos: sandwichCombinations[sandwich.Sum],
				cells:  sandwich.Cells,
				area:   sudoku.NewArea(sandwich.Cells...),
			})
		}
	}
	return solvers
}
