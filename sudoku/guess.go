package sudoku

//type Guesser func(s Sudoku) (CellLocation, Values)
//
//func DefaultGuesser(s Sudoku) (CellLocation, Values) {
//	bestCell := CellLocation{}
//	bestDigits := Digits(0)
//	for _, cell := range s.SolvedArea().Not().Locations {
//		d := s.Get(cell)
//		if d.Count() == 1 {
//			continue
//		}
//		if bestDigits.Count() == 0 || d.Count() < bestDigits.Count() {
//			bestCell = cell
//			bestDigits = d
//		}
//	}
//	return bestCell, bestDigits.Values
//}

//func (s *sudoku) GuessSolutions(ctx context.Context, g Guesser) func(func(Sudoku) bool) {
//	return s.GuessSolutionsWith(ctx, g, defaultSolverFactories...)
//}
//
//func (s *sudoku) GuessSolutionsWith(ctx context.Context, guesser Guesser, factories ...StrategyFactory) func(func(Sudoku) bool) {
//	return func(yield func(Sudoku) bool) {
//		solvers := s.createSolvers(factories)
//
//		var err error
//		solvers, err = s.solve(solvers, ctx)
//		if err != nil {
//			return
//		}
//		if s.IsSolved() {
//			yield(s)
//			return
//		}
//
//		if guesser == nil {
//			guesser = DefaultGuesser
//		}
//
//		for s := range s.guessSolutions(solvers, guesser, make([]CellLocation, 0, 81), ctx) {
//			if !yield(&s) {
//				return
//			}
//		}
//	}
//}
//
//func (s *sudoku) guessSolutions(solvers []Strategy, guesser Guesser, path []CellLocation, ctx context.Context) func(yield func(sudoku) bool) {
//	s.stats.GuesserRuns++
//	return func(yield func(sudoku) bool) {
//		baseClone := *s
//		baseClone.logger = voidLogger{}
//
//		cell, values := guesser(&baseClone)
//
//		guessPath := append(path, cell)
//		for v := range values {
//			clone := baseClone
//
//			if clone.Set(cell, v) != nil {
//				continue
//			}
//
//			if clone.Validate() != nil {
//				continue
//			}
//
//			nextSolvers, err := clone.solve(solvers, ctx)
//			if err != nil {
//				if errors.Is(err, context.DeadlineExceeded) {
//					return
//				}
//				continue
//			}
//
//			if clone.IsSolved() {
//				if !yield(clone) {
//					return
//				}
//			} else {
//				solutionCount := 0
//				for solution := range clone.guessSolutions(nextSolvers, guesser, guessPath, ctx) {
//					solutionCount++
//					if !yield(solution) {
//						return
//					}
//				}
//				if solutionCount == 0 {
//					s.stats.GuessMisses++
//					_ = s.RemoveOption(cell, v)
//				}
//			}
//		}
//	}
//}
