package configuration

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"testing"
)

func exampleNewConfigManager() (*configManager, error) {
	manager, err := NewNginxConfigurationManager("test", "F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		return nil, err
	}
	return manager.(*configManager), nil
}

func TestConfigManager_SaveWithCheck(t *testing.T) {

	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	config := manager.configuration
	c, err := config.Query("comment:sep: :reg: .*")
	if err != nil {
		t.Fatal(err)
	}
	newC := parser.NewComment("testConfigManager", false, c.Self().GetIndention())
	err = config.InsertByQueryer(newC, c)
	if err != nil {
		t.Fatal(err)
	}
	//jsonData := config.Json()
	//fmt.Println(string(jsonData))
	//fmt.Println(string(config.View()))
	err = manager.SaveWithCheck()
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigManager_Backup(t *testing.T) {
	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	err = manager.Backup("", 7, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigManager_Read(t *testing.T) {

	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(manager.Read()))
}

func TestConfigManager_ReadJson(t *testing.T) {

	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(manager.ReadJson()))
}

func TestConfigManager_ReadStatistics(t *testing.T) {

	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(manager.ReadStatistics()))
}

func TestConfigManager_UpdateFromJsonBytes(t *testing.T) {

	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}

	config := manager.configuration
	c, err := config.Query("comment:sep: :reg: .*")
	if err != nil {
		t.Fatal(err)
	}
	newC := parser.NewComment("testConfigManager_update...", false, c.Self().GetIndention())
	err = config.InsertByQueryer(newC, c)
	if err != nil {
		t.Fatal(err)
	}

	manager2, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}

	err = manager2.UpdateFromJsonBytes(config.Json())
	if err != nil {
		t.Fatal(err)
	}
}
