package sudoku

import (
	"context"
	"fmt"
)

type Restriction[D Digits, A Area] interface {
	Name() string
}

type Strategy[D Digits, A Area] interface {
	Name() string
	Solve(s Sudoku[D, A]) ([]Strategy[D, A], error)
	AreaFilter() A
}

type SolveProcessor[D Digits, A Area] interface {
	Name() string
	ProcessSolve(s Sudoku[D, A], cell CellLocation) error
}

type StrategyFactory[D Digits, A Area] = func(s Sudoku[D, A]) []Strategy[D, A]

func (s *sudoku[D, A, G, S]) IsSolved() bool {
	return s.solved.Size() == 81
}

func (s *sudoku[D, A, G, S]) SolveWith(ctx context.Context, factories ...StrategyFactory[D, A]) error {
	_, err := s.solve(s.createSolvers(factories), ctx)
	return err
}

func (s *sudoku[D, A, G, S]) createSolvers(factories []StrategyFactory[D, A]) []Strategy[D, A] {
	restrictions := make([]Restriction[D, A], 0, len(s.restrictions))
	restrictions = append(restrictions, s.restrictions...)

	solvers := make([]Strategy[D, A], 0, len(restrictions))
	for _, factory := range factories {
		solvers = append(solvers, factory(s)...)
	}
	return solvers
}

func (s *sudoku[D, A, G, S]) solve(solvers []Strategy[D, A], ctx context.Context) ([]Strategy[D, A], error) {
	for !s.IsSolved() {
		if err := ctx.Err(); err != nil {
			return solvers, err
		}

		var err error
		solvers, err = s.runSolvers(solvers)
		if err != nil {
			return solvers, err
		}

		if s.nextChanged.Empty() && s.chainLimit > 0 {
			s.stats.ExclusionChainRuns++
			for limit := 1; limit <= s.chainLimit; limit++ {
				if err := s.solveExclusionChain(s.InvertArea(s.solved), limit); err != nil {
					return solvers, err
				}
				if !s.nextChanged.Empty() {
					break
				}
			}
		}

		if s.nextChanged.Empty() {
			break
		}
	}
	return solvers, s.Validate()
}

func (s *sudoku[D, A, G, S]) runSolvers(solvers []Strategy[D, A]) ([]Strategy[D, A], error) {
	s.changed = s.nextChanged
	s.nextChanged = *new(A)

	newSolvers := make([]Strategy[D, A], 0, len(solvers)*2)
	for _, slv := range solvers {
		if !s.IntersectAreas(slv.AreaFilter(), s.changed).Empty() {
			s.stats.SolverRuns++
			cellUpdatesBefore := s.stats.CellUpdates
			s.logger.EnterContext(slv)
			s2, err := slv.Solve(s)
			s.logger.ExitContext()
			if err != nil {
				return nil, err
			}
			if s.stats.CellUpdates > cellUpdatesBefore {
				s.stats.SolverHits++
			}
			newSolvers = append(newSolvers, s2...)
		} else {
			newSolvers = append(newSolvers, slv)
		}
	}
	return newSolvers, nil
}

type ExclusionChainError struct {
	cell   CellLocation
	errors [9]error
}

func (s *sudoku[D, A, G, S]) solveExclusionChain(area Area, levels int) error {
	s.logger.EnterContext(StringContext("solveExclusionChain"))
	defer s.logger.ExitContext()

	for _, cell := range area.Locations {
		d := s.Get(cell)
		errs := [9]error{}
		for v := range d.Values {
			clone := *s
			clone.logger = voidLogger[D]{}
			clone.nextChanged = *new(A)
			err := clone.Set(cell, v)
			if err == nil {
				err = clone.Validate()
			}
			if err == nil && levels > 1 {
				err = clone.solveExclusionChain(s.IntersectAreas(clone.nextChanged, s.InvertArea(clone.solved)), levels-1)
			}

			if err != nil {
				errs[v-1] = err
				if removeErr := s.RemoveOption(cell, v); removeErr != nil {
					return fmt.Errorf("%+v breaks with %d(%w) and without %d(%w)", cell, v, err, v, removeErr)
				}
			}
		}
		if levels == s.chainLimit && !s.nextChanged.Empty() {
			//fmt.Println(cell, errs)
			return nil
		}
	}
	return nil
}

var processSolveContext = StringContext("processSolve")

func (s *sudoku[D, A, G, S]) processSolve(l CellLocation, mask D) error {
	s.logger.EnterContext(processSolveContext)
	defer s.logger.ExitContext()

	s.AreaWith(&s.solved, l)
	for _, cell := range s.exclusionAreas[l.Row][l.Col].Locations {
		if err := s.RemoveMask(cell, mask); err != nil {
			return err
		}
	}
	for _, sp := range s.solveProcessors {
		s.logger.EnterContext(sp)
		if err := sp.ProcessSolve(s, l); err != nil {
			return err
		}
		s.logger.ExitContext()
	}
	return nil
}
