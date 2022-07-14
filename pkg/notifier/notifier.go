package notifier

import (
	"errors"
	"time"

	"github.com/kmdkuk/gote/pkg/option"
	"github.com/kmdkuk/gote/pkg/provider"
)

type Notifier interface {
	NotifyStatus(status bool) error
	Notify(msg string) error
}

func NewNotifier(opts option.Options) (Notifier, error) {
	var p provider.Provider
	switch opts.Notification {
	case "twitter":
		p = provider.NewTwitter(opts.Twitter)
	case "slack":
		p = provider.NewSlack(opts.Slack)
	default:
		return nil, errors.New("invalid notification destination")
	}
	return &notifier{
		provider:      p,
		msgConnect:    opts.MsgConnect,
		msgDisconnect: opts.MsgDisconnect,
		msgSuffix:     opts.MsgSuffix,
	}, nil
}

type notifier struct {
	provider      provider.Provider
	msgConnect    string
	msgDisconnect string
	msgSuffix     string
}

func (n *notifier) NotifyStatus(connection bool) error {
	msg := n.buildMessage(connection)
	return n.Notify(msg)
}

func (n *notifier) Notify(msg string) error {
	return n.provider.Post(msg)
}

func (n *notifier) buildMessage(connection bool) string {
	t := time.Now()
	layout := "2006-01-02 15:04"
	if connection {
		return "[" + t.Format(layout) + "] " + n.msgConnect + " " + n.msgSuffix
	}
	return "[" + t.Format(layout) + "] " + n.msgDisconnect + " " + n.msgSuffix
}
