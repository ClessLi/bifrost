package resolv

import "testing"

func TestSave(t *testing.T) {
	conf, err := Load("../../test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	serr := Save(conf)

	if serr != nil {
		t.Log(serr)
	}
}
