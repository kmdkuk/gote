package twitter

import (
	"log"
	"os"
	"time"

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

func Tweet(recentStatus bool) {
	t := time.Now()
	layout := "2006-01-02 15:04"
	log.Println("ping失敗")
	message := ""
	hashtag := "#kmdkukのネット回線"
	if recentStatus {
		message += "[" + t.Format(layout) + "] 切断されました． " + hashtag
	} else {
		message += "[" + t.Format(layout) + "] 復旧されました． " + hashtag
	}
	tweet, err := twitterApi.PostTweet(message, nil)
	if err != nil {
		log.Printf("Tweet: %v", err)
	} else {
		log.Printf("Tweet success: %v", tweet)
	}
}

func init() {
	twitterApi = connectTwitterAPI()
}
