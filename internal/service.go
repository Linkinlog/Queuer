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

	params := &dbParams{
		dsn: fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			queue.QueueDatabaseUser,
			queue.QueueDatabasePassword,
			queue.QueueDatabaseHost,
			queue.QueueDatabasePort,
			queue.QueueDatabaseName,
		),
	}

	queueReader, err := OpenQueue(params)
	if err != nil {
		logger.Error("failed to open queue", "error", err)
		return
	}

	events, err := queueReader.Read(queue.Service)
	if err != nil {
		logger.Error("failed to read queue", "error", err)
		return
	}

	if len(events) == 0 {
		logger.Info("no events to process", "queue", queue.Name)
		return
	}

	for _, event := range events {
		srv := ToService(queue.Service)
		if srv == nil {
			logger.Error("unknown service", "service", queue.Service)
			return
		}

		logger.Info("processing queue entry", "queue", queue.Name, "timeout", queue.Timeout, "service", srv.String(), "data", event.Data)

		if err := srv.UnmarshalJSON([]byte(event.Data)); err != nil {
			logger.Error("failed to unmarshal data", "data", event.Data, "error", err)
			continue
		}

		for i := 0; i <= queue.Retries; i++ {
			res, err := processItem(srv, queue.Timeout)
			if err != nil {
				logger.Error("failed to process service", "name", queue.Name, "service", srv.String(), "error", err)
				continue
			}
			if res == nil {
				logger.Error("service returned nil result", "service", srv.String())
				continue
			}

			logger.Info("service processed", "service", srv.String(), "result", string(res))

			if err := queueReader.MarkProcessed(event.EventID); err != nil {
				logger.Error("failed to mark event as processed", "eventID", event.EventID, "error", err)
				continue
			}

			break
		}
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
