package checker

import "fmt"

type Checker interface {
	Check() (bool, error)
	Close() error
}

func NewChecker(mode, dest string) (Checker, error) {
	switch mode {
	case "http":
		return newHttpChecker(dest), nil
	case "ping":
		return newPingChecker(dest)
	default:
		return nil, fmt.Errorf("invalid check mode")
	}
}
