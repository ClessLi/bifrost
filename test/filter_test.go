package test

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"testing"
)

func TestFilter(t *testing.T) {
	conf, err := nginx.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	keykw := nginx.NewKeyWords("key", "server_name", `^.*com.*`, true, false)
	svrkw := nginx.NewKeyWords("server", "", "", true, false, keykw)
	servers := conf.QueryAll(svrkw)
	for _, server := range servers {
		t.Log(server.String())
	}
}

func TestParams(t *testing.T) {
	conf, err := nginx.Load("test_circle_load/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	for _, param := range conf.Params() {
		t.Log(param.String())
	}

}
