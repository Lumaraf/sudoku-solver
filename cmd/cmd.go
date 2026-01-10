package main

import (
	"context"
	"fmt"
	"github.com/lumaraf/sudoku-solver/rule"
	"github.com/lumaraf/sudoku-solver/strategy"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"time"
)

func main_test() {
	sb := sudoku.NewSudokuBuilder9x9()
	sb.Use(
		rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
		boxGroupsRule[sudoku.Digits9, sudoku.Area9x9]{},
		rule.AntiKnightRule[sudoku.Digits9, sudoku.Area9x9]{},
	)

	//s.SetLogger(sudoku.NewLogger[sudoku.Digits6]())
	s, _ := sb.Build()
	g := s.NewGuesser()
	g.SetChainLimit(1)
	g.Use(
		strategy.AllStrategies[sudoku.Digits9, sudoku.Area9x9](),
	)
	start := time.Now()
	for solution := range g.Guess(sudoku.DefaultGuessSelector[sudoku.Digits9, sudoku.Area9x9], context.Background()) {
		solution.Print()
		//
		time.Sleep(1 * time.Second)
	}
	//s.Print()
	fmt.Printf("time: %v\nerr: %w\n", time.Since(start))
}

type boxGroupsRule[D sudoku.Digits, A sudoku.Area] struct{}

func (r boxGroupsRule[D, A]) Name() string {
	return "box groups"
}

func (r boxGroupsRule[D, A]) Apply(sb sudoku.SudokuBuilder[D, A]) error {
	rowGroups := []D{
		sb.NewDigits(1, 2, 3),
		sb.NewDigits(4, 5, 6),
		sb.NewDigits(7, 8, 9),
	}
	colGroups := []D{
		sb.NewDigits(1, 2, 3),
		sb.NewDigits(4, 5, 6),
		sb.NewDigits(7, 8, 9),
	}

	boxRows, boxCols := sb.BoxSize()
	boxesPerRow := sb.Size() / boxRows
	for box := 0; box < sb.Size(); box++ {
		boxRowOffset := box / boxesPerRow * boxRows
		boxColOffset := box % boxesPerRow * boxCols

		for row := 0; row < boxRows; row++ {
			a := sb.NewArea()
			for col := 0; col < boxCols; col++ {
				sb.AreaWith(&a, sudoku.CellLocation{boxRowOffset + row, boxColOffset + col})
			}
			sb.AddValidator(antiGroupValidator[D, A]{
				area:   a,
				groups: rowGroups,
			})
		}

		for col := 0; col < boxCols; col++ {
			a := sb.NewArea()
			for row := 0; row < boxRows; row++ {
				sb.AreaWith(&a, sudoku.CellLocation{boxRowOffset + row, boxColOffset + col})
			}
			sb.AddValidator(antiGroupValidator[D, A]{
				area:   a,
				groups: colGroups,
			})
		}
	}
	return nil
}

type antiGroupValidator[D sudoku.Digits, A sudoku.Area] struct {
	area   A
	groups []D
}

func (v antiGroupValidator[D, A]) Name() string {
	return "anti-group"
}

func (v antiGroupValidator[D, A]) Validate(s sudoku.Sudoku[D, A]) error {
	groupMap := make(map[D]bool)
	for _, group := range v.groups {
		groupMap[group] = false
	}
	for _, cell := range v.area.Locations {
		d := s.Get(cell)
		for _, group := range v.groups {
			//fmt.Println("group", group, "cell", cell, "d", d)
			if s.IntersectDigits(d, group).Empty() {
				continue
			}
			if s.IntersectDigits(d, s.InvertDigits(group)).Empty() {
				if groupMap[group] {
					return fmt.Errorf("digit %d in cell %s is in group %d twice", d, cell, group)
				}
				groupMap[group] = true
			}
		}
	}
	return nil
}

func main() {
	{
		s, err := sudoku.NewSudoku6x6(
			rule.ClassicRules[sudoku.Digits6, sudoku.Area6x6]{},
			rule.GivenDigitsFromString[sudoku.Digits6, sudoku.Area6x6](
				"  5  2",
				"6     ",
				"4    5",
				"5   4 ",
				"  12  ",
				"     1",
			),
		)
		if err != nil {
			panic(err)
		}
		s.SetLogger(sudoku.NewLogger[sudoku.Digits6]())
		slv := s.NewSolver()
		slv.SetChainLimit(0)
		slv.Use(
			strategy.AllStrategies[sudoku.Digits6, sudoku.Area6x6](),
		)
		start := time.Now()
		err = slv.Solve(context.Background())
		s.Print()
		fmt.Printf("time: %v\nerr: %w\n", time.Since(start), err)
	}

	{
		s, err := sudoku.NewSudoku9x9(
			rule.ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			rule.GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](
				" 3       ",
				"   195   ",
				"  8    6 ",
				"8   6    ",
				"4  8    1",
				"    2    ",
				" 6    28 ",
				"   419  5",
				"       7 ",
			),
		)
		if err != nil {
			panic(err)
		}
		s.SetLogger(sudoku.NewLogger[sudoku.Digits9]())
		slv := s.NewSolver()
		slv.Use(
			strategy.AllStrategies[sudoku.Digits9, sudoku.Area9x9](),
		)
		start := time.Now()
		slv.Solve(context.Background())
		s.Print()
		fmt.Printf("time: %v\nerr: %w\n", time.Since(start), err)
	}

	//parsers := map[string]ruleParser{}
	//parsers["given"] = newRestrictionParser[givenDigits]("given")
}

//type specParser struct {
//	parsers map[string]ruleParser
//}
//
//type rule interface {
//}
//
//type ruleParser func(node *yaml.Node) (rule, error)
//
//func newRestrictionParser[R rule](name string) ruleParser {
//	return func(node *yaml.Node) (rule, error) {
//		r := new(R)
//		err := node.Decode(r)
//		return *r, err
//	}
//}
//
//type givenDigits struct {
//	Rows []string `yaml:"rows"`
//}
