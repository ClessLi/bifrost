package test

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"path/filepath"
	"testing"
)

func TestRel(t *testing.T) {
	conf, err := resolv.Load("./config_test/nginx.conf")
	if err != nil {
		t.Log(err)
	}

	fileList, err := conf.List()
	if err != nil {
		t.Log(err)
	}

	for _, s := range fileList {
		t.Log(filepath.Rel("F:\\GO_Project\\src\\go-nginx-conf-parser\\test", s))
	}
}
