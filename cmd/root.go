package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kmdkuk/gote/network"
	"github.com/spf13/cobra"
)

var (
	opt Options
)

type Options struct {
	mode         string
	host         string
	notification string
}

func init() {
	opt = Options{}

	rootCmd.Flags().StringVarP(&opt.mode, "mode", "m", "ping", "How to do a health check. ping or http")
	rootCmd.Flags().StringVarP(&opt.host, "target", "t", "127.0.0.1", "Target for health check. domain or ip or URL")
	rootCmd.Flags().StringVarP(&opt.notification, "notification", "n", "slack", "Destination to notify when health check fails. slack or twitter")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintf(os.Stderr, "%v\n", rootCmd.UsageString())
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:           "gote",
	Short:         "Service health check and notification",
	Long:          "Service health check and notification",
	SilenceErrors: true,
	SilenceUsage:  true,
	Run:           run,
}

func run(cmd *cobra.Command, args []string) {
	c := network.NewChecker(opt.mode, opt.host, opt.notification)
	err := c.Check()
	if err != nil {
		log.Fatalln(err)
	}
}
