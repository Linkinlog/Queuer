package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/linkinlog/queuer/internal"
	"github.com/spf13/cobra"
)

var (
	Verbosity      int
	ConfigFilePath string
)

var RootCmd = &cobra.Command{
	Use:     "queuer",
	Short:   "Queuer is a simple task queue",
	Example: "queuer -f config.json",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		slogOpts := &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelWarn,
		}

		if Verbosity > 0 {
			slogOpts.Level = slog.LevelInfo
		}

		if Verbosity > 1 {
			slogOpts.Level = slog.LevelDebug
		}

		if Verbosity > 2 {
			slogOpts.AddSource = true
		}

		logger := slog.New(slog.NewJSONHandler(os.Stdout, slogOpts))

		internal.Start(logger, ConfigFilePath)

		logger.Info("Finished", "Duration", time.Since(start).String())
	},
}

func Execute() error {
	RootCmd.PersistentFlags().CountVarP(&Verbosity, "verbosity", "v", "Set the verbosity level (e.g., -v for verbose, -vv for very verbose)")
	RootCmd.PersistentFlags().StringVarP(&ConfigFilePath, "file", "f", "config.json", "Path to the configuration file")
	return RootCmd.Execute()
}
