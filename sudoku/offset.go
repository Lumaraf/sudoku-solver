package sudoku

type Offset struct {
	Row int
	Col int
}

type Offsets []Offset

func (o Offsets) locations(size int, cell CellLocation) func(yield func(cell CellLocation) bool) {
	return func(yield func(cell CellLocation) bool) {
		for _, offset := range o {
			row := cell.Row + offset.Row
			if row < 0 || row >= size {
				continue
			}
			col := cell.Col + offset.Col
			if col < 0 || col >= size {
				continue
			}
			if !yield(CellLocation{row, col}) {
				return
			}
		}
	}
}
