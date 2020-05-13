package test

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/statistics"
	"testing"
)

func TestStatistics(t *testing.T) {
	conf, err := resolv.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	t.Log(statistics.HTTPServers(conf))
	t.Log(statistics.HTTPServerNames(conf))
	t.Log(statistics.HTTPPorts(conf))
	//t.Log(statistics.HTTPPortsSTR(conf))
	t.Log(statistics.HTTPLocations(conf))
	t.Log(statistics.StreamServers(conf))
	t.Log(statistics.StreamPorts(conf))
}
