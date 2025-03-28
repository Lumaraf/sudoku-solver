package sudoku

import (
	"errors"
	"fmt"
	"strings"
)

func LoadGridFromString(s Sudoku, grid string) (err error) {
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

func PrintGrid(s Sudoku) {
	lines := [][]rune{
		[]rune("╔═══════╤═══════╤═══════╦═══════╤═══════╤═══════╦═══════╤═══════╤═══════╗"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("╠═══════╪═══════╪═══════╬═══════╪═══════╪═══════╬═══════╪═══════╪═══════╣"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("╠═══════╪═══════╪═══════╬═══════╪═══════╪═══════╬═══════╪═══════╪═══════╣"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("╟───────┼───────┼───────╫───────┼───────┼───────╫───────┼───────┼───────╢"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("║       │       │       ║       │       │       ║       │       │       ║"),
		[]rune("╚═══════╧═══════╧═══════╩═══════╧═══════╧═══════╩═══════╧═══════╧═══════╝"),
	}

	count := 0
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			top := 1 + row*4
			left := 2 + col*8
			d := s.Get(CellLocation{row, col})
			for v := range d.Values {
				lines[top+(v-1)/3][left+((v-1)%3)*2] = rune(v + '0')
				count++
			}
		}
	}

	for _, l := range lines {
		fmt.Println(string(l))
	}
}
