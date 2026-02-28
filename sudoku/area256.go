package sudoku

import (
	"math"
	"math/bits"
	"math/rand"
)

type area256[AS interface {
	allCells() [4]uint64
	Size() int
}] [4]uint64

func (a area256[AS]) And(b area256[AS]) area256[AS] {
	return area256[AS]{
		a[0] & b[0],
		a[1] & b[1],
		a[2] & b[2],
		a[3] & b[3],
	}
}
func (a area256[AS]) Or(b area256[AS]) area256[AS] {
	return area256[AS]{
		a[0] | b[0],
		a[1] | b[1],
		a[2] | b[2],
		a[3] | b[3],
	}
}
func (a area256[AS]) Not() area256[AS] {
	var spec AS
	all := spec.allCells()
	return area256[AS]{
		^a[0] & all[0],
		^a[1] & all[1],
		^a[2] & all[2],
		^a[3] & all[3],
	}
}

func (a area256[AS]) All() area256[AS] {
	var spec AS
	return spec.allCells()
}

func (a area256[AS]) Size() int {
	var spec AS
	return spec.Size()
}

func (a area256[AS]) With(l CellLocation) area256[AS] {
	index, mask := a.getMask(l)
	a[index] = a[index] | mask
	return a
}

func (a area256[AS]) Without(l CellLocation) area256[AS] {
	index, mask := a.getMask(l)
	a[index] = a[index] & ^mask
	return a
}

func (a area256[AS]) getMask(l CellLocation) (int, uint64) {
	var spec AS
	idx := l.Row*spec.Size() + l.Col
	return idx / 64, 1 << (uint64(idx) % 64)
}

func (a area256[AS]) ShiftLeft(n int) area256[AS] {
	var spec AS
	size := spec.Size()
	if n >= size {
		return area256[AS]{0, 0, 0, 0}
	}
	shift := uint64(n)
	if shift == 0 {
		return a
	}

	a = area256[AS]{
		(a[0] >> shift) | (a[1] << (64 - shift)),
		(a[1] >> shift) | (a[2] << (64 - shift)),
		(a[2] >> shift) | (a[3] << (64 - shift)),
		a[3] >> shift,
	}

	// TODO find faster way to clear bits that wrap around to the next row
	for row := 0; row < spec.Size(); row++ {
		for col := spec.Size() - n; col < spec.Size(); col++ {
			a = a.Without(CellLocation{row, col})
		}
	}

	return a
}

func (a area256[AS]) ShiftRight(n int) area256[AS] {
	var spec AS
	size := spec.Size()
	if n >= size {
		return area256[AS]{0, 0, 0, 0}
	}
	shift := uint64(n)
	if shift == 0 {
		return a
	}

	a = area256[AS]{
		a[0] << shift,
		(a[1] << shift) | (a[0] >> (64 - shift)),
		(a[2] << shift) | (a[1] >> (64 - shift)),
		(a[3] << shift) | (a[2] >> (64 - shift)),
	}

	// TODO find faster way to clear bits that wrap around to the next row
	for row := 0; row < spec.Size(); row++ {
		for col := 0; col < n; col++ {
			a = a.Without(CellLocation{row, col})
		}
	}

	return a.And(a.All())
}

func (a area256[AS]) ShiftUp(n int) area256[AS] {
	var spec AS
	size := spec.Size()
	if n >= size {
		return area256[AS]{0, 0, 0, 0}
	}
	shift := uint64(n * size)
	if shift == 0 {
		return a
	}
	if shift < 64 {
		return area256[AS]{
			(a[0] >> shift) | (a[1] << (64 - shift)),
			(a[1] >> shift) | (a[2] << (64 - shift)),
			(a[2] >> shift) | (a[3] << (64 - shift)),
			a[3] >> shift,
		}
	} else if shift < 128 {
		s := shift - 64
		return area256[AS]{
			(a[1] >> s) | (a[2] << (64 - s)),
			(a[2] >> s) | (a[3] << (64 - s)),
			a[3] >> s,
			0,
		}
	} else if shift < 192 {
		s := shift - 128
		return area256[AS]{
			(a[2] >> s) | (a[3] << (64 - s)),
			a[3] >> s,
			0,
			0,
		}
	} else {
		s := shift - 192
		return area256[AS]{
			a[3] >> s,
			0,
			0,
			0,
		}
	}
}

func (a area256[AS]) ShiftDown(n int) area256[AS] {
	var spec AS
	size := spec.Size()
	if n >= size {
		return area256[AS]{0, 0, 0, 0}
	}
	shift := uint64(n * size)
	if shift == 0 {
		return a
	}
	if shift < 64 {
		return area256[AS]{
			a[0] << shift,
			(a[1] << shift) | (a[0] >> (64 - shift)),
			(a[2] << shift) | (a[1] >> (64 - shift)),
			(a[3] << shift) | (a[2] >> (64 - shift)),
		}.And(a.All())
	} else if shift < 128 {
		s := shift - 64
		return area256[AS]{
			0,
			a[0] << s,
			(a[1] << s) | (a[0] >> (64 - s)),
			(a[2] << s) | (a[1] >> (64 - s)),
		}.And(a.All())
	} else if shift < 192 {
		s := shift - 128
		return area256[AS]{
			0,
			0,
			a[0] << s,
			(a[1] << s) | (a[0] >> (64 - s)),
		}.And(a.All())
	} else {
		s := shift - 192
		return area256[AS]{
			0,
			0,
			0,
			a[0] << s,
		}.And(a.All())
	}

}

func (a area256[AS]) ShiftBy(offset Offset) area256[AS] {
	if offset.Row > 0 {
		a = a.ShiftDown(offset.Row)
	} else if offset.Row < 0 {
		a = a.ShiftUp(-offset.Row)
	}
	if offset.Col > 0 {
		a = a.ShiftRight(offset.Col)
	} else if offset.Col < 0 {
		a = a.ShiftLeft(-offset.Col)
	}
	return a
}

func (a area256[AS]) Get(l CellLocation) bool {
	index, mask := a.getMask(l)
	return a[index]&mask != 0
}

func (a area256[AS]) Set(l CellLocation, v bool) area256[AS] {
	if v {
		return a.With(l)
	}
	return a.Without(l)
}

func (a area256[AS]) Locations(yield func(int, CellLocation) bool) {
	var spec AS
	index := 0
	size := spec.Size()
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

func (a area256[AS]) RandomLocation() CellLocation {
	var spec AS
	size := spec.Size()
	return a.nextCell(rand.Intn(size * size))
}

func (a area256[AS]) nextCell(index int) CellLocation {
	var maskedArea area256[AS]
	if index < 64 {
		maskedArea[0] = a[0] & ^(math.MaxUint64 >> (64 - index))
		maskedArea[1] = a[1]
		maskedArea[2] = a[2]
		maskedArea[3] = a[3]
	} else if index < 128 {
		maskedArea[0] = 0
		maskedArea[1] = a[1] & ^(math.MaxUint64 >> (128 - index))
		maskedArea[2] = a[2]
		maskedArea[3] = a[3]
	} else if index < 192 {
		maskedArea[0] = 0
		maskedArea[1] = 0
		maskedArea[2] = a[2] & ^(math.MaxUint64 >> (192 - index))
		maskedArea[3] = a[3]
	} else {
		maskedArea[0] = 0
		maskedArea[1] = 0
		maskedArea[2] = 0
		maskedArea[3] = a[3] & ^(math.MaxUint64 >> (256 - index))
	}
	found := false
	var cell CellLocation
	a.Locations(func(i int, l CellLocation) bool {
		if !found && maskedArea.Get(l) {
			cell = l
			found = true
			return false
		}
		return true
	})
	if found {
		return cell
	}
	// fallback: first cell in a
	found = false
	a.Locations(func(i int, l CellLocation) bool {
		if !found {
			cell = l
			found = true
			return false
		}
		return true
	})
	if found {
		return cell
	}
	return CellLocation{}
}

func (a area256[AS]) Count() int {
	return bits.OnesCount64(a[0]) + bits.OnesCount64(a[1]) + bits.OnesCount64(a[2]) + bits.OnesCount64(a[3])
}

func (a area256[AS]) Empty() bool {
	return a[0] == 0 && a[1] == 0 && a[2] == 0 && a[3] == 0
}

func (a area256[AS]) String() string {
	var spec AS
	size := spec.Size()
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
