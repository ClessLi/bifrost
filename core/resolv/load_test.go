package resolv

import "testing"

func TestLoad(t *testing.T) {
	conf, err := Load("../../test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	t.Log(conf.String())
}
