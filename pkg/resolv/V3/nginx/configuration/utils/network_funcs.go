package utils

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"

	"github.com/marmotedu/errors"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type NetworkType int

const (
	TCP NetworkType = iota
	UDP
)

var ipAddrMustCompile = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)

func ResolveDomainNameToIPv4(name string) ([]string, error) {
	// return itself, if the domain name is an IPv4 string
	if ipAddrMustCompile.MatchString(name) {
		return []string{name}, nil
	}
	dnResolverRWLock.RLock()
	defer dnResolverRWLock.RUnlock()

	if dnResolver == nil {
		return nil, errors.WithCode(code.ErrV3InvalidOperation, "Domain Name Resolver not initialized")
	}

	return dnResolver.ResolveToIPv4s(name)
}

// SocketConnectivityCheck check TCP and UDP socket connectivity.
func SocketConnectivityCheck(nettype NetworkType, host, port string, timeout time.Duration) v1.NetworkConnectivity {
	if nettype == TCP {
		return tcpConnectivity(host, port, timeout)
	} else if nettype == UDP {
		return udpConnectivity(host, port, timeout)
	}

	return v1.NetUnknown
}

func tcpConnectivity(host, port string, timeout time.Duration) v1.NetworkConnectivity {
	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return v1.NetUnreachable
	}
	defer conn.Close()

	return v1.NetReachable
}

func listenICMPUnreachable(d time.Duration, ready context.CancelFunc) ([]byte, error) {
	// 监听ipv4的icmp报文
	c, _ := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	defer c.Close()
	buf := make([]byte, 1500)
	// 设置超时时间
	_ = c.SetReadDeadline(time.Now().Add(d))
	// listening ICMP is ready
	ready()

	// 读取会话中信息
	n, _, _ := c.ReadFrom(buf)
	if n == 0 {
		// 如果在会话中没有读取到任何内容则返回空
		return nil, nil
	}
	// 解析报文内容
	msg, err := icmp.ParseMessage(ipv4.ICMPTypeDestinationUnreachable.Protocol(), buf[:n])
	if err != nil {
		return nil, err
	}
	// 如果报文类型为icmp不可达类型的报文则返回报文内容
	if msg.Type == ipv4.ICMPTypeDestinationUnreachable {
		// TODO: 判断类型，避免异常退出
		body := msg.Body.(*icmp.DstUnreach)

		return body.Data, nil
	}
	// 如果在会话中没有读取到任何内容则返回空
	return nil, nil
}

func parseUnreachUDP(unreachData []byte) (ip, port string) {
	// 解析udp不可达的报文
	ipHeader, err := ipv4.ParseHeader(unreachData)
	if err != nil {
		fmt.Printf("Failed to parse IP header:%s", err.Error())

		return
	}

	ip = ipHeader.Dst.String() // 头部的目标地址则为我们探测的目标地址
	// 创建一个切片存放数据，长度为ipv4报文包头的长度
	dataBytes := unreachData[ipv4.HeaderLen:]
	// 解析 UDP 数据包
	var udpHeader []byte
	// 出去前面8个字节的ipv4头
	udpHeader = append(udpHeader, dataBytes[:8]...)
	// 目的端口的位置，实测得到
	port = strconv.Itoa(int(binary.BigEndian.Uint16(udpHeader[2:4])))

	return
}
