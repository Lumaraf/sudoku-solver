package sudoku

import (
	"math"
	"math/bits"
	"math/rand"
)

type Area [2]uint64

var (
	rowAreas [9]Area
	colAreas [9]Area
	boxAreas [9]Area
)

func init() {
	for i := 0; i < 9; i++ {
		for c := 0; c < 9; c++ {
			rowAreas[i].Set(CellLocation{i, c}, true)
			colAreas[i].Set(CellLocation{c, i}, true)
			boxAreas[i].Set(CellLocation{3*(i/3) + c/3, 3*(i%3) + c%3}, true)
		}
	}
}

func NewArea(locations ...CellLocation) (a Area) {
	for _, l := range locations {
		a.Set(l, true)
	}
	return
}

func RowArea(row int) Area {
	if row < 0 || row >= 9 {
		return Area{}
	}
	return rowAreas[row]
}

func ColArea(col int) Area {
	if col < 0 || col >= 9 {
		return Area{}
	}
	return colAreas[col]
}

func BoxArea(box int) Area {
	if box < 0 || box >= 9 {
		return Area{}
	}
	return boxAreas[box]
}

func (a *Area) Get(l CellLocation) bool {
	index, mask := a.getMask(l)
	return a[index]&mask != 0
}

func (a *Area) Set(l CellLocation, v bool) {
	index, mask := a.getMask(l)
	if v {
		a[index] = a[index] | mask
	} else {
		a[index] = a[index] & ^mask
	}
}

func (a Area) And(o Area) Area {
	return Area{
		a[0] & o[0],
		a[1] & o[1],
	}
}

func (a Area) Or(o Area) Area {
	return Area{
		a[0] | o[0],
		a[1] | o[1],
	}
}

func (a Area) Not() Area {
	return Area{
		^a[0],
		(^a[1]) & 0b11111111111111111,
	}
}

func (a *Area) getMask(l CellLocation) (int, uint64) {
	idx := l.Row*9 + l.Col
	return idx / 64, 1 << (uint64(idx) % 64)
}

func (a Area) Locations(yield func(int, CellLocation) bool) {
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

func (a Area) RandomLocation() CellLocation {
	return a.nextCell(rand.Intn(81))
}

func (a Area) nextCell(index int) CellLocation {
	var maskedArea Area
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

func (a Area) Size() int {
	return bits.OnesCount64(a[0]) + bits.OnesCount64(a[1])
}

func (a Area) Empty() bool {
	return a[0] == 0 && a[1] == 0
}

func (a Area) String() string {
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
