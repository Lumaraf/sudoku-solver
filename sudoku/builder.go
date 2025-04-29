package sudoku

type RestrictionFactory[D Digits, A Area] func(c SudokuBuilder[D, A]) error

type SudokuBuilder[D Digits, A Area] interface {
	baseSpec
	digitsSpec[D]
	areaSpec[A]

	Row(row int) A
	Column(col int) A
	Box(box int) A

	SetCell(row, col, value int) error

	AddRestrictionFactory(r RestrictionFactory[D, A]) error

	AddRestriction(r Restriction[D, A])
	AddValidator(r Validator[D, A])
	AddSolveProcessor(sp SolveProcessor[D, A])
	AddExclusionArea(l CellLocation, a A)

	Build() Sudoku[D, A]
}

func NewSudokuBuilder9x9() SudokuBuilder[Digits9, Area9x9] {
	return newSudokuBuilder[Digits9, Area9x9, grid9x9[Digits9], size9]()
}

type sudokuBuilder[D Digits, A Area, G comparable, S size[D, A, G]] struct {
	*sudoku[D, A, G, S]
}

func newSudokuBuilder[D Digits, A Area, G comparable, S size[D, A, G]]() SudokuBuilder[D, A] {
	return &sudokuBuilder[D, A, G, S]{
		sudoku: newSudoku[D, A, G, S](),
	}
}

func (s *sudokuBuilder[D, A, G, S]) Row(row int) (a A) {
	for col := 0; col < s.Size(); col++ {
		s.AreaWith(&a, CellLocation{row, col})
	}
	return
}

func (s *sudokuBuilder[D, A, G, S]) Column(col int) (a A) {
	for row := 0; row < s.Size(); row++ {
		s.AreaWith(&a, CellLocation{row, col})
	}
	return
}

func (s *sudokuBuilder[D, A, G, S]) Box(box int) (a A) {
	boxRows, boxCols := s.BoxSize()
	boxesPerRow := s.Size() / boxRows
	rowOffset := (box / boxesPerRow) * boxRows
	colOffset := (box % boxesPerRow) * boxCols
	for row := 0; row < boxRows; row++ {
		for col := 0; col < boxCols; col++ {
			s.AreaWith(&a, CellLocation{rowOffset + row, colOffset + col})
		}
	}
	return
}

func (s *sudokuBuilder[D, A, G, S]) SetCell(row, col, value int) error {
	return s.Set(CellLocation{row, col}, value)
}

func (s *sudokuBuilder[D, A, G, S]) AddRestrictionFactory(r RestrictionFactory[D, A]) error {
	return r(s)
}

func (s *sudokuBuilder[D, A, G, S]) AddRestriction(r Restriction[D, A]) {
	s.restrictions = append(s.restrictions, r)
}

func (s *sudokuBuilder[D, A, G, S]) AddValidator(v Validator[D, A]) {
	s.validators = append(s.validators, v)
}

func (s *sudokuBuilder[D, A, G, S]) AddSolveProcessor(sp SolveProcessor[D, A]) {
	s.solveProcessors = append(s.solveProcessors, sp)
}

func (s *sudokuBuilder[D, A, G, S]) AddExclusionArea(l CellLocation, a A) {
	s.AreaWithout(&a, l)
	s.exclusionAreas[l.Row][l.Col] = s.UnionAreas(s.exclusionAreas[l.Row][l.Col], a)
}

func (s *sudokuBuilder[D, A, G, S]) Build() Sudoku[D, A] {
	return s.sudoku
}
