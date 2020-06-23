package test

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"path/filepath"
	"testing"
)

func TestRel(t *testing.T) {
	conf, err := nginx.Load("./config_test/nginx.conf")
	if err != nil {
		t.Log(err)
	}

	fileList, err := conf.List()
	if err != nil {
		t.Log(err)
	}

	for _, s := range fileList {
		t.Log(filepath.Rel("F:\\GO_Project\\src\\bifrost\\test", s))
	}
}
