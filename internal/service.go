package internal

import (
	"fmt"
	"log/slog"

	"github.com/linkinlog/queuer/internal/config"
)

func Start(logger *slog.Logger, configPath string) {
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		logger.Error("failed to parse config", "error", err)
		return
	}

	for _, queue := range cfg.Queues {
		fmt.Println(queue)
		logger.Info("processing queue", "queue", queue.Name)
	}
}
