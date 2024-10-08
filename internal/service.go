package internal

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/linkinlog/queuer/internal/config"
	"github.com/linkinlog/queuer/internal/services"
)

type Service interface {
	io.Reader
	fmt.Stringer
}

func Start(logger *slog.Logger, configPath string) {
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		logger.Error("failed to parse config", "error", err)
		return
	}

	for _, queue := range cfg.Queues {
		logger.Info("fetching queue data", "queue", queue.Name)
		srv := ToService(queue.Service)
		if srv == nil {
			logger.Error("unknown service", "service", queue.Service)
			continue
		}
		logger.Info("processing queue", "queue", queue.Name)
	}
}

func ToService(s string) Service {
	switch s {
	case "example":
		return services.NewExampleService(2, 3)
	case "example2":
		return services.NewExampleService2(2, 3)
	default:
		return nil
	}
}
