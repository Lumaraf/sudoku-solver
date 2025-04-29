package strategy

import (
	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type SwordfishSolver struct {
	rows []int
	cols []int
}

func (slv SwordfishSolver) Name() string {
	return "SwordfishSolver"
}

func (slv SwordfishSolver) Solve(s sudoku.Sudoku) ([]sudoku.Strategy, error) {
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

	if len(slv.rows) < 3 && len(slv.cols) < 3 {
		return nil, nil
	}

	for digit := 1; digit <= 9; digit++ {
		slv.findSwordfish(s, digit)
	}
	return []sudoku.Strategy{slv}, nil
}

func (slv SwordfishSolver) findSwordfish(s sudoku.Sudoku, digit int) bool {
	// Check rows for Swordfish pattern
	for row1 := range slv.rows[:len(slv.rows)-2] {
		cols1 := slv.findCandidateCols(s, slv.rows[row1], digit)
		if len(cols1) != 3 {
			continue
		}
		for row2 := range slv.rows[row1+1 : len(slv.rows)-1] {
			cols2 := slv.findCandidateCols(s, slv.rows[row1+1+row2], digit)
			if len(cols2) != 3 || !equalSets(cols1, cols2) {
				continue
			}
			for row3 := range slv.rows[row1+2+row2:] {
				cols3 := slv.findCandidateCols(s, slv.rows[row1+2+row2+row3], digit)
				if len(cols3) == 3 && equalSets(cols1, cols3) {
					if slv.eliminateInCols(s, cols1, slv.rows[row1], slv.rows[row1+1+row2], slv.rows[row1+2+row2+row3], digit) {
						return true
					}
				}
			}
		}
	}

	// Check columns for Swordfish pattern
	for col1 := range slv.cols[:len(slv.cols)-2] {
		rows1 := slv.findCandidateRows(s, slv.cols[col1], digit)
		if len(rows1) != 3 {
			continue
		}
		for col2 := range slv.cols[col1+1 : len(slv.cols)-1] {
			rows2 := slv.findCandidateRows(s, slv.cols[col1+1+col2], digit)
			if len(rows2) != 3 || !equalSets(rows1, rows2) {
				continue
			}
			for col3 := range slv.cols[col1+2+col2:] {
				rows3 := slv.findCandidateRows(s, slv.cols[col1+2+col2+col3], digit)
				if len(rows3) == 3 && equalSets(rows1, rows3) {
					if slv.eliminateInRows(s, rows1, slv.cols[col1], slv.cols[col1+1+col2], slv.cols[col1+2+col2+col3], digit) {
						return true
					}
				}
			}
		}
	}

	return false
}

func (slv SwordfishSolver) findCandidateCols(s sudoku.Sudoku, row, digit int) []int {
	cols := make([]int, 0, 3)
	for _, l := range sudoku.RowArea(row).And(s.SolvedArea().Not()).Locations {
		if s.Get(l).CanContain(digit) {
			cols = append(cols, l.Col)
		}
	}
	return cols
}

func (slv SwordfishSolver) findCandidateRows(s sudoku.Sudoku, col, digit int) []int {
	rows := make([]int, 0, 3)
	for _, l := range sudoku.ColArea(col).And(s.SolvedArea().Not()).Locations {
		if s.Get(l).CanContain(digit) {
			rows = append(rows, l.Row)
		}
	}
	return rows
}

func (slv SwordfishSolver) eliminateInCols(s sudoku.Sudoku, cols []int, row1, row2, row3, digit int) bool {
	changed := false
	for row := 0; row < 9; row++ {
		if row == row1 || row == row2 || row == row3 {
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

func (slv SwordfishSolver) eliminateInRows(s sudoku.Sudoku, rows []int, col1, col2, col3, digit int) bool {
	changed := false
	for col := 0; col < 9; col++ {
		if col == col1 || col == col2 || col == col3 {
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

func (slv SwordfishSolver) AreaFilter() sudoku.Area {
	return sudoku.Area{}.Not()
}

func SwordfishSolverFactory(restrictions []sudoku.Restriction) []sudoku.Strategy {
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
	if len(rows) > 2 || len(cols) > 2 {
		return []sudoku.Strategy{SwordfishSolver{
			rows: rows,
			cols: cols,
		}}
	}
	return nil
}

func equalSets(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	set := make(map[int]bool, len(a))
	for _, v := range a {
		set[v] = true
	}
	for _, v := range b {
		if !set[v] {
			return false
		}
	}
	return true
}
