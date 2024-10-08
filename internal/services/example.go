package services

import "strconv"

func NewExampleService(factor int, base int) *ExampleService {
	return &ExampleService{
		factor: factor,
		base:   base,
	}
}

type ExampleService struct {
	factor int
	base   int
}

func (s *ExampleService) Read(p []byte) (n int, err error) {
	solution := s.factor * s.base
	return copy(p, []byte(strconv.Itoa(solution))), nil
}

func (s *ExampleService) String() string {
	return "example"
}
