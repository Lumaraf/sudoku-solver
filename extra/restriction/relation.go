package restriction

import (
	"errors"
	"github.com/lumaraf/sudoku-solver/sudoku"
)

type RelationRestriction struct {
	Area    sudoku.Area
	Offsets sudoku.Offsets
	Masks   map[int]sudoku.Digits
}

func (r RelationRestriction) Name() string {
	return "Relation"
}

func (r RelationRestriction) Validate(s sudoku.Sudoku) error {
outer:
	for _, cell := range s.SolvedArea().And(r.Area).Locations {
		v, _ := s.Get(cell).Single()

		mask, found := r.Masks[v]
		if !found {
			continue
		}
		for offsetCell := range r.Offsets.Locations(cell) {
			if s.Get(offsetCell)&mask != 0 {
				continue outer
			}
		}
		return errors.New("Relation rule is not satisfied")
	}
	return nil
}
