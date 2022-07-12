package twitter

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
	"github.com/kmdkuk/gote/cmd/option"
	"go.uber.org/zap"
)

var twitterApi *anaconda.TwitterApi

func connectTwitterAPI() *anaconda.TwitterApi {
	at := option.Opt.Twitter.AccessToken
	ats := option.Opt.Twitter.AccessTokenSecret
	ck := option.Opt.Twitter.ConsumerKey
	cs := option.Opt.Twitter.ConsumerSecret
	return anaconda.NewTwitterApiWithCredentials(at, ats, ck, cs)
}

func Post(message string) error {
	logger := zap.L()
	if twitterApi == nil {
		twitterApi = connectTwitterAPI()
	}
	tweet, err := twitterApi.PostTweet(message, nil)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("tweet success: %v", tweet))
	return nil
}
