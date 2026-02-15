package rule

import (
	"errors"
	"fmt"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

var (
	ErrTooFewDigits = errors.New("too few available digits in unique set")
)

type UniqueAreaRule[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	name string
	area A
}

func NewUniqueAreaRule[D sudoku.Digits[D], A sudoku.Area[A]](name string, area A) UniqueAreaRule[D, A] {
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
		name: r.name,
		area: r.area,
	})
	sb.AddValidator(UniqueValidator[D, A]{
		name: r.name,
		area: r.area,
	})
	for _, cell := range r.area.Locations {
		sb.AddExclusionArea(cell, r.area)
	}
	return nil
}

type ClassicRules[D sudoku.Digits[D], A sudoku.Area[A]] struct{}

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

type DiagonalRule[D sudoku.Digits[D], A sudoku.Area[A]] struct{}

func (r DiagonalRule[D, A]) Name() string {
	return "diagonal"
}

func (r DiagonalRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	falling := sb.NewArea()
	rising := sb.NewArea()
	for n := 0; n < sb.Size(); n++ {
		falling = falling.With(sudoku.CellLocation{n, n})
		rising = rising.With(sudoku.CellLocation{sb.Size() - 1 - n, n})
	}
	return sb.Use(
		NewUniqueAreaRule[D, A]("falling diagonal", falling),
		NewUniqueAreaRule[D, A]("rising diagonal", rising),
	)
}

type DisjointGroupsRule[D sudoku.Digits[D], A sudoku.Area[A]] struct{}

func (r DisjointGroupsRule[D, A]) Name() string {
	return "disjoint groups"
}

func (r DisjointGroupsRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	groups := make([]A, sb.Size())
	for box := 0; box < sb.Size(); box++ {
		for n, l := range sb.Box(box).Locations {
			groups[n] = groups[n].With(l)
		}
	}

	rules := make(sudoku.Rules[D, A], 0, len(groups))
	for n, group := range groups {
		rules = append(rules, NewUniqueAreaRule[D, A](fmt.Sprintf("disjoint group %d", n+1), group))
	}
	return sb.Use(rules...)
}

type UniqueRestriction[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	name string
	area A
}

func (r UniqueRestriction[D, A]) Name() string {
	return fmt.Sprintf("Unique %s", r.name)
}

func (r UniqueRestriction[D, A]) Area() A {
	return r.area
}

type UniqueValidator[D sudoku.Digits[D], A sudoku.Area[A]] struct {
	name string
	area A
}

func (v UniqueValidator[D, A]) Name() string {
	return fmt.Sprintf("Unique %s", v.name)
}

func (v UniqueValidator[D, A]) Validate(s sudoku.Sudoku[D, A]) error {
	mask := s.NewDigits()
	for _, cell := range v.area.Locations {
		d := s.Get(cell)
		mask = mask.Or(d)
	}
	if mask.Count() < v.area.Size() {
		return ErrTooFewDigits
	}
	return nil
}
