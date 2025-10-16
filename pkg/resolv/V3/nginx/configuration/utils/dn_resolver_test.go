package utils

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestNewDNSClient(t *testing.T) {
	type args struct {
		dnsIP string
	}
	tests := []struct {
		name string
		args args
		want DomainNameResolver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDNSClient(tt.args.dnsIP); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDNSClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIPv4Hosts(t *testing.T) {
	type args struct {
		hosts map[string][]net.IP
	}
	tests := []struct {
		name string
		args args
		want DomainNameResolver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIPv4Hosts(tt.args.hosts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIPv4Hosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDomainNameResolver(t *testing.T) {
	type args struct {
		r DomainNameResolver
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDomainNameResolver(tt.args.r)
		})
	}
}

func Test_dnsClient_ResolveToIPv4s(t *testing.T) {
	type fields struct {
		isTCP        bool
		dailTimeout  time.Duration
		readTimeout  time.Duration
		writeTimeout time.Duration
		dnsHost      string
		dnsPort      int
	}
	type args struct {
		domainName string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantIpv4s []string
		wantErr   bool
	}{
		{
			name: "resolve test, without root domain dot",
			fields: fields{
				isTCP:        false,
				dailTimeout:  time.Second * 10,
				readTimeout:  time.Second * 10,
				writeTimeout: time.Second * 10,
				dnsHost:      "114.114.114.114",
				dnsPort:      53,
			},
			args: args{
				domainName: "baidu.com",
			},
			wantIpv4s: []string{"39.156.70.37", "220.181.7.203"},
		},
		{
			name: "resolve test",
			fields: fields{
				isTCP:        false,
				dailTimeout:  time.Second * 10,
				readTimeout:  time.Second * 10,
				writeTimeout: time.Second * 10,
				dnsHost:      "114.114.114.114",
				dnsPort:      53,
			},
			args: args{
				domainName: "baidu.com.",
			},
			wantIpv4s: []string{"220.181.7.203", "39.156.70.37"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dnsClient{
				isTCP:        tt.fields.isTCP,
				dailTimeout:  tt.fields.dailTimeout,
				readTimeout:  tt.fields.readTimeout,
				writeTimeout: tt.fields.writeTimeout,
				dnsHost:      tt.fields.dnsHost,
				dnsPort:      tt.fields.dnsPort,
			}
			gotIpv4s, err := d.ResolveToIPv4s(tt.args.domainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveToIPv4s() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(gotIpv4s, tt.wantIpv4s) {
				t.Errorf("ResolveToIPv4s() gotIpv4s = %v, want %v", gotIpv4s, tt.wantIpv4s)
			}
		})
	}
}

func Test_hostList_ResolveToIPv4s(t *testing.T) {
	type args struct {
		domainName string
	}
	tests := []struct {
		name      string
		h         hostList
		args      args
		wantIpv4s []string
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIpv4s, err := tt.h.ResolveToIPv4s(tt.args.domainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveToIPv4s() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(gotIpv4s, tt.wantIpv4s) {
				t.Errorf("ResolveToIPv4s() gotIpv4s = %v, want %v", gotIpv4s, tt.wantIpv4s)
			}
		})
	}
}

func Test_newDNSClient(t *testing.T) {
	type args struct {
		dnsIP   string
		dnsPort int
		isTCP   bool
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want *dnsClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDNSClient(tt.args.dnsIP, tt.args.dnsPort, tt.args.isTCP, tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDNSClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
