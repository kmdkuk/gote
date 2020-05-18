package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kmdkuk/gote/network"
	"github.com/kmdkuk/gote/notification/twitter"
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
)

var (
	recentPingResult bool
	recentStatus     bool
	count            int
)

func init() {
	recentPingResult = true
	recentStatus = true
	count = 0
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintf(os.Stderr, "%v\n", rootCmd.UsageString())
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:           "gote",
	Short:         "Service health check and notification",
	Long:          "Service health check and notification",
	SilenceErrors: true,
	SilenceUsage:  true,
	Run:           run,
}

func run(cmd *cobra.Command, args []string) {
	var sleep time.Duration
	var timeout time.Duration

	flag.DurationVar(&sleep, "s", 2*time.Second, "sleep")
	flag.DurationVar(&timeout, "t", 1*time.Second, "timeout")
	flag.Parse()

	proto := "ip4"
	host := "minecraft.kmdkuk.com"

	conn, err := icmp.ListenPacket(proto+":icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket: %v", err)
	}
	defer conn.Close()

	for {
		if network.SendPing(conn, proto, host, timeout) {
			if count > 0 {
				log.Printf("pingが復旧するまで %d 回エラー", count)
			}
			count = 0
			recentPingResult = true
			if isStatusToggled() {
				twitter.Tweet(recentStatus)
				recentStatus = true
			}
		} else {
			count++
			recentPingResult = false
			if isStatusToggled() == true {
				twitter.Tweet(recentStatus)
				recentStatus = false
			}
		}
		time.Sleep(sleep)
	}
}

func isStatusToggled() bool {
	result := false
	if recentStatus {
		if count > 5 && recentPingResult == false {
			result = true
		}
	} else {
		if recentPingResult == true {
			result = true
		}
	}
	return result
}
