package sudoku

import (
	"errors"
	"math/bits"
)

var ErrValueOutOfRange = errors.New("value out of range")

type DigitsOps[D Digits[D]] interface {
	NewDigits(values ...int) D
	AllDigits() D
	//IntersectDigits(d1 D, d2 D) D
	//UnionDigits(d1 D, d2 D) D
	//InvertDigits(d D) D
}

type Digits[D Digits[D]] interface {
	comparable

	And(other D) D
	Or(other D) D
	Not() D

	All() D

	With(v int) D
	Without(v int) D

	CanContain(v int) bool
	Empty() bool
	Count() int
	Single() (v int, isSingle bool)
	Min() int
	Max() int
	Values(func(int) bool)
}

type Values func(func(int) bool)

type digits_16[AD interface{ allDigits() uint16 }] uint16

func (c digits_16[AD]) And(d digits_16[AD]) digits_16[AD] {
	return c & d
}

func (c digits_16[AD]) Or(d digits_16[AD]) digits_16[AD] {
	return c | d
}

func (c digits_16[AD]) Not() digits_16[AD] {
	var ad AD
	return digits_16[AD](ad.allDigits() & ^uint16(c))
}

func (c digits_16[AD]) All() digits_16[AD] {
	var ad AD
	return digits_16[AD](ad.allDigits())
}

func (c digits_16[AD]) With(v int) digits_16[AD] {
	return digits_16[AD](uint16(c) | c.getBit(v))
}

func (c digits_16[AD]) Without(v int) digits_16[AD] {
	return digits_16[AD](uint16(c) & ^c.getBit(v))
}

func (c digits_16[AD]) CanContain(v int) bool {
	return uint16(c)&c.getBit(v) != 0
}

func (c digits_16[AD]) Empty() bool {
	return c == 0
}

func (c digits_16[AD]) Count() int {
	return bits.OnesCount16(uint16(c))
}

func (c digits_16[AD]) Single() (v int, isSingle bool) {
	isSingle = c.Count() == 1
	if isSingle {
		v = bits.TrailingZeros16(uint16(c)) + 1
	}
	return
}

func (c digits_16[AD]) getBit(v int) uint16 {
	return 1 << (v - 1)
}

func (c digits_16[AD]) Values(yield func(int) bool) {
	mask := uint16(c)
	for mask != 0 {
		lz := bits.TrailingZeros16(mask)
		mask = mask & ^(1 << lz)
		if !yield(int(lz) + 1) {
			return
		}
	}
}

func (c digits_16[AD]) Min() int {
	return bits.TrailingZeros16(uint16(c)) + 1
}

func (c digits_16[AD]) Max() int {
	return 16 - bits.LeadingZeros16(uint16(c))
}

func (c digits_16[AD]) String() string {
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
