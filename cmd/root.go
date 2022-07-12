package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kmdkuk/gote/cmd/option"
	"github.com/kmdkuk/gote/pkg/network"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
)

func init() {
	cobra.OnInitialize(initconf)

	rootCmd.Flags().StringVarP(&configFile, "file", "f", ".gote.yaml", "Gote config file")

	rootCmd.Flags().StringVarP(&option.Opt.Mode, "mode", "m", "ping", "How to do a health check. ping or http")
	rootCmd.Flags().StringVarP(&option.Opt.Host, "target", "t", "127.0.0.1", "Target for health check. domain or ip or URL")
	rootCmd.Flags().StringVarP(&option.Opt.Notification, "notification", "n", "slack", "Destination to notify when health check fails. slack or twitter")

	rootCmd.Flags().StringVar(&option.Opt.MsgDisconnect, "msgdisconnect", "disconnected", "Message when disconnecting")
	rootCmd.Flags().StringVar(&option.Opt.MsgConnect, "msgconnect", "connected", "Message when connecting")
	rootCmd.Flags().StringVar(&option.Opt.MsgSuffix, "msgsuffix", "", "Suffix of common message")

	rootCmd.AddCommand(configCmd)
}

func initconf() {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return
		} else {
			// Config file was found but another error was produced
		}
	}
	if err := viper.Unmarshal(&option.Opt); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
	c := network.NewChecker()
	err := c.Check()
	if err != nil {
		log.Fatalln(err)
	}
}
