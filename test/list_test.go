package test

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"testing"
)

func TestList(t *testing.T) {
	conf, err := resolv.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	l, lErr := conf.List()
	if lErr != nil {
		t.Log(lErr)
	}

	t.Log(l)
}
