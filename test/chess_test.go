package test

import (
	"testing"
)

func TestChess(t *testing.T) {
	sudokuTests{
		"anti knight": {
			rows: []string{
				" 5   9   ",
				"8        ",
				"     3 4 ",
				"7 8   1 9",
				"         ",
				"    3    ",
				"         ",
				"  3 1   8",
				"   9   2 ",
			},
			anitKnight: true,
		},
		"miracle": {
			rows: []string{
				"         ",
				"         ",
				"         ",
				"         ",
				"  1      ",
				"      2  ",
				"         ",
				"         ",
				"         ",
			},
			anitKing:       true,
			anitKnight:     true,
			nonConsecutive: true,
		},
		"159": {
			rows: []string{
				"         ",
				"         ",
				"    E    ",
				"E   E    ",
				"E   E    ",
				"E   E    ",
				"E        ",
				"         ",
				"         ",
			},
			anitKnight: true,
			rule159:    true,
		},
	}.Run(t)
}
