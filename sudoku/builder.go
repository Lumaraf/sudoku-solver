package sudoku

type Rule[D Digits[D], A Area[A]] interface {
	Apply(s SudokuBuilder[D, A]) error
}

type Rules[D Digits[D], A Area[A]] []Rule[D, A]

func (rs Rules[D, A]) Apply(s SudokuBuilder[D, A]) error {
	for _, r := range rs {
		if err := r.Apply(s); err != nil {
			return err
		}
	}
	return nil
}

type Restriction[D Digits[D], A Area[A]] interface {
	Name() string
}

type ChangeProcessor[D Digits[D], A Area[A]] interface {
	Name() string
	ProcessChange(s Sudoku[D, A], cell CellLocation, mask D) error
}

type SolveProcessor[D Digits[D], A Area[A]] interface {
	Name() string
	ProcessSolve(s Sudoku[D, A], cell CellLocation, mask D) error
}

type SolveProcessors[D Digits[D], A Area[A]] []SolveProcessor[D, A]

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

type ExclusionChainSolveProcessor[D Digits[D], A Area[A]] struct{}

func (e ExclusionChainSolveProcessor[D, A]) Name() string {
	return "Exclusion Chain"
}

func (e ExclusionChainSolveProcessor[D, A]) ProcessSolve(s Sudoku[D, A], cell CellLocation, mask D) error {
	v, _ := mask.Single()
	area := s.PossibleLocations(v).Without(cell)
	for _, excludedCell := range area.And(s.GetExclusionArea(cell)).Locations {
		//for _, excludedCell := range s.GetExclusionArea(cell).Locations {
		if err := s.RemoveMask(excludedCell, mask); err != nil {
			return err
		}
	}
	return nil
}

type SudokuBuilder[D Digits[D], A Area[A]] interface {
	BaseSpec

	AreaOps[A]
	DigitsOps[D]

	buildTarget() Sudoku[D, A]

	Row(row int) A
	Column(col int) A
	Box(box int) A

	SetCell(row, col, value int) error
	MaskCell(row, col int, mask D) error

	Use(rules ...Rule[D, A]) error

	AddRestriction(rf Restriction[D, A])
	AddValidator(r Validator[D, A])
	AddChangeProcessor(cp ChangeProcessor[D, A])
	AddSolveProcessor(sp SolveProcessor[D, A])
	AddExclusionArea(l CellLocation, a A)

	Build() (Sudoku[D, A], error)
}

type sudokuBuilder[D Digits[D], A Area[A], G comparable, S size[D, A, G], GO gridOps[D, A, G]] struct {
	*sudoku[D, A, G, S, GO]
	solveProcessors SolveProcessors[D, A]
}

func newSudokuBuilder[D Digits[D], A Area[A], G comparable, S size[D, A, G], GO gridOps[D, A, G]]() SudokuBuilder[D, A] {
	s := newSudoku[D, A, G, S, GO]()
	s.changeProcessors = append(s.changeProcessors, SolveProcessors[D, A]{
		ExclusionChainSolveProcessor[D, A]{},
	})
	return &sudokuBuilder[D, A, G, S, GO]{
		sudoku: s,
		solveProcessors: SolveProcessors[D, A]{
			ExclusionChainSolveProcessor[D, A]{},
		},
	}
}

func (s *sudokuBuilder[D, A, G, S, GO]) buildTarget() Sudoku[D, A] {
	return s.sudoku
}

func (s *sudokuBuilder[D, A, G, S, GO]) SetCell(row, col, value int) error {
	return s.Set(CellLocation{row, col}, value)
}

func (s *sudokuBuilder[D, A, G, S, GO]) MaskCell(row, col int, mask D) error {
	return s.Mask(CellLocation{row, col}, mask)
}

func (s *sudokuBuilder[D, A, G, S, GO]) Use(rules ...Rule[D, A]) error {
	for _, r := range rules {
		if err := r.Apply(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *sudokuBuilder[D, A, G, S, GO]) AddRestriction(r Restriction[D, A]) {
	s.restrictions = append(s.restrictions, r)
}

func (s *sudokuBuilder[D, A, G, S, GO]) AddValidator(v Validator[D, A]) {
	s.validators = append(s.validators, v)
}

func (s *sudokuBuilder[D, A, G, S, GO]) AddChangeProcessor(cp ChangeProcessor[D, A]) {
	s.changeProcessors = append(s.changeProcessors, cp)
}

func (s *sudokuBuilder[D, A, G, S, GO]) AddSolveProcessor(sp SolveProcessor[D, A]) {
	s.solveProcessors = append(s.solveProcessors, sp)
}

func (s *sudokuBuilder[D, A, G, S, GO]) AddExclusionArea(l CellLocation, a A) {
	a = a.Without(l)
	s.exclusionAreas[l.Row][l.Col] = s.exclusionAreas[l.Row][l.Col].Or(a)
}

func (s *sudokuBuilder[D, A, G, S, GO]) Build() (Sudoku[D, A], error) {
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
