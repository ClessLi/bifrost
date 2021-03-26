package configuration

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"sync"
	"testing"
	"time"
)

func exampleNewConfigManager() (*configManager, error) {
	c, err := NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		return nil, err
	}
	manager := NewNginxConfigurationManager(loader.NewLoader(), c, ".", "", 1, 7, new(sync.RWMutex))
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
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		err = manager.regularlyBackup(time.Second, make(chan int))
		if err != nil {
			t.Fatal(err)
		}
		wg.Done()
	}()
	wg.Wait()
	//err = manager.Stop()
	//if err != nil {
	//	t.Fatal(err)
	//}
}
