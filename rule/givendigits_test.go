package rule

import (
	"github.com/lumaraf/sudoku-solver/sudoku"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGivenDigits(t *testing.T) {
	for _, test := range []struct {
		name         string
		rows         []string
		expectSolved int
	}{
		{
			name:         "empty",
			rows:         []string{},
			expectSolved: 0,
		},
		{
			name: "row",
			rows: []string{
				"123456789",
			},
			expectSolved: 9,
		},
		{
			name: "column",
			rows: []string{
				"1",
				"2",
				"3",
				"4",
				"5",
				"6",
				"7",
				"8",
				"9",
			},
			expectSolved: 9,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			s, err := sudoku.NewSudoku9x9(
				GivenDigitsFromString[sudoku.Digits9, sudoku.Area9x9](test.rows...),
			)
			assert.NoError(t, err)
			assert.Equal(t, test.expectSolved, s.SolvedArea().Size())
		})
	}
}
