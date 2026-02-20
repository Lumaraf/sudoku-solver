package sudoku

func NewSudoku9x9(rules ...Rule[Digits9, Area9x9]) (Sudoku[Digits9, Area9x9], error) {
	builder := NewSudokuBuilder9x9()
	if err := builder.Use(rules...); err != nil {
		return nil, err
	}
	return builder.Build()
}

func NewSudokuBuilder9x9() SudokuBuilder[Digits9, Area9x9] {
	return newSudokuBuilder[Digits9, Area9x9, grid9x9, size9]()
}

//type gridSize9 struct{}
//
//func (gridSize9) gridSize() int { return 9 }
//
//func (gridSize9) allCells() [2]uint64 {
//	return [2]uint64{
//		0xFFFFFFFFFFFFFFFF,
//		0b11111111111111111,
//	}
//}

type Area9x9 = area128[size9]

//type allDigits9 struct{}
//
//func (allDigits9) allDigits() uint16 {
//	return 0b111111111
//}

type Digits9 = digits_16[size9]

type grid9x9 [9 * 9]Digits9

type size9 struct{}

func (size9) allCells() [2]uint64 {
	return [2]uint64{
		0xFFFFFFFFFFFFFFFF,
		0b11111111111111111,
	}
}

func (size9) allDigits() uint16 {
	return 0b111111111
}

func (size9) Size() int {
	return 9
}

func (size9) BoxSize() (int, int) {
	return 3, 3
}

func (size9) GridCell(g *grid9x9, row, col int) *Digits9 {
	return &g[row*9+col]
}

func (s size9) PossibleLocations(g grid9x9, d Digits9) (a Area9x9) {
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
