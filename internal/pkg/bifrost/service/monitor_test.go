package service

import (
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service/web_server_manager"
	"testing"
	"time"
)

func TestNewSysInfo(t *testing.T) {
	info := web_server_manager.WebServerConfigInfo{
		Name:           "config_test",
		Type:           web_server_manager.NGINX,
		BackupCycle:    7,
		BackupSaveTime: 1,
		BackupDir:      "../../../../test/config_test",
		ConfPath:       "F:\\GO_Project\\src\\bifrost\\test\\config_test\\nginx.conf",
		VerifyExecPath: "",
	}
	svcsController := web_server_manager.NewWebServerConfigServicesController(info)
	err := svcsController.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer svcsController.Stop()
	handler := svcsController.GetServicesHandler()
	monitor := NewSysInfo(handler)
	err = monitor.Start()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 6)
	data, err := monitor.DisplayStatus()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	err = monitor.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
