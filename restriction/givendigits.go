package restriction

import "github.com/lumaraf/sudoku-solver/sudoku"

func GivenDigits[D sudoku.Digits, A sudoku.Area](rows ...string) sudoku.RestrictionFactory[D, A] {
	return func(sb sudoku.SudokuBuilder[D, A]) error {
		for row, rowContent := range rows {
			for col, cellContent := range rowContent {
				if cellContent < '1' || cellContent > '9' {
					continue
				}
				if err := sb.SetCell(row, col, int(cellContent-'0')); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
