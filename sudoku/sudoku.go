package sudoku

import (
	"context"
	"errors"
	"fmt"
)

var mutipleSolutionsError = errors.New("multiple solutions")

type Sudoku interface {
	// SetChainLimit sets the limit for the chain length in solving algorithms.
	SetChainLimit(limit int)

	// Get retrieves the digits at the specified cell location.
	Get(l CellLocation) Digits

	// Set assigns a value to the specified cell location.
	Set(l CellLocation, v int) error

	// Try attempts to apply a function to a clone of the Sudoku puzzle.
	Try(f func(s Sudoku) error) error

	// Mask restricts the possible digits at the specified cell location.
	Mask(l CellLocation, d Digits) error

	// RemoveMask removes the restriction of possible digits at the specified cell location.
	RemoveMask(l CellLocation, d Digits) error

	// RemoveOption removes a specific digit option from the specified cell location.
	RemoveOption(l CellLocation, v int) error

	// SolvedArea returns the area of the grid that is solved.
	SolvedArea() Area

	// ChangedArea returns the area of the grid that has changed.
	ChangedArea() Area

	// NextChangedArea returns the area of the grid that will change next.
	NextChangedArea() Area

	// GetExclusionArea returns the exclusion area for the specified cell location.
	GetExclusionArea(l CellLocation) Area

	// SetLogger sets the logger for the Sudoku puzzle.
	SetLogger(logger Logger)

	// Validate checks the validity of the current Sudoku puzzle state.
	Validate() error

	// Print prints the current state of the Sudoku puzzle.
	Print() error

	// Stats returns the statistics of the Sudoku puzzle.
	Stats() Stats

	// AddRestriction adds a restriction to the Sudoku puzzle.
	AddRestriction(r Restriction)

	// Solve attempts to solve the Sudoku puzzle.
	Solve(ctx context.Context) error

	// SolveWith attempts to solve the Sudoku puzzle using the provided solver factories.
	SolveWith(ctx context.Context, factories ...SolverFactory) error

	// GuessSolutions generates possible solutions by guessing.
	GuessSolutions(ctx context.Context, g Guesser) func(func(Sudoku) bool)

	// IsSolved checks if the Sudoku puzzle is solved.
	IsSolved() bool
}

type CellLocation struct {
	Row int
	Col int
}

func (l CellLocation) Box() int {
	return l.Row/3*3 + l.Col/3
}

type sudoku struct {
	grid            [9][9]Digits
	exlusionAreas   [9][9]Area
	restrictions    []Restriction
	solveProcessors []SolveProcessor
	findAllOptions  bool
	chainLimit      int
	changed         Area
	nextChanged     Area
	solved          Area
	stats           Stats
	logger          Logger
}

type Stats struct {
	CellUpdates        int
	SolverRuns         int
	SolverHits         int
	ExclusionChainRuns int
	GuesserRuns        int
	GuessMisses        int
}

func NewSudoku() Sudoku {
	//switch size {
	//case 6:
	//	return newSudoku[grid6, area6]()
	//case 9:
	//	return newSudoku[grid9, area9]()
	//}
	return newSudoku()
}

func newSudoku() *sudoku {
	s := sudoku{
		chainLimit:  2,
		nextChanged: Area{}.Not(),
		logger:      voidLogger{},
	}
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			s.grid[row][col] = AllDigits
		}
	}
	return &s
}

func (s *sudoku) SetChainLimit(limit int) {
	s.chainLimit = limit
}

func (s *sudoku) Get(l CellLocation) Digits {
	return s.grid[l.Row][l.Col]
}

func (s *sudoku) Set(l CellLocation, v int) error {
	if err := checkValue(v); err != nil {
		return err
	}
	cell := &s.grid[l.Row][l.Col]
	if !cell.CanContain(v) {
		return errors.New("cell doesn't allow value")
	}
	if cell.Count() > 1 {
		s.nextChanged.Set(l, true)
		s.stats.CellUpdates++
		oldCell := *cell
		cell.ForceValue(v)
		s.logger.UpdateCell(l, oldCell, *cell)
		return s.processSolve(l, *cell)
	}
	return nil
}

func (s *sudoku) Try(f func(s Sudoku) error) error {
	clone := *s
	clone.logger = voidLogger{}
	clone.nextChanged = Area{}
	return f(&clone)
}

func (s *sudoku) Mask(l CellLocation, d Digits) error {
	target := &s.grid[l.Row][l.Col]
	if *target&^d != 0 {
		s.nextChanged.Set(l, true)
		s.stats.CellUpdates++
		oldDigits := *target
		newDigits := *target & d
		if newDigits == 0 {
			return ErrEmptyCell(l)
		}
		*target = newDigits
		s.logger.UpdateCell(l, oldDigits, newDigits)
		if target.Count() == 1 {
			return s.processSolve(l, newDigits)
		}
	}
	return nil
}

func (s *sudoku) RemoveMask(l CellLocation, d Digits) error {
	target := &s.grid[l.Row][l.Col]
	if *target&d != 0 {
		s.nextChanged.Set(l, true)
		s.stats.CellUpdates++
		oldDigits := *target
		newDigits := (*target & ^d) & AllDigits
		if newDigits == 0 {
			s.solved.Set(l, false)
			return ErrEmptyCell(l)
		}
		*target = newDigits
		s.logger.UpdateCell(l, oldDigits, newDigits)
		if target.Count() == 1 {
			return s.processSolve(l, *target)
		}
	}
	return nil
}

func (s *sudoku) RemoveOption(l CellLocation, v int) error {
	if err := checkValue(v); err != nil {
		return err
	}
	mask := Digits(0)
	mask.AddOption(v)
	return s.RemoveMask(l, mask)
}

func (s *sudoku) RequireValueInArea(v int, a Area) error {
	overlap := Area{}.Not()
	for _, l := range a.Locations {
		overlap = overlap.And(s.GetExclusionArea(l))
	}
	for _, l := range overlap.Locations {
		if err := s.RemoveOption(l, v); err != nil {
			return err
		}
	}
	return nil
}

func (s *sudoku) SolvedArea() Area {
	return s.solved
}

func (s *sudoku) ChangedArea() Area {
	return s.changed
}

func (s *sudoku) NextChangedArea() Area {
	return s.nextChanged
}

func (s *sudoku) GetExclusionArea(l CellLocation) Area {
	return s.exlusionAreas[l.Row][l.Col]
}

func (s *sudoku) SetLogger(logger Logger) {
	s.logger = logger
}

func (s *sudoku) Validate() error {
	for _, r := range s.restrictions {
		if err := r.Validate(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *sudoku) Print() error {
	//lines := [][]rune{
	//	[]rune("╔═══════╤═══════╤═══════╦═══════╤═══════╤═══════╦═══════╤═══════╤═══════╗"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("╠═══════╪═══════╪═══════╬═══════╪═══════╪═══════╬═══════╪═══════╪═══════╣"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("╠═══════╪═══════╪═══════╬═══════╪═══════╪═══════╬═══════╪═══════╪═══════╣"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
	//	[]rune("╚═══════╧═══════╧═══════╩═══════╧═══════╧═══════╩═══════╧═══════╧═══════╝"),
	//}
	//
	//count := 0
	//for row := 0; row < 9; row++ {
	//	for col := 0; col < 9; col++ {
	//		top := 1 + row*4
	//		left := 2 + col*8
	//		d := s.grid[row][col]
	//		for v := range d.Values {
	//			lines[top+(v-1)/3][left+((v-1)%3)*2] = rune(v + '0')
	//			count++
	//		}
	//	}
	//}
	//
	//for _, l := range lines {
	//	fmt.Println(string(l))
	//}

	count := 0
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			d := s.grid[row][col]
			count += d.Count()
		}
	}

	PrintGrid(s)
	fmt.Printf("%.2f%% solved\n", (1-float64(count-81)/(81*8))*100)
	return nil
}

func (s *sudoku) Stats() Stats {
	return s.stats
}
