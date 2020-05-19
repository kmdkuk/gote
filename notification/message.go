package notification

import (
	"time"

	"github.com/kmdkuk/gote/cmd/option"
)

func BuildMessage(recentStatus bool) string {
	t := time.Now()
	layout := "2006-01-02 15:04"
	message := ""
	if recentStatus {
		message += "[" + t.Format(layout) + "] " + option.Opt.MsgDisconnect + " " + option.Opt.MsgSuffix
	} else {
		message += "[" + t.Format(layout) + "] " + option.Opt.MsgConnect + " " + option.Opt.MsgSuffix
	}
	return message
}
