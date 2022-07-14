package controller

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/kmdkuk/gote/pkg/constants"
	"github.com/kmdkuk/gote/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Controller", func() {
	logging.SetDebugLogging()
	connectMsg := strconv.FormatBool(true)
	disConnectMsg := strconv.FormatBool(false)
	trueCheck := checkResult{true, nil}
	falseCheck := checkResult{false, nil}
	errCheck := checkResult{false, fmt.Errorf("check error occurred")}
	testCases := []struct {
		name         string
		input        []checkResult
		expectOutput []string
	}{
		{
			name: "First check true to not notify",
			input: []checkResult{
				trueCheck,
			},
			expectOutput: []string{
				constants.MsgFirst,
			},
		},
		{
			name: "10 count true to not notify",
			input: []checkResult{
				trueCheck,
				trueCheck,
				trueCheck,
				trueCheck,
				trueCheck,
				trueCheck,
				trueCheck,
				trueCheck,
				trueCheck,
				trueCheck,
			},
			expectOutput: []string{
				constants.MsgFirst,
			},
		},
		{
			name: "6 count false to notify false",
			input: []checkResult{
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				trueCheck,
			},
			expectOutput: []string{
				constants.MsgFirst,
				disConnectMsg,
				connectMsg,
			},
		},
		{
			name: "6 count false to notify false",
			input: []checkResult{
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
			},
			expectOutput: []string{
				constants.MsgFirst,
				disConnectMsg,
			},
		},
		{
			name: "5 count false to not notify false",
			input: []checkResult{
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
			},
			expectOutput: []string{
				constants.MsgFirst,
			},
		},
		{
			name: "check results are unstable if 5 count false",
			input: []checkResult{
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				trueCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				trueCheck,
			},
			expectOutput: []string{
				constants.MsgFirst,
			},
		},
		{
			name: "check results are unstable if 6 count false",
			input: []checkResult{
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				trueCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				falseCheck,
				trueCheck,
			},
			expectOutput: []string{
				constants.MsgFirst,
				disConnectMsg,
				connectMsg,
				disConnectMsg,
				connectMsg,
			},
		},
		{
			name: "check results are unstable if 6 count error",
			input: []checkResult{
				errCheck,
				errCheck,
				errCheck,
				errCheck,
				errCheck,
				errCheck,
				trueCheck,
				errCheck,
				errCheck,
				errCheck,
				errCheck,
				errCheck,
				errCheck,
				trueCheck,
			},
			expectOutput: []string{
				constants.MsgFirst,
				disConnectMsg,
				connectMsg,
				disConnectMsg,
				connectMsg,
			},
		},
	}
	for _, testcase := range testCases {
		tc := testcase
		It(fmt.Sprintf("should check status, %s", tc.name), func() {
			By("prepare controller")
			checkerMock := CheckerMock{
				Pattern: tc.input,
				Count:   0,
			}
			notifierMock := NotifierMock{}
			interval := time.Millisecond
			controller := &controller{
				checker:      &checkerMock,
				notifier:     &notifierMock,
				count:        0,
				recentStatus: true,
				interval:     time.Nanosecond,
			}
			go func() {
				controller.Run()
			}()

			time.Sleep(time.Duration(len(tc.input))*interval + time.Millisecond)
			By("checking output")
			notifierMock.mu.Lock()
			defer notifierMock.mu.Unlock()
			Expect(notifierMock.Output).To(Equal(tc.expectOutput), tc.name)
		})
	}
})

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)

	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(time.Second)
	SetDefaultConsistentlyDuration(10 * time.Second)
	SetDefaultConsistentlyPollingInterval(time.Second)

	RunSpecs(t, "Runner Suite")
}

type CheckerMock struct {
	Pattern []checkResult
	Count   int
}

type checkResult struct {
	status bool
	err    error
}

func (c *CheckerMock) Check() (bool, error) {
	logger := zap.L()
	if len(c.Pattern) <= c.Count {
		logger.Debug("wait")
		time.Sleep(999 * time.Hour)
		return false, fmt.Errorf("Not reaching this point during the test.")
	}
	status := c.Pattern[c.Count]
	c.Count++
	logger.Debug(fmt.Sprintf("check count %d", c.Count))
	return status.status, status.err
}

func (c CheckerMock) Close() error {
	return nil
}

type NotifierMock struct {
	Output []string
	mu     sync.Mutex
}

func (n *NotifierMock) Notify(msg string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Output = append(n.Output, msg)
	return nil
}

func (n *NotifierMock) NotifyStatus(connection bool) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Output = append(n.Output, strconv.FormatBool(connection))
	return nil
}
