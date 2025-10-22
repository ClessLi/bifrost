package utils

import (
	"context"
	"net"
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"

	"golang.org/x/sync/errgroup"
)

func udpConnectivity(ipv4, port string, timeout time.Duration) v1.NetworkConnectivity {
	address := net.JoinHostPort(ipv4, port)
	conn, err := net.DialTimeout("udp", address, timeout)
	if err != nil {
		return v1.NetUnreachable
	}
	defer conn.Close()

	listenReadyCtx, listenReady := context.WithCancel(context.Background())
	var unreach []byte
	eg := new(errgroup.Group)
	eg.Go(func() (err error) {
		unreach, err = listenICMPUnreachable(timeout, listenReady)

		return
	})

	<-listenReadyCtx.Done()
	// udp ping
	pingdata := []byte("")
	_, err = conn.Write(pingdata)
	if err != nil {
		return v1.NetUnreachable
	}

	// listen ICMP unreachable
	err = eg.Wait()
	if err != nil {
		return v1.NetUnknown
	}

	if unreach == nil {
		return v1.NetReachable
	}

	// parse unreachable package
	dstip, dstport := parseUnreachUDP(unreach)

	if ipv4 == dstip && port == dstport {
		return v1.NetUnreachable
	}

	return v1.NetUnknown
}
