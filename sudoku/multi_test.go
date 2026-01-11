package sudoku

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultiSudokuBuilder_Overlap(t *testing.T) {
	sb1 := NewSudokuBuilder9x9()
	sb2 := NewSudokuBuilder9x9()
	mb := MultiSudokuBuilder[Digits9, Area9x9]{}
	assert.NoError(t, mb.Overlap(sb1, CellLocation{0, 0}, 3, 3, sb2))
	//assert.NoError(t, mb.Overlap(sb2, CellLocation{8, 8}, 3, 3, sb1))

	s1, err := sb1.Build()
	assert.NoError(t, err)
	s2, err := sb2.Build()
	assert.NoError(t, err)

	assert.NoError(t, s1.Set(CellLocation{0, 0}, 1))
	assert.Equal(t, s2.NewDigits(1), s2.Get(CellLocation{6, 6}))

	assert.NoError(t, s2.Mask(CellLocation{7, 8}, s2.NewDigits(3, 4, 5)))
	assert.Equal(t, s1.NewDigits(3, 4, 5), s1.Get(CellLocation{1, 2}))

	s1.Print()
	s2.Print()
}
