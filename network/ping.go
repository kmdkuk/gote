package network

import (
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func sendPing(conn *icmp.PacketConn, proto, host string, timeout time.Duration) bool {
	ip, err := net.ResolveIPAddr(proto, host)
	if err != nil {
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
	if _, err := conn.WriteTo(wb, &net.IPAddr{IP: ip.IP}); err != nil {
		log.Printf("WriteTo: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(timeout))
	rb := make([]byte, 1500)
	n, _, err := conn.ReadFrom(rb)
	if err != nil {
		return false
	}
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:n])
	if err == nil && rm.Type == ipv4.ICMPTypeEchoReply {
		// log.Println(ip.IP.String() + " ping成功")
		return true
	}
	return false
}
