package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var configCmd = &cobra.Command{
	Use:           "config",
	Short:         "Print config",
	Long:          "print config",
	SilenceErrors: true,
	SilenceUsage:  true,
	Run:           runConfig,
}

func runConfig(c *cobra.Command, args []string) {
	logger := zap.L()
	logger.Info(fmt.Sprintf("config: %+v", opts))
}
