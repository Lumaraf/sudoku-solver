package sudoku

func NewSudoku9x9(rules ...Rule[Digits9, Area9x9]) (Sudoku[Digits9, Area9x9], error) {
	builder := NewSudokuBuilder9x9()
	if err := builder.Use(rules...); err != nil {
		return nil, err
	}
	return builder.Build()
}

func NewSudokuBuilder9x9() SudokuBuilder[Digits9, Area9x9] {
	return newSudokuBuilder[Digits9, Area9x9, grid9x9[Digits9], size9]()
}

type gridSize9 struct{}

func (gridSize9) gridSize() int { return 9 }

type Area9x9 = area128[gridSize9]

type Digits9 = digits_16

type grid9x9[D Digits] [9][9]D

type size9 struct{ digitsOps_16 }

func (s size9) Size() int {
	return 9
}

func (s size9) BoxSize() (int, int) {
	return 3, 3
}

func (s size9) GridCell(g *grid9x9[digits_16], row, col int) *Digits9 {
	return &g[row][col]
}

func (s size9) NewDigits(values ...int) Digits9 {
	d := Digits9(0)
	for _, v := range values {
		d = d.withOption(v)
	}
	return d & s.AllDigits()
}

func (s size9) AllDigits() Digits9 {
	return digits_16(0b111111111)
}

func (s size9) NewArea(locs ...CellLocation) Area9x9 {
	a := Area9x9{}
	for _, loc := range locs {
		a = a.with(loc, true)
	}
	return a
}

func (s size9) NewAreaFromOffsets(center CellLocation, o Offsets) Area9x9 {
	a := Area9x9{}
	for loc := range o.locations(s.Size(), center) {
		a = a.with(loc, true)
	}
	return a
}

func (s size9) AreaWith(a *Area9x9, l CellLocation) {
	*a = a.with(l, true)
}

func (s size9) AreaWithout(a *Area9x9, l CellLocation) {
	*a = a.with(l, false)
}

func (s size9) IntersectAreas(a1 Area9x9, a2 Area9x9) (i Area9x9) {
	i.bits[0] = a1.bits[0] & a2.bits[0]
	i.bits[1] = a1.bits[1] & a2.bits[1]
	return
}

func (s size9) UnionAreas(a1 Area9x9, a2 Area9x9) (u Area9x9) {
	u.bits[0] = a1.bits[0] | a2.bits[0]
	u.bits[1] = a1.bits[1] | a2.bits[1]
	return
}

func (s size9) InvertArea(a Area9x9) Area9x9 {
	a.bits[0] = ^a.bits[0]
	a.bits[1] = (^a.bits[1]) & 0b11111111111111111
	return a
}
