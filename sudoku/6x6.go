package sudoku

func NewSudoku6x6(rules ...Rule[Digits6, Area6x6]) (Sudoku[Digits6, Area6x6], error) {
	builder := NewSudokuBuilder6x6()
	if err := builder.Use(rules...); err != nil {
		return nil, err
	}
	return builder.Build()
}

func NewSudokuBuilder6x6() SudokuBuilder[Digits6, Area6x6] {
	return newSudokuBuilder[Digits6, Area6x6, grid6x6[Digits6], size6]()
}

type gridSize6 struct{}

func (gridSize6) gridSize() int { return 6 }

type Area6x6 = area128[gridSize6]

type Digits6 = digits_16

type grid6x6[D Digits] [6][6]D

type size6 struct {
	digitsOps_16
}

func (s size6) Size() int {
	return 6
}

func (s size6) BoxSize() (int, int) {
	return 2, 3
}

func (s size6) GridCell(g *grid6x6[Digits6], row, col int) *Digits6 {
	return &g[row][col]
}

func (s size6) NewDigits(values ...int) Digits6 {
	d := Digits6(0)
	for _, v := range values {
		d = d.withOption(v)
	}
	return d & s.AllDigits()
}

func (s size6) AllDigits() Digits6 {
	return digits_16(0b111111)
}

func (s size6) NewArea(locs ...CellLocation) Area6x6 {
	a := Area6x6{}
	for _, loc := range locs {
		a = a.with(loc, true)
	}
	return a
}

func (s size6) NewAreaFromOffsets(center CellLocation, o Offsets) Area6x6 {
	a := Area6x6{}
	for loc := range o.locations(s.Size(), center) {
		a = a.with(loc, true)
	}
	return a
}

func (s size6) AreaWith(a *Area6x6, l CellLocation) {
	*a = a.with(l, true)
}

func (s size6) AreaWithout(a *Area6x6, l CellLocation) {
	*a = a.with(l, false)
}

func (s size6) IntersectAreas(a1 Area6x6, a2 Area6x6) Area6x6 {
	return And(a1, a2)
}

func (s size6) UnionAreas(a1 Area6x6, a2 Area6x6) Area6x6 {
	return Or(a1, a2)
}

func (s size6) InvertArea(a Area6x6) Area6x6 {
	return And(Not(a), Area6x6{
		gs:   a.gs,
		bits: [2]uint64{0b111111_111111_111111_111111_111111_111111, 0},
	})
}
