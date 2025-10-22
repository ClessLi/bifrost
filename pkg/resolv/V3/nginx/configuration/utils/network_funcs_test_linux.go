package utils

import (
	"context"
	"net"
	"testing"
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"

	"github.com/marmotedu/errors"
	"golang.org/x/sync/errgroup"
)

type udpTest struct {
	addr *net.UDPAddr
	conn *net.UDPConn
	ctx  context.Context
	done context.CancelFunc
}

func (u *udpTest) Run() error {
	var err error
	u.conn, err = net.ListenUDP("udp", u.addr)
	defer u.conn.Close()
	if err != nil {
		return err
	}
	for {
		select {
		case <-u.ctx.Done():
			return nil
		default:
			_, err := u.conn.WriteToUDP([]byte("Hello, client!"), u.addr)
			if err != nil {
				return err
			}
		}
	}
}

func (u *udpTest) Close() error {
	u.done()
	timeout := time.After(5 * time.Second)
	select {
	case <-u.ctx.Done():
		return nil
	case <-timeout:
		return errors.New("timeout")
	}
}

func testCheckSocketConnectivity(t *testing.T) {
	udpT := &udpTest{}
	var err error
	udpip, udpport := "127.0.0.1", "5000"
	udpT.addr, err = net.ResolveUDPAddr("udp", net.JoinHostPort(udpip, udpport))
	if err != nil {
		t.Fatal(err)
	}
	udpT.ctx, udpT.done = context.WithCancel(context.Background())
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return udpT.Run()
	})
	defer eg.Wait()
	defer udpT.Close()
	type args struct {
		nettype NetworkType
		host    string
		port    string
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want v1.NetworkConnectivity
	}{
		{
			name: "test tcp connectivity",
			args: args{
				nettype: TCP,
				host:    "www.baidu.com",
				port:    "443",
				timeout: time.Second,
			},
			want: v1.NetReachable,
		},
		{
			name: "test udp connectivity, which is unreachable",
			args: args{
				nettype: UDP,
				host:    "127.0.0.1",
				port:    "443",
				timeout: time.Second,
			},
			want: v1.NetUnreachable,
		},
		{
			name: "test tcp connectivity, which is reachable",
			args: args{
				nettype: UDP,
				host:    udpip,
				port:    udpport,
				timeout: time.Second,
			},
			want: v1.NetReachable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SocketConnectivityCheck(tt.args.nettype, tt.args.host, tt.args.port, tt.args.timeout); got != tt.want {
				t.Errorf("SocketConnectivityCheck() = %v, want %v", got, tt.want)
			} else {
				t.Logf("SocketConnectivityCheck() = %v", got)
			}
		})
	}
}
