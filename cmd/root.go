package cmd

import (
	"os"

	"github.com/kmdkuk/gote/pkg/controller"
	"github.com/kmdkuk/gote/pkg/logging"
	"github.com/kmdkuk/gote/pkg/option"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	configFile string
	opts       option.Options
)

func init() {
	cobra.OnInitialize(initconf)

	rootCmd.Flags().StringVarP(&configFile, "file", "f", ".gote.yaml", "Gote config file")

	option.AddOptionFlags(rootCmd, &opts)
	logging.AddLoggingFlags(rootCmd)

	rootCmd.AddCommand(configCmd)
}

func initconf() {
	logger := zap.L()
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return
		}
		// Config file was found but another error was produced
	}
	if err := viper.Unmarshal(&opts); err != nil {
		logger.Error("config unmarshal failed", zap.Error(err))
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger := zap.L()
		logger.Error("error occurred", zap.Error(err))
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:           "gote",
	Short:         "Service health check and notification",
	Long:          "Service health check and notification",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          run,
}

func run(cmd *cobra.Command, args []string) (err error) {
	c, err := controller.NewController(opts)
	if err != nil {
		return err
	}
	defer func() {
		err2 := c.Close()
		if err2 != nil {
			err = err2
		}
	}()
	c.Run()
	return nil
}
