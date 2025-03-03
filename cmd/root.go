package cmd

import (
	"fmt"
	"time"

	"github.com/linkinlog/queuer/internal"
	"github.com/linkinlog/queuer/internal/config"
	"github.com/spf13/cobra"
)

var (
	Verbosity      int
	ConfigFilePath string
)

var RootCmd = &cobra.Command{
	Use:     "queuer",
	Short:   "Queuer is a simple task queue",
	Example: "queuer -f example.json",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		cfg, err := config.ParseConfig(ConfigFilePath, Verbosity)
		if err != nil {
			fmt.Println("failed to parse config", "error", err)
			return
		}

		internal.Start(cfg)

		fmt.Println("Queuer finished in", time.Since(start))
	},
}

func Execute() error {
	RootCmd.PersistentFlags().CountVarP(&Verbosity, "verbosity", "v", "Set the verbosity level (e.g., -v for verbose, -vv for very verbose)")
	RootCmd.PersistentFlags().StringVarP(&ConfigFilePath, "file", "f", "example.json", "Path to the configuration file")
	return RootCmd.Execute()
}
