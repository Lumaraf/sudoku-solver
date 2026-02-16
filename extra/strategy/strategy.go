package strategy

import "github.com/lumaraf/sudoku-solver/sudoku"

func AllExtraStrategies[D sudoku.Digits[D], A sudoku.Area[A]]() sudoku.StrategyFactories[D, A] {
	return sudoku.StrategyFactories[D, A]{
		// KillerCageStrategy:
		// Utilizes the sum constraints of killer sudoku cages to limit candidate placements.
		sudoku.StrategyFactoryFunc[D, A](KillerCageStrategyFactory[D, A]),

		// HiddenKillerCageStrategy:
		// Identifies hidden killer cages by analyzing the grid for areas that must sum to specific values based on existing cages.
		sudoku.StrategyFactoryFunc[D, A](HiddenKillerCageStrategyFactory[D, A]),
	}
}
