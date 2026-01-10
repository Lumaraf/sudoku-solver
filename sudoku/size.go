package sudoku

type BaseSpec interface {
	Size() int
	BoxSize() (int, int)
}

type DigitsSpec[D Digits] interface {
	NewDigits(values ...int) D
	AllDigits() D
	IntersectDigits(d1 D, d2 D) D
	UnionDigits(d1 D, d2 D) D
	InvertDigits(d D) D
}

type AreaSpec[A Area] interface {
	NewArea(locs ...CellLocation) A
	NewAreaFromOffsets(center CellLocation, o Offsets) A
	AreaWith(a *A, l CellLocation)
	AreaWithout(a *A, l CellLocation)
	IntersectAreas(a1 A, a2 A) A
	UnionAreas(a1 A, a2 A) A
	InvertArea(a A) A
}

type size[D Digits, A Area, G comparable] interface {
	BaseSpec
	DigitsSpec[D]
	AreaSpec[A]

	GridCell(g *G, row, col int) *D
}
