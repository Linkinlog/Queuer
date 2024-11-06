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

// Program internals beyond here, no touchy

func processQueue(logger *slog.Logger, queue *config.Queue, wg *sync.WaitGroup) {
	defer wg.Done()

	loggerParams := &dbParams{
		dsn: fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			queue.LogDatabaseUser,
			queue.LogDatabasePassword,
			queue.LogDatabaseHost,
			queue.LogDatabasePort,
			queue.LogDatabaseName,
		),
	}

	logWriter, err := OpenLog(loggerParams)
	if err != nil {
		logger.Error("failed to open log", "error", err)
		return
	}

	logger.Info("fetching queue data", "queue", queue.Name)

	queueParams := &dbParams{
		dsn: fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			queue.QueueDatabaseUser,
			queue.QueueDatabasePassword,
			queue.QueueDatabaseHost,
			queue.QueueDatabasePort,
			queue.QueueDatabaseName,
		),
	}

	queueReader, err := OpenQueue(queueParams)
	if err != nil {
		if err := logWriter.WriteLog([]byte(fmt.Sprintf("failed to open queue: %v", err))); err != nil {
			logger.Error("failed to write to log", "error", err)
		}
		logger.Error("failed to open queue", "error", err)
		return
	}

	defer queueReader.Close()

	events, err := queueReader.Read(queue.Service)
	if err != nil {
		if err := logWriter.WriteLog([]byte(fmt.Sprintf("failed to read queue: %v", err))); err != nil {
			logger.Error("failed to write to log", "error", err)
		}
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
			if err := logWriter.WriteLog([]byte(fmt.Sprintf("unknown service: %s", queue.Service))); err != nil {
				logger.Error("failed to write to log", "error", err)
			}
			logger.Error("unknown service", "service", queue.Service)
			return
		}

		logger.Info("processing queue entry", "queue", queue.Name, "timeout", queue.Timeout, "service", srv.String(), "data", event.Data)

		if err := srv.UnmarshalJSON([]byte(event.Data)); err != nil {
			if err := logWriter.WriteLog([]byte(fmt.Sprintf("failed to unmarshal data: %v", err))); err != nil {
				logger.Error("failed to write to log", "error", err)
			}
			logger.Error("failed to unmarshal data", "data", event.Data, "error", err)
			continue
		}

		for i := 0; i <= queue.Retries; i++ {
			res, err := processItem(srv, queue.Timeout)
			if err != nil {
				if err := logWriter.WriteLog([]byte(fmt.Sprintf("failed to process service: %v", err))); err != nil {
					logger.Error("failed to write to log", "error", err)
				}
				logger.Error("failed to process service", "name", queue.Name, "service", srv.String(), "error", err)
				continue
			}
			if res == nil {
				if err := logWriter.WriteLog([]byte("service returned nil result")); err != nil {
					logger.Error("failed to write to log", "error", err)
				}
				logger.Error("service returned nil result", "service", srv.String())
				continue
			}

			logger.Info("service processed", "service", srv.String(), "result", string(res))

			params := &dbParams{
				dsn: fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
					queue.TargetDatabaseUser,
					queue.TargetDatabasePassword,
					queue.TargetDatabaseHost,
					queue.TargetDatabasePort,
					queue.TargetDatabaseName,
				),
			}

			targetWriter, err := OpenTarget(params)
			if err != nil {
				if err := logWriter.WriteLog([]byte(fmt.Sprintf("failed to open target: %v", err))); err != nil {
					logger.Error("failed to write to log", "error", err)
				}
				logger.Error("failed to open target", "error", err)
				continue
			}
			defer targetWriter.Close()

			if err := targetWriter.Write(event.EventID, res); err != nil {
				if err := logWriter.WriteLog([]byte(fmt.Sprintf("failed to write to target: %v", err))); err != nil {
					logger.Error("failed to write to log", "error", err)
				}
				logger.Error("failed to write to target", "error", err)
				continue
			}

			if err := queueReader.MarkProcessed(event.EventID); err != nil {
				if err := logWriter.WriteLog([]byte(fmt.Sprintf("failed to mark event as processed: %v", err))); err != nil {
					logger.Error("failed to write to log", "error", err)
				}
				logger.Error("failed to mark event as processed", "eventID", event.EventID, "error", err)
				continue
			}

			if i > 0 {
				logger.Info("service processed after retries", "service", srv.String(), "retries", i)
			}

			break
		}
	}
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
