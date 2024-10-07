package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Queues []*Queue
}

type Queue struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Service     string `json:"service"`

	QueueDatabaseHost     string `json:"queueDatabaseHost"`
	QueueDatabasePort     string `json:"queueDatabasePort"`
	QueueDatabaseUser     string `json:"queueDatabaseUser"`
	QueueDatabasePassword string `json:"queueDatabasePassword"`
	QueueDatabaseTable    string `json:"queueDatabaseTable"`
	QueueDatabaseDriver   string `json:"queueDatabaseDriver"`

	TargetDatabaseHost     string `json:"targetDatabaseHost"`
	TargetDatabasePort     string `json:"targetDatabasePort"`
	TargetDatabaseUser     string `json:"targetDatabaseUser"`
	TargetDatabasePassword string `json:"targetDatabasePassword"`
	TargetDatabaseTable    string `json:"targetDatabaseTable"`
	TargetDatabaseDriver   string `json:"targetDatabaseDriver"`
}

func ParseConfig(filename string) (*Config, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	err = json.NewDecoder(file).Decode(&cfg.Queues)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
