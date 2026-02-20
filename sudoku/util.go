package sudoku

import (
	"fmt"
)

func PrintGrid[D Digits[D], A Area[A]](s Sudoku[D, A]) {
	boxRows, boxCols := s.BoxSize()
	gridSize := s.Size()

	rowLength := boxCols*(boxCols*2+1) + 1

	createLine := func(start, edge, fill, boxEdge, end rune) string {
		line := make([]rune, 0, rowLength)
		for col := 0; col < gridSize; col++ {
			if col == 0 {
				line = append(line, start)
			} else if col%boxCols == 0 {
				line = append(line, boxEdge)
			} else {
				line = append(line, edge)
			}
			for i := 0; i < boxCols*2+1; i++ {
				line = append(line, fill)
			}
		}
		line = append(line, end)
		return string(line)
	}

	topLine := createLine('╔', '╤', '═', '╦', '╗')
	midLine := createLine('╟', '┼', '─', '╫', '╢')
	boxLine := createLine('╠', '╪', '═', '╬', '╣')
	bottomLine := createLine('╚', '╧', '═', '╩', '╝')

	symbols := make([]rune, 0, 36)
	for r := '1'; r <= '9'; r++ {
		symbols = append(symbols, r)
	}
	for r := 'A'; r <= 'Z'; r++ {
		symbols = append(symbols, r)
	}

	fmt.Println(topLine)
	for row := 0; row < gridSize; row++ {
		for subRow := 0; subRow < boxRows; subRow++ {
			line := make([]rune, 0, rowLength)
			for col := 0; col < gridSize; col++ {
				if col%boxCols == 0 {
					line = append(line, '║')
				} else {
					line = append(line, '│')
				}
				line = append(line, ' ')
				cell := s.Get(CellLocation{Row: row, Col: col})
				for digit := subRow * boxCols; digit < (subRow+1)*boxCols; digit++ {
					if cell.CanContain(digit + 1) {
						line = append(line, symbols[digit])
					} else {
						line = append(line, ' ')
					}
					line = append(line, ' ')
				}
			}
			line = append(line, '║')
			fmt.Println(string(line))
		}
		if (row+1)%boxRows == 0 && row+1 < gridSize {
			fmt.Println(boxLine)
		} else if row+1 < gridSize {
			fmt.Println(midLine)
		}
	}
	fmt.Println(bottomLine)
}

func GetRestrictions[D Digits[D], A Area[A], R Restriction[D, A]](s Sudoku[D, A]) func(yield func(R) bool) {
	return func(yield func(R) bool) {
		for _, restriction := range s.getRestrictions() {
			if r, ok := restriction.(R); ok {
				if !yield(r) {
					return
				}
			}
		}
	}
}

var bitMasks = [65536][]uint8{}

func init() {
	c := 0
	for i := 0; i < len(bitMasks); i++ {
		var bits []uint8
		for j := 0; j < 16; j++ {
			if (i & (1 << j)) != 0 {
				bits = append(bits, uint8(j))
				c++
			}
		}
		bitMasks[i] = bits
	}
}

func iterateBits[UI ~uint16 | ~uint32 | ~uint64](v UI) func(yield func(uint8) bool) {
	return func(yield func(uint8) bool) {
		offset := uint8(0)
		for v != 0 {
			for _, bit := range bitMasks[v&0xFFFF] {
				if !yield(bit + offset) {
					return
				}
			}
			offset += 16
			v = v >> 16
		}
	}
}
