package web_server_manager

import (
	"bytes"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/apsdehal/go-logger"
	"os"
	"testing"
	"time"
)

func TestNewNginxManager(t *testing.T) {
	utils.InitLogger(os.Stdout, logger.DebugLevel)
	info := WebServerConfigInfo{
		Name:           "config_test",
		Type:           NGINX,
		BackupCycle:    7,
		BackupSaveTime: 1,
		BackupDir:      "../../../../../test/config_test",
		ConfPath:       "F:\\GO_Project\\src\\bifrost\\test\\config_test\\nginx.conf",
		VerifyExecPath: "",
	}
	controller := newWebServerConfigServiceController(info)
	controller.SetState(Normal)
	time.Sleep(time.Second * 35)
	service := controller.GetService()
	t.Log(service.DisplayWebServerVersion())
	data, err := service.DisplayConfig()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(data))
	}
	controller.SetState(Disabled)
}

func TestBytesBuffer(t *testing.T) {
	buf := bytes.NewBuffer([]byte("test"))
	t.Log(string(buf.Bytes()))
	t.Log(string(buf.Bytes()))
}
