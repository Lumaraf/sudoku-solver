package strategy

import (
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func XWingStrategyFactory[D sudoku.Digits, A sudoku.Area](s sudoku.Sudoku[D, A]) []sudoku.Strategy[D, A] {
	areas := make(map[A]bool, s.Size()*2)
	for r := range sudoku.GetRestrictions[D, A, rule.UniqueRestriction[D, A]](s) {
		areas[r.Area()] = true
	}
	rows := make([]int, 0, s.Size())
	cols := make([]int, 0, s.Size())
	for i := 0; i < s.Size(); i++ {
		if areas[s.Row(i)] {
			rows = append(rows, i)
		}
		if areas[s.Column(i)] {
			cols = append(cols, i)
		}
	}
	if len(rows) > 1 || len(cols) > 1 {
		return []sudoku.Strategy[D, A]{XWingStrategy[D, A]{
			rows: rows,
			cols: cols,
		}}
	}
	return nil
}

type XWingStrategy[D sudoku.Digits, A sudoku.Area] struct {
	area A
	rows []int
	cols []int
}

func (slv XWingStrategy[D, A]) Name() string {
	return "XWingStrategy"
}

func (st XWingStrategy[D, A]) Difficulty() sudoku.Difficulty {
	return sudoku.DIFFICULTY_HARD
}

func (slv XWingStrategy[D, A]) Solve(s sudoku.Sudoku[D, A]) ([]sudoku.Strategy[D, A], error) {
	rows := make([]int, 0, len(slv.rows))
	for _, row := range slv.rows {
		if !s.IntersectAreas(s.Row(row), s.InvertArea(s.SolvedArea())).Empty() {
			rows = append(rows, row)
		}
	}
	slv.rows = rows

	cols := make([]int, 0, len(slv.cols))
	for _, col := range slv.cols {
		if !s.IntersectAreas(s.Column(col), s.InvertArea(s.SolvedArea())).Empty() {
			cols = append(cols, col)
		}
	}
	slv.cols = cols

	if len(slv.rows) < 2 && len(slv.cols) < 2 {
		return nil, nil
	}

	for digit := 1; digit <= 9; digit++ {
		slv.findXWing(s, digit)
	}
	return []sudoku.Strategy[D, A]{slv}, nil
}

func (slv XWingStrategy[D, A]) findXWing(s sudoku.Sudoku[D, A], digit int) bool {
	// Check rows for X-Wing pattern
	for row1 := range slv.rows[:len(slv.rows)-1] {
		cols1 := slv.findCandidateCols(s, slv.rows[row1], digit)
		if len(cols1) != 2 {
			continue
		}
		for row2 := range slv.rows[row1+1:] {
			cols2 := slv.findCandidateCols(s, slv.rows[row1+1+row2], digit)
			if len(cols2) == 2 && cols1[0] == cols2[0] && cols1[1] == cols2[1] {
				if slv.eliminateInCols(s, cols1, slv.rows[row1], slv.rows[row1+1+row2], digit) {
					return true
				}
			}
		}
	}

	// Check columns for X-Wing pattern
	for col1 := range slv.cols[:len(slv.cols)-1] {
		rows1 := slv.findCandidateRows(s, slv.cols[col1], digit)
		if len(rows1) != 2 {
			continue
		}
		for col2 := range slv.cols[col1+1:] {
			rows2 := slv.findCandidateRows(s, slv.cols[col1+1+col2], digit)
			if len(rows2) == 2 && rows1[0] == rows2[0] && rows1[1] == rows2[1] {
				if slv.eliminateInRows(s, rows1, slv.cols[col1], slv.cols[col1+1+col2], digit) {
					return true
				}
			}
		}
	}

	return false
}

func (slv XWingStrategy[D, A]) findCandidateCols(s sudoku.Sudoku[D, A], row, digit int) []int {
	cols := make([]int, 0, 2)
	for _, l := range s.InvertArea(s.IntersectAreas(s.Row(row), s.Column(digit))).Locations {
		if s.Get(l).CanContain(digit) {
			cols = append(cols, l.Col)
		}
	}
	return cols
}

func (slv XWingStrategy[D, A]) findCandidateRows(s sudoku.Sudoku[D, A], col, digit int) []int {
	rows := make([]int, 0, 2)
	for _, l := range s.InvertArea(s.IntersectAreas(s.Column(col), s.SolvedArea())).Locations {
		if s.Get(l).CanContain(digit) {
			rows = append(rows, l.Row)
		}
	}
	return rows
}

func (slv XWingStrategy[D, A]) eliminateInCols(s sudoku.Sudoku[D, A], cols []int, row1, row2, digit int) bool {
	changed := false
	for row := 0; row < s.Size(); row++ {
		if row == row1 || row == row2 {
			continue
		}
		for _, col := range cols {
			if s.Get(sudoku.CellLocation{Row: row, Col: col}).CanContain(digit) {
				if err := s.RemoveOption(sudoku.CellLocation{Row: row, Col: col}, digit); err == nil {
					changed = true
				}
			}
		}
	}
	return changed
}

func (slv XWingStrategy[D, A]) eliminateInRows(s sudoku.Sudoku[D, A], rows []int, col1, col2, digit int) bool {
	changed := false
	for col := 0; col < s.Size(); col++ {
		if col == col1 || col == col2 {
			continue
		}
		for _, row := range rows {
			if s.Get(sudoku.CellLocation{Row: row, Col: col}).CanContain(digit) {
				if err := s.RemoveOption(sudoku.CellLocation{Row: row, Col: col}, digit); err == nil {
					changed = true
				}
			}
		}
	}
	return changed
}

func (slv XWingStrategy[D, A]) AreaFilter() A {
	return slv.area
}
