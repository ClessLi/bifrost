package utils

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ClessLi/bifrost/internal/pkg/code"

	"github.com/marmotedu/errors"
	"github.com/miekg/dns"
)

type DomainNameResolver interface {
	ResolveToIPv4s(domainName string) (ipv4s []string, err error)
}

type hostList map[string][]string

func (h hostList) ResolveToIPv4s(domainName string) (ipv4s []string, err error) {
	if h[domainName] == nil {
		return nil, errors.WithCode(code.ErrV3DomainNameResolutionFailed, "the domain name resolution record for `%s` does not exist", domainName)
	}

	return h[domainName], nil
}

type dnsClient struct {
	isTCP        bool
	dailTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	dnsHost      string
	dnsPort      int
}

func (d *dnsClient) ResolveToIPv4s(domainName string) (ipv4s []string, err error) {
	client := &dns.Client{
		DialTimeout:  d.dailTimeout,
		ReadTimeout:  d.readTimeout,
		WriteTimeout: d.writeTimeout,
	}
	if d.isTCP {
		client.Net = "tcp"
	}
	if domainName[len(domainName)-1] != '.' {
		domainName = domainName + "."
	}
	r, _, err := client.Exchange(new(dns.Msg).SetQuestion(domainName, dns.TypeA), fmt.Sprintf("%s:%d", d.dnsHost, d.dnsPort))
	if err != nil {
		return ipv4s, errors.WithCode(code.ErrV3DomainNameResolutionFailed, err.Error())
	}
	for _, rr := range r.Answer {
		if ar, ok := rr.(*dns.A); ok {
			ipv4 := ar.A.To4().String()
			if ipv4 != "<nil>" {
				ipv4s = append(ipv4s, ipv4)
			}
		}
	}

	return
}

func NewIPv4Hosts(hosts map[string][]net.IP) DomainNameResolver {
	h := make(hostList)
	for dn, ips := range hosts {
		for _, ip := range ips {
			if ip.To4().String() != "<nil>" {
				h[dn] = append(h[dn], ip.To4().String())
			}
		}
	}

	return h
}

func newDNSClient(dnsIP string, dnsPort int, isTCP bool, timeout time.Duration) *dnsClient {
	client := &dnsClient{
		isTCP:   isTCP,
		dnsHost: dnsIP,
		dnsPort: dnsPort,
	}
	if timeout > 0 {
		client.readTimeout = timeout
		client.writeTimeout = timeout
		client.dailTimeout = timeout
	}

	return client
}

func NewDNSClient(dnsIP string) DomainNameResolver {
	return newDNSClient(dnsIP, 53, false, time.Second*10)
}

var dnResolver DomainNameResolver

var dnResolverRWLock = new(sync.RWMutex)

func SetDomainNameResolver(r DomainNameResolver) {
	dnResolverRWLock.Lock()
	defer dnResolverRWLock.Unlock()
	dnResolver = r
}
