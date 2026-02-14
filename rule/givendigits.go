package rule

import "github.com/lumaraf/sudoku-solver/sudoku"

type GivenDigits[D sudoku.Digits[D], A sudoku.Area] struct {
	givenDigits map[sudoku.CellLocation]int
}

func GivenDigitsFromString[D sudoku.Digits[D], A sudoku.Area](rows ...string) GivenDigits[D, A] {
	givenDigits := make(map[sudoku.CellLocation]int)
	for row, rowContent := range rows {
		for col, cellContent := range rowContent {
			if cellContent < '1' || cellContent > '9' {
				continue
			}
			givenDigits[sudoku.CellLocation{Row: row, Col: col}] = int(cellContent - '0')
		}
	}
	return GivenDigits[D, A]{givenDigits: givenDigits}
}

func (r GivenDigits[D, A]) Apply(s sudoku.SudokuBuilder[D, A]) error {
	for loc, digit := range r.givenDigits {
		if err := s.SetCell(loc.Row, loc.Col, digit); err != nil {
			return err
		}
	}
	return nil
}
