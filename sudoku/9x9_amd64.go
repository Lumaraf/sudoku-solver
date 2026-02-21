package sudoku

import (
	"simd/archsimd"
)

func newSudokuBuilder9x9_simd() SudokuBuilder[Digits9, Area9x9] {
	if archsimd.X86.AVX2() {
		return newSudokuBuilder[Digits9, Area9x9, grid9x9, size9, gridOps9x9_avx2]()
	}
	return nil
}

type gridOps9x9_avx2 struct {
	genericGridOps[Digits9, Area9x9, grid9x9, size9]
}

func (o gridOps9x9_avx2) PossibleLocations(g grid9x9, d Digits9) (a Area9x9) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			cell := g[row*9+col]
			if !cell.And(d).Empty() {
				a = a.With(CellLocation{row, col})
			}
		}
	}
	return
}
