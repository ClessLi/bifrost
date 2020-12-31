package web_server_manager

import (
	"bytes"
	"testing"
	"time"
)

func TestNewNginxManager(t *testing.T) {
	manager := NewNginxManager("test", 7, 1, "../../../../../test/config_test", "../../../../../test/config_test/nginx.conf", "")
	err := manager.ManagementStart()
	if err != nil {
		t.Fatal(err)
		return
	}
	time.Sleep(time.Minute)
	err = manager.ManagementStop()
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestBytesBuffer(t *testing.T) {
	buf := bytes.NewBuffer([]byte("test"))
	t.Log(buf.Bytes())
	t.Log(buf.Bytes())
}
