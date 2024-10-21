package services

import (
	"context"
	"encoding/json"
	"strconv"
)

func NewSquarer() *Squarer {
	return &Squarer{}
}

type Squarer struct {
	Factor int `json:"factor"`
	Base   int `json:"base"`
}

func (s *Squarer) UnmarshalJSON(data []byte) error {
	temp := &Squarer{}

	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	s.Factor = temp.Factor
	s.Base = temp.Base

	return nil
}

func (s *Squarer) String() string {
	return "squarer"
}

func (s *Squarer) Run(ctx context.Context) (result []byte, err error) {
	resInt := s.Factor * s.Base
	resStr := strconv.Itoa(resInt)
	res := []byte(resStr)

	return res, nil
}
