package web_server_manager

import "testing"

func TestWebServerConfigServicesController_GetServicesHandler(t *testing.T) {
	info := WebServerConfigInfo{
		Name:           "config_test",
		Type:           NGINX,
		BackupCycle:    7,
		BackupSaveTime: 1,
		BackupDir:      "../../../../../test/config_test",
		ConfPath:       "F:\\GO_Project\\src\\bifrost\\test\\config_test\\nginx.conf",
		VerifyExecPath: "",
	}
	svcsController := NewWebServerConfigServicesController(info)
	handler := svcsController.GetServicesHandler()
	s, err := handler.DisplayStatus("config_test")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(s)
	}
	err = svcsController.Start()
	if err != nil {
		t.Fatal(err)
	}

	s, err = handler.DisplayStatus("config_test")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(s)
	}
	data, err := handler.DisplayConfig("config_test")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(data))
	}
	err = svcsController.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
