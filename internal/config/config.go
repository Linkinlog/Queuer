package config

import (
	"encoding/json"
	"os"
)

type Credentials struct {
	QueueDatabaseUser     string
	QueueDatabasePassword string

	TargetDatabaseUser     string
	TargetDatabasePassword string

	LogDatabaseUser     string
	LogDatabasePassword string
}

type Config struct {
	Queues []*Queue
	Creds  *Credentials
}

type Queue struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Service     string `json:"service"`
	Timeout     int    `json:"timeout"`
	Retries     int    `json:"retries"`

	QueueDatabaseHost string `json:"queueDatabaseHost"`
	QueueDatabasePort string `json:"queueDatabasePort"`
	QueueDatabaseName string `json:"queueDatabaseName"`

	TargetDatabaseHost string `json:"targetDatabaseHost"`
	TargetDatabasePort string `json:"targetDatabasePort"`
	TargetDatabaseName string `json:"targetDatabaseName"`

	LogDatabaseHost string `json:"logDatabaseHost"`
	LogDatabasePort string `json:"logDatabasePort"`
	LogDatabaseName string `json:"logDatabaseName"`
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

	cfg.Creds = &Credentials{
		QueueDatabaseUser:      queueDatabaseUser(),
		QueueDatabasePassword:  queueDatabasePassword(),
		TargetDatabaseUser:     targetDatabaseUser(),
		TargetDatabasePassword: targetDatabasePassword(),
		LogDatabaseUser:        logDatabaseUser(),
		LogDatabasePassword:    logDatabasePassword(),
	}

	return cfg, nil
}

func queueDatabaseUser() string {
	return os.Getenv("QUEUE_DATABASE_USER")
}

func queueDatabasePassword() string {
	return os.Getenv("QUEUE_DATABASE_PASSWORD")
}

func targetDatabaseUser() string {
	return os.Getenv("TARGET_DATABASE_USER")
}

func targetDatabasePassword() string {
	return os.Getenv("TARGET_DATABASE_PASSWORD")
}

func logDatabaseUser() string {
	return os.Getenv("LOG_DATABASE_USER")
}

func logDatabasePassword() string {
	return os.Getenv("LOG_DATABASE_PASSWORD")
}
