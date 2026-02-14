package rule

import "github.com/lumaraf/sudoku-solver/sudoku"

var (
	parityEven = 0
	parityOdd  = 1
)

type ParityRule[D sudoku.Digits[D], A sudoku.Area] struct {
	parities map[sudoku.CellLocation]int
}

func ParityFromString[D sudoku.Digits[D], A sudoku.Area](rows ...string) ParityRule[D, A] {
	parities := make(map[sudoku.CellLocation]int)
	for row, rowContent := range rows {
		for col, cellContent := range rowContent {
			switch cellContent {
			case 'E':
				parities[sudoku.CellLocation{Row: row, Col: col}] = parityEven
			case 'O':
				parities[sudoku.CellLocation{Row: row, Col: col}] = parityOdd
			}
		}
	}
	return ParityRule[D, A]{parities: parities}
}

func (r ParityRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	evenValues := make([]int, 0, sb.Size()/2)
	oddValues := make([]int, 0, sb.Size()/2)
	for v := 1; v <= sb.Size(); v++ {
		if v%2 == 0 {
			evenValues = append(evenValues, v)
		} else {
			oddValues = append(oddValues, v)
		}
	}

	even := sb.NewDigits(evenValues...)
	odd := sb.NewDigits(oddValues...)

	for l, parity := range r.parities {
		switch parity {
		case parityEven:
			if err := sb.MaskCell(l.Row, l.Col, even); err != nil {
				return err
			}
		case parityOdd:
			if err := sb.MaskCell(l.Row, l.Col, odd); err != nil {
				return err
			}
		}
	}
	return nil
}
