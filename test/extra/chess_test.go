package extra

import (
	"github.com/lumaraf/sudoku-solver/test"
	"testing"
)

func TestChess(t *testing.T) {
	test.SudokuTests{
		"anti knight": {
			Rows: []string{
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
			AntiKnight: true,
		},
		"miracle": {
			Rows: []string{
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
			AntiKing:       true,
			AntiKnight:     true,
			NonConsecutive: true,
		},
		"159": {
			Rows: []string{
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
			AntiKnight: true,
			Rule159:    true,
		},
	}.Run(t)
}
