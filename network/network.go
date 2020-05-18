package network

import (
	"fmt"
	"log"
	"time"

	"github.com/kmdkuk/gote/notification"
	"golang.org/x/net/icmp"
)

type Checker interface {
	Check() error
}

type checker struct {
	mode         string
	host         string
	notification string
}

func NewChecker(mode, host, notification string) Checker {
	return &checker{
		mode,
		host,
		notification,
	}
}

func (c *checker) Check() error {
	if c.mode == "ping" {
		return c.checkPing()
	} else if c.mode == "http" {
		return c.checkHttp()
	} else {
		return fmt.Errorf("invalid mode")
	}
}

func (c *checker) checkPing() error {
	recentStatus := true
	count := 0

	sleep := 2 * time.Second
	timeout := 1 * time.Second

	proto := "ip4"

	conn, err := icmp.ListenPacket(proto+":icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket: %v", err)
	}
	defer conn.Close()

	hashtag := "#kmdkukのネット回線"

	for {
		if pingResult := sendPing(conn, proto, c.host, timeout); pingResult {
			if count > 0 {
				log.Printf("pingが復旧するまで %d 回エラー", count)
			}
			count = 0
			if isStatusToggled(count, recentStatus, pingResult) {
				message := notification.BuildMessage(recentStatus, hashtag)
				notification.Notification(c.notification, message)
				recentStatus = true
			}
		} else {
			count++
			log.Println(count, "here")
			if isStatusToggled(count, recentStatus, pingResult) {
				log.Println("send notification")
				message := notification.BuildMessage(recentStatus, hashtag)
				notification.Notification(c.notification, message)
				recentStatus = false
			}
		}
		time.Sleep(sleep)
	}
}

func isStatusToggled(count int, recentStatus, recentPingResult bool) bool {
	result := false
	if recentStatus {
		if count > 5 && recentPingResult == false {
			result = true
		}
	} else {
		if recentPingResult == true {
			result = true
		}
	}
	return result
}

func (c *checker) checkHttp() error {
	return fmt.Errorf("not implemented")
}
