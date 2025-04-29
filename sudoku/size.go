package sudoku

type baseSpec interface {
	Size() int
	BoxSize() (int, int)
}

type digitsSpec[D Digits] interface {
	NewDigits(values ...int) D
	AllDigits() D
	IntersectDigits(d1 D, d2 D) D
	UnionDigits(d1 D, d2 D) D
	InvertDigits(d D) D
}

type areaSpec[A Area] interface {
	NewArea() A
	AreaWith(a *A, l CellLocation)
	AreaWithout(a *A, l CellLocation)
	IntersectAreas(a1 A, a2 A) A
	UnionAreas(a1 A, a2 A) A
	InvertArea(a A) A
}

type size[D Digits, A Area, G comparable] interface {
	baseSpec
	digitsSpec[D]
	areaSpec[A]

	GridCell(g *G, row, col int) *D
}

func newSize9() size[digits16, area128, grid9x9[digits16]] {
	return size9{}
}

type size9 struct{}

func (s size9) Size() int {
	return 9
}

func (s size9) BoxSize() (int, int) {
	return 3, 3
}

func (s size9) GridCell(g *grid9x9[digits16], row, col int) *digits16 {
	return &g[row][col]
}

func (s size9) NewDigits(values ...int) digits16 {
	d := digits16(0)
	for _, v := range values {
		d = d.withOption(v)
	}
	return d & s.AllDigits()
}

func (s size9) AllDigits() digits16 {
	return digits16(0b111111111)
}

func (s size9) IntersectDigits(d1 digits16, d2 digits16) digits16 {
	return d1 & d2
}

func (s size9) UnionDigits(d1 digits16, d2 digits16) digits16 {
	return d1 | d2
}

func (s size9) InvertDigits(d digits16) digits16 {
	return d ^ s.AllDigits()
}

func (s size9) NewArea() area128 {
	return area128{}
}

func (s size9) AreaWith(a *area128, l CellLocation) {
	*a = a.with(l, true)
}

func (s size9) AreaWithout(a *area128, l CellLocation) {
	*a = a.with(l, false)
}

func (s size9) IntersectAreas(a1 area128, a2 area128) (i area128) {
	i[0] = a1[0] & a2[0]
	i[1] = a1[1] & a2[1]
	return
}

func (s size9) UnionAreas(a1 area128, a2 area128) (u area128) {
	u[0] = a1[0] | a2[0]
	u[1] = a1[1] | a2[1]
	return
}

func (s size9) InvertArea(a area128) area128 {
	a[0] = ^a[0]
	a[1] = (^a[1]) & 0b11111111111111111
	return a
}
