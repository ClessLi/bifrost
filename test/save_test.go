package test

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"testing"
)

func TestSave(t *testing.T) {
	conf, err := nginx.Load("./config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	serr := nginx.Save(conf)

	if serr != nil {
		t.Log(serr)
	}
}

func TestVerifyAndSave(t *testing.T) {
	conf, err := nginx.Load("./config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	tmpConf := conf

	conf = nginx.NewConf(nil, "./config_test/nginx.conf")
	t.Log(tmpConf.String())
	t.Log(conf.String())
}
