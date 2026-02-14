package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/strategy"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

func main() {
	//   EEEE
	// EEE   E
	//E E    EE
	//E   EE  E
	//E   EE  E
	//EE     EE
	// E   EEE
	// EEE  E
	//  EEE E

	for a := range tetrisPlacements[sudoku.Digits9, sudoku.Area9x9](sudoku.NewSudokuBuilder9x9()) {
		fmt.Println("-------")
		fmt.Println(a)

		sb := sudoku.NewSudokuBuilder9x9()
		rows := make([]string, 9)
		for row := 0; row < sb.Size(); row++ {
			rowStr := ""
			for col := 0; col < sb.Size(); col++ {
				if a.Get(sudoku.CellLocation{row, col}) {
					rowStr += "E"
				} else {
					rowStr += " "
				}
			}
			rows[row] = rowStr
		}

		sb.Use(
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			//rule.DiagonalRule[sudoku.Digits9, sudoku.Area9x9]{},
			//rule.DisjointGroupsRule[sudoku.Digits9, sudoku.Area9x9]{},
			rule.ParityFromString[sudoku.Digits9, sudoku.Area9x9](rows...),
			NeighborMaskRule[sudoku.Digits9, sudoku.Area9x9]{
				fn: func(a, b int) bool {
					return a%2 == 0 || b%2 == 0 || (a+b) != 10
				},
			},
		)

		s, _ := sb.Build()
		g := s.NewGuesser()
		g.Use(strategy.AllStrategies[sudoku.Digits9, sudoku.Area9x9]())
		for s := range g.Guess(nil, context.Background()) {
			s.Print()
			return
		}
	}

	//tetrisPlacements[sudoku.Digits9, sudoku.Area9x9](sb)

	//sb := sudoku.NewSudokuBuilder9x9()
	//sb.Use(
	//	rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
	//	//NeighborMaskRule[sudoku.Digits9, sudoku.Area9x9]{
	//	//	fn: func(a, b int) bool {
	//	//		return (a + b) != 13
	//	//	},
	//	//},
	//	AntiMagicRule[sudoku.Digits9, sudoku.Area9x9]{},
	//	//rule.NonConsecutiveRule[sudoku.Digits9, sudoku.Area9x9]{},
	//	rule.AntiKnightRule[sudoku.Digits9, sudoku.Area9x9]{},
	//)
	//
	//s, _ := sb.Build()
	//g := s.NewGuesser()
	////g.Use(strategy.AllStrategies[sudoku.Digits9, sudoku.Area9x9]())
	//for s := range g.Guess(nil, context.Background()) {
	//	s.Print()
	//	return
	//}
}

type NeighborMaskRule[D sudoku.Digits[D], A sudoku.Area] struct {
	fn func(a, b int) bool
}

func (r NeighborMaskRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	masks := make(map[int]D, sb.Size())
	for a := 1; a <= sb.Size(); a++ {
		validValues := make([]int, 0, sb.Size())
		for b := 1; b <= sb.Size(); b++ {
			if r.fn(a, b) && r.fn(b, a) {
				validValues = append(validValues, b)
			}
		}
		masks[a] = sb.NewDigits(validValues...)
	}
	sb.AddChangeProcessor(NeighborMaskChangeProcessor[D, A]{
		masks: masks,
	})
	return nil
}

type NeighborMaskChangeProcessor[D sudoku.Digits[D], A sudoku.Area] struct {
	masks map[int]D
}

func (cp NeighborMaskChangeProcessor[D, A]) Name() string {
	return "NeighborMaskChangeProcessor"
}

func (cp NeighborMaskChangeProcessor[D, A]) ProcessChange(s sudoku.Sudoku[D, A], cell sudoku.CellLocation, mask D) error {
	combinedMask := s.NewDigits()
	for v := range mask.Values {
		combinedMask = s.UnionDigits(combinedMask, cp.masks[v])
	}
	targetArea := s.NewAreaFromOffsets(cell, sudoku.Offsets{
		{Row: -1, Col: 0},
		{Row: 1, Col: 0},
		{Row: 0, Col: -1},
		{Row: 0, Col: 1},
		//{Row: -1, Col: -1},
		//{Row: -1, Col: 1},
		//{Row: 1, Col: -1},
		//{Row: 1, Col: 1},
	})
	for _, targetCell := range targetArea.Locations {
		if err := s.Mask(targetCell, combinedMask); err != nil {
			return err
		}
	}
	return nil
}

type AntiMagicRule[D sudoku.Digits[D], A sudoku.Area] struct{}

func (r AntiMagicRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	sb.AddValidator(AntiMagicValidator[D, A]{})
	return nil
}

type AntiMagicValidator[D sudoku.Digits[D], A sudoku.Area] struct{}

func (v AntiMagicValidator[D, A]) Name() string {
	return "AntiMagicValidator"
}

func (v AntiMagicValidator[D, A]) Validate(s sudoku.Sudoku[D, A]) error {
	for a := range iterateMagicSets(s) {
		a = s.IntersectAreas(a, s.SolvedArea())
		if a.Size() != 3 {
			continue
		}

		sum := 0
		for _, l := range a.Locations {
			sum += s.Get(l).Min()
		}

		if sum == 15 {
			return fmt.Errorf("AntiMagicValidator: area %s sums to 15", a)
		}
	}
	return nil
}

func iterateMagicSets[D sudoku.Digits[D], A sudoku.Area](s sudoku.Sudoku[D, A]) func(yield func(area A) bool) {
	return func(yield func(area A) bool) {
		for row := 0; row < s.Size(); row++ {
			for col := 0; col < s.Size(); col++ {
				if !yield(s.NewAreaFromOffsets(sudoku.CellLocation{row, col}, sudoku.Offsets{
					{0, 0},
					{0, 1},
					{0, 2},
				})) {
					return
				}
				if !yield(s.NewAreaFromOffsets(sudoku.CellLocation{row, col}, sudoku.Offsets{
					{0, 0},
					{1, 0},
					{2, 0},
				})) {
					return
				}
				if !yield(s.NewAreaFromOffsets(sudoku.CellLocation{row, col}, sudoku.Offsets{
					{0, 0},
					{1, 1},
					{2, 2},
				})) {
					return
				}
				if !yield(s.NewAreaFromOffsets(sudoku.CellLocation{row, col}, sudoku.Offsets{
					{0, 0},
					{1, -1},
					{2, -2},
				})) {
					return
				}
			}
		}
	}
}

func getTetrisShapes[D sudoku.Digits[D], A sudoku.Area](sb sudoku.SudokuBuilder[D, A]) []sudoku.Offsets {
	shapes := map[A]sudoku.Offsets{}

	shapeFromString := func(rows ...string) {
		offsets := sudoku.Offsets{}
		w := 0
		h := len(rows)
		for r, row := range rows {
			if len(row) > w {
				w = len(row)
			}
			for c, ch := range row {
				if ch != ' ' {
					offsets = append(offsets, sudoku.Offset{Row: r, Col: c})
				}
			}
		}

		key := sb.NewAreaFromOffsets(sudoku.CellLocation{Row: 0, Col: 0}, offsets)
		shapes[key] = offsets

		for n := 0; n < 4; n++ {
			rotatedOffsets := sudoku.Offsets{}
			for _, off := range offsets {
				rotatedOffsets = append(rotatedOffsets, sudoku.Offset{Row: off.Col, Col: h - off.Row - 1})
			}

			key := sb.NewAreaFromOffsets(sudoku.CellLocation{Row: 0, Col: 0}, offsets)
			shapes[key] = offsets

			offsets = rotatedOffsets
			w, h = h, w
		}
	}

	shapeFromString(
		"XXXX",
	)
	shapeFromString(
		"XX",
		"XX",
	)
	shapeFromString(
		"X ",
		"X ",
		"XX",
	)
	shapeFromString(
		" X",
		" X",
		"XX",
	)
	shapeFromString(
		" X",
		"XX",
		"X ",
	)
	shapeFromString(
		"X ",
		"XX",
		" X",
	)
	shapeFromString(
		" X ",
		"XXX",
	)

	result := make([]sudoku.Offsets, 0, len(shapes))
	for _, off := range shapes {
		result = append(result, off)
	}
	return result
}

func tetrisPlacements[D sudoku.Digits[D], A sudoku.Area](sb sudoku.SudokuBuilder[D, A]) func(func(A) bool) {
	shapes := getTetrisShapes(sb)
	rand.Shuffle(len(shapes), func(i, j int) {
		shapes[i], shapes[j] = shapes[j], shapes[i]
	})

	checkAreas := make([]A, 0, 9*3+2)
	for n := 0; n < sb.Size(); n++ {
		checkAreas = append(checkAreas, sb.Row(n))
		checkAreas = append(checkAreas, sb.Column(n))
		checkAreas = append(checkAreas, sb.Box(n))
	}
	falling := sb.NewArea()
	rising := sb.NewArea()
	disjointGroups := make([]A, sb.Size())
	for n := 0; n < sb.Size(); n++ {
		sb.AreaWith(&falling, sudoku.CellLocation{n, n})
		sb.AreaWith(&rising, sudoku.CellLocation{sb.Size() - 1 - n, n})
		for i, l := range sb.Box(n).Locations {
			sb.AreaWith(&disjointGroups[i], l)
		}
	}
	//checkAreas = append(checkAreas, disjointGroups...)
	//checkAreas = append(checkAreas, falling, rising)
	check := func(area A) bool {
		for _, ca := range checkAreas {
			if sb.IntersectAreas(area, ca).Size() > 4 {
				return false
			}
		}
		return true
	}

	usedShapes := make(map[int]bool, len(shapes))
	var findPlacements func(a, blocked A, searchRow int) func(func(A) bool)
	findPlacements = func(a, blocked A, searchRow int) func(func(A) bool) {
		return func(yield func(A) bool) {
			for r := searchRow; r < sb.Size(); r++ {
				row := sb.Row(r)
				if sb.IntersectAreas(a, row).Size() < 4 {
					searchRow = r
					break
				}
			}

			positions := sb.IntersectAreas(sb.Row(searchRow), sb.InvertArea(blocked))
			for !positions.Empty() {
				l := positions.RandomLocation()
				sb.AreaWithout(&positions, l)
				for idx, shape := range shapes {
					if usedShapes[idx] {
						continue
					}
					placement := sb.NewAreaFromOffsets(l, shape)
					if sb.UnionAreas(blocked, placement).Size() != blocked.Size()+4 {
						continue
					}
					newArea := sb.UnionAreas(a, placement)
					if !check(newArea) {
						continue
					}
					if newArea.Size() == 9*4 {
						if !yield(newArea) {
							return
						}
					} else {
						usedShapes[idx] = true
						newBlocked := sb.UnionAreas(blocked, placement)
						for result := range findPlacements(newArea, newBlocked, searchRow) {
							if !yield(result) {
								return
							}
						}
						usedShapes[idx] = false
					}
				}
			}
		}
	}

	blocked := sb.NewArea()
	sb.AreaWith(&blocked, sudoku.CellLocation{0, 0})
	sb.AreaWith(&blocked, sudoku.CellLocation{0, 8})
	sb.AreaWith(&blocked, sudoku.CellLocation{8, 0})
	sb.AreaWith(&blocked, sudoku.CellLocation{8, 8})

	//sb.AreaWith(&blocked, sudoku.CellLocation{3, 4})
	//sb.AreaWith(&blocked, sudoku.CellLocation{4, 3})
	//sb.AreaWith(&blocked, sudoku.CellLocation{4, 4})
	//sb.AreaWith(&blocked, sudoku.CellLocation{4, 5})
	//sb.AreaWith(&blocked, sudoku.CellLocation{5, 4})
	return findPlacements(sb.NewArea(), blocked, 0)
}
