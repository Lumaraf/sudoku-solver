package sudoku

import (
	"errors"
	"math/bits"
)

//const AllDigits = digits_16(0b111111111)

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

type digits_16 uint16

func (c digits_16) CanContain(v int) bool {
	return uint16(c)&c.getBit(v) != 0
}

func (c digits_16) withOption(v int) digits_16 {
	return c | digits_16(c.getBit(v))
}

func (c digits_16) Empty() bool {
	return c == 0
}

func (c digits_16) Count() int {
	return bits.OnesCount16(uint16(c))
}

func (c digits_16) Single() (v int, isSingle bool) {
	isSingle = c.Count() == 1
	if isSingle {
		v = bits.TrailingZeros16(uint16(c)) + 1
	}
	return
}

func (c digits_16) getBit(v int) uint16 {
	return 1 << (v - 1)
}

func (c digits_16) Values(yield func(int) bool) {
	mask := uint16(c)
	for mask != 0 {
		lz := bits.TrailingZeros16(mask)
		mask = mask & ^(1 << lz)
		if !yield(int(lz) + 1) {
			return
		}
	}
}

func (c digits_16) Min() int {
	return bits.TrailingZeros16(uint16(c)) + 1
}

func (c digits_16) Max() int {
	return 64 - bits.LeadingZeros16(uint16(c))
}

func (c digits_16) and(d digits_16) digits_16 {
	return c & d
}
func (c digits_16) or(d digits_16) digits_16 {
	return c | d
}
func (c digits_16) not() digits_16 {
	return ^c
}

func (c digits_16) String() string {
	str := ""
	for i := 1; i <= 16; i++ {
		if c.CanContain(i) {
			str += string(rune(i + '0'))
			str += ","
		}
	}
	return str
}

type digitsOps_16 struct{}

func (digitsOps_16) IntersectDigits(d1, d2 digits_16) digits_16 {
	return d1.and(d2)
}

func (digitsOps_16) UnionDigits(d1 digits_16, d2 digits_16) digits_16 {
	return d1.or(d2)
}

func (digitsOps_16) InvertDigits(d digits_16) digits_16 {
	return d.not()
}
