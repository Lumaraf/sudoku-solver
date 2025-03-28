package sudoku

import "fmt"

type ErrEmptyCell CellLocation

func (e ErrEmptyCell) Error() string {
	return fmt.Sprintf("empty cell %d,%d", e.Row, e.Col)
}
