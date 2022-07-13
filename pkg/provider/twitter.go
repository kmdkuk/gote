package provider

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
	"github.com/kmdkuk/gote/pkg/option"
	"go.uber.org/zap"
)

func NewTwitter(twitterOpts option.Twitter) Provider {
	return twitter{
		api: connectTwitterAPI(twitterOpts),
	}
}

type twitter struct {
	api *anaconda.TwitterApi
}

func connectTwitterAPI(twitterOpts option.Twitter) *anaconda.TwitterApi {
	at := twitterOpts.AccessToken
	ats := twitterOpts.AccessTokenSecret
	ck := twitterOpts.ConsumerKey
	cs := twitterOpts.ConsumerSecret
	return anaconda.NewTwitterApiWithCredentials(at, ats, ck, cs)
}

func (t twitter) Post(message string) error {
	logger := zap.L()
	tweet, err := t.api.PostTweet(message, nil)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("tweet success: %v", tweet))
	return nil
}
