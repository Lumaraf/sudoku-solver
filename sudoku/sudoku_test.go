package sudoku

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	s := newSudoku()
	s.nextChanged = Area{}

	loc := CellLocation{Row: 0, Col: 0}

	assert.Equal(t, AllDigits, s.Get(loc))

	assert.Error(t, s.Set(loc, 0))
	assert.Error(t, s.Set(loc, 10))

	assert.Equal(t, Area{}, s.NextChangedArea())
	assert.NoError(t, s.Set(loc, 5))
	assert.Equal(t, NewArea(loc), s.NextChangedArea())

	assert.Equal(t, NewDigits(5), s.Get(loc))

	assert.NoError(t, s.Set(loc, 5))

	assert.Error(t, s.Set(loc, 4))
	assert.Equal(t, NewDigits(5), s.Get(loc))
}

func TestRemoveOption(t *testing.T) {
	s := newSudoku()
	s.nextChanged = Area{}

	loc := CellLocation{Row: 0, Col: 0}

	assert.Error(t, s.RemoveOption(loc, 0))
	assert.Error(t, s.RemoveOption(loc, 10))

	assert.Equal(t, Area{}, s.NextChangedArea())
	assert.NoError(t, s.RemoveOption(loc, 5))
	assert.Equal(t, NewArea(loc), s.NextChangedArea())

	assert.Equal(t, NewDigits(1, 2, 3, 4, 6, 7, 8, 9), s.Get(loc))

	assert.NoError(t, s.RemoveOption(loc, 3))

	assert.Equal(t, NewDigits(1, 2, 4, 6, 7, 8, 9), s.Get(loc))
}

func TestMask(t *testing.T) {
	s := newSudoku()
	s.nextChanged = Area{}

	loc := CellLocation{Row: 0, Col: 0}

	assert.Equal(t, Area{}, s.NextChangedArea())
	assert.NoError(t, s.Mask(loc, NewDigits(1, 2, 3)))
	assert.Equal(t, NewDigits(1, 2, 3), s.Get(loc))
	assert.Equal(t, NewArea(loc), s.NextChangedArea())

	assert.NoError(t, s.Mask(loc, NewDigits(2, 3, 4, 5)))
	assert.Equal(t, NewDigits(2, 3), s.Get(loc))

	assert.Equal(t, Area{}, s.SolvedArea())
	assert.NoError(t, s.Mask(loc, NewDigits(3, 4, 5)))
	assert.Equal(t, NewDigits(3), s.Get(loc))
	assert.Equal(t, NewArea(loc), s.SolvedArea())

	assert.Error(t, s.Mask(loc, NewDigits(1, 4, 5)))
	assert.Equal(t, NewDigits(3), s.Get(loc))
}

func TestRemoveMask(t *testing.T) {
	s := newSudoku()
	s.nextChanged = Area{}

	loc := CellLocation{Row: 0, Col: 0}

	assert.NoError(t, s.RemoveMask(loc, NewDigits(1, 2)))
	assert.Equal(t, NewDigits(3, 4, 5, 6, 7, 8, 9), s.Get(loc))
	assert.Equal(t, NewArea(loc), s.NextChangedArea())

	assert.NoError(t, s.RemoveMask(loc, NewDigits(8, 9)))
	assert.Equal(t, NewDigits(3, 4, 5, 6, 7), s.Get(loc))

	assert.Equal(t, Area{}, s.SolvedArea())
	assert.NoError(t, s.RemoveMask(loc, NewDigits(3, 4, 6, 7)))
	assert.Equal(t, NewDigits(5), s.Get(loc))

	assert.Error(t, s.RemoveMask(loc, AllDigits))
	assert.Equal(t, NewDigits(5), s.Get(loc))
}
