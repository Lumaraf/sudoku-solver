package rule

import (
	"errors"
	"fmt"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

var (
	ErrTooFewDigits = errors.New("too few available digits in unique set")
)

type ClassicRules[D sudoku.Digits, A sudoku.Area] struct{}

func (r ClassicRules[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	for row := 0; row < sb.Size(); row++ {
		a := sb.Row(row)
		sb.AddRestriction(UniqueRestriction[D, A]{
			name: fmt.Sprintf("row %d", row+1),
			area: a,
		})
		sb.AddValidator(UniqueValidator[D, A]{
			name: fmt.Sprintf("row %d", row+1),
			area: a,
		})
		for _, cell := range a.Locations {
			sb.AddExclusionArea(cell, a)
		}
	}
	for col := 0; col < sb.Size(); col++ {
		a := sb.Column(col)
		sb.AddRestriction(UniqueRestriction[D, A]{
			name: fmt.Sprintf("col %d", col+1),
			area: a,
		})
		sb.AddValidator(UniqueValidator[D, A]{
			name: fmt.Sprintf("col %d", col+1),
			area: a,
		})
		for _, cell := range a.Locations {
			sb.AddExclusionArea(cell, a)
		}
	}
	for box := 0; box < sb.Size(); box++ {
		a := sb.Box(box)
		sb.AddRestriction(UniqueRestriction[D, A]{
			name: fmt.Sprintf("box %d", box+1),
			area: a,
		})
		sb.AddValidator(UniqueValidator[D, A]{
			name: fmt.Sprintf("box %d", box+1),
			area: a,
		})
		for _, cell := range a.Locations {
			sb.AddExclusionArea(cell, a)
		}
	}
	return nil
}

type DiagonalRule[D sudoku.Digits, A sudoku.Area] struct{}

func (r DiagonalRule[D, A]) Name() string {
	return "diagonal"
}

func (r DiagonalRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	falling := sb.NewArea()
	rising := sb.NewArea()
	for n := 0; n < sb.Size(); n++ {
		sb.AreaWith(&falling, sudoku.CellLocation{n, n})
		sb.AreaWith(&rising, sudoku.CellLocation{sb.Size() - 1 - n, n})
	}
	sb.AddRestriction(UniqueRestriction[D, A]{
		name: "falling diagonal",
		area: falling,
	})
	sb.AddValidator(UniqueValidator[D, A]{
		name: "falling diagonal",
		area: falling,
	})
	sb.AddRestriction(UniqueRestriction[D, A]{
		name: "rising diagonal",
		area: rising,
	})
	sb.AddValidator(UniqueValidator[D, A]{
		name: "rising diagonal",
		area: rising,
	})
	return nil
}

type UniqueRestriction[D sudoku.Digits, A sudoku.Area] struct {
	name string
	area A
}

func (r UniqueRestriction[D, A]) Name() string {
	return fmt.Sprintf("Unique %s", r.name)
}

func (r UniqueRestriction[D, A]) Area() A {
	return r.area
}

type UniqueValidator[D sudoku.Digits, A sudoku.Area] struct {
	name string
	area A
}

func (r UniqueValidator[D, A]) Name() string {
	return fmt.Sprintf("Unique %s", r.name)
}

func (v UniqueValidator[D, A]) Validate(s sudoku.Sudoku[D, A]) error {
	mask := s.NewDigits()
	for _, cell := range v.area.Locations {
		d := s.Get(cell)
		mask = s.UnionDigits(mask, d)
	}
	if mask.Count() < v.area.Size() {
		return ErrTooFewDigits
	}
	return nil
}

//func (r UniqueRestriction[D, A]) Validate(s sudoku.Sudoku[D, A]) error {
//	mask := *new(D)
//	count := 0
//	for _, cell := range s.InvertArea(s.IntersectAreas(r.area, s.SolvedArea())).Locations {
//		//for _, cell := range r.area.And(s.SolvedArea().Not()).Locations {
//		mask = s.UnionDigits(mask, s.Get(cell))
//		//mask = mask | s.Get(cell)
//		count++
//	}
//
//	if mask.Count() < count {
//		return ErrTooFewDigits
//	}
//	return nil
//}
//
//func (r UniqueRestriction[D, A]) ExclusionAreas() map[sudoku.CellLocation]A {
//	area := sudoku.Area{}
//	for _, cell := range r.area.Locations {
//		area.Set(cell, true)
//	}
//	areas := map[sudoku.CellLocation]sudoku.Area{}
//	for _, cell := range r.area.Locations {
//		cellArea := area
//		cellArea.Set(cell, false)
//		areas[cell] = cellArea
//	}
//	return areas
//}
//
//func AddClassicRestrictions(s sudoku.Sudoku) {
//	for row := 0; row < 9; row++ {
//		area := sudoku.Area{}
//		for col := 0; col < 9; col++ {
//			area.Set(sudoku.CellLocation{row, col}, true)
//		}
//		s.AddRestriction(UniqueRestriction{
//			name: fmt.Sprintf("row %d", row+1),
//			area: area,
//		})
//	}
//
//	for col := 0; col < 9; col++ {
//		area := sudoku.Area{}
//		for row := 0; row < 9; row++ {
//			area.Set(sudoku.CellLocation{row, col}, true)
//		}
//		s.AddRestriction(UniqueRestriction{
//			name: fmt.Sprintf("col %d", col+1),
//			area: area,
//		})
//	}
//
//	for box := 0; box < 9; box++ {
//		rowOffset, colOffset := box/3*3, box%3*3
//		area := sudoku.Area{}
//		for row := 0; row < 3; row++ {
//			for col := 0; col < 3; col++ {
//				area.Set(sudoku.CellLocation{rowOffset + row, colOffset + col}, true)
//			}
//		}
//		s.AddRestriction(UniqueRestriction{
//			name: fmt.Sprintf("box %d", box+1),
//			area: area,
//		})
//	}
//}
//
//func AddDiagonalRestrictions(s sudoku.Sudoku) {
//	area := sudoku.Area{}
//	for n := 0; n < 9; n++ {
//		area.Set(sudoku.CellLocation{n, n}, true)
//	}
//	s.AddRestriction(UniqueRestriction{
//		name: "falling diagonal",
//		area: area,
//	})
//
//	area = sudoku.Area{}
//	for n := 0; n < 9; n++ {
//		area.Set(sudoku.CellLocation{8 - n, n}, true)
//	}
//	s.AddRestriction(UniqueRestriction{
//		name: "rising diagonal",
//		area: area,
//	})
//}
//
//func AddColorRestrictions(s sudoku.Sudoku) {
//	for color := 0; color < 9; color++ {
//		rowOffset, colOffset := color/3, color%3
//		area := sudoku.Area{}
//		for row := 0; row < 3; row++ {
//			for col := 0; col < 3; col++ {
//				area.Set(sudoku.CellLocation{rowOffset + row*3, colOffset + col*3}, true)
//			}
//		}
//		s.AddRestriction(UniqueRestriction{
//			name: fmt.Sprintf("color %d:%d", rowOffset, colOffset),
//			area: area,
//		})
//	}
//}
