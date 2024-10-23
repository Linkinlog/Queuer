package services

import (
	"encoding/json"
	"time"
)

func NewLongRunner() *LongRunner {
	return &LongRunner{
		TimeToRun: 5000, // TODO remove this once we are feeding from DB
	}
}

type LongRunner struct {
	TimeToRun int `json:"time_to_run"`
}

func (lr *LongRunner) UnmarshalJSON(data []byte) error {
	temp := &LongRunner{}

	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	return nil
}

func (lr *LongRunner) String() string {
	return "longrunner"
}

func (lr *LongRunner) Run() (chan []byte, chan error) {
	results := make(chan []byte)
	errs := make(chan error)
	go func() {
		<-time.After(time.Duration(lr.TimeToRun) * time.Millisecond)
		result, err := json.Marshal("longrunner done")
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	return results, errs
}
