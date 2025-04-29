package restriction

//import (
//	"fmt"
//	"github.com/lumaraf/sudoku-solver/sudoku"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestAddClassicRestrictions(t *testing.T) {
//	s := sudoku.NewSudoku()
//	AddClassicRestrictions(s)
//
//	for row := 0; row < 9; row++ {
//		for col := 0; col < 9; col++ {
//			l := sudoku.CellLocation{row, col}
//
//			expectedArea := sudoku.RowArea(row).
//				Or(sudoku.ColArea(col)).
//				Or(sudoku.BoxArea(l.Box()))
//			expectedArea.Set(l, false)
//
//			assert.Equal(t, expectedArea, s.GetExclusionArea(l))
//		}
//	}
//}
//
//func TestUniqueRestriction_Validate(t *testing.T) {
//	for _, test := range []struct {
//		name          string
//		masks         map[int]sudoku.Digits
//		expectedError error
//	}{
//		{
//			"default",
//			map[int]sudoku.Digits{},
//			nil,
//		},
//		{
//			"invalid set",
//			map[int]sudoku.Digits{
//				0: sudoku.NewDigits(1, 2, 3),
//				1: sudoku.NewDigits(1, 2, 3),
//				2: sudoku.NewDigits(1, 2, 3),
//				3: sudoku.NewDigits(1, 2, 3),
//				4: sudoku.NewDigits(5),
//				5: sudoku.NewDigits(6),
//				6: sudoku.NewDigits(7),
//				7: sudoku.NewDigits(8),
//				8: sudoku.NewDigits(9),
//			},
//			ErrTooFewDigits,
//		},
//	} {
//		t.Run(test.name, func(t *testing.T) {
//			s := sudoku.NewSudoku()
//			r := UniqueRestriction{
//				name: "test",
//				area: sudoku.NewArea([]sudoku.CellLocation{
//					{0, 0},
//					{0, 1},
//					{0, 2},
//					{0, 3},
//					{0, 4},
//					{0, 5},
//					{0, 6},
//					{0, 7},
//					{0, 8},
//				}...),
//			}
//			s.AddRestriction(r)
//
//			for col, mask := range test.masks {
//				assert.NoError(t, s.Mask(sudoku.CellLocation{0, col}, mask))
//			}
//
//			err := r.Validate(s)
//			if test.expectedError == nil {
//				assert.NoError(t, err)
//			} else {
//				fmt.Println(err)
//				assert.ErrorIs(t, err, test.expectedError)
//			}
//		})
//	}
//}
