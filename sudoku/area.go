package sudoku

import (
	"math"
	"math/bits"
	"math/rand"
)

type Area interface {
	Get(l CellLocation) bool
	Locations(yield func(int, CellLocation) bool)
	RandomLocation() CellLocation
	Size() int
	Empty() bool
	String() string
}

type Area9x9 = area128

type area128 [2]uint64

//var (
//	rowAreas [9]Area
//	colAreas [9]Area
//	boxAreas [9]Area
//)
//
//func init() {
//	for i := 0; i < 9; i++ {
//		for c := 0; c < 9; c++ {
//			rowAreas[i].Set(CellLocation{i, c}, true)
//			colAreas[i].Set(CellLocation{c, i}, true)
//			boxAreas[i].Set(CellLocation{3*(i/3) + c/3, 3*(i%3) + c%3}, true)
//		}
//	}
//}

//func NewArea[A Area](locations ...CellLocation) (a A) {
//	a = *new(A)
//	for _, l := range locations {
//		a.Set(l, true)
//	}
//	return
//}

//func RowArea(row int) Area {
//	if row < 0 || row >= 9 {
//		return Area{}
//	}
//	return rowAreas[row]
//}
//
//func ColArea(col int) Area {
//	if col < 0 || col >= 9 {
//		return Area{}
//	}
//	return colAreas[col]
//}
//
//func BoxArea(box int) Area {
//	if box < 0 || box >= 9 {
//		return Area{}
//	}
//	return boxAreas[box]
//}

func newArea() Area {
	return area128{}
}

func (a area128) Get(l CellLocation) bool {
	index, mask := a.getMask(l)
	return a[index]&mask != 0
}

func (a area128) with(l CellLocation, v bool) area128 {
	index, mask := a.getMask(l)
	if v {
		a[index] = a[index] | mask
	} else {
		a[index] = a[index] & ^mask
	}
	return a
}

func (a area128) getMask(l CellLocation) (int, uint64) {
	idx := l.Row*9 + l.Col
	return idx / 64, 1 << (uint64(idx) % 64)
}

func (a area128) Locations(yield func(int, CellLocation) bool) {
	index := 0
	for idx, b := range a {
		for b != 0 {
			lz := bits.TrailingZeros64(b)
			b = b & ^(1 << lz)
			pos := idx*64 + lz
			if !yield(index, CellLocation{pos / 9, pos % 9}) {
				return
			}
			index++
		}
	}
}

func (a area128) RandomLocation() CellLocation {
	return a.nextCell(rand.Intn(81))
}

func (a area128) nextCell(index int) CellLocation {
	var maskedArea area128
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

func (a area128) Size() int {
	return bits.OnesCount64(a[0]) + bits.OnesCount64(a[1])
}

func (a area128) Empty() bool {
	return a[0] == 0 && a[1] == 0
}

func (a area128) String() string {
	grid := make([]rune, 0, 9*9+9)
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
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
