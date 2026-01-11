package strategy

import "github.com/lumaraf/sudoku-solver/sudoku"

// AllStrategies returns all available strategies used by the solver.
// Each strategy is documented directly above its registration for clarity.
func AllStrategies[D sudoku.Digits, A sudoku.Area]() sudoku.StrategyFactories[D, A] {
	return sudoku.StrategyFactories[D, A]{
		// UniqueSetStrategy:
		// Detects sets of cells within a unit (row, column or box) that contain exactly N candidates among N cells.
		// Removes these candidates from all other cells in the same unit. This is commonly known as the "naked set" technique.
		sudoku.StrategyFactoryFunc[D, A](UniqueSetStrategyFactory[D, A]),

		// UniqueIntersectionStrategy:
		// Identifies intersections between units (e.g. row and box) where candidates are restricted to a shared subset of cells.
		// Eliminates these candidates from other cells in the intersecting unit. This is also called "pointing pairs/triples" or "box-line reduction".
		sudoku.StrategyFactoryFunc[D, A](UniqueIntersectionStrategyFactory[D, A]),

		// LogicChainStrategy:
		// Uses chains of logical implications to deduce eliminations. It simulates placing a candidate and follows the consequences,
		// ruling out candidates that would lead to contradictions. This covers techniques like "simple coloring" and "forcing chains".
		sudoku.StrategyFactoryFunc[D, A](LogicChainStrategyFactory[D, A]),

		// XWingStrategy:
		// Searches for the X-Wing pattern: a candidate appears exactly twice in two different rows and the same columns (or vice versa).
		// This allows elimination of the candidate from other cells in those columns/rows.
		sudoku.StrategyFactoryFunc[D, A](XWingStrategyFactory[D, A]),

		// UniqueExclusionStrategy:
		// Examines all possible placements of a candidate in a unit and excludes candidates that cannot appear in any valid solution.
		// This is related to "hidden singles" and advanced exclusion logic.
		sudoku.StrategyFactoryFunc[D, A](UniqueExclusionStrategyFactory[D, A]),

		// PatternOverlayStrategy:
		// Finds all possible placement patterns for every digit and overlays them to eliminate impossible options.
		sudoku.StrategyFactoryFunc[D, A](PatternOverlayStrategyFactory[D, A]),

		// KillerCageStrategy:
		// Utilizes the sum constraints of killer sudoku cages to limit candidate placements.
		sudoku.StrategyFactoryFunc[D, A](KillerCageStrategyFactory[D, A]),

		// HiddenKillerCageStrategy:
		// Identifies hidden killer cages by analyzing the grid for areas that must sum to specific values based on existing cages.
		sudoku.StrategyFactoryFunc[D, A](HiddenKillerCageStrategyFactory[D, A]),
	}
}
