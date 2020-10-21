package test

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	nginxStatistics "github.com/ClessLi/bifrost/pkg/statistics/nginx"
	"os"
	"testing"
)

func TestStatistics(t *testing.T) {
	path, caches, err := nginx.Load("config_test/nginx.conf")
	//conf, err := nginx.Load("filter_bug_test/nginx.conf")

	if err != nil {
		t.Log(err)
		os.Exit(1)
	}

	conf, err := caches.GetConfig(path)
	if err != nil {
		t.Log(err)
		os.Exit(2)
	}
	//t.Log(nginxStatistics.HTTPServersNum(conf))
	//t.Log(nginxStatistics.HTTPServerNames(conf))
	t.Log(nginxStatistics.HTTPPorts(conf))
	//t.Log(statistics.HTTPPortsSTR(conf))
	//t.Log(nginxStatistics.HTTPLocationsNum(conf))
	//t.Log(nginxStatistics.StreamServersNum(conf))
	//t.Log(nginxStatistics.StreamServers(conf))
	t.Log(nginxStatistics.HTTPServers(conf))
}
