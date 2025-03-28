package test

import (
	"context"
	"fmt"
	"github.com/lumaraf/sudoku-solver/restriction"
	_ "github.com/lumaraf/sudoku-solver/solver"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	rule159        bool
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

	if spec.rule159 {
		restriction.Add159Restriction(s)
	}

	for row, rowContent := range spec.rows {
		for col, cellContent := range rowContent {
			cell := sudoku.CellLocation{row, col}
			switch cellContent {
			case 'O':
				if err := s.Mask(cell, sudoku.NewDigits(1, 3, 5, 7, 9)); err != nil {
					s.Print()
					panic(err)
				}
			case 'E':
				if err := s.Mask(cell, sudoku.NewDigits(2, 4, 6, 8)); err != nil {
					s.Print()
					panic(err)
				}
			default:
				if cellContent < '1' || cellContent > '9' {
					continue
				}
				if err := s.Set(cell, int(cellContent-'0')); err != nil {
					s.Print()
					panic(err)
				}
			}
		}
	}
}

type sudokuTests map[string]sudokuSpec

func (tests sudokuTests) Run(t *testing.T) {
	for name, spec := range tests {
		name, spec := name, spec
		t.Run(name, func(t *testing.T) {
			//t.Parallel()

			s := sudoku.NewSudoku()
			s.SetChainLimit(5)

			restriction.AddClassicRestrictions(s)
			spec.setup(s)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			s.SetLogger(sudoku.NewLogger())
			assert.NoError(t, s.Solve(ctx))
			assert.True(t, s.IsSolved())
			fmt.Printf("Stats: %+v\n", s.Stats())
		})
	}
}
