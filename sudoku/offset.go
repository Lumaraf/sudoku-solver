package sudoku

type Offset struct {
	Row int
	Col int
}

type Offsets []Offset

var AdjacentOffsets = Offsets{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

func (o Offsets) Locations(cell CellLocation) func(yield func(cell CellLocation) bool) {
	return func(yield func(cell CellLocation) bool) {
		for _, offset := range o {
			row := cell.Row + offset.Row
			if row < 0 || row >= 9 {
				continue
			}
			col := cell.Col + offset.Col
			if col < 0 || col >= 9 {
				continue
			}
			if !yield(CellLocation{row, col}) {
				return
			}
		}
	}
}
