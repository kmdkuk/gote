package checker

import (
	"net/http"
)

func newHttpChecker(url string) Checker {
	return &httpChecker{
		url: url,
	}
}

type httpChecker struct {
	url string
}

func (h httpChecker) Check() (bool, error) {
	resp, err := http.Get(h.url)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}

func (h httpChecker) Close() error {
	return nil
}
