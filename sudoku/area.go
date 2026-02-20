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

	With(l CellLocation) A
	Without(l CellLocation) A

	Get(l CellLocation) bool
	Locations(yield func(int, CellLocation) bool)
	RandomLocation() CellLocation
	Size() int
	Empty() bool
	String() string
}
