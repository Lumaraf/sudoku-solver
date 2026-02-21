package sudoku

func NewSudoku6x6(rules ...Rule[Digits6, Area6x6]) (Sudoku[Digits6, Area6x6], error) {
	builder := NewSudokuBuilder6x6()
	if err := builder.Use(rules...); err != nil {
		return nil, err
	}
	return builder.Build()
}

func NewSudokuBuilder6x6() SudokuBuilder[Digits6, Area6x6] {
	return newSudokuBuilder[Digits6, Area6x6, grid6x6, size6, genericGridOps[Digits6, Area6x6, grid6x6, size6]]()
}

type Area6x6 = area128[size6]

type Digits6 = digits_16[size6]

type grid6x6 [6 * 6]Digits6

type size6 struct{}

func (size6) allCells() [2]uint64 {
	return [2]uint64{0b111111_111111_111111_111111_111111_111111, 0}
}

func (size6) allDigits() uint16 {
	return 0b111111
}

func (s size6) Size() int {
	return 6
}

func (s size6) BoxSize() (int, int) {
	return 2, 3
}

func (s size6) GridCell(g *grid6x6, row, col int) *Digits6 {
	return &g[row*6+col]
}

func (s size6) NewArea(locs ...CellLocation) Area6x6 {
	a := Area6x6{}
	for _, loc := range locs {
		a = a.With(loc)
	}
	return a
}

func (s size6) NewAreaFromOffsets(center CellLocation, o Offsets) Area6x6 {
	a := Area6x6{}
	for loc := range o.locations(s.Size(), center) {
		a = a.With(loc)
	}
	return a
}
