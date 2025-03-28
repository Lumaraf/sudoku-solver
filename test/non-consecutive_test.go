package test

import (
	"testing"
)

func TestNonConsecutive(t *testing.T) {
	sudokuTests{
		"one": {
			rows: []string{
				"        5",
				" 1    7  ",
				"7        ",
				"    7  59",
				"         ",
				"42  9    ",
				"        8",
				"  1    7 ",
				"8        ",
			},
			nonConsecutive: true,
		},
	}.Run(t)
}
