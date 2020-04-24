package test

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"testing"
)

func TestFilter(t *testing.T) {
	conf, err := resolv.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	keykw := resolv.NewKeyWords("key", "server_name", `^.*com.*`, true)
	srvkw := resolv.NewKeyWords("server", "", "", true, keykw)
	servers := conf.Filter(srvkw)
	for _, server := range servers {
		t.Log(server.String())
	}
}
