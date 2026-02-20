package strategy

import (
	"github.com/lumaraf/sudoku-solver/sudoku"
)

// compares areas that have to contain the same set of values and removes values that are not in both areas
func SetEquivalenceStrategyFactory[D sudoku.Digits[D], A sudoku.Area[A]](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	if s.Size() != 9 {
		return nil
	}

	strategies := make([]sudoku.Strategy[D, A], 0)
	for row1 := 0; row1 < s.Size()-1; row1++ {
		for row2 := row1 + 1; row2 < s.Size(); row2++ {
			if s.BoxAt(sudoku.CellLocation{row1, 0}) == s.BoxAt(sudoku.CellLocation{row2, 0}) {
				continue
			}

			for col1 := 0; col1 < s.Size()-1; col1++ {
				for col2 := col1 + 1; col2 < s.Size(); col2++ {
					if s.BoxAt(sudoku.CellLocation{0, col1}) == s.BoxAt(sudoku.CellLocation{0, col2}) {
						continue
					}

					strategies = append(strategies, createSetEquivalenceStrategy(s, []int{row1, row2}, []int{col1, col2}))
				}
			}
		}
	}
	return strategies
}

func createSetEquivalenceStrategy[D sudoku.Digits[D], A sudoku.Area[A]](s sudoku.Sudoku[D, A], rows, cols []int) SetEquivalenceStrategy[D, A] {
	var rowsArea A
	var colsArea A
	var boxes A
	for _, row := range rows {
		rowsArea = rowsArea.Or(s.Row(row))
		for _, col := range cols {
			box := s.BoxAt(sudoku.CellLocation{row, col})
			boxes = boxes.Or(s.Box(box))
		}
	}
	for _, col := range cols {
		colsArea = colsArea.Or(s.Column(col))
	}

	// area 0: boxes at intersections of rows and cols without cells contained by the rows and cols
	// area 1: cells contained by the rows and cols without boxes at intersections, but including the intersection cell
	areas := [2]A{
		boxes.And(rowsArea.Or(colsArea).Not()),
		rowsArea.Or(colsArea).And(boxes.Not()).Or(rowsArea.And(colsArea)),
	}
	return SetEquivalenceStrategy[D, A]{areas: areas}
}

type SetEquivalenceStrategy[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	areas [2]A
}

func (slv SetEquivalenceStrategy[D, A]) Name() string {
	return "SetEquivalenceStrategy"
}

func (slv SetEquivalenceStrategy[D, A]) Difficulty() sudoku.Difficulty {
	return sudoku.DIFFICULTY_IMPOSSIBLE
}

func (slv SetEquivalenceStrategy[D, A]) AreaFilter() A {
	return slv.areas[0].Or(slv.areas[1])
}

func (slv SetEquivalenceStrategy[D, A]) Solve(s sudoku.Sudoku[D, A], push func(sudoku.Strategy[D, A])) error {
	masks := make([]D, 2)
	for i, area := range slv.areas {
		for _, l := range area.Locations {
			masks[i] = masks[i].Or(s.Get(l))
		}
	}

	if !masks[1].And(masks[0].Not()).Empty() {
		for _, l := range slv.areas[0].Locations {
			if err := s.Mask(l, masks[1]); err != nil {
				return err
			}
		}
	}

	if !masks[0].And(masks[1].Not()).Empty() {
		for _, l := range slv.areas[1].Locations {
			if err := s.Mask(l, masks[0]); err != nil {
				return err
			}
		}
	}

	push(slv)
	return nil
}
