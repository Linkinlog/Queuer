package services

import (
	"encoding/json"
	"errors"
)

func NewAdder() *Adder {
	return &Adder{}
}

type Adder struct {
	Addends []int `json:"addends"`
}

func (a *Adder) UnmarshalJSON(data []byte) error {
	temp := struct {
		Addends []*int `json:"addends"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	for _, addend := range temp.Addends {
		if addend == nil {
			return errors.New("addend is required")
		}
		a.Addends = append(a.Addends, *addend)
	}

	return nil
}

func (a *Adder) String() string {
	return "adder"
}

func (a *Adder) Run() (chan []byte, chan error) {
	results := make(chan []byte)
	errs := make(chan error)
	go func() {
		sum := 0

		for _, addend := range a.Addends {
			sum += addend
		}

		result, err := json.Marshal(sum)
		if err != nil {
			errs <- err
			return
		}

		results <- result
	}()

	return results, errs
}
