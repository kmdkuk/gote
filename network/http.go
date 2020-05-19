package network

import (
	"log"
	"net/http"
)

func sendHttp(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("[Err]http.Get", err)
		return false
	}
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}
