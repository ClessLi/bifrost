package test

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"testing"
)

func TestList(t *testing.T) {
	path, caches, err := nginx.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	conf, err := caches.GetConfig(path)
	if err != nil {
		t.Log(err)
	}

	l, lErr := conf.List()
	if lErr != nil {
		t.Log(lErr)
	}

	t.Log(l)
}
