package sudoku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func runAreaTests[A Area[A]](t *testing.T) {
	t.Parallel()

	t.Run("With and Without", func(t *testing.T) {
		t.Parallel()

		var a A
		assert.True(t, a.Empty())
		a = a.With(CellLocation{Row: 0, Col: 0})
		assert.Equal(t, 1, a.Count())
		a = a.With(CellLocation{Row: 0, Col: 0})
		assert.Equal(t, 1, a.Count())
		a = a.With(CellLocation{Row: 1, Col: 1})
		assert.Equal(t, 2, a.Count())
		a = a.Without(CellLocation{Row: 1, Col: 1})
		assert.Equal(t, 1, a.Count())
	})

	t.Run("And", func(t *testing.T) {
		t.Parallel()

		var a, b A
		a = a.With(CellLocation{Row: 0, Col: 0})
		a = a.With(CellLocation{Row: 1, Col: 1})
		b = b.With(CellLocation{Row: 1, Col: 1})
		b = b.With(CellLocation{Row: 2, Col: 2})

		c := a.And(b)
		assert.Equal(t, 1, c.Count())
		assert.True(t, c.Get(CellLocation{Row: 1, Col: 1}))
	})

	t.Run("Or", func(t *testing.T) {
		t.Parallel()

		var a, b A
		a = a.With(CellLocation{Row: 0, Col: 0})
		a = a.With(CellLocation{Row: 1, Col: 1})
		b = b.With(CellLocation{Row: 1, Col: 1})
		b = b.With(CellLocation{Row: 2, Col: 2})

		c := a.Or(b)
		assert.Equal(t, 3, c.Count())
		assert.True(t, c.Get(CellLocation{Row: 0, Col: 0}))
		assert.True(t, c.Get(CellLocation{Row: 1, Col: 1}))
		assert.True(t, c.Get(CellLocation{Row: 2, Col: 2}))
	})

	t.Run("Not", func(t *testing.T) {
		t.Parallel()

		var a A
		a = a.With(CellLocation{Row: 0, Col: 0})

		b := a.Not()
		assert.False(t, b.Get(CellLocation{Row: 0, Col: 0}))
	})

	t.Run("ShiftLeft", func(t *testing.T) {
		t.Parallel()

		var a A
		for n := 0; n < a.Size(); n++ {
			a = a.With(CellLocation{Row: n, Col: n})
		}

		for s := 0; s <= a.Size(); s++ {
			b := a.ShiftLeft(s)
			assert.Equal(t, b.Size()-s, b.Count())
		}
	})

	t.Run("ShiftRight", func(t *testing.T) {
		t.Parallel()

		var a A
		for n := 0; n < a.Size(); n++ {
			a = a.With(CellLocation{Row: n, Col: n})
		}

		for s := 0; s <= a.Size(); s++ {
			b := a.ShiftRight(s)
			assert.Equal(t, b.Size()-s, b.Count())
		}
	})

	t.Run("ShiftUp", func(t *testing.T) {
		t.Parallel()

		var a A
		for n := 0; n < a.Size(); n++ {
			a = a.With(CellLocation{Row: n, Col: n})
		}

		for s := 0; s <= a.Size(); s++ {
			b := a.ShiftUp(s)
			assert.Equal(t, b.Size()-s, b.Count())
		}
	})

	t.Run("ShiftDown", func(t *testing.T) {
		t.Parallel()

		var a A
		for n := 0; n < a.Size(); n++ {
			a = a.With(CellLocation{Row: n, Col: n})
		}

		for s := 0; s <= a.Size(); s++ {
			b := a.ShiftDown(s)
			assert.Equal(t, b.Size()-s, b.Count())
		}
	})
}

func TestArea6x6(t *testing.T) {
	runAreaTests[Area6x6](t)
}

func TestArea9x9(t *testing.T) {
	runAreaTests[Area9x9](t)
}

func TestArea12x12(t *testing.T) {
	runAreaTests[Area12x12](t)
}
