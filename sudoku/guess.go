package sudoku

import (
	"context"
	"errors"
)

type GuessSelector[D Digits[D], A Area[A]] func(s Sudoku[D, A]) (CellLocation, Values)

func DefaultGuessSelector[D Digits[D], A Area[A]](s Sudoku[D, A]) (CellLocation, Values) {
	bestCell := CellLocation{}
	bestDigits := s.NewDigits()
	for _, cell := range s.SolvedArea().Not().Locations {
		d := s.Get(cell)
		if d.Count() == 1 {
			continue
		}
		if bestDigits.Count() == 0 || d.Count() < bestDigits.Count() {
			bestCell = cell
			bestDigits = d
		}
	}
	return bestCell, bestDigits.Values
}

type Guesser[D Digits[D], A Area[A]] interface {
	Solver[D, A]
	Guess(g GuessSelector[D, A], ctx context.Context) func(func(Sudoku[D, A]) bool)
}

type guesser[D Digits[D], A Area[A], G comparable, S size[D, A, G]] struct {
	*solver[D, A, G, S]
}

func (g *guesser[D, A, G, S]) Guess(gs GuessSelector[D, A], ctx context.Context) func(func(Sudoku[D, A]) bool) {
	return func(yield func(Sudoku[D, A]) bool) {
		strategies := g.createStrategies()

		var err error
		strategies, err = g.solve(g.sudoku, strategies, ctx)
		if err != nil {
			return
		}
		if g.sudoku.IsSolved() {
			yield(g.sudoku)
			return
		}

		if gs == nil {
			gs = DefaultGuessSelector
		}

		for s := range g.guessSolutions(g.sudoku, strategies, gs, make([]CellLocation, 0, g.sudoku.Size()*g.sudoku.Size()), ctx) {
			if !yield(&s) {
				return
			}
		}
	}
}

func (g *guesser[D, A, G, S]) guessSolutions(s *sudoku[D, A, G, S], solvers []Strategy[D, A], gs GuessSelector[D, A], path []CellLocation, ctx context.Context) func(yield func(sudoku[D, A, G, S]) bool) {
	g.sudoku.stats.GuesserRuns++
	return func(yield func(sudoku[D, A, G, S]) bool) {
		baseClone := *s
		baseClone.logger = voidLogger[D]{}

		cell, values := gs(&baseClone)

		guessPath := append(path, cell)
		for v := range values {
			clone := baseClone

			if clone.Set(cell, v) != nil {
				continue
			}

			if clone.Validate() != nil {
				continue
			}

			nextSolvers, err := g.solve(&clone, solvers, ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return
				}
				continue
			}

			if clone.IsSolved() {
				if !yield(clone) {
					return
				}
			} else {
				solutionCount := 0
				for solution := range g.guessSolutions(&clone, nextSolvers, gs, guessPath, ctx) {
					solutionCount++
					if !yield(solution) {
						return
					}
				}
				if solutionCount == 0 {
					s.stats.GuessMisses++
					_ = s.RemoveOption(cell, v)
				}
			}
		}
	}
}

func (s *sudoku[D, A, G, S]) NewGuesser() Guesser[D, A] {
	return &guesser[D, A, G, S]{
		solver: &solver[D, A, G, S]{
			sudoku: s,
		},
	}
}
