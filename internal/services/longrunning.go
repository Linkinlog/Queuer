package services

import (
	"encoding/json"
	"errors"
	"log/slog"
	"time"
)

var timesRan = 0

func NewLongRunner() *LongRunner {
	return &LongRunner{}
}

type LongRunner struct {
	TimeToRun int `json:"time_to_run"`
	s         *slog.Logger
}

func (lr *LongRunner) SetLogger(s *slog.Logger) {
	lr.s = s
}

func (lr *LongRunner) UnmarshalJSON(data []byte) error {
	lr.s.Debug("Unmarshalling data", "data", string(data))
	temp := struct {
		TimeToRun *int `json:"time_to_run"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.TimeToRun == nil {
		return errors.New("time_to_run is required")
	}

	lr.TimeToRun = *temp.TimeToRun

	lr.s.Debug("Unmarshalled data", "time_to_run", lr.TimeToRun)

	return nil
}

func (lr *LongRunner) String() string {
	return "longrunner"
}

func (lr *LongRunner) Run() (chan []byte, chan error) {
	results := make(chan []byte)
	errs := make(chan error)

	// we do this so we can simulate a service which requires a single retry
	tempTTR := lr.TimeToRun
	if timesRan < 1 {
		lr.s.Debug("LongRunner: timesRan < 1")
		tempTTR = lr.TimeToRun * 20
	}
	timesRan++

	go func() {
		<-time.After(time.Duration(tempTTR) * time.Millisecond)

		lr.s.Debug("LongRunner: done", "time_to_run", lr.TimeToRun)

		result, err := json.Marshal("longrunner done")
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	return results, errs
}
