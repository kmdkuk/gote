package cmd

import (
	"fmt"

	"github.com/kmdkuk/gote/cmd/option"
	"github.com/spf13/cobra"
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
	fmt.Printf("config: %v\n", option.Opt)
}
