package utils

import (
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

func udpConnectivity(ipv4, port string, timeout time.Duration) v1.NetworkConnectivity {
	// TODO: Improve the udp network connectivity detection mechanism suitable for the "Windows" OS
	return v1.NetUnknown
}
