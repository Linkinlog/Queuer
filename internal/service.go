package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/linkinlog/queuer/internal/config"
	"github.com/linkinlog/queuer/internal/services"
)

type Service interface {
	json.Unmarshaler
	fmt.Stringer
	Run(ctx context.Context) (result []byte, err error)
}

func Start(logger *slog.Logger, configPath string) {
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		logger.Error("failed to parse config", "error", err)
		return
	}

	for _, queue := range cfg.Queues {
		// todo handle context timeout and cancel
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(queue.Timeout) * time.Millisecond)
		defer cancel()

		if err := process(queue, logger, ctx); err != nil {
			logger.Error("failed to process queue", "queue", queue.Name, "error", err)
		}

		if ctx.Done() != nil {
			fmt.Println("TODO")
			logger.Error("queue processing timeout", "queue", queue.Name)
		}
	}
}

func process(queue *config.Queue, logger *slog.Logger, ctx context.Context) error {
	logger.Info("fetching queue data", "queue", queue.Name)
	// todo fetch from databases
	srv := ToService(queue.Service)
	if srv == nil {
		return fmt.Errorf("unknown service: %s", queue.Service)
	}
	logger.Info("processing queue", "queue", queue.Name)

	// todo get data from queue, marshal to service unmarhsaler, run service

	result, err := srv.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to process queue: %w", err)
	}

	logger.Info("queue processed", "queue", queue.Name, "result", string(result))
	// todo write to databases

	return nil
}

func ToService(s string) Service {
	switch s {
	case "squarer":
		return services.NewSquarer()
	case "adder":
		return services.NewAdder()
	default:
		return nil
	}
}
