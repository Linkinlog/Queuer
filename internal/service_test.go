package internal_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/linkinlog/queuer/internal"
)

type testLogger struct {
	slog.JSONHandler
}

func (t *testLogger) Handle(ctx context.Context, record slog.Record) error {
	if record.Level == slog.LevelError {
		string := record.Message

		record.Attrs(func(a slog.Attr) bool {
			if a.Key == "error" {
				string += fmt.Sprintf(" %v", a.Value)
				return false
			}

			return true
		})

		panic(string)
	}

	t.JSONHandler.Handle(ctx, record)

	return nil
}

func TestStart(t *testing.T) {
	filePath := "../example.json"

	logger := slog.New(&testLogger{
		*slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
	})

	internal.Start(logger, filePath)
}
