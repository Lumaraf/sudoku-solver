# sudoku solver

## Implementing a Strategy

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

#### Methods for Working with Digits (Candidates)

- `NewDigits(values ...int) D`: Create a digit set from a list of values.
- `AllDigits() D`: Get a digit set containing all possible digits.
- `IntersectDigits(d1 D, d2 D) D`: Intersection of two digit sets.
- `UnionDigits(d1 D, d2 D) D`: Union of two digit sets.
- `InvertDigits(d D) D`: Invert a digit set (all digits not present in `d`).

#### Methods for Working with Areas (Sets of Cells)

- `NewArea(locs ...CellLocation) A`: Create an area from a list of cell locations.
- `NewAreaFromOffsets(center CellLocation, o Offsets) A`: Create an area from a center cell and a set of offsets.
- `AreaWith(a *A, l CellLocation)`: Add a cell to an area.
- `AreaWithout(a *A, l CellLocation)`: Remove a cell from an area.
- `IntersectAreas(a1 A, a2 A) A`: Intersection of two areas.
- `UnionAreas(a1 A, a2 A) A`: Union of two areas.
- `InvertArea(a A) A`: Invert an area (all cells not present in `a`).

- **Digits** (`sudoku/digits.go`):  
  Represents possible values in a cell. Key methods:
  - `CanContain(v int) bool`
  - `Empty() bool`
  - `Count() int`
  - `Single() (int, bool)`
  - `Values(func(int) bool)`

- **Area** (`sudoku/area.go`):  
  Represents a set of cells. Key methods:
  - `Get(l CellLocation) bool`
  - `Locations(func(int, CellLocation) bool)`
  - `Size() int`
  - `Empty() bool`

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

This documentation should help you and contributors understand how to add new strategies and interact with the Sudoku API. If you need a code example or more details on a specific part, see the provided strategy files or ask for further clarification.
