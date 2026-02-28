package rule

import (
	"errors"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

// the digits in column 1, 5, and 9 contain the column in the same row which contains the digit 1, 5, and 9 respectively.
type Rule159[D sudoku.Digits[D], A sudoku.Area[A]] struct{}

func (r Rule159[D, A]) Name() string {
	return "159"
}

func (r Rule159[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	if sb.Size() != 9 {
		return errors.New("rule 159 only works with 9x9 sudoku")
	}

	sb.AddValidator(Rule159Validator[D, A]{})
	sb.AddChangeProcessor(Rule159ChangeProcessor[D, A]{})
	return nil
}

type Rule159Validator[D sudoku.Digits[D], A sudoku.Area[A]] struct{}

func (v Rule159Validator[D, A]) Name() string {
	return "159"
}

func (v Rule159Validator[D, A]) Validate(s sudoku.Sudoku[D, A]) error {
	for _, l := range s.SolvedArea().Locations {
		d := s.Get(l)
		v, ok := d.Single()
		if !ok {
			continue
		}
		switch v {
		case 1, 5, 9:
			other := s.Get(sudoku.CellLocation{Row: l.Row, Col: v - 1})
			if !other.CanContain(l.Col + 1) {
				return errors.New("rule 159 violation")
			}
		}
	}
	return nil
}

type Rule159ChangeProcessor[D sudoku.Digits[D], A sudoku.Area[A]] struct{}

func (cp Rule159ChangeProcessor[D, A]) Name() string {
	return "159"
}

func (cp Rule159ChangeProcessor[D, A]) ProcessChanges(s sudoku.Sudoku[D, A]) error {
	for _, l := range s.ChangedArea().Locations {
		d := s.Get(l)
		switch l.Col {
		case 0, 4, 8:
			for v := range s.AllDigits().And(d.Not()).Values {
				if err := s.RemoveOption(sudoku.CellLocation{Row: l.Row, Col: v - 1}, l.Col+1); err != nil {
					return err
				}
			}
		}

		if !d.CanContain(1) {
			if err := s.RemoveOption(sudoku.CellLocation{Row: l.Row, Col: 0}, l.Col+1); err != nil {
				return err
			}
		}
		if !d.CanContain(5) {
			if err := s.RemoveOption(sudoku.CellLocation{Row: l.Row, Col: 4}, l.Col+1); err != nil {
				return err
			}
		}
		if !d.CanContain(9) {
			if err := s.RemoveOption(sudoku.CellLocation{Row: l.Row, Col: 8}, l.Col+1); err != nil {
				return err
			}
		}
	}
	return nil
}
