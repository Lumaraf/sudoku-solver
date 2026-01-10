package rule

import (
	"fmt"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClassicRules(t *testing.T) {
	s, err := sudoku.NewSudoku9x9(
		ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
	)
	assert.NoError(t, err)

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			l := sudoku.CellLocation{row, col}
			assert.Equal(t, 20, s.GetExclusionArea(l).Size())
		}
	}
}

func TestUniqueRestriction_Validate(t *testing.T) {
	for _, test := range []struct {
		name          string
		masks         map[int][]int
		expectedError error
	}{
		{
			"default",
			map[int][]int{},
			nil,
		},
		{
			"invalid set",
			map[int][]int{
				0: {1, 2, 3},
				1: {1, 2, 3},
				2: {1, 2, 3},
				3: {1, 2, 3},
				4: {5},
				5: {6},
				6: {7},
				7: {8},
				8: {9},
			},
			ErrTooFewDigits,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			s, err := sudoku.NewSudoku9x9(
				ClassicRules[sudoku.Digits9, sudoku.Area9x9]{},
			)
			assert.NoError(t, err)

			for col, mask := range test.masks {
				assert.NoError(t, s.Mask(
					sudoku.CellLocation{0, col},
					s.NewDigits(mask...),
				))
			}

			s.NewDigits()

			err = s.Validate()
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				fmt.Println(err)
				assert.ErrorIs(t, err, test.expectedError)
			}
		})
	}
}
