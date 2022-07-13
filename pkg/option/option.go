package option

import "github.com/spf13/cobra"

type Options struct {
	Mode          string
	Host          string
	Notification  string
	Twitter       `mapstructure:"Twitter"`
	Slack         `mapstructure:"Slack"`
	MsgDisconnect string
	MsgConnect    string
	MsgSuffix     string
}

type Twitter struct {
	AccessToken       string `mapstructure:"access_token"`
	AccessTokenSecret string `mapstructure:"access_token_secret"`
	ConsumerKey       string `mapstructure:"consumer_key"`
	ConsumerSecret    string `mapstructure:"consumer_secret"`
}

type Slack struct {
	WebhookURL string `mapstructure:"webhook_url"`
	Channel    string
}

func AddOptionFlags(cmd *cobra.Command, opts *Options) {
	cmd.PersistentFlags().StringVarP(&opts.Mode, "mode", "m", "ping", "How to do a health check. ping or http")
	cmd.PersistentFlags().StringVarP(&opts.Host, "target", "t", "127.0.0.1", "Target for health check. domain or ip or URL")
	cmd.PersistentFlags().StringVarP(&opts.Notification, "notification", "n", "slack", "Destination to notify when health check fails. slack or twitter")
	cmd.PersistentFlags().StringVar(&opts.MsgDisconnect, "msgdisconnect", "disconnected", "Message when disconnecting")
	cmd.PersistentFlags().StringVar(&opts.MsgConnect, "msgconnect", "connected", "Message when connecting")
	cmd.PersistentFlags().StringVar(&opts.MsgSuffix, "msgsuffix", "", "Suffix of common message")
}
