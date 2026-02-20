package sudoku

func NewSudoku12x12(rules ...Rule[Digits12, Area12x12]) (Sudoku[Digits12, Area12x12], error) {
	builder := NewSudokuBuilder12x12()
	if err := builder.Use(rules...); err != nil {
		return nil, err
	}
	return builder.Build()
}

func NewSudokuBuilder12x12() SudokuBuilder[Digits12, Area12x12] {
	return newSudokuBuilder[Digits12, Area12x12, grid12x12, size12]()
}

type Area12x12 = area256[size12]

type Digits12 = digits_16[size12]

type grid12x12 [12 * 12]Digits12

type size12 struct{}

func (size12) allCells() [4]uint64 {
	return [4]uint64{
		0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF,
		0b1111111111111111,
	}
}

func (size12) allDigits() uint16 {
	return 0b111111111111
}

func (s size12) Size() int {
	return 12
}

func (s size12) BoxSize() (int, int) {
	return 3, 4
}

func (s size12) GridCell(g *grid12x12, row, col int) *Digits12 {
	return &g[row*12+col]
}

func (s size12) PossibleLocations(g grid12x12, d Digits12) (a Area12x12) {
	for row := 0; row < s.Size(); row++ {
		for col := 0; col < s.Size(); col++ {
			cell := s.GridCell(&g, row, col)
			if !(*cell).And(d).Empty() {
				a = a.With(CellLocation{row, col})
			}
		}
	}
	return
}
