package services

import (
	"encoding/json"
)

func NewAdder() *Adder {
	return &Adder{
		Addends: []int{1, 3}, // TODO remove this once we are feeding from DB
	}
}

type Adder struct {
	Addends []int `json:"addends"`
}

func (a *Adder) UnmarshalJSON(data []byte) error {
	temp := &Adder{}

	if err := json.Unmarshal(data, &temp.Addends); err != nil {
		return err
	}

	a.Addends = temp.Addends

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
