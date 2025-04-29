package test

import (
	"context"
	"fmt"
	restriction2 "github.com/lumaraf/sudoku-solver/extra/restriction"
	"github.com/lumaraf/sudoku-solver/restriction"
	_ "github.com/lumaraf/sudoku-solver/strategy"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type SudokuSpec struct {
	Rows           []string
	Cages          []string
	CageSums       map[rune]int
	SandwichRows   []int
	SandwichCols   []int
	Thermometers   [][]sudoku.CellLocation
	Palindromes    [][]sudoku.CellLocation
	AntiKing       bool
	AntiKnight     bool
	NonConsecutive bool
	Diagonals      bool
	Colors         bool
	Rule159        bool
}

func (spec SudokuSpec) setup(s sudoku.Sudoku) {
	if spec.CageSums != nil {
		cageAreas := map[rune]sudoku.Area{}
		for row, rowContent := range spec.Cages {
			for col, cellContent := range rowContent {
				area := cageAreas[cellContent]
				area.Set(sudoku.CellLocation{row, col}, true)
				cageAreas[cellContent] = area
			}
		}

		for key, sum := range spec.CageSums {
			restriction2.AddKillerCageRestriction(s, cageAreas[key], sum)
		}
	}

	if spec.SandwichRows != nil {
		for row, sum := range spec.SandwichRows {
			if sum >= 0 {
				s.AddRestriction(restriction2.NewRowSandwichRestriction(row, sum))
			}
		}
	}
	if spec.SandwichCols != nil {
		for col, sum := range spec.SandwichCols {
			if sum >= 0 {
				s.AddRestriction(restriction2.NewColSandwichRestriction(col, sum))
			}
		}
	}

	if spec.Thermometers != nil {
		for _, thermo := range spec.Thermometers {
			restriction2.AddThermometerRestriction(s, thermo)
			s.AddRestriction(restriction.NewUniqueRestriction(
				"thermometer",
				thermo...,
			))
		}
	}

	if spec.Palindromes != nil {
		for _, p := range spec.Palindromes {
			for index := 0; index < len(p)/2; index++ {
				s.AddRestriction(restriction2.EqualRestriction{
					Cells: []sudoku.CellLocation{
						p[index],
						p[len(p)-index-1],
					},
				})
			}
		}
	}

	if spec.AntiKing {
		restriction2.AddAntiKingRestriction(s)
	}

	if spec.AntiKnight {
		restriction2.AddAntiKnighRestriction(s)
	}

	if spec.NonConsecutive {
		restriction2.AddNonConsecutiveRestriction(s)
	}

	if spec.Diagonals {
		restriction.AddDiagonalRestrictions(s)
	}

	if spec.Colors {
		restriction.AddColorRestrictions(s)
	}

	if spec.Rule159 {
		restriction2.Add159Restriction(s)
	}

	for row, rowContent := range spec.Rows {
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

type SudokuTests map[string]SudokuSpec

func (tests SudokuTests) Run(t *testing.T) {
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

			slv := sudoku.NewSolver(s)
			assert.NoError(t, slv.Solve(ctx))
			assert.True(t, slv.IsSolved())
			fmt.Printf("Stats: %+v\n", s.Stats())
		})
	}
}
