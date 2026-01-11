package test

import (
	"context"
	"fmt"
	"github.com/lumaraf/sudoku-solver/strategy"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type SudokuTests[D sudoku.Digits, A sudoku.Area] map[string][]sudoku.Rule[D, A]

func (tests SudokuTests[D, A]) Run(t *testing.T, builderFunc func() sudoku.SudokuBuilder[D, A]) {
	for name, rules := range tests {
		t.Run(name, func(t *testing.T) {
			//t.Parallel()

			b := builderFunc()
			assert.NoError(t, b.Use(rules...))

			s, err := b.Build()
			assert.NoError(t, err)

			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()

			s.SetLogger(sudoku.NewLogger[D]())

			//g := s.NewGuesser()
			//g.Use(strategy.AllStrategies[D, A]())
			//for s := range g.Guess(nil, ctx) {
			//	s.Print()
			//}

			slv := s.NewSolver()
			slv.Use(strategy.AllStrategies[D, A]())
			slv.SetChainLimit(0)
			assert.NoError(t, slv.Solve(ctx))
			if !assert.True(t, s.IsSolved()) {
				s.Print()
			}
			fmt.Printf("Stats: %+v\n", s.Stats())
		})
	}
}
