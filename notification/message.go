package notification

import "time"

func BuildMessage(recentStatus bool, suffix string) string {
	t := time.Now()
	layout := "2006-01-02 15:04"
	message := ""
	if recentStatus {
		message += "[" + t.Format(layout) + "] 切断されました． " + suffix
	} else {
		message += "[" + t.Format(layout) + "] 復旧されました． " + suffix
	}
	return message
}
