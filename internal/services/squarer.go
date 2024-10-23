package services

import (
	"encoding/json"
)

func NewSquarer() *Squarer {
	return &Squarer{
		// TODO remove this once we are feeding from DB
		Factor: 2,
		Base:   3,
	}
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

func (s *Squarer) Run() (chan []byte, chan error) {
	results := make(chan []byte)
	errs := make(chan error)

	go func() {
		res := s.Factor * s.Base

		result, err := json.Marshal(res)
		if err != nil {
			errs <- err
			return
		}

		results <- result
	}()

	return results, errs
}
