package rule

import (
	"errors"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

var (
	ErrInvalidAreaSum = errors.New("invalid area sum")
)

func KillerCageRulesFromString[D sudoku.Digits[D], A sudoku.Area](grid []string, sums map[rune]int) sudoku.Rules[D, A] {
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

type KillerCageRule[D sudoku.Digits[D], A sudoku.Area] struct {
	Area []sudoku.CellLocation
	Sum  int
}

func (r KillerCageRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	area := sb.NewArea(r.Area...)
	sb.AddRestriction(AreaSumRestriction[D, A]{
		area: area,
		sum:  r.Sum,
	})
	return sb.Use(NewUniqueAreaRule[D, A]("killer cage", area))
}

type AreaSumRule[D sudoku.Digits[D], A sudoku.Area] struct {
	Area []sudoku.CellLocation
	Sum  int
}

func (r AreaSumRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	area := sb.NewArea(r.Area...)
	sb.AddRestriction(AreaSumRestriction[D, A]{
		area: area,
		sum:  r.Sum,
	})
	return nil
}

type AreaSumRestriction[D sudoku.Digits[D], A sudoku.Area] struct {
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

func (r AreaSumRestriction[D, A]) Validate(s sudoku.Sudoku[D, A]) error {
	areaMin := 0
	areaMax := 0
	for _, cell := range r.Area().Locations {
		d := s.Get(cell)
		areaMin += d.Min()
		areaMax += d.Max()
	}
	if areaMin > r.Sum() || areaMax < r.Sum() {
		return ErrInvalidAreaSum
	}
	return nil
}
