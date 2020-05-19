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
	count        int
	recentStatus bool
}

func NewChecker(mode, host, notification string) Checker {
	return &checker{
		mode,
		host,
		notification,
		0,
		true,
	}
}

func (c *checker) Check() error {
	if c.mode == "ping" {
		return c.checkPing()
	} else if c.mode == "http" {
		return c.checkHTTP()
	} else {
		return fmt.Errorf("invalid mode")
	}
}

func (c *checker) checkPing() error {

	sleep := 2 * time.Second
	timeout := 1 * time.Second

	proto := "ip4"

	conn, err := icmp.ListenPacket(proto+":icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket: %v", err)
	}
	defer conn.Close()

	for {
		c.statusUpdate(sendPing(conn, proto, c.host, timeout))
		time.Sleep(sleep)
	}
}

func (c *checker) isStatusToggled(recentCheckResult bool) bool {
	result := false
	if c.recentStatus {
		if c.count > 5 && recentCheckResult == false {
			result = true
		}
	} else {
		if recentCheckResult == true {
			result = true
		}
	}
	return result
}

func (c *checker) statusUpdate(checkResult bool) {
	suffix := "#kmdkukのネット回線"
	if checkResult {
		if c.count > 0 {
			log.Printf("pingが復旧するまで %d 回エラー", c.count)
		}
		c.count = 0
		if c.isStatusToggled(checkResult) {
			message := notification.BuildMessage(c.recentStatus, suffix)
			notification.Notification(c.notification, message)
			c.recentStatus = true
		}
	} else {
		c.count++
		log.Println(c.count, "here")
		if c.isStatusToggled(checkResult) {
			log.Println("send notification")
			message := notification.BuildMessage(c.recentStatus, suffix)
			notification.Notification(c.notification, message)
			c.recentStatus = false
		}
	}
}

func (c *checker) checkHTTP() error {
	sleep := 2 * time.Second
	for {
		c.statusUpdate(sendHttp(c.host))
		time.Sleep(sleep)
	}
	return fmt.Errorf("not implemented")
}
