package services

import (
	"encoding/json"
	"errors"
)

func NewSquarer() *Squarer {
	return &Squarer{}
}

type Squarer struct {
	Factor int `json:"factor"`
	Base   int `json:"base"`
}

func (s *Squarer) UnmarshalJSON(data []byte) error {
	temp := struct {
		Factor *int `json:"factor"`
		Base   *int `json:"base"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.Factor == nil {
		return errors.New("factor is required")
	}
	if temp.Base == nil {
		return errors.New("base is required")
	}

	s.Factor = *temp.Factor
	s.Base = *temp.Base

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
