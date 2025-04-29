package extra

import (
	"github.com/lumaraf/sudoku-solver/test"
	"testing"
)

func TestNonConsecutive(t *testing.T) {
	test.SudokuTests{
		"one": {
			Rows: []string{
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
			NonConsecutive: true,
		},
	}.Run(t)
}
