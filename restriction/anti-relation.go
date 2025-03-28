package restriction

import (
	"errors"

	"github.com/lumaraf/sudoku-solver/sudoku"
)

var errBrokenRelation = errors.New("broken relation")

type AntiRelationRestriction struct {
	Offsets    sudoku.Offsets
	Masks      map[int]sudoku.Digits
	Exceptions map[sudoku.Area]bool
}

func (r AntiRelationRestriction) Name() string {
	return "AntiRelation"
}

func (r AntiRelationRestriction) Validate(s sudoku.Sudoku) error {
	for _, cell := range s.SolvedArea().Locations {
		v, _ := s.Get(cell).Single()

		mask, found := r.Masks[v]
		if !found {
			continue
		}
		mask = ^mask
		for offsetCell := range r.Offsets.Locations(cell) {
			if _, found := r.Exceptions[sudoku.NewArea(cell, offsetCell)]; found {
				continue
			}
			if s.Get(offsetCell)&mask == 0 {
				return errBrokenRelation
			}
		}
	}
	return nil
}

func (r AntiRelationRestriction) ProcessSolvea(s sudoku.Sudoku, cell sudoku.CellLocation) error {
	v, _ := s.Get(cell).Single()
	mask, found := r.Masks[v]
	if !found {
		return nil
	}
	for offsetCell := range r.Offsets.Locations(cell) {
		if _, found := r.Exceptions[sudoku.NewArea(cell, offsetCell)]; found {
			continue
		}
		if err := s.RemoveMask(offsetCell, mask); err != nil {
			return err
		}
	}
	return nil
}
