package twitter

import (
	"log"

	"github.com/ChimeraCoder/anaconda"
	"github.com/kmdkuk/gote/cmd/option"
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
	tweet, err := twitterApi.PostTweet(message, nil)
	if err != nil {
		log.Printf("Tweet: %v", err)
		return err
	} else {
		log.Printf("Tweet success: %v", tweet)
	}
	return nil
}

func init() {
	twitterApi = connectTwitterAPI()
}
