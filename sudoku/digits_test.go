package sudoku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDigits_16_Min(t *testing.T) {
	d := Digits9(0)
	for i := 1; i <= 9; i++ {
		d = d.With(i)
		assert.Equal(t, d.Min(), 1)
	}

	d = Digits9(0)
	for i := 9; i >= 1; i-- {
		d = d.With(i)
		assert.Equal(t, d.Min(), i)
	}
}

func TestDigits_16_Max(t *testing.T) {
	d := Digits9(0)
	for i := 1; i <= 9; i++ {
		d = d.With(i)
		assert.Equal(t, d.Max(), i)
	}

	d = Digits9(0)
	for i := 9; i >= 1; i-- {
		d = d.With(i)
		assert.Equal(t, d.Max(), 9)
	}
}
