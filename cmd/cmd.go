package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lumaraf/sudoku-solver/restriction"
	"github.com/lumaraf/sudoku-solver/solver"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type sudokuSpec struct {
	rows           []string
	cages          []string
	cageSums       map[rune]int
	sandwichRows   []int
	sandwichCols   []int
	thermometers   [][]sudoku.CellLocation
	palindromes    [][]sudoku.CellLocation
	anitKing       bool
	anitKnight     bool
	nonConsecutive bool
	diagonals      bool
	colors         bool
}

func (spec sudokuSpec) setup(s sudoku.Sudoku) {
	if spec.cageSums != nil {
		cageAreas := map[rune]sudoku.Area{}
		for row, rowContent := range spec.cages {
			for col, cellContent := range rowContent {
				area := cageAreas[cellContent]
				area.Set(sudoku.CellLocation{row, col}, true)
				cageAreas[cellContent] = area
			}
		}

		for key, sum := range spec.cageSums {
			restriction.AddKillerCageRestriction(s, cageAreas[key], sum)
		}
	}

	if spec.sandwichRows != nil {
		for row, sum := range spec.sandwichRows {
			if sum >= 0 {
				s.AddRestriction(restriction.NewRowSandwichRestriction(row, sum))
			}
		}
	}
	if spec.sandwichCols != nil {
		for col, sum := range spec.sandwichCols {
			if sum >= 0 {
				s.AddRestriction(restriction.NewColSandwichRestriction(col, sum))
			}
		}
	}

	if spec.thermometers != nil {
		for _, thermo := range spec.thermometers {
			restriction.AddThermometerRestriction(s, thermo)
			s.AddRestriction(restriction.NewUniqueRestriction(
				"thermometer",
				thermo...,
			))
		}
	}

	if spec.palindromes != nil {
		for _, p := range spec.palindromes {
			for index := 0; index < len(p)/2; index++ {
				s.AddRestriction(restriction.EqualRestriction{
					Cells: []sudoku.CellLocation{
						p[index],
						p[len(p)-index-1],
					},
				})
			}
		}
	}

	if spec.anitKing {
		restriction.AddAntiKingRestriction(s)
	}

	if spec.anitKnight {
		restriction.AddAntiKnighRestriction(s)
	}

	if spec.nonConsecutive {
		restriction.AddNonConsecutiveRestriction(s)
	}

	if spec.diagonals {
		restriction.AddDiagonalRestrictions(s)
	}

	if spec.colors {
		restriction.AddColorRestrictions(s)
	}

	for row, rowContent := range spec.rows {
		for col, cellContent := range rowContent {
			if cellContent < '1' || cellContent > '9' {
				continue
			}
			if err := s.Set(sudoku.CellLocation{row, col}, int(cellContent-'0')); err != nil {
				s.Print()
				panic(err)
			}
		}
	}
}

func main() {
	specs := map[string]sudokuSpec{
		"non consecutive": {
			rows: []string{
				"        5",
				" 1    7  ",
				"7        ",
				"    7  59",
				"         ",
				"42  9    ",
				"        8",
				"  1    7 ",
				"8        ",
			},
			nonConsecutive: true,
		},
		"sandwich": {
			rows: []string{
				"  8      ",
				"         ",
				"1        ",
				"         ",
				" 1     7 ",
				"         ",
				"   5 9   ",
				"         ",
				"         ",
			},
			sandwichRows: []int{0, 0, 4, 10, 32, 4, 7, 0, 15},
			sandwichCols: []int{27, 0, 13, 6, 32, 17, 19, 30, 17},
		},
		"killer cage": {
			cages: []string{
				"ABCCDDEFG",
				"AHHIIIEFG",
				"JJKLLLEMG",
				"JJKLNMMMO",
				"JPPPNQROO",
				"STUVVQRWW",
				"SSUXVVYYZ",
				"SSaXXbbcd",
				"eeaXffccd",
			},
			cageSums: map[rune]int{
				'A': 10, 'B': 6, 'C': 5, 'D': 15, 'E': 20, 'F': 7,
				'G': 16, 'H': 13, 'I': 15, 'J': 17, 'K': 12, 'L': 19,
				'M': 19, 'N': 7, 'O': 18, 'P': 23, 'Q': 8, 'R': 3,
				'S': 25, 'T': 5, 'U': 4, 'V': 22, 'W': 12, 'X': 22,
				'Y': 13, 'Z': 2, 'a': 16, 'b': 15, 'c': 12, 'd': 10,
				'e': 10, 'f': 4,
			},
		},
		"thermometer": {
			rows: []string{
				"     8   ",
				"       9 ",
				"       3 ",
				"         ",
				"    7    ",
				"         ",
				" 1       ",
				" 5       ",
				"   9     ",
			},
			thermometers: [][]sudoku.CellLocation{
				[]sudoku.CellLocation{{2, 0}, {3, 1}, {3, 2}, {2, 3}, {1, 3}, {0, 2}},
				[]sudoku.CellLocation{{0, 6}, {1, 5}, {2, 5}, {3, 6}, {3, 7}, {2, 8}},
				[]sudoku.CellLocation{{8, 2}, {7, 3}, {6, 3}, {5, 2}, {5, 1}, {6, 0}},
				[]sudoku.CellLocation{{6, 8}, {5, 7}, {5, 6}, {6, 5}, {7, 5}, {8, 6}},
			},
		},
		"diagonal": {
			rows: []string{
				"   2 8   ",
				" 1  5  6 ",
				"8  1 3  9",
				"         ",
				" 93   61 ",
				"5  912  8",
				"  5   3  ",
				"         ",
				"3       7",
			},
			diagonals: true,
		},
		"color": {
			rows: []string{
				"   3     ",
				"6  5   9 ",
				"     9  4",
				"   9     ",
				" 8  5  7 ",
				"     8   ",
				"2  4     ",
				" 4   2  3",
				"     6   ",
			},
			colors: true,
		},
		"diagonal sandwich thermos": {
			rows: []string{
				"        6",
				"         ",
				"  2 6    ",
				"     9   ",
				"    8 5  ",
				"         ",
				"      3  ",
				"         ",
				"         ",
			},
			sandwichRows: []int{18, -1, -1, 5, 33, 27, -1, -1, 20},
			sandwichCols: []int{28, -1, -1, 20, 0, 16, -1, -1, 0},
			thermometers: [][]sudoku.CellLocation{
				[]sudoku.CellLocation{{0, 2}, {1, 1}, {2, 0}},
				[]sudoku.CellLocation{{2, 6}, {2, 5}, {1, 4}},
				[]sudoku.CellLocation{{2, 6}, {3, 6}, {4, 7}},
				[]sudoku.CellLocation{{3, 1}, {4, 2}, {5, 2}},
				[]sudoku.CellLocation{{5, 3}, {6, 2}, {6, 1}},
				[]sudoku.CellLocation{{5, 3}, {6, 2}, {7, 1}},
				[]sudoku.CellLocation{{5, 3}, {6, 2}, {7, 2}},
				[]sudoku.CellLocation{{7, 5}, {6, 4}, {6, 3}},
				[]sudoku.CellLocation{{6, 8}, {7, 7}, {8, 6}},
			},
			diagonals: true,
		},
		"palindrome sandwich thermos": {
			rows: []string{
				"         ",
				"         ",
				"         ",
				"         ",
				"8   5   6",
				"         ",
				"         ",
				"         ",
				"         ",
			},
			sandwichRows: []int{0, 6, 8, -1, -1, -1, 20, 4, 3},
			sandwichCols: []int{15, 35, 15, -1, -1, -1, 5, 10, 18},
			thermometers: [][]sudoku.CellLocation{
				[]sudoku.CellLocation{{3, 3}, {2, 2}, {1, 2}, {0, 1}},
				[]sudoku.CellLocation{{3, 5}, {2, 6}, {1, 6}, {0, 7}},
				[]sudoku.CellLocation{{4, 2}, {5, 2}, {6, 3}, {7, 3}, {8, 3}},
				[]sudoku.CellLocation{{4, 6}, {5, 6}, {6, 5}, {7, 5}, {8, 5}},
			},
			palindromes: [][]sudoku.CellLocation{
				[]sudoku.CellLocation{{3, 1}, {4, 1}, {5, 1}, {6, 2}, {7, 2}, {8, 2}},
				[]sudoku.CellLocation{{3, 7}, {4, 7}, {5, 7}, {6, 6}, {7, 6}, {8, 6}},
			},
		},
	}

	for key, spec := range specs {
		//if key != "classic impossible" {
		//	continue
		//}

		fmt.Println(key)

		s := sudoku.NewSudoku()
		s.SetChainLimit(3)
		//s.EnableGuessing()

		restriction.AddClassicRestrictions(s)
		spec.setup(s)

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 5*time.Second)
		err := s.SolveWith(
			ctx,
			solver.UniqueSetSolverFactory,
			solver.UniqueExclusionSolverFactory,
			solver.SandwichSolverFactory,
			solver.AreaSumSolverFactory,
			solver.IncreaseSolverFactory,
			solver.EqualSolverFactory,
			//solver.AntiRelationSolverFactory,
		)

		s.Print()
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println()
	}
}
