package sudoku

import (
	"math"
	"math/bits"
	"math/rand"
)

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

type area128[AS interface {
	gridSize() int
	allCells() [2]uint64
}] [2]uint64

func (a area128[AS]) And(b area128[AS]) area128[AS] {
	return area128[AS]{
		a[0] & b[0],
		a[1] & b[1],
	}
}
func (a area128[AS]) Or(b area128[AS]) area128[AS] {
	return area128[AS]{
		a[0] | b[0],
		a[1] | b[1],
	}
}
func (a area128[AS]) Not() area128[AS] {
	var spec AS
	all := spec.allCells()
	return area128[AS]{
		^a[0] & all[0],
		^a[1] & all[1],
	}
}

func (a area128[AS]) All() area128[AS] {
	var spec AS
	return spec.allCells()
}

func (a area128[AS]) With(l CellLocation) area128[AS] {
	index, mask := a.getMask(l)
	a[index] = a[index] | mask
	return a
}

func (a area128[AS]) Without(l CellLocation) area128[AS] {
	index, mask := a.getMask(l)
	a[index] = a[index] & ^mask
	return a
}

func (a area128[AS]) getMask(l CellLocation) (int, uint64) {
	var spec AS
	idx := l.Row*spec.gridSize() + l.Col
	return idx / 64, 1 << (uint64(idx) % 64)
}

func (a area128[AS]) Get(l CellLocation) bool {
	index, mask := a.getMask(l)
	return a[index]&mask != 0
}

func (a area128[AS]) Set(l CellLocation, v bool) area128[AS] {
	if v {
		return a.With(l)
	}
	return a.Without(l)
}

func (a area128[AS]) Locations(yield func(int, CellLocation) bool) {
	var spec AS
	index := 0
	size := spec.gridSize()
	for idx, b := range a {
		for b != 0 {
			lz := bits.TrailingZeros64(b)
			b = b & ^(1 << lz)
			pos := idx*64 + lz
			if !yield(index, CellLocation{pos / size, pos % size}) {
				return
			}
			index++
		}
	}
}

func (a area128[AS]) RandomLocation() CellLocation {
	var spec AS
	size := spec.gridSize()
	return a.nextCell(rand.Intn(size * size))
}

func (a area128[AS]) nextCell(index int) CellLocation {
	var maskedArea area128[AS]
	if index < 64 {
		maskedArea[0] = a[0] & ^(math.MaxUint64 >> (64 - index))
		maskedArea[1] = a[1]
	} else {
		maskedArea[0] = 0
		maskedArea[1] = a[1] & ^(math.MaxUint64 >> (128 - index))
	}
	for _, cell := range maskedArea.Locations {
		return cell
	}
	for _, cell := range a.Locations {
		return cell
	}
	return CellLocation{}
}

func (a area128[AS]) Size() int {
	return bits.OnesCount64(a[0]) + bits.OnesCount64(a[1])
}

func (a area128[AS]) Empty() bool {
	return a[0] == 0 && a[1] == 0
}

func (a area128[AS]) String() string {
	var spec AS
	size := spec.gridSize()
	grid := make([]rune, 0, size*size+size)
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			if a.Get(CellLocation{row, col}) {
				grid = append(grid, 'X')
			} else {
				grid = append(grid, ' ')
			}
		}
		grid = append(grid, '\n')
	}
	return string(grid)
}
