package strategy

import (
	//strategy2 "github.com/lumaraf/sudoku-solver/extra/strategy"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func init() {
	// classic/general solvers
	sudoku.RegisterStrategyFactory(UniqueSetSolverFactory)
	sudoku.RegisterStrategyFactory(UniqueExclusionSolverFactory)
	sudoku.RegisterStrategyFactory(UniqueIntersectionSolverFactory)
	//sudoku.RegisterStrategyFactory(XWingSolverFactory)
	//sudoku.RegisterStrategyFactory(SwordfishSolverFactory)

	// special rule solvers
	//sudoku.RegisterStrategyFactory(strategy2.EqualSolverFactory)
	//sudoku.RegisterStrategyFactory(strategy2.IncreaseSolverFactory)
	//sudoku.RegisterStrategyFactory(strategy2.SandwichSolverFactory)
	//sudoku.RegisterStrategyFactory(strategy2.AreaSumSolverFactory)
	//sudoku.RegisterStrategyFactory(strategy2.KillerCageSolverFactory)

	// expensive solvers
	sudoku.RegisterStrategyFactory(LogicChainSolverFactory)
}
