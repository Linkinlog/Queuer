package internal_test

import (
	"log/slog"
	"testing"

	"github.com/linkinlog/queuer/internal"
	"github.com/linkinlog/queuer/internal/config"
)

func TestStart(t *testing.T) {
	cfg := &config.Config{
		Queues: []*config.Queue{
			{
				Name:        "Addition Service",
				Environment: "development",
				Service:     "adder",
				Timeout:     1000,
				Retries:     3,

				QueueDatabaseHost: "queue",
				QueueDatabasePort: "5432",
				QueueDatabaseName: "queue",

				TargetDatabaseHost: "transaction",
				TargetDatabasePort: "5432",
				TargetDatabaseName: "target",

				LogDatabaseHost: "logs",
				LogDatabasePort: "5432",
				LogDatabaseName: "logs",
			},
			{
				Name:        "Square Service",
				Environment: "prod",
				Service:     "squarer",
				Timeout:     1000,
				Retries:     3,

				QueueDatabaseHost: "queue",
				QueueDatabasePort: "5432",
				QueueDatabaseName: "queue",

				TargetDatabaseHost: "transaction",
				TargetDatabasePort: "5432",
				TargetDatabaseName: "target",

				LogDatabaseHost: "logs",
				LogDatabasePort: "5432",
				LogDatabaseName: "logs",
			},
			{
				Name:        "Long Running Service",
				Environment: "prod",
				Service:     "longrunner",
				Timeout:     1000,
				Retries:     1,

				QueueDatabaseHost: "queue",
				QueueDatabasePort: "5432",
				QueueDatabaseName: "queue",

				TargetDatabaseHost: "transaction",
				TargetDatabasePort: "5432",
				TargetDatabaseName: "target",

				LogDatabaseHost: "logs",
				LogDatabasePort: "5432",
				LogDatabaseName: "logs",
			},
		},
		SlogOpts: &slog.HandlerOptions{},
	}

	internal.Start(cfg)
}
