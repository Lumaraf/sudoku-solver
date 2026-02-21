package sudoku

type Validator[D Digits[D], A Area[A]] interface {
	Name() string
	Validate(s Sudoku[D, A]) error
}

func (s *sudoku[D, A, G, S, GO]) Validate() error {
	for _, v := range s.validators {
		if err := v.Validate(s); err != nil {
			return err
		}
	}
	return nil
}
