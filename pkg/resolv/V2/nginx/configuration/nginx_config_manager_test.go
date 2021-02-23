package configuration

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"sync"
	"testing"
	"time"
)

func exampleNewConfigManager() (*configManager, error) {
	l := loader.NewLoader()
	ctx, loopPreventer, err := l.LoadFromFilePath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		return nil, err
	}
	c := NewConfiguration(ctx.(*parser.Config), loopPreventer, new(sync.RWMutex))
	manager := NewNginxConfigurationManager(l, c, ".", "", 1, 7, new(sync.RWMutex))
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
	err = manager.regularlyBackup(time.Second, make(chan int))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	err = manager.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
