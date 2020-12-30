package web_server_manager

import (
	"testing"
	"time"
)

func TestNewNginxManager(t *testing.T) {
	manager := NewNginxManager("test", 7, 1, "../../../../../test/config_test", "../../../../../test/config_test/nginx.conf", "")
	err := manager.Start()
	if err != nil {
		t.Fatal(err)
		return
	}
	time.Sleep(time.Minute)
	err = manager.Stop()
	if err != nil {
		t.Fatal(err)
		return
	}
}
