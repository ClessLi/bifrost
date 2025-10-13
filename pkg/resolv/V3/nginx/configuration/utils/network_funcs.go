package utils

import (
	"regexp"

	"github.com/ClessLi/bifrost/internal/pkg/code"

	"github.com/marmotedu/errors"
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

// TODO: check TCP and UDP socket connectivity
// func CheckSocketConnectivity(nettype, address string, timeout time.Duration) bool {
// 	conn, err := net.DialTimeout(nettype, address, timeout)
// 	if err != nil {
// 		return false
// 	}
// 	defer conn.Close()
// 	return true
// }
