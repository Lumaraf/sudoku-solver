# Sudoku Solver

This is a Sudoku solver implemented in Go, designed to be flexible and extensible.

- support for different grid sizes (e.g., 6x6, 9x9, 16x16) through generics
- a modular strategy system for implementing various solving techniques
- optional backtracking via the `Guesser` api for puzzles that have multiple solutions or can't be solved with the available strategies
- multi solver support for puzzles that overlay multiple grids (e.g., Samurai Sudoku)

## Concepts

TODO

## Areas

Areas represent groups of cells in the Sudoku grid, such as rows, columns, boxes, or custom-defined areas. They are used to apply rules and strategies that involve multiple cells. Each area type implements the `Area` interface, which provides methods for iterating over cells, checking for the presence of specific cells, changing and combining areas,

## Digits

Digits represent the possible values that can be placed in a cell. The `Digits` interface provides methods for managing candidate digits in a cell, such as adding, removing, checking for specific digits and combining sets of digits.

## Rules

### Available Rules

- **ClassicRules** - The standard Sudoku rules: each digit must appear exactly once in each row, column, and box.
- **GivenDigits** - The initial clues provided in the puzzle.
- **DiagonalRule** - For Sudoku variants with diagonal constraints, digits must also be unique along the main diagonals.
- **DisjointAreaRule** - Digits in the same location in each box must be unique. Also known as "Color Sudoku".
- **UniqueAreaRule** - Defines that a set of cells (an area) must contain unique digits.
- **KillerCageRule** - For Killer Sudoku, defines that a cage of cells must sum to a specific value without repeating digits.
- **AreaSumRule** - All cells in an area must sum to a specific value. Digits may repeat.
- **NonConsecutiveRule** - No two adjacent cells may contain consecutive digits.
- **ParityRule** - Cells must contain either only odd or only even digits.
- **AntiKingRule** - No two cells that are a king's move apart may contain the same digit.
- **AntiKnightRule** - No two cells that are a knight's move apart may contain the same digit.

### Implementing Custom Rules

TODO

## Strategies

A **strategy** in this solver is a modular technique for deducing new information about the puzzle state, such as eliminating candidates or solving cells. Strategies are registered with the solver and applied automatically.

### How to Implement a Strategy

1. **Define the Strategy Type**  
   Implement the `Strategy` interface (see `sudoku/solve.go`):
   ```go
   type Strategy[D Digits, A Area] interface {
       Name() string
       Difficulty() Difficulty
       Solve(s Sudoku[D, A]) ([]Strategy[D, A], error)
       AreaFilter() A
   }
   ```
   - `Name()`: Returns the name of the strategy.
   - `Difficulty()`: Returns a difficulty rating.
   - `Solve(s)`: Applies the strategy to the given Sudoku. Returns additional strategies to apply, if any.
   - `AreaFilter()`: Returns the area of the grid this strategy is interested in.

2. **Create a Factory**  
   Implement a factory function that returns instances of your strategy for a given puzzle. This is typically a function matching:
   ```go
   func MyStrategyFactory[D Digits, A Area](s Sudoku[D, A]) []Strategy[D, A]
   ```
   Register your factory in the `AllStrategies` function in `strategy/strategy.go`.

3. **Register the Strategy**  
   Add your factory to the list in `AllStrategies`:
   ```go
   return sudoku.StrategyFactories[D, A]{
       // ...existing strategies...
       sudoku.StrategyFactoryFunc[D, A](MyStrategyFactory[D, A]),
   }
   ```

#### Example: Unique Set Strategy

The Unique Set strategy looks for sets of N cells in a unit (row, column, or box) that together contain exactly N candidates, and removes those candidates from other cells in the unit.

See `strategy/unqiue.set.go` for a concrete implementation.

---

## Relevant Parts of the Sudoku API

The solver is highly generic and type-safe, using Go generics for digits and area representations. Here are the most important interfaces and types:

### Core Methods of the Sudoku Interface

- `Get(l CellLocation) D`: Get the candidates at a cell.
- `Set(l CellLocation, v int) error`: Set a value at a cell.
- `Mask(l, d) error`: Restrict candidates at a cell.
- `RemoveOption(l, v) error`: Remove a candidate from a cell.
- `Row/Column/Box(int) A`: Get an area (row, column, or box).
- `SolvedArea() A`, `ChangedArea() A`, `NextChangedArea() A`: Track progress.
- `Try(func(Sudoku[D, A]) error) error`: Clone and test a hypothetical change.
- `Validate() error`: Check puzzle validity.
- `Print() error`: Print the current state.
  
- **StrategyFactory** (`sudoku/solve.go`):  
  Used to create strategies for a puzzle.

### Supporting Types

- **CellLocation**:  
  Struct with `Row` and `Col` fields, identifies a cell.

- **Logger** (`sudoku/logger.go`):  
  For debugging and tracing strategy application.

- **Stats**:  
  Tracks statistics like cell updates and solver runs.

### Example Usage

To implement a strategy, you typically:
- Iterate over areas (rows, columns, boxes) using the `Sudoku` interface.
- Query and update candidates using `Get`, `Set`, `Mask`, and `RemoveOption`.
- Use the `Area` and `Digits` interfaces to manipulate sets of cells and candidates.
- Use the digits and area methods above for set operations on digits and areas.

---
