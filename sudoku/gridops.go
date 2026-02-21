package sudoku

type gridOps[D Digits[D], A Area[A], G comparable] interface {
	PossibleLocations(g G, d D) A
}

type genericGridOps[D Digits[D], A Area[A], G comparable, S size[D, A, G]] struct{}

func (o genericGridOps[D, A, G, S]) PossibleLocations(g G, d D) (a A) {
	var s S
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
