package network

import (
	"fmt"
	"time"

	"github.com/kmdkuk/gote/cmd/option"
	"github.com/kmdkuk/gote/pkg/notification"
	"go.uber.org/zap"
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
	logger := zap.L()

	sleep := 2 * time.Second
	timeout := 1 * time.Second

	proto := "ip4"

	conn, err := icmp.ListenPacket(proto+":icmp", "0.0.0.0")
	if err != nil {
		logger.Fatal("ListenPacket", zap.Error(err))
	}
	defer conn.Close()

	for {
		c.statusUpdate(sendPing(conn, proto, option.Opt.Host, timeout))
		time.Sleep(sleep)
	}
}

func (c *checker) isStatusToggled(recentCheckResult bool) bool {
	if c.recentStatus {
		if c.count > 5 && !recentCheckResult {
			return true
		}
	} else {
		if recentCheckResult {
			return true
		}
	}
	return false
}

func (c *checker) statusUpdate(checkResult bool) {
	logger := zap.L()
	if checkResult {
		if c.count > 0 {
			logger.Info(fmt.Sprintf("復旧するまで %d 回エラー", c.count))
		}
		c.count = 0
		logger.Debug("check succeeded")
		if c.isStatusToggled(checkResult) {
			message := notification.BuildMessage(c.recentStatus)
			notification.Notification(option.Opt.Notification, message)
			c.recentStatus = true
		}
	} else {
		c.count++
		logger.Debug("check failed")
		if c.isStatusToggled(checkResult) {
			logger.Info("send notification")
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
