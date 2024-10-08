package services

import "strconv"

func NewExampleService2(addendX, addendY int) *ExampleService2 {
	return &ExampleService2{
		addendX: addendX,
		addendY: addendY,
	}
}

type ExampleService2 struct {
	addendX, addendY int
}

func (s *ExampleService2) Read(p []byte) (n int, err error) {
	sum := s.addendX + s.addendY
	return copy(p, []byte(strconv.Itoa(sum))), nil
}

func (s *ExampleService2) String() string {
	return "example2"
}
