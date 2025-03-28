package solver

import "github.com/lumaraf/sudoku-solver/sudoku"

func init() {
	// classic/general solvers
	sudoku.RegisterSolverFactory(UniqueSetSolverFactory)
	sudoku.RegisterSolverFactory(UniqueExclusionSolverFactory)
	sudoku.RegisterSolverFactory(UniqueIntersectionSolverFactory)
	//sudoku.RegisterSolverFactory(XWingSolverFactory)
	//sudoku.RegisterSolverFactory(SwordfishSolverFactory)

	// special rule solvers
	sudoku.RegisterSolverFactory(EqualSolverFactory)
	sudoku.RegisterSolverFactory(IncreaseSolverFactory)
	sudoku.RegisterSolverFactory(SandwichSolverFactory)
	sudoku.RegisterSolverFactory(AreaSumSolverFactory)
	sudoku.RegisterSolverFactory(KillerCageSolverFactory)

	// expensive solvers
	sudoku.RegisterSolverFactory(LogicChainSolverFactory)
}
