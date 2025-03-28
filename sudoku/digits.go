package sudoku

import (
	"errors"
	"math/bits"
)

const AllDigits = Digits(0b111111111)

var ErrValueOutOfRange = errors.New("value out of range")

type Digits uint16

type Values func(func(int) bool)

func NewDigits(values ...int) Digits {
	d := Digits(0)
	for _, v := range values {
		d.AddOption(v)
	}
	return d
}

func (c Digits) CanContain(v int) bool {
	return uint16(c)&c.getBit(v) != 0
}

func (c *Digits) AddOption(v int) {
	*c = *c | Digits(c.getBit(v))
}

func (c *Digits) RemoveOption(v int) {
	*c = *c & (^Digits(c.getBit(v)))
}

func (c *Digits) ForceValue(v int) {
	*c = Digits(c.getBit(v))
}

func (c Digits) Count() int {
	return bits.OnesCount16(uint16(c))
}

func (c Digits) Single() (v int, isSingle bool) {
	isSingle = c.Count() == 1
	if isSingle {
		v = bits.TrailingZeros16(uint16(c)) + 1
	}
	return
}

func (c Digits) getBit(v int) uint16 {
	return 1 << (v - 1)
}

func (c Digits) Values(yield func(int) bool) {
	mask := uint16(c)
	for mask != 0 {
		lz := bits.TrailingZeros16(mask)
		mask = mask & ^(1 << lz)
		if !yield(int(lz) + 1) {
			return
		}
	}
}

func (c Digits) Min() int {
	return bits.TrailingZeros16(uint16(c)) + 1
}

func (c Digits) Max() int {
	return 16 - bits.LeadingZeros16(uint16(c))
}

func checkValue(v int) error {
	if v < 1 || v > 9 {
		return ErrValueOutOfRange
	}
	return nil
}
