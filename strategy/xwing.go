package strategy

import (
	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type XWingSolver struct {
	rows []int
	cols []int
}

func (slv XWingSolver) Name() string {
	return "XWingSolver"
}

func (slv XWingSolver) Solve(s sudoku.Sudoku) ([]sudoku.Strategy, error) {
	rows := make([]int, 0, len(slv.rows))
	for _, row := range slv.rows {
		if !sudoku.RowArea(row).And(s.SolvedArea().Not()).Empty() {
			rows = append(rows, row)
		}
	}
	slv.rows = rows

	cols := make([]int, 0, len(slv.cols))
	for _, col := range slv.cols {
		if !sudoku.ColArea(col).And(s.SolvedArea().Not()).Empty() {
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
	return []sudoku.Strategy{slv}, nil
}

func (slv XWingSolver) findXWing(s sudoku.Sudoku, digit int) bool {
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

func (slv XWingSolver) findCandidateCols(s sudoku.Sudoku, row, digit int) []int {
	cols := make([]int, 0, 2)
	for _, l := range sudoku.RowArea(row).And(s.SolvedArea().Not()).Locations {
		if s.Get(l).CanContain(digit) {
			cols = append(cols, l.Col)
		}
	}
	return cols
}

func (slv XWingSolver) findCandidateRows(s sudoku.Sudoku, col, digit int) []int {
	rows := make([]int, 0, 2)
	for _, l := range sudoku.ColArea(col).And(s.SolvedArea().Not()).Locations {
		if s.Get(l).CanContain(digit) {
			rows = append(rows, l.Row)
		}
	}
	return rows
}

func (slv XWingSolver) eliminateInCols(s sudoku.Sudoku, cols []int, row1, row2, digit int) bool {
	changed := false
	for row := 0; row < 9; row++ {
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

func (slv XWingSolver) eliminateInRows(s sudoku.Sudoku, rows []int, col1, col2, digit int) bool {
	changed := false
	for col := 0; col < 9; col++ {
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

func (slv XWingSolver) AreaFilter() sudoku.Area {
	return sudoku.Area{}.Not()
}

func XWingSolverFactory(restrictions []sudoku.Restriction) []sudoku.Strategy {
	areas := make(map[sudoku.Area]bool, 18)
	for _, r := range restrictions {
		if unique, ok := r.(restriction.UniqueRestriction); ok {
			areas[unique.Area()] = true
		}
	}
	rows := make([]int, 0, 9)
	cols := make([]int, 0, 9)
	for i := 0; i < 9; i++ {
		if areas[sudoku.RowArea(i)] {
			rows = append(rows, i)
		}
		if areas[sudoku.ColArea(i)] {
			cols = append(cols, i)
		}
	}
	if len(rows) > 1 || len(cols) > 1 {
		return []sudoku.Strategy{XWingSolver{
			rows: rows,
			cols: cols,
		}}
	}
	return nil
}
