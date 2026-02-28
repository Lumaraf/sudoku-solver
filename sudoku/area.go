package sudoku

type AreaOps[A Area[A]] interface {
	NewArea(locs ...CellLocation) A
	NewAreaFromOffsets(center CellLocation, o Offsets) A
}

type Area[A Area[A]] interface {
	comparable

	And(other A) A
	Or(other A) A
	Not() A

	All() A
	Size() int

	With(l CellLocation) A
	Without(l CellLocation) A

	ShiftLeft(n int) A
	ShiftRight(n int) A
	ShiftUp(n int) A
	ShiftDown(n int) A
	ShiftBy(offset Offset) A

	Get(l CellLocation) bool
	Locations(yield func(int, CellLocation) bool)
	RandomLocation() CellLocation
	Count() int
	Empty() bool
	String() string
}
