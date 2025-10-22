package utils

import (
	"testing"
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

func testCheckSocketConnectivity(t *testing.T) {
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
			name: "test udp connectivity",
			args: args{
				nettype: UDP,
				host:    "127.0.0.1",
				port:    "443",
				timeout: time.Second,
			},
			want: v1.NetUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SocketConnectivityCheck(tt.args.nettype, tt.args.host, tt.args.port, tt.args.timeout); got != tt.want {
				t.Errorf("SocketConnectivityCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
