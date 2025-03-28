package sudoku

type Restriction interface {
	Name() string
	Validate(s Sudoku) error
}

func (s *sudoku) AddRestriction(r Restriction) {
	s.restrictions = append(s.restrictions, r)

	if ea, ok := r.(interface{ ExclusionAreas() map[CellLocation]Area }); ok {
		for l, a := range ea.ExclusionAreas() {
			s.exlusionAreas[l.Row][l.Col] = s.exlusionAreas[l.Row][l.Col].Or(a)
		}
	}

	if sp, ok := r.(SolveProcessor); ok {
		s.solveProcessors = append(s.solveProcessors, sp)
	}
}

type HiddenRestrictionDetector interface {
	Name() string
	FindHiddenRestrictions([]Restriction) []Restriction
}

var hiddenRestrictionDetectors = make([]HiddenRestrictionDetector, 0, 10)

func RegisterHiddenRestrictionDetector(d HiddenRestrictionDetector) {
	hiddenRestrictionDetectors = append(hiddenRestrictionDetectors, d)
}

func findHiddenRestrictions(restrictions []Restriction) []Restriction {
	r := make([]Restriction, 0, 10)
	for _, d := range hiddenRestrictionDetectors {
		r = append(r, d.FindHiddenRestrictions(restrictions)...)
	}
	return r
}
