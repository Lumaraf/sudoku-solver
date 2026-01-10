package sudoku

type Rule[D Digits, A Area] interface {
	Apply(s SudokuBuilder[D, A]) error
}

type Restriction[D Digits, A Area] interface {
	Name() string
}

type ChangeProcessor[D Digits, A Area] interface {
	Name() string
	ProcessChange(s Sudoku[D, A], cell CellLocation, mask D) error
}

type SolveProcessor[D Digits, A Area] interface {
	Name() string
	ProcessSolve(s Sudoku[D, A], cell CellLocation, mask D) error
}

type SolveProcessors[D Digits, A Area] []SolveProcessor[D, A]

func (sps SolveProcessors[D, A]) Name() string {
	return "Solve Processors"
}

func (sps SolveProcessors[D, A]) ProcessChange(s Sudoku[D, A], cell CellLocation, mask D) error {
	if mask.Count() != 1 {
		return nil
	}
	s.setSolved(cell)
	for _, sp := range sps {
		if err := sp.ProcessSolve(s, cell, mask); err != nil {
			return err
		}
	}
	return nil
}

type ExclusionChainSolveProcessor[D Digits, A Area] struct{}

func (e ExclusionChainSolveProcessor[D, A]) Name() string {
	return "Exclusion Chain"
}

func (e ExclusionChainSolveProcessor[D, A]) ProcessSolve(s Sudoku[D, A], cell CellLocation, mask D) error {
	for _, excludedCell := range s.GetExclusionArea(cell).Locations {
		if err := s.RemoveMask(excludedCell, mask); err != nil {
			return err
		}
	}
	return nil
}

type SudokuBuilder[D Digits, A Area] interface {
	BaseSpec
	DigitsSpec[D]
	AreaSpec[A]

	buildTarget() Sudoku[D, A]

	Row(row int) A
	Column(col int) A
	Box(box int) A

	SetCell(row, col, value int) error

	Use(rules ...Rule[D, A]) error

	AddRestriction(rf Restriction[D, A])
	AddValidator(r Validator[D, A])
	AddChangeProcessor(cp ChangeProcessor[D, A])
	AddSolveProcessor(sp SolveProcessor[D, A])
	AddExclusionArea(l CellLocation, a A)

	Build() (Sudoku[D, A], error)
}

type sudokuBuilder[D Digits, A Area, G comparable, S size[D, A, G]] struct {
	*sudoku[D, A, G, S]
	solveProcessors SolveProcessors[D, A]
}

func newSudokuBuilder[D Digits, A Area, G comparable, S size[D, A, G]]() SudokuBuilder[D, A] {
	s := newSudoku[D, A, G, S]()
	s.changeProcessors = append(s.changeProcessors, SolveProcessors[D, A]{
		ExclusionChainSolveProcessor[D, A]{},
	})
	return &sudokuBuilder[D, A, G, S]{
		sudoku: s,
		solveProcessors: SolveProcessors[D, A]{
			ExclusionChainSolveProcessor[D, A]{},
		},
	}
}

func (s *sudokuBuilder[D, A, G, S]) buildTarget() Sudoku[D, A] {
	return s.sudoku
}

func (s *sudokuBuilder[D, A, G, S]) SetCell(row, col, value int) error {
	return s.Set(CellLocation{row, col}, value)
}

func (s *sudokuBuilder[D, A, G, S]) Use(rules ...Rule[D, A]) error {
	for _, r := range rules {
		if err := r.Apply(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *sudokuBuilder[D, A, G, S]) AddRestriction(r Restriction[D, A]) {
	s.restrictions = append(s.restrictions, r)
}

func (s *sudokuBuilder[D, A, G, S]) AddValidator(v Validator[D, A]) {
	s.validators = append(s.validators, v)
}

func (s *sudokuBuilder[D, A, G, S]) AddChangeProcessor(cp ChangeProcessor[D, A]) {
	s.changeProcessors = append(s.changeProcessors, cp)
}

func (s *sudokuBuilder[D, A, G, S]) AddSolveProcessor(sp SolveProcessor[D, A]) {
	s.solveProcessors = append(s.solveProcessors, sp)
}

func (s *sudokuBuilder[D, A, G, S]) AddExclusionArea(l CellLocation, a A) {
	s.AreaWithout(&a, l)
	//if v, isSingle := s.Get(l).Single(); !isSingle {
	//	for _, cell := range a.Locations {
	//		_ = s.RemoveOption(cell, v)
	//	}
	//}
	s.exclusionAreas[l.Row][l.Col] = s.UnionAreas(s.exclusionAreas[l.Row][l.Col], a)
}

func (s *sudokuBuilder[D, A, G, S]) Build() (Sudoku[D, A], error) {
	s.changeProcessors[0] = s.solveProcessors
	for row := 0; row < s.Size(); row++ {
		for col := 0; col < s.Size(); col++ {
			l := CellLocation{row, col}
			mask := s.sudoku.Get(l)
			if mask == s.sudoku.AllDigits() {
				continue
			}
			if err := s.sudoku.processChange(l, mask); err != nil {
				return nil, err
			}
		}
	}
	return s.sudoku, nil
}
