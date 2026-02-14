package sudoku

import (
	"errors"
	"math/bits"
)

var ErrValueOutOfRange = errors.New("value out of range")

type Digits[D Digits[D]] interface {
	comparable

	And(other D) D
	Or(other D) D
	Not() D

	CanContain(v int) bool
	Empty() bool
	Count() int
	Single() (v int, isSingle bool)
	Min() int
	Max() int
	Values(func(int) bool)
}

type Values func(func(int) bool)

type digits_16 struct {
	v uint16
}

func (c digits_16) CanContain(v int) bool {
	return c.v&c.getBit(v) != 0
}

func (c digits_16) withOption(v int) digits_16 {
	return digits_16{v: c.v | c.getBit(v)}
}

func (c digits_16) Empty() bool {
	return c.v == 0
}

func (c digits_16) Count() int {
	return bits.OnesCount16(c.v)
}

func (c digits_16) Single() (v int, isSingle bool) {
	isSingle = c.Count() == 1
	if isSingle {
		v = bits.TrailingZeros16(c.v) + 1
	}
	return
}

func (c digits_16) getBit(v int) uint16 {
	return 1 << (v - 1)
}

func (c digits_16) Values(yield func(int) bool) {
	mask := c.v
	for mask != 0 {
		lz := bits.TrailingZeros16(mask)
		mask = mask & ^(1 << lz)
		if !yield(int(lz) + 1) {
			return
		}
	}
}

func (c digits_16) Min() int {
	return bits.TrailingZeros16(c.v) + 1
}

func (c digits_16) Max() int {
	return 64 - bits.LeadingZeros16(c.v)
}

func (c digits_16) And(d digits_16) digits_16 {
	return digits_16{v: c.v & d.v}
}
func (c digits_16) Or(d digits_16) digits_16 {
	return digits_16{v: c.v | d.v}
}
func (c digits_16) Not() digits_16 {
	return digits_16{v: ^c.v}
}

func (c digits_16) String() string {
	str := ""
	for i := 1; i <= 16; i++ {
		if c.CanContain(i) {
			if str != "" {
				str += ","
			}
			str += string(rune(i + '0'))
		}
	}
	return str
}

type digitsOps_16 struct{}

func (digitsOps_16) IntersectDigits(d1, d2 digits_16) digits_16 {
	return d1.And(d2)
}

func (digitsOps_16) UnionDigits(d1 digits_16, d2 digits_16) digits_16 {
	return d1.Or(d2)
}

func (digitsOps_16) InvertDigits(d digits_16) digits_16 {
	return d.Not()
}
