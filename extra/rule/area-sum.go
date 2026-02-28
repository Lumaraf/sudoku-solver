package rule

import (
	"errors"

	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

var (
	ErrInvalidAreaSum = errors.New("invalid area sum")
)

func KillerCageRulesFromString[D sudoku.Digits[D], A sudoku.Area[A]](grid []string, sums map[rune]int) sudoku.Rules[D, A] {
	cages := make(map[rune][]sudoku.CellLocation)
	for row, rowContent := range grid {
		for col, cellContent := range rowContent {
			if cellContent < 'A' || cellContent > 'Z' {
				continue
			}
			cages[cellContent] = append(cages[cellContent], sudoku.CellLocation{
				Row: row,
				Col: col,
			})
		}
	}

	rules := make(sudoku.Rules[D, A], 0, len(cages))
	for cageLabel, locations := range cages {
		rules = append(rules, KillerCageRule[D, A]{
			Area: locations,
			Sum:  sums[cageLabel],
		})
	}
	return rules
}

type KillerCageRule[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	Area []sudoku.CellLocation
	Sum  int
}

func (r KillerCageRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	area := sb.NewArea(r.Area...)
	//sb.AddRestriction(AreaSumRestriction[D, A]{
	//	area: area,
	//	sum:  r.Sum,
	//})
	//sb.AddValidator(AreaSumValidator[D, A]{
	//	area: area,
	//	sum:  r.Sum,
	//})
	return sb.Use(
		rule.NewUniqueAreaRule[D, A]("killer cage", area),
		AreaSumRule[D, A]{
			Area: r.Area,
			Sum:  r.Sum,
		},
	)
}

type AreaSumRule[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	Area []sudoku.CellLocation
	Sum  int
}

func (r AreaSumRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	area := sb.NewArea(r.Area...)
	sb.AddRestriction(AreaSumRestriction[D, A]{
		area: area,
		sum:  r.Sum,
	})
	sb.AddValidator(AreaSumValidator[D, A]{
		area: area,
		sum:  r.Sum,
	})
	sb.AddSolveProcessor(AreaSumSolveProcessor[D, A]{
		area: area,
		sum:  r.Sum,
	})
	return nil
}

type AreaSumRestriction[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	area A
	sum  int
}

func (r AreaSumRestriction[D, A]) Name() string {
	return "AreaSumRestriction"
}

func (r AreaSumRestriction[D, A]) Area() A {
	return r.area
}

func (r AreaSumRestriction[D, A]) Sum() int {
	return r.sum
}

type AreaSumValidator[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	area A
	sum  int
}

func (v AreaSumValidator[D, A]) Name() string {
	return "AreaSumValidator"
}

func (v AreaSumValidator[D, A]) Validate(s sudoku.Sudoku[D, A]) error {
	areaMin := 0
	areaMax := 0
	for _, cell := range v.area.Locations {
		d := s.Get(cell)
		areaMin += d.Min()
		areaMax += d.Max()
	}
	if areaMin > v.sum || areaMax < v.sum {
		return ErrInvalidAreaSum
	}
	return nil
}

type AreaSumSolveProcessor[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	area A
	sum  int
}

func (p AreaSumSolveProcessor[D, A]) Name() string {
	return "AreaSumSolveProcessor"
}

func (p AreaSumSolveProcessor[D, A]) ProcessSolve(s sudoku.Sudoku[D, A], cell sudoku.CellLocation, mask D) error {
	if !p.area.Get(cell) {
		return nil
	}

	// if all except for one cell in the area are solved, we can determine the value of the last cell
	solved := p.area.And(s.SolvedArea())
	if solved.Count() == p.area.Count()-1 {
		sum := 0
		for _, solvedCell := range solved.Locations {
			d, _ := s.Get(solvedCell).Single()
			sum += d
		}
		v := p.sum - sum
		if v <= 0 || v > s.Size() {
			return ErrInvalidAreaSum
		}
		unsolved := p.area.And(solved.Not())
		for _, unsolvedCell := range unsolved.Locations {
			return s.Set(unsolvedCell, p.sum-sum)
		}
	}
	return nil
}
