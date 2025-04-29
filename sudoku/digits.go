package sudoku

import (
	"errors"
	"math/bits"
)

//const AllDigits = digits16(0b111111111)

var ErrValueOutOfRange = errors.New("value out of range")

type Digits interface {
	comparable

	CanContain(v int) bool
	Empty() bool
	Count() int
	Single() (v int, isSingle bool)
	Min() int
	Max() int
	Values(func(int) bool)
}

type Values func(func(int) bool)

type Digits4 = digits16
type Digits6 = digits16
type Digits9 = digits16
type Digits16 = digits16

//type Digits25 = digits32

type digits16 uint16

func (c digits16) CanContain(v int) bool {
	return uint16(c)&c.getBit(v) != 0
}

func (c digits16) withOption(v int) digits16 {
	return c | digits16(c.getBit(v))
}

func (c digits16) Empty() bool {
	return c == 0
}

func (c digits16) Count() int {
	return bits.OnesCount16(uint16(c))
}

func (c digits16) Single() (v int, isSingle bool) {
	isSingle = c.Count() == 1
	if isSingle {
		v = bits.TrailingZeros16(uint16(c)) + 1
	}
	return
}

func (c digits16) getBit(v int) uint16 {
	return 1 << (v - 1)
}

func (c digits16) Values(yield func(int) bool) {
	mask := uint16(c)
	for mask != 0 {
		lz := bits.TrailingZeros16(mask)
		mask = mask & ^(1 << lz)
		if !yield(int(lz) + 1) {
			return
		}
	}
}

func (c digits16) Min() int {
	return bits.TrailingZeros16(uint16(c)) + 1
}

func (c digits16) Max() int {
	return 64 - bits.LeadingZeros16(uint16(c))
}
