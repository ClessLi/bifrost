package test

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"testing"
)

func TestSave(t *testing.T) {
	conf, err := resolv.Load("nginx.conf")

	if err != nil {
		t.Log(err)
	}

	serr := resolv.Save(conf)

	if serr != nil {
		t.Log(serr)
	}
}
