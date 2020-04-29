package test

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/filter"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"testing"
)

func TestStatistics(t *testing.T) {
	conf, err := resolv.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	t.Log(filter.HTTPServers(conf))
	t.Log(filter.HTTPServerNames(conf))
	t.Log(filter.HTTPPorts(conf))
	t.Log(filter.HTTPLocations(conf))
	t.Log(filter.StreamServers(conf))
	t.Log(filter.StreamPorts(conf))
}
