package sudoku

import (
	"math"
	"math/bits"
	"math/rand"
)

type Area interface {
	comparable

	Get(l CellLocation) bool
	Locations(yield func(int, CellLocation) bool)
	RandomLocation() CellLocation
	Size() int
	Empty() bool
	String() string
}

type gridSize interface {
	gridSize() int
}

type area128[GS gridSize] struct {
	gs   GS
	bits [2]uint64
}

func (a area128[GS]) Get(l CellLocation) bool {
	index, mask := a.getMask(l)
	return a.bits[index]&mask != 0
}

func (a area128[GS]) and(b area128[GS]) area128[GS] {
	return area128[GS]{
		a.gs,
		[2]uint64{
			a.bits[0] & b.bits[0],
			a.bits[1] & b.bits[1],
		},
	}
}
func (a area128[GS]) or(b area128[GS]) area128[GS] {
	return area128[GS]{
		a.gs,
		[2]uint64{
			a.bits[0] | b.bits[0],
			a.bits[1] | b.bits[1],
		},
	}
}
func (a area128[GS]) not() area128[GS] {
	return area128[GS]{
		a.gs,
		[2]uint64{
			^a.bits[0],
			^a.bits[1],
		},
	}
}

func (a area128[GS]) with(l CellLocation, v bool) area128[GS] {
	index, mask := a.getMask(l)
	if v {
		a.bits[index] = a.bits[index] | mask
	} else {
		a.bits[index] = a.bits[index] & ^mask
	}
	return a
}

func (a area128[GS]) getMask(l CellLocation) (int, uint64) {
	idx := l.Row*a.gs.gridSize() + l.Col
	return idx / 64, 1 << (uint64(idx) % 64)
}

func (a area128[GS]) Locations(yield func(int, CellLocation) bool) {
	index := 0
	size := a.gs.gridSize()
	for idx, b := range a.bits {
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

func (a area128[GS]) RandomLocation() CellLocation {
	size := a.gs.gridSize()
	return a.nextCell(rand.Intn(size * size))
}

func (a area128[GS]) nextCell(index int) CellLocation {
	var maskedArea area128[GS]
	if index < 64 {
		maskedArea.bits[0] = a.bits[0] & ^(math.MaxUint64 >> (64 - index))
		maskedArea.bits[1] = a.bits[1]
	} else {
		maskedArea.bits[0] = 0
		maskedArea.bits[1] = a.bits[1] & ^(math.MaxUint64 >> (128 - index))
	}
	for _, cell := range maskedArea.Locations {
		return cell
	}
	for _, cell := range a.Locations {
		return cell
	}
	return CellLocation{}
}

func (a area128[GS]) Size() int {
	return bits.OnesCount64(a.bits[0]) + bits.OnesCount64(a.bits[1])
}

func (a area128[GS]) Empty() bool {
	return a.bits[0] == 0 && a.bits[1] == 0
}

func (a area128[GS]) String() string {
	size := a.gs.gridSize()
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
