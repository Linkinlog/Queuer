package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/linkinlog/queuer/internal/config"
	"github.com/linkinlog/queuer/internal/db"
	"github.com/linkinlog/queuer/internal/logger"
	"github.com/linkinlog/queuer/internal/services"
)

type Service interface {
	json.Unmarshaler
	fmt.Stringer
	SetLogger(*slog.Logger)
	Run() (results chan []byte, errs chan error)
}

func Start(cfg *config.Config) {
	var wg sync.WaitGroup
	for _, queue := range cfg.Queues {
		wg.Add(1)
		go processQueue(queue, cfg.Creds, cfg.SlogOpts, &wg)
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

func processQueue(queue *config.Queue, creds *config.Credentials, s *slog.HandlerOptions, wg *sync.WaitGroup) {
	defer wg.Done()

	loggerDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		creds.LogDatabaseUser,
		creds.LogDatabasePassword,
		queue.LogDatabaseHost,
		queue.LogDatabasePort,
		queue.LogDatabaseName,
	)

	loggerParams := &logger.LoggerParams{
		W:   os.Stdout,
		DSN: loggerDSN,
	}

	logWriter := logger.New(loggerParams)

	logger := slog.New(slog.NewJSONHandler(logWriter, s))

	logger.Info("fetching queue data", "queue", queue.Name)

	queueDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		creds.QueueDatabaseUser,
		creds.QueueDatabasePassword,
		queue.QueueDatabaseHost,
		queue.QueueDatabasePort,
		queue.QueueDatabaseName,
	)

	queueReader, err := db.OpenQueue(queueDSN)
	if err != nil {
		logger.Error("failed to open queue", "error", err)
		return
	}

	defer queueReader.Close()

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

		srv.SetLogger(logger)

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

			targetDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
				creds.TargetDatabaseUser,
				creds.TargetDatabasePassword,
				queue.TargetDatabaseHost,
				queue.TargetDatabasePort,
				queue.TargetDatabaseName,
			)

			targetWriter, err := db.OpenTarget(targetDSN)
			if err != nil {
				logger.Error("failed to open target", "error", err)
				continue
			}
			defer targetWriter.Close()

			if err := targetWriter.Write(event.EventID, res); err != nil {
				logger.Error("failed to write to target", "error", err)
				continue
			}

			if err := queueReader.MarkProcessed(event.EventID); err != nil {
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
