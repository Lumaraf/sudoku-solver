package sudoku

import (
	"context"
	"errors"
)

type MultiSudokuBuilder[D Digits, A Area] struct {
	sudokus  []Sudoku[D, A]
	builders []SudokuBuilder[D, A]
}

func (b *MultiSudokuBuilder[D, A]) ensureAdded(sb SudokuBuilder[D, A]) {
	s := sb.buildTarget()
	for _, item := range b.sudokus {
		if item == s {
			return
		}
	}
	b.sudokus = append(b.sudokus, s)
	b.builders = append(b.builders, sb)
}

func (b *MultiSudokuBuilder[D, A]) Overlap(sb1 SudokuBuilder[D, A], corner CellLocation, w, h int, sb2 SudokuBuilder[D, A]) error {
	b.ensureAdded(sb1)
	b.ensureAdded(sb2)

	offset := Offset{}
	if corner.Col == 0 {
		offset.Col = sb1.Size() - w
	} else if corner.Col == sb1.Size()-1 {
		offset.Col = w - sb1.Size()
	} else {
		return errors.New("invalid corner")
	}

	if corner.Row == 0 {
		offset.Row = sb1.Size() - h
	} else if corner.Row == sb1.Size()-1 {
		offset.Row = h - sb1.Size()
	} else {
		return errors.New("invalid corner")
	}

	sb1.AddChangeProcessor(overlapChangeProcessor[D, A]{
		targetSudoku: sb2.buildTarget(),
		offset:       offset,
	})
	sb2.AddChangeProcessor(overlapChangeProcessor[D, A]{
		targetSudoku: sb1.buildTarget(),
		offset: Offset{
			Row: -offset.Row,
			Col: -offset.Col,
		},
	})

	return nil
}

type MultiSudoku[D Digits, A Area] struct {
	sudokus []Sudoku[D, A]
}

func (b *MultiSudokuBuilder[D, A]) Build() (MultiSudoku[D, A], error) {
	for idx, sb := range b.builders {
		s, err := sb.Build()
		if err != nil {
			return MultiSudoku[D, A]{}, err
		}
		if s != b.sudokus[idx] {
			panic("invalid sudoku")
		}
	}
	return MultiSudoku[D, A]{
		sudokus: b.sudokus,
	}, nil
}

func (ms *MultiSudoku[D, A]) Solve(ctx context.Context, factories StrategyFactories[D, A]) error {
	// TODO repeat until there is no more progress
	for _, s := range ms.sudokus {
		slv := s.NewSolver()
		slv.Use(factories...)
		slv.Solve(ctx)
	}
	for _, s := range ms.sudokus {
		s.Print()
	}
	return nil
}

type overlapChangeProcessor[D Digits, A Area] struct {
	targetSudoku Sudoku[D, A]
	offset       Offset
}

func (o overlapChangeProcessor[D, A]) Name() string {
	return "OverlapChangeProcessor"
}

func (o overlapChangeProcessor[D, A]) ProcessChange(s Sudoku[D, A], cell CellLocation, mask D) error {
	targetCell := CellLocation{
		Row: cell.Row + o.offset.Row,
		Col: cell.Col + o.offset.Col,
	}
	if targetCell.Row < 0 || targetCell.Row >= s.Size() {
		return nil
	}
	if targetCell.Col < 0 || targetCell.Col >= s.Size() {
		return nil
	}
	return o.targetSudoku.Mask(targetCell, mask)
}
