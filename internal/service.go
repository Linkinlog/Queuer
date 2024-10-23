package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/linkinlog/queuer/internal/config"
	"github.com/linkinlog/queuer/internal/services"
)

type Service interface {
	json.Unmarshaler
	fmt.Stringer
	Run() (results chan []byte, errs chan error)
}

func Start(logger *slog.Logger, configPath string) {
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		logger.Error("failed to parse config", "error", err)
		return
	}

	var wg sync.WaitGroup
	for _, queue := range cfg.Queues {
		wg.Add(1)
		go processQueue(logger, queue, &wg)
	}

	wg.Wait()
}

func processQueue(logger *slog.Logger, queue *config.Queue, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info("fetching queue data", "queue", queue.Name)

	// todo fetch from databases
	srv := ToService(queue.Service)
	if srv == nil {
		logger.Error("unknown service", "service", queue.Service)
		return
	}
	logger.Info("processing queue entry", "queue", queue.Name, "timeout", queue.Timeout, "service", srv.String())

	if res, err := processItem(srv, queue.Timeout); err != nil {
		logger.Error("failed to process service", "service", srv.String(), "error", err)
	} else if res != nil {
		logger.Info("service processed", "service", srv.String(), "result", string(res))
	}

	// todo write to databases
}

func processItem(srv Service, timeout int) (result []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	resChan, errChan := srv.Run()
	select {
	case res := <-resChan:
		return res, nil
	case err := <-errChan:
		return nil, fmt.Errorf("service failed: %w", err)
	case <-ctx.Done():
		return nil, fmt.Errorf("service timed out")
	}
}

func ToService(s string) Service {
	switch s {
	case "squarer":
		return services.NewSquarer()
	case "adder":
		return services.NewAdder()
	case "longrunner":
		return services.NewLongRunner()
	default:
		return nil
	}
}
