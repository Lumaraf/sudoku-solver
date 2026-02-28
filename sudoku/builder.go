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
	ProcessChanges(s Sudoku[D, A]) error
}

type SolveProcessor[D Digits[D], A Area[A]] interface {
	Name() string
	ProcessSolve(s Sudoku[D, A], cell CellLocation, mask D) error
}

type SolveProcessors[D Digits[D], A Area[A]] []SolveProcessor[D, A]

func (sps SolveProcessors[D, A]) Name() string {
	return "Solve Processors"
}

func (sps SolveProcessors[D, A]) ProcessChanges(s Sudoku[D, A]) error {
	for _, l := range s.ChangedArea().Locations {
		mask := s.Get(l)
		if mask.Count() != 1 {
			continue
		}
		s.setSolved(l)
		for _, sp := range sps {
			if err := sp.ProcessSolve(s, l, mask); err != nil {
				return err
			}
		}
	}
	return nil
}

type ExclusionAreaSolveProcessor[D Digits[D], A Area[A]] struct{}

func (e ExclusionAreaSolveProcessor[D, A]) Name() string {
	return "Exclusion Chain"
}

func (e ExclusionAreaSolveProcessor[D, A]) ProcessSolve(s Sudoku[D, A], cell CellLocation, mask D) error {
	v, _ := mask.Single()
	area := s.PossibleLocations(v).Without(cell)
	for _, excludedCell := range area.And(s.GetExclusionArea(cell)).Locations {
		if err := s.RemoveMask(excludedCell, mask); err != nil {
			return err
		}
	}
	return nil
}

type OffsetMaskChangeProcessor[D Digits[D], A Area[A]] struct {
	offsetMasks map[int]map[Offset]D
}

func (cp OffsetMaskChangeProcessor[D, A]) Name() string {
	return "Offset Mask"
}

func (cp OffsetMaskChangeProcessor[D, A]) ProcessChanges(s Sudoku[D, A]) error {
	for _, cell := range s.ChangedArea().Locations {
		mask := s.Get(cell)
		cellMasks := make(map[Offset]D)
		for v := range mask.Values {
			if offsetMasks, ok := cp.offsetMasks[v]; ok {
				for offset, offsetMask := range offsetMasks {
					cellMasks[offset] = cellMasks[offset].Or(offsetMask)
				}
			}
		}

		for offset, combinedMask := range cellMasks {
			offsetCell := CellLocation{cell.Row + offset.Row, cell.Col + offset.Col}
			if offsetCell.Row < 0 || offsetCell.Row >= s.Size() || offsetCell.Col < 0 || offsetCell.Col >= s.Size() {
				continue
			}

			if err := s.Mask(offsetCell, combinedMask); err != nil {
				return err
			}
		}
	}
	return nil
}

type OffsetMaskRestriction[D Digits[D], A Area[A]] struct {
	offsetMasks map[int]map[Offset]D
}

func (r OffsetMaskRestriction[D, A]) Name() string {
	return "Offset Mask Restriction"
}

func (r OffsetMaskRestriction[D, A]) MasksForValue(v int) map[Offset]D {
	return r.offsetMasks[v]
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
	AddOffsetMask(v int, offset Offset, mask D)

	Build() (Sudoku[D, A], error)
}

type sudokuBuilder[D Digits[D], A Area[A], G comparable, S size[D, A, G], GO gridOps[D, A, G]] struct {
	*sudoku[D, A, G, S, GO]
	solveProcessors SolveProcessors[D, A]
	offsetMasks     map[int]map[Offset]D
}

func newSudokuBuilder[D Digits[D], A Area[A], G comparable, S size[D, A, G], GO gridOps[D, A, G]]() SudokuBuilder[D, A] {
	s := newSudoku[D, A, G, S, GO]()
	s.changeProcessors = append(s.changeProcessors, SolveProcessors[D, A]{
		ExclusionAreaSolveProcessor[D, A]{},
	})
	return &sudokuBuilder[D, A, G, S, GO]{
		sudoku: s,
		solveProcessors: SolveProcessors[D, A]{
			ExclusionAreaSolveProcessor[D, A]{},
		},
		offsetMasks: make(map[int]map[Offset]D),
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

func (s *sudokuBuilder[D, A, G, S, GO]) AddOffsetMask(v int, offset Offset, mask D) {
	if s.offsetMasks[v] == nil {
		s.offsetMasks[v] = make(map[Offset]D)
	} else if existingMask, ok := s.offsetMasks[v][offset]; ok {
		mask = mask.And(existingMask)
	}
	s.offsetMasks[v][offset] = mask
}

func (s *sudokuBuilder[D, A, G, S, GO]) Build() (Sudoku[D, A], error) {
	s.changeProcessors[0] = s.solveProcessors
	if len(s.offsetMasks) > 0 {
		s.changeProcessors = append(s.changeProcessors, OffsetMaskChangeProcessor[D, A]{
			offsetMasks: s.offsetMasks,
		})
		s.restrictions = append(s.restrictions, OffsetMaskRestriction[D, A]{
			offsetMasks: s.offsetMasks,
		})
	}
	for row := 0; row < s.Size(); row++ {
		for col := 0; col < s.Size(); col++ {
			l := CellLocation{row, col}
			mask := s.sudoku.Get(l)
			if mask == s.sudoku.AllDigits() {
				continue
			}
		}
	}
	if err := s.sudoku.ProcessChanges(); err != nil {
		return nil, err
	}
	return s.sudoku, nil
}
