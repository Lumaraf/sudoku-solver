package sudoku

import (
	"errors"
	"fmt"
	"strings"
)

func LoadGridFromString[D Digits, A Area](s Sudoku[D, A], grid string) (err error) {
	rows := strings.Split(grid, "\n")
	for row, rowContent := range rows {
		for col, cellContent := range rowContent {
			if cellContent < '1' || cellContent > '9' {
				if cellContent != ' ' {
					err = errors.Join(err, errors.New("invalid cell content"))
				}
				continue
			}
			err = errors.Join(err, s.Set(CellLocation{row, col}, int(cellContent-'0')))
		}
	}
	return
}

func PrintGrid[D Digits, A Area](s Sudoku[D, A]) {
	boxRows, boxCols := s.BoxSize()
	gridSize := s.Size()

	rowLength := boxCols*(boxCols*2+1) + 1

	createLine := func(start, edge, fill, boxEdge, end rune) string {
		line := make([]rune, 0, rowLength)
		line = append(line, start)
		for col := 0; col < gridSize; col++ {
			if col > 0 {
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
						line = append(line, rune('0'+digit+1))
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

//func GetRestrictions[D Digits, A Area, R Restriction[D, A]](s Sudoku[D, A]) func(yield func(R) bool) {
//	return func(yield func(R) bool) {
//		for _, restriction := range s.getRestrictions() {
//			if r, ok := restriction.(R); ok {
//				if !yield(r) {
//					return
//				}
//			}
//		}
//	}
//}
