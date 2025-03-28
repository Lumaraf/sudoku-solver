package sudoku

import (
	"context"
	"fmt"
)

type Solver interface {
	Name() string
	Solve(s Sudoku) ([]Solver, error)
	AreaFilter() Area
}

type SolveProcessor interface {
	Name() string
	ProcessSolve(s Sudoku, cell CellLocation) error
}

type SolverFactory = func(restrictions []Restriction) []Solver

var defaultSolverFactories = make([]SolverFactory, 0, 10)

func RegisterSolverFactory(factory SolverFactory) {
	defaultSolverFactories = append(defaultSolverFactories, factory)
}

func (s *sudoku) IsSolved() bool {
	return s.solved.Size() == 81
}

func (s *sudoku) Solve(ctx context.Context) error {
	return s.SolveWith(ctx, defaultSolverFactories...)
}

func (s *sudoku) SolveWith(ctx context.Context, factories ...SolverFactory) error {
	_, err := s.solve(s.createSolvers(factories), ctx)
	return err
}

func (s *sudoku) createSolvers(factories []SolverFactory) []Solver {
	restrictions := make([]Restriction, 0, len(s.restrictions))
	restrictions = append(restrictions, s.restrictions...)
	for {
		hiddenRestrictions := findHiddenRestrictions(restrictions)
		if len(hiddenRestrictions) == 0 {
			break
		}
		restrictions = append(restrictions, hiddenRestrictions...)
	}

	solvers := make([]Solver, 0, len(restrictions))
	for _, factory := range factories {
		solvers = append(solvers, factory(restrictions)...)
	}
	return solvers
}

func (s *sudoku) solve(solvers []Solver, ctx context.Context) ([]Solver, error) {
	for !s.IsSolved() {
		if err := ctx.Err(); err != nil {
			return solvers, err
		}

		var err error
		solvers, err = s.runSolvers(solvers)
		if err != nil {
			return solvers, err
		}

		if s.nextChanged == (Area{}) && s.chainLimit > 0 {
			s.stats.ExclusionChainRuns++
			for limit := 1; limit <= s.chainLimit; limit++ {
				if err := s.solveExclusionChain(s.solved.Not(), limit); err != nil {
					return solvers, err
				}
				if s.nextChanged != (Area{}) {
					break
				}
			}
		}

		if s.nextChanged == (Area{}) {
			break
		}
	}
	return solvers, s.Validate()
}

func (s *sudoku) runSolvers(solvers []Solver) ([]Solver, error) {
	s.changed = s.nextChanged
	s.nextChanged = Area{}

	newSolvers := make([]Solver, 0, len(solvers)*2)
	for _, slv := range solvers {
		if slv.AreaFilter().And(s.changed) != (Area{}) {
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

func (s *sudoku) solveExclusionChain(area Area, levels int) error {
	s.logger.EnterContext(StringContext("solveExclusionChain"))
	defer s.logger.ExitContext()

	for _, cell := range area.Locations {
		d := s.Get(cell)
		errs := [9]error{}
		for v := range d.Values {
			clone := *s
			clone.logger = voidLogger{}
			clone.nextChanged = Area{}
			err := clone.Set(cell, v)
			if err == nil {
				err = clone.Validate()
			}
			if err == nil && levels > 1 {
				err = clone.solveExclusionChain(clone.nextChanged.And(clone.solved.Not()), levels-1)
			}

			if err != nil {
				errs[v-1] = err
				if removeErr := s.RemoveOption(cell, v); removeErr != nil {
					return fmt.Errorf("%+v breaks with %d(%w) and without %d(%w)", cell, v, err, v, removeErr)
				}
			}
		}
		if levels == s.chainLimit && s.nextChanged != (Area{}) {
			//fmt.Println(cell, errs)
			return nil
		}
	}
	return nil
}

var processSolveContext = StringContext("processSolve")

func (s *sudoku) processSolve(l CellLocation, mask Digits) error {
	s.logger.EnterContext(processSolveContext)
	defer s.logger.ExitContext()

	s.solved.Set(l, true)
	for _, cell := range s.exlusionAreas[l.Row][l.Col].Locations {
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
