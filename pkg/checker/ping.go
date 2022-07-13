package checker

import (
	"fmt"
	"net"
	"os"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func newPingChecker(host string) (Checker, error) {
	proto := "ip4"
	conn, err := icmp.ListenPacket(proto+":icmp", "0.0.0.0")
	if err != nil {
		return nil, err
	}

	return &pingChecker{
		proto: "ip4",
		conn:  conn,
		host:  host,
	}, nil
}

type pingChecker struct {
	proto string
	conn  *icmp.PacketConn
	host  string
}

func (p pingChecker) Check() (bool, error) {
	logger := zap.L()

	timeout := 1 * time.Second

	ip, err := net.ResolveIPAddr(p.proto, p.host)
	if err != nil {
		return false, err
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
		return false, err
	}
	if _, err := p.conn.WriteTo(wb, &net.IPAddr{IP: ip.IP}); err != nil {
		return false, err
	}

	p.conn.SetReadDeadline(time.Now().Add(timeout))
	rb := make([]byte, 1500)
	n, _, err := p.conn.ReadFrom(rb)
	if err != nil {
		return false, err
	}
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:n])
	if err != nil {
		return false, err
	}
	if rm.Type != ipv4.ICMPTypeEchoReply {
		logger.Info(fmt.Sprintf("returned non-echo reply type: %s", rm.Type))
		return false, nil
	}
	return true, nil
}

func (p pingChecker) Close() error {
	return p.conn.Close()
}
