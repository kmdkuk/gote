package network

import (
	"fmt"
	"log"
	"time"

	"github.com/kmdkuk/gote/cmd/option"
	"github.com/kmdkuk/gote/notification"
	"golang.org/x/net/icmp"
)

type Checker interface {
	Check() error
}

type checker struct {
	count        int
	recentStatus bool
}

func NewChecker() Checker {
	return &checker{
		0,
		true,
	}
}

func (c *checker) Check() error {
	if option.Opt.Mode == "ping" {
		return c.checkPing()
	} else if option.Opt.Mode == "http" {
		return c.checkHTTP()
	} else {
		return fmt.Errorf("invalid Mode")
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
		c.statusUpdate(sendPing(conn, proto, option.Opt.Host, timeout))
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
	if checkResult {
		if c.count > 0 {
			log.Printf("pingが復旧するまで %d 回エラー", c.count)
		}
		c.count = 0
		if c.isStatusToggled(checkResult) {
			message := notification.BuildMessage(c.recentStatus)
			notification.Notification(option.Opt.Notification, message)
			c.recentStatus = true
		}
	} else {
		c.count++
		log.Println(c.count, "here")
		if c.isStatusToggled(checkResult) {
			log.Println("send notification")
			message := notification.BuildMessage(c.recentStatus)
			notification.Notification(option.Opt.Notification, message)
			c.recentStatus = false
		}
	}
}

func (c *checker) checkHTTP() error {
	sleep := 2 * time.Second
	for {
		c.statusUpdate(sendHttp(option.Opt.Host))
		time.Sleep(sleep)
	}
}
