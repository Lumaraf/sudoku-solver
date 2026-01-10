package sudoku

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	s := newSudoku[Digits9, Area9x9, grid9x9[Digits9], size9]()
	s.nextChanged = Area9x9{}

	loc := CellLocation{Row: 0, Col: 0}

	assert.Equal(t, s.AllDigits(), s.Get(loc))

	assert.Error(t, s.Set(loc, 0))
	assert.Error(t, s.Set(loc, 10))

	assert.Equal(t, Area9x9{}, s.NextChangedArea())
	assert.NoError(t, s.Set(loc, 5))
	assert.Equal(t, s.NewArea().with(loc, true), s.NextChangedArea())

	assert.Equal(t, s.NewDigits(5), s.Get(loc))

	assert.NoError(t, s.Set(loc, 5))

	assert.Error(t, s.Set(loc, 4))
	assert.Equal(t, s.NewDigits(5), s.Get(loc))
}

func TestRemoveOption(t *testing.T) {
	s := newSudoku[Digits9, Area9x9, grid9x9[Digits9], size9]()
	s.nextChanged = Area9x9{}

	loc := CellLocation{Row: 0, Col: 0}

	assert.Error(t, s.RemoveOption(loc, 0))
	assert.Error(t, s.RemoveOption(loc, 10))

	assert.Equal(t, Area9x9{}, s.NextChangedArea())
	assert.NoError(t, s.RemoveOption(loc, 5))
	assert.Equal(t, s.NewArea().with(loc, true), s.NextChangedArea())

	assert.Equal(t, s.NewDigits(1, 2, 3, 4, 6, 7, 8, 9), s.Get(loc))

	assert.NoError(t, s.RemoveOption(loc, 3))

	assert.Equal(t, s.NewDigits(1, 2, 4, 6, 7, 8, 9), s.Get(loc))
}

func TestMask(t *testing.T) {
	s := newSudoku[Digits9, Area9x9, grid9x9[Digits9], size9]()
	s.changeProcessors = append(s.changeProcessors, SolveProcessors[Digits9, Area9x9]{})
	s.nextChanged = Area9x9{}

	loc := CellLocation{Row: 0, Col: 0}

	assert.Equal(t, Area9x9{}, s.NextChangedArea())
	assert.NoError(t, s.Mask(loc, s.NewDigits(1, 2, 3)))
	assert.Equal(t, s.NewDigits(1, 2, 3), s.Get(loc))
	assert.Equal(t, s.NewArea().with(loc, true), s.NextChangedArea())

	assert.NoError(t, s.Mask(loc, s.NewDigits(2, 3, 4, 5)))
	assert.Equal(t, s.NewDigits(2, 3), s.Get(loc))

	assert.Equal(t, Area9x9{}, s.SolvedArea())
	assert.NoError(t, s.Mask(loc, s.NewDigits(3, 4, 5)))
	assert.Equal(t, s.NewDigits(3), s.Get(loc))
	assert.Equal(t, s.NewArea().with(loc, true), s.SolvedArea())

	assert.Error(t, s.Mask(loc, s.NewDigits(1, 4, 5)))
	assert.Equal(t, s.NewDigits(3), s.Get(loc))
}

func TestRemoveMask(t *testing.T) {
	s := newSudoku[Digits9, Area9x9, grid9x9[Digits9], size9]()
	s.nextChanged = Area9x9{}

	loc := CellLocation{Row: 0, Col: 0}

	assert.NoError(t, s.RemoveMask(loc, s.NewDigits(1, 2)))
	assert.Equal(t, s.NewDigits(3, 4, 5, 6, 7, 8, 9), s.Get(loc))
	assert.Equal(t, s.NewArea().with(loc, true), s.NextChangedArea())

	assert.NoError(t, s.RemoveMask(loc, s.NewDigits(8, 9)))
	assert.Equal(t, s.NewDigits(3, 4, 5, 6, 7), s.Get(loc))

	assert.Equal(t, Area9x9{}, s.SolvedArea())
	assert.NoError(t, s.RemoveMask(loc, s.NewDigits(3, 4, 6, 7)))
	assert.Equal(t, s.NewDigits(5), s.Get(loc))

	assert.Error(t, s.RemoveMask(loc, s.AllDigits()))
	assert.Equal(t, s.NewDigits(5), s.Get(loc))
}
