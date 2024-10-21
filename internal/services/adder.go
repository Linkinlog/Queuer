package services

import (
	"context"
	"encoding/json"
)

func NewAdder() *Adder {
	return &Adder{}
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

func (a *Adder) Run(ctx context.Context) (result []byte, err error) {
	sum := 0

	for _, addend := range a.Addends {
		sum += addend
	}

	return json.Marshal(sum)
}
