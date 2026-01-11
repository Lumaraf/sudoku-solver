package rule

import (
	"errors"
	"fmt"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

var (
	ErrTooFewDigits = errors.New("too few available digits in unique set")
)

type UniqueAreaRule[D sudoku.Digits, A sudoku.Area] struct {
	name string
	area A
}

func NewUniqueAreaRule[D sudoku.Digits, A sudoku.Area](name string, area A) UniqueAreaRule[D, A] {
	return UniqueAreaRule[D, A]{
		name: name,
		area: area,
	}
}

func (r UniqueAreaRule[D, A]) Name() string {
	return "unique area"
}

func (r UniqueAreaRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	sb.AddRestriction(UniqueRestriction[D, A]{
		name: "custom area",
		area: r.area,
	})
	sb.AddValidator(UniqueValidator[D, A]{
		name: "custom area",
		area: r.area,
	})
	for _, cell := range r.area.Locations {
		sb.AddExclusionArea(cell, r.area)
	}
	return nil
}

type ClassicRules[D sudoku.Digits, A sudoku.Area] struct{}

func (r ClassicRules[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	rules := make(sudoku.Rules[D, A], 0, sb.Size()*3)
	for row := 0; row < sb.Size(); row++ {
		a := sb.Row(row)
		rules = append(rules, NewUniqueAreaRule[D, A](fmt.Sprintf("row %d", row+1), a))
	}
	for col := 0; col < sb.Size(); col++ {
		a := sb.Column(col)
		rules = append(rules, NewUniqueAreaRule[D, A](fmt.Sprintf("col %d", col+1), a))
	}
	for box := 0; box < sb.Size(); box++ {
		a := sb.Box(box)
		rules = append(rules, NewUniqueAreaRule[D, A](fmt.Sprintf("box %d", box+1), a))
	}
	return sb.Use(rules...)
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
	return sb.Use(
		NewUniqueAreaRule[D, A]("falling diagonal", falling),
		NewUniqueAreaRule[D, A]("rising diagonal", rising),
	)
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
