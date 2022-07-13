package option

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
