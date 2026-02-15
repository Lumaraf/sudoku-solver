package sudoku

import (
	"context"
	"fmt"
	"sort"
)

type Difficulty uint

const (
	DIFFICULTY_EASY Difficulty = iota
	DIFFICULTY_NORMAL
	DIFFICULTY_HARD
	DIFFICULTY_IMPOSSIBLE
)

type Strategy[D Digits[D], A Area[A]] interface {
	Name() string
	Difficulty() Difficulty
	Solve(s Sudoku[D, A]) ([]Strategy[D, A], error)
	AreaFilter() A
}

type Strategies[D Digits[D], A Area[A]] []Strategy[D, A]

func (s Strategies[D, A]) Len() int {
	return len(s)
}

func (s Strategies[D, A]) Less(i, j int) bool {
	return s[i].Difficulty() < s[j].Difficulty()
}

func (s Strategies[D, A]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type StrategyFactory[D Digits[D], A Area[A]] interface {
	For(s Sudoku[D, A]) []Strategy[D, A]
}

type StrategyFactoryFunc[D Digits[D], A Area[A]] func(s Sudoku[D, A]) []Strategy[D, A]

func (sff StrategyFactoryFunc[D, A]) For(s Sudoku[D, A]) []Strategy[D, A] {
	return sff(s)
}

type StrategyFactories[D Digits[D], A Area[A]] []StrategyFactory[D, A]

func (sf StrategyFactories[D, A]) For(s Sudoku[D, A]) []Strategy[D, A] {
	strategies := make([]Strategy[D, A], 0, len(sf))
	for _, factory := range sf {
		strategies = append(strategies, factory.For(s)...)
	}
	return strategies
}

type Solver[D Digits[D], A Area[A]] interface {
	SetChainLimit(limit int)
	Use(factories ...StrategyFactory[D, A])
	Solve(ctx context.Context) error
}

type solver[D Digits[D], A Area[A], G comparable, S size[D, A, G]] struct {
	sudoku            *sudoku[D, A, G, S]
	chainLimit        int
	strategyFactories []StrategyFactory[D, A]
}

func (slv *solver[D, A, G, S]) SetChainLimit(limit int) {
	slv.chainLimit = limit
}

func (slv *solver[D, A, G, S]) Use(factories ...StrategyFactory[D, A]) {
	slv.strategyFactories = append(slv.strategyFactories, factories...)
}

func (slv *solver[D, A, G, S]) Solve(ctx context.Context) error {
	_, err := slv.solve(slv.sudoku, slv.createStrategies(), ctx)
	return err
}

func (slv *solver[D, A, G, S]) createStrategies() Strategies[D, A] {
	strategies := make(Strategies[D, A], 0, len(slv.strategyFactories))
	for _, factory := range slv.strategyFactories {
		strategies = append(strategies, factory.For(slv.sudoku)...)
	}
	sort.Stable(strategies)
	return strategies
}

func (slv *solver[D, A, G, S]) solve(s *sudoku[D, A, G, S], solvers []Strategy[D, A], ctx context.Context) ([]Strategy[D, A], error) {
	for !s.IsSolved() {
		if err := ctx.Err(); err != nil {
			return solvers, err
		}

		var err error
		solvers, err = slv.runSolvers(s, solvers)
		if err != nil {
			return solvers, err
		}

		if s.nextChanged.Empty() && slv.chainLimit > 0 {
			s.stats.ExclusionChainRuns++
			for limit := 1; limit <= slv.chainLimit; limit++ {
				if err := slv.solveExclusionChain(s, s.solved.Not(), limit); err != nil {
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

func (slv *solver[D, A, G, S]) runSolvers(s *sudoku[D, A, G, S], strategies []Strategy[D, A]) ([]Strategy[D, A], error) {
	s.changed = s.nextChanged
	s.nextChanged = *new(A)

	var lastDifficulty Difficulty
	newSolvers := make([]Strategy[D, A], 0, len(strategies)*2)
	for n, strategy := range strategies {
		if strategy.Difficulty() > lastDifficulty && !s.nextChanged.Empty() {
			newSolvers = append(newSolvers, strategies[n:]...)
			break
		}
		lastDifficulty = strategy.Difficulty()

		if !strategy.AreaFilter().And(s.changed).Empty() {
			s.stats.SolverRuns++
			cellUpdatesBefore := s.stats.CellUpdates
			s.logger.EnterContext(strategy)
			s2, err := strategy.Solve(s)
			s.logger.ExitContext()
			if err != nil {
				return nil, err
			}
			if s.stats.CellUpdates > cellUpdatesBefore {
				s.stats.SolverHits++
			}
			newSolvers = append(newSolvers, s2...)
		} else {
			newSolvers = append(newSolvers, strategy)
		}
	}
	return newSolvers, nil
}

type ExclusionChainError struct {
	cell   CellLocation
	errors [9]error
}

func (slv *solver[D, A, G, S]) solveExclusionChain(s *sudoku[D, A, G, S], area A, levels int) error {
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
				err = slv.solveExclusionChain(&clone, clone.nextChanged.And(clone.solved.Not()), levels-1)
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

func (s *sudoku[D, A, G, S]) setSolved(l CellLocation) {
	s.solved = s.solved.With(l)
}

func (s *sudoku[D, A, G, S]) IsSolved() bool {
	return s.solved.Size() == s.Size()*s.Size()
}

func (s *sudoku[D, A, G, S]) NewSolver() Solver[D, A] {
	return &solver[D, A, G, S]{sudoku: s}
}

var processChangeContext = StringContext("processChange")

func (s *sudoku[D, A, G, S]) processChange(l CellLocation, mask D) error {
	s.logger.EnterContext(processChangeContext)
	defer s.logger.ExitContext()

	for _, cp := range s.changeProcessors {
		s.logger.EnterContext(cp)
		if err := cp.ProcessChange(s, l, mask); err != nil {
			return err
		}
		s.logger.ExitContext()
	}
	return nil
}
