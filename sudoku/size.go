package sudoku

type BaseSpec interface {
	Size() int
	BoxSize() (int, int)
}

type size[D Digits[D], A Area[A], G comparable] interface {
	BaseSpec

	GridCell(g *G, row, col int) *D
	PossibleLocations(g G, d D) A
}
