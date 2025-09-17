package test

import (
	"testing"

	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
)

func TestSave(t *testing.T) {
	path, caches, err := nginx.Load("./config_test/nginx.conf")
	if err != nil {
		t.Log(err)

		return
	}

	conf, err := caches.GetConfig(path)
	if err != nil {
		t.Log(err)

		return
	}

	_, serr := nginx.Save(conf)

	if serr != nil {
		t.Log(serr)
	}
}
