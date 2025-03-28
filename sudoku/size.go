package sudoku

type size interface {
	Size() int
	BoxSize() (int, int)
}

type Size4 struct{}

func (s Size4) Size() int {
	return 4
}

func (s Size4) BoxSize() (int, int) {
	return 2, 2
}
