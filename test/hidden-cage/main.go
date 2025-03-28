package main

import (
	"fmt"
	"github.com/lumaraf/sudoku-solver/sudoku"
	"time"
)

func main() {
	areas := make([]sudoku.Area, 0, 27)
	for i := 0; i < 9; i++ {
		areas = append(areas, sudoku.RowArea(i))
		areas = append(areas, sudoku.ColArea(i))
		areas = append(areas, sudoku.BoxArea(i))
	}

	start := time.Now()
	mergedAreas := make(map[sudoku.Area]bool)
	for area := range iterateAreas(areas) {
		mergedAreas[area] = true
	}
	fmt.Println(len(mergedAreas))
	fmt.Println(time.Since(start))
}

func iterateAreas(areas []sudoku.Area) func(yield func(a sudoku.Area) bool) {
	return func(yield func(a sudoku.Area) bool) {
		for i, area := range areas {
			if !yield(area) {
				return
			}
			for other := range iterateMergedAreas(area, areas[i+1:]) {
				if !yield(other) {
					return
				}
			}
		}
	}
}

func iterateMergedAreas(baseArea sudoku.Area, areas []sudoku.Area) func(yield func(a sudoku.Area) bool) {
	return func(yield func(a sudoku.Area) bool) {
		for i, area := range areas {
			if area.And(baseArea).Size() != 0 {
				continue
			}
			merged := baseArea.Or(area)
			if !yield(merged) {
				return
			}
			for other := range iterateMergedAreas(merged, areas[i+1:]) {
				if !yield(other) {
					return
				}
			}
		}
	}
}
