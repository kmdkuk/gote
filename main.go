package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/net/ipv4"

	"golang.org/x/net/icmp"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

var recentPingResult bool
var recentStatus bool
var count int
var db *sql.DB

// DB の接続情報
const (
	DRIVER_NAME = "mysql" // ドライバ名(mysql固定)
	// user:password@tcp(container-name:port)/dbname ※mysql はデフォルトで用意されているDB
	DATA_SOURCE_NAME = "docker:docker@tcp(db:3306)/network-monitoring?parseTime=true"
)

type DisconnectDate struct {
	ID      int
	StartAt time.Time
	EndAt   time.Time
}

func RecordStartAt() {
	// insert
	ins, err := db.Prepare("INSERT INTO `network-monitoring`.disconnect_dates(start_at) VALUES(?)")
	if err != nil {
		log.Printf("[Err] Prepare %v", err)
	}
	log.Println(time.Now())
	_, err = ins.Exec(time.Now())
	if err != nil {
		log.Printf("[Err] Insert %v", err)
	}
}

func RecordEndAt() {
	query := "select * from disconnect_dates where end_at=0;"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("[Err] %s %v", query, err)
		return
	}
	record := DisconnectDate{}
	rows.Next()
	if err := rows.Scan(&record.ID, &record.StartAt, &record.EndAt); err != nil {
		log.Printf("[Err] %v", err)
		return
	}
	// set end_at
	upd, err := db.Prepare("UPDATE disconnect_dates set end_at = ? where id = ?")
	if err != nil {
		log.Printf("[Err] Update Prepare: %v", err)
		return
	}
	upd.Exec(time.Now(), record.ID)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func connectTwitterApi() *anaconda.TwitterApi {
	at := os.Getenv("ACCESS_TOKEN")
	ats := os.Getenv("ACCESS_TOKEN_SECRET")
	ck := os.Getenv("CONSUMER_KEY")
	cs := os.Getenv("CONSUMER_SECRET")
	fmt.Println(at)
	fmt.Println(ats)
	fmt.Println(ck)
	fmt.Println(cs)
	return anaconda.NewTwitterApiWithCredentials(at, ats, ck, cs)
}

func initialize() *anaconda.TwitterApi {
	// loadEnv()
	recentPingResult = true
	recentStatus = true
	count = 0
	var err error
	db, err = sql.Open(DRIVER_NAME, DATA_SOURCE_NAME)
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
	return connectTwitterApi()
}

func tweet(api *anaconda.TwitterApi) {
	t := time.Now()
	layout := "2006-01-02 15:04"
	log.Println("ping失敗")
	message := ""
	hashtag := "#kmdkukのネット回線"
	if recentStatus {
		message += "[" + t.Format(layout) + "] 切断されました． " + hashtag
	} else {
		message += "[" + t.Format(layout) + "] 復旧されました． " + hashtag
	}
	tweet, err := api.PostTweet(message, nil)
	if err != nil {
		log.Printf("Tweet: %v", err)
	} else {
		log.Printf("Tweet success: %v", tweet)
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

func sendPing(c *icmp.PacketConn, proto, host string, timeout time.Duration) bool {
	ip, err := net.ResolveIPAddr(proto, host)
	if err != nil {
		log.Printf("ping失敗: %d", count)
		log.Printf("ResolveIPAddr: %v", err)
		return false
	}
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Printf("Marshal: %v", err)
	}
	if _, err := c.WriteTo(wb, &net.IPAddr{IP: ip.IP}); err != nil {
		log.Printf("WriteTo: %v", err)
	}

	c.SetReadDeadline(time.Now().Add(timeout))
	rb := make([]byte, 1500)
	n, _, err := c.ReadFrom(rb)
	if err != nil {
		// log.Printf(ip.IP.String()+" ping失敗: %d", count)
		// log.Printf("err: %v", err)
		return false
	}
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:n])
	if err == nil && rm.Type == ipv4.ICMPTypeEchoReply {
		// log.Println(ip.IP.String() + " ping成功")
		return true
	}
	// log.Printf(ip.IP.String()+" ping失敗: %d", count)
	// log.Printf("err: %v", err)
	return false
}

func main() {
	api := initialize()
	defer db.Close()
	fmt.Println("Hello, world!")
	RecordStartAt()
	RecordEndAt()
	var sleep time.Duration
	var timeout time.Duration

	flag.DurationVar(&sleep, "s", 2*time.Second, "sleep")
	flag.DurationVar(&timeout, "t", 1*time.Second, "timeout")
	flag.Parse()

	proto := "ip4"
	host := "minecraft.kmdkuk.com"

	c, err := icmp.ListenPacket(proto+":icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket: %v", err)
	}
	defer c.Close()

	for {
		if sendPing(c, proto, host, timeout) {
			if count > 0 {
				log.Printf("pingが復旧するまで %d 回エラー", count)
			}
			count = 0
			recentPingResult = true
			if isStatusToggled() {
				tweet(api)
				recentStatus = true
			}
		} else {
			count++
			recentPingResult = false
			if isStatusToggled() == true {
				tweet(api)
				recentStatus = false
			}
		}
		time.Sleep(sleep)
	}
}
