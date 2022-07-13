package controller

import (
	"fmt"
	"time"

	"github.com/kmdkuk/gote/pkg/checker"
	"github.com/kmdkuk/gote/pkg/notifier"
	"github.com/kmdkuk/gote/pkg/option"
	"go.uber.org/zap"
)

type Controller interface {
	Run()
	Close() error
}

func NewController(opts option.Options) (Controller, error) {
	c, err := checker.NewChecker(opts.Mode, opts.Host)
	if err != nil {
		return nil, err
	}

	n, err := notifier.NewNotifier(opts)
	if err != nil {
		return nil, err
	}
	return &controller{
		checker:      c,
		notifier:     n,
		count:        0,
		recentStatus: true,
	}, nil
}

type controller struct {
	checker      checker.Checker
	notifier     notifier.Notifier
	count        int
	recentStatus bool
}

func (c controller) Run() {
	logger := zap.L()
	c.notifier.Notify("Start Health Checker")
	for {
		time.Sleep(2 * time.Second)
		status, err := c.checker.Check()
		if err != nil {
			logger.Error("check error occurred", zap.Error(err))
			continue
		}
		c.statusUpdate(status)
	}
}

func (c *controller) statusUpdate(checkResult bool) {
	logger := zap.L()
	if checkResult {
		if c.count > 0 {
			logger.Info(fmt.Sprintf("復旧するまで %d 回エラー", c.count))
		}
		c.count = 0
		logger.Debug("check succeeded")
		if c.isStatusToggled(checkResult) {
			c.notifier.NotifyStatus(c.recentStatus)
			c.recentStatus = true
		}
		return
	}
	c.count++
	logger.Debug("check failed")
	if c.isStatusToggled(checkResult) {
		c.notifier.NotifyStatus(c.recentStatus)
		c.recentStatus = false
	}
}

func (c controller) isStatusToggled(recentCheckResult bool) bool {
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

func (c controller) Close() error {
	return c.checker.Close()
}
