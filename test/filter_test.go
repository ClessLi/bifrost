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

	keykw := nginx.NewKeyWords("key", "server_name", `^.*com.*`, true)
	svrkw := nginx.NewKeyWords("server", "", "", true, keykw)
	servers := conf.QueryAll(svrkw)
	for _, server := range servers {
		t.Log(server.String())
	}
}
