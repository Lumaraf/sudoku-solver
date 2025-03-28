package restriction

import (
	"errors"
	"fmt"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

var (
	ErrTooFewDigits = errors.New("too few available digits in unique set")
)

type UniqueRestriction struct {
	name string
	area sudoku.Area
}

func NewUniqueRestriction(name string, cells ...sudoku.CellLocation) UniqueRestriction {
	return UniqueRestriction{
		name: name,
		area: sudoku.NewArea(cells...),
	}
}

func (r UniqueRestriction) Name() string {
	return "Unique"
}

func (r UniqueRestriction) Area() sudoku.Area {
	return r.area
}

func (r UniqueRestriction) Validate(s sudoku.Sudoku) error {
	mask := sudoku.Digits(0)
	count := 0
	for _, cell := range r.area.And(s.SolvedArea().Not()).Locations {
		mask = mask | s.Get(cell)
		count++
	}

	if mask.Count() < count {
		return ErrTooFewDigits
	}
	return nil
}

func (r UniqueRestriction) ExclusionAreas() map[sudoku.CellLocation]sudoku.Area {
	area := sudoku.Area{}
	for _, cell := range r.area.Locations {
		area.Set(cell, true)
	}
	areas := map[sudoku.CellLocation]sudoku.Area{}
	for _, cell := range r.area.Locations {
		cellArea := area
		cellArea.Set(cell, false)
		areas[cell] = cellArea
	}
	return areas
}

func AddClassicRestrictions(s sudoku.Sudoku) {
	for row := 0; row < 9; row++ {
		area := sudoku.Area{}
		for col := 0; col < 9; col++ {
			area.Set(sudoku.CellLocation{row, col}, true)
		}
		s.AddRestriction(UniqueRestriction{
			name: fmt.Sprintf("row %d", row+1),
			area: area,
		})
	}

	for col := 0; col < 9; col++ {
		area := sudoku.Area{}
		for row := 0; row < 9; row++ {
			area.Set(sudoku.CellLocation{row, col}, true)
		}
		s.AddRestriction(UniqueRestriction{
			name: fmt.Sprintf("col %d", col+1),
			area: area,
		})
	}

	for box := 0; box < 9; box++ {
		rowOffset, colOffset := box/3*3, box%3*3
		area := sudoku.Area{}
		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				area.Set(sudoku.CellLocation{rowOffset + row, colOffset + col}, true)
			}
		}
		s.AddRestriction(UniqueRestriction{
			name: fmt.Sprintf("box %d", box+1),
			area: area,
		})
	}
}

func AddDiagonalRestrictions(s sudoku.Sudoku) {
	area := sudoku.Area{}
	for n := 0; n < 9; n++ {
		area.Set(sudoku.CellLocation{n, n}, true)
	}
	s.AddRestriction(UniqueRestriction{
		name: "falling diagonal",
		area: area,
	})

	area = sudoku.Area{}
	for n := 0; n < 9; n++ {
		area.Set(sudoku.CellLocation{8 - n, n}, true)
	}
	s.AddRestriction(UniqueRestriction{
		name: "rising diagonal",
		area: area,
	})
}

func AddColorRestrictions(s sudoku.Sudoku) {
	for color := 0; color < 9; color++ {
		rowOffset, colOffset := color/3, color%3
		area := sudoku.Area{}
		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				area.Set(sudoku.CellLocation{rowOffset + row*3, colOffset + col*3}, true)
			}
		}
		s.AddRestriction(UniqueRestriction{
			name: fmt.Sprintf("color %d:%d", rowOffset, colOffset),
			area: area,
		})
	}
}
