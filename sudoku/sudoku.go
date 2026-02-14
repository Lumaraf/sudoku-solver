package sudoku

import (
	"errors"
	"fmt"
)

var mutipleSolutionsError = errors.New("multiple solutions")

type Sudoku[D Digits[D], A Area] interface {
	BaseSpec
	DigitsSpec[D]
	AreaSpec[A]

	Row(row int) A
	Column(col int) A
	Box(box int) A

	// Get retrieves the digits at the specified cell location.
	Get(l CellLocation) D

	// Set assigns a value to the specified cell location.
	Set(l CellLocation, v int) error

	// Try attempts to apply a function to a clone of the Sudoku puzzle.
	Try(f func(s Sudoku[D, A]) error) error

	// Mask restricts the possible digits at the specified cell location.
	Mask(l CellLocation, d D) error

	// RemoveMask removes the rule of possible digits at the specified cell location.
	RemoveMask(l CellLocation, d D) error

	// RemoveOption removes a specific digit option from the specified cell location.
	RemoveOption(l CellLocation, v int) error

	// PossibleLocations returns the area of possible locations for the specified digit.
	PossibleLocations(v int) A

	// SolvedArea returns the area of the grid that is solved.
	SolvedArea() A

	// ChangedArea returns the area of the grid that has changed.
	ChangedArea() A

	// NextChangedArea returns the area of the grid that will change next.
	NextChangedArea() A

	// GetExclusionArea returns the exclusion area for the specified cell location.
	GetExclusionArea(l CellLocation) A

	// IsUniqueArea checks if all cells in the area require unique digits.
	IsUniqueArea(area A) bool

	// SetLogger sets the logger for the Sudoku puzzle.
	SetLogger(logger Logger[D])

	// Validate checks the validity of the current Sudoku puzzle state.
	Validate() error

	// Print prints the current state of the Sudoku puzzle.
	Print() error

	// Stats returns the statistics of the Sudoku puzzle.
	Stats() Stats

	// IsSolved checks if the Sudoku puzzle is solved.
	IsSolved() bool

	NewSolver() Solver[D, A]

	NewGuesser() Guesser[D, A]

	getRestrictions() []any

	setSolved(l CellLocation)
}

type CellLocation struct {
	Row int
	Col int
}

func (l CellLocation) Box() int {
	return l.Row/3*3 + l.Col/3
}

type sudoku[D Digits[D], A Area, G comparable, S size[D, A, G]] struct {
	size[D, A, G]
	grid             G
	exclusionAreas   [][]A
	restrictions     []any
	validators       []Validator[D, A]
	changeProcessors []ChangeProcessor[D, A]
	findAllOptions   bool
	chainLimit       int
	changed          A
	nextChanged      A
	solved           A
	stats            Stats
	logger           Logger[D]
}

type Stats struct {
	CellUpdates        int
	SolverRuns         int
	SolverHits         int
	ExclusionChainRuns int
	GuesserRuns        int
	GuessMisses        int
}

func newSudoku[D Digits[D], A Area, G comparable, S size[D, A, G]]() *sudoku[D, A, G, S] {
	sizeSpec := *new(S)
	s := sudoku[D, A, G, S]{
		size:        sizeSpec,
		chainLimit:  2,
		nextChanged: sizeSpec.InvertArea(*new(A)),
		logger:      voidLogger[D]{},
	}
	s.exclusionAreas = make([][]A, s.Size())
	for row := 0; row < s.Size(); row++ {
		s.exclusionAreas[row] = make([]A, s.Size())
		for col := 0; col < s.Size(); col++ {
			cell := s.GridCell(&s.grid, row, col)
			*cell = sizeSpec.AllDigits()
		}
	}
	return &s
}

func (s *sudoku[D, A, G, S]) Row(row int) (a A) {
	for col := 0; col < s.Size(); col++ {
		s.AreaWith(&a, CellLocation{row, col})
	}
	return
}

func (s *sudoku[D, A, G, S]) Column(col int) (a A) {
	for row := 0; row < s.Size(); row++ {
		s.AreaWith(&a, CellLocation{row, col})
	}
	return
}

func (s *sudoku[D, A, G, S]) Box(box int) (a A) {
	boxRows, boxCols := s.BoxSize()
	boxesPerCol := s.Size() / boxCols
	rowOffset := (box / boxesPerCol) * boxRows
	colOffset := (box % boxesPerCol) * boxCols
	for row := 0; row < boxRows; row++ {
		for col := 0; col < boxCols; col++ {
			s.AreaWith(&a, CellLocation{rowOffset + row, colOffset + col})
		}
	}
	return
}

func (s *sudoku[D, A, G, S]) Get(l CellLocation) D {
	cell := s.GridCell(&s.grid, l.Row, l.Col)
	return *cell
}

func (s *sudoku[D, A, G, S]) Set(l CellLocation, v int) error {
	if err := s.checkValue(v); err != nil {
		return err
	}
	cell := s.GridCell(&s.grid, l.Row, l.Col)
	if !(*cell).CanContain(v) {
		return errors.New("cell doesn't allow value")
	}
	if (*cell).Count() > 1 {
		s.AreaWith(&s.nextChanged, l)
		s.stats.CellUpdates++
		oldCell := *cell
		*cell = s.NewDigits(v)
		s.logger.UpdateCell(l, oldCell, *cell)
		return s.processChange(l, *cell)
	}
	return nil
}

func (s *sudoku[D, A, G, S]) Try(f func(s Sudoku[D, A]) error) error {
	clone := *s
	clone.logger = voidLogger[D]{}
	clone.nextChanged = *new(A)
	return f(&clone)
}

func (s *sudoku[D, A, G, S]) Mask(l CellLocation, d D) error {
	target := s.GridCell(&s.grid, l.Row, l.Col)
	if !s.IntersectDigits(*target, s.InvertDigits(d)).Empty() {
		s.AreaWith(&s.nextChanged, l)
		s.stats.CellUpdates++
		oldDigits := *target
		newDigits := s.IntersectDigits(*target, d)
		if newDigits.Empty() {
			return ErrEmptyCell(l)
		}
		*target = newDigits
		s.logger.UpdateCell(l, oldDigits, newDigits)
		return s.processChange(l, newDigits)
	}
	return nil
}

func (s *sudoku[D, A, G, S]) RemoveMask(l CellLocation, d D) error {
	target := s.GridCell(&s.grid, l.Row, l.Col)
	if !s.IntersectDigits(*target, d).Empty() {
		s.AreaWith(&s.nextChanged, l)
		s.stats.CellUpdates++
		oldDigits := *target
		newDigits := s.IntersectDigits(*target, s.InvertDigits(d))
		if newDigits.Empty() {
			s.AreaWithout(&s.solved, l)
			return ErrEmptyCell(l)
		}
		*target = newDigits
		s.logger.UpdateCell(l, oldDigits, newDigits)
		return s.processChange(l, newDigits)
	}
	return nil
}

func (s *sudoku[D, A, G, S]) RemoveOption(l CellLocation, v int) error {
	if err := s.checkValue(v); err != nil {
		return err
	}
	mask := s.NewDigits(v)
	return s.RemoveMask(l, mask)
}

func (s *sudoku[D, A, G, S]) checkValue(v int) error {
	if v < 1 || v > s.Size() {
		return ErrValueOutOfRange
	}
	return nil
}

func (s *sudoku[D, A, G, S]) PossibleLocations(v int) A {
	return s.size.PossibleLocations(s.grid, s.NewDigits(v))
}

func (s *sudoku[D, A, G, S]) SolvedArea() A {
	return s.solved
}

func (s *sudoku[D, A, G, S]) ChangedArea() A {
	return s.changed
}

func (s *sudoku[D, A, G, S]) NextChangedArea() A {
	return s.nextChanged
}

func (s *sudoku[D, A, G, S]) GetExclusionArea(l CellLocation) A {
	return s.exclusionAreas[l.Row][l.Col]
}

func (s *sudoku[D, A, G, S]) IsUniqueArea(area A) bool {
	for _, l := range area.Locations {
		expectedArea := area
		s.AreaWithout(&expectedArea, l)
		if s.IntersectAreas(area, s.GetExclusionArea(l)) != expectedArea {
			return false
		}
	}
	return true
}

func (s *sudoku[D, A, G, S]) SetLogger(logger Logger[D]) {
	s.logger = logger
}

func (s *sudoku[D, A, G, S]) Print() error {
	count := 0
	for row := 0; row < s.Size(); row++ {
		for col := 0; col < s.Size(); col++ {
			d := s.GridCell(&s.grid, row, col)
			count += (*d).Count()
		}
	}

	size := s.Size()
	cells := size * size

	PrintGrid[D, A](s)
	fmt.Printf("%.2f%% solved\n", (1-float64(count-cells)/float64(cells*(size-1)))*100)
	fmt.Println(s.Stats())
	return nil
}

func (s *sudoku[D, A, G, S]) Stats() Stats {
	return s.stats
}

func (s *sudoku[D, A, G, S]) getRestrictions() []any {
	return s.restrictions
}
