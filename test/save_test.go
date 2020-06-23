package test

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"testing"
)

func TestSave(t *testing.T) {
	conf, err := resolv.Load("./config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	serr := resolv.Save(conf)

	if serr != nil {
		t.Log(serr)
	}
}

func TestVerifyAndSave(t *testing.T) {
	conf, err := resolv.Load("./config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	tmpConf := conf

	conf = resolv.NewConf(nil, "./config_test/nginx.conf")
	t.Log(tmpConf.String())
	t.Log(conf.String())
}
