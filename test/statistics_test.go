package test

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	nginxStatistics "github.com/ClessLi/bifrost/pkg/statistics/nginx"
	"testing"
)

func TestStatistics(t *testing.T) {
	//conf, err := resolv.Load("config_test/nginx.conf")
	conf, err := nginx.Load("filter_bug_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	t.Log(nginxStatistics.HTTPServersNum(conf))
	t.Log(nginxStatistics.HTTPServerNames(conf))
	t.Log(nginxStatistics.HTTPPorts(conf))
	//t.Log(statistics.HTTPPortsSTR(conf))
	t.Log(nginxStatistics.HTTPLocationsNum(conf))
	t.Log(nginxStatistics.StreamServersNum(conf))
	t.Log(nginxStatistics.StreamPorts(conf))
}
