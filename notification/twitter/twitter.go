package twitter

import (
	"log"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

var twitterApi *anaconda.TwitterApi

func connectTwitterAPI() *anaconda.TwitterApi {
	at := os.Getenv("ACCESS_TOKEN")
	ats := os.Getenv("ACCESS_TOKEN_SECRET")
	ck := os.Getenv("CONSUMER_KEY")
	cs := os.Getenv("CONSUMER_SECRET")
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
