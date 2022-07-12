package network

import (
	"net/http"

	"go.uber.org/zap"
)

func sendHttp(url string) bool {
	logger := zap.L()
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("htto.Get", zap.Error(err))
		return false
	}
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}
