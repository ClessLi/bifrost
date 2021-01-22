package configuration

import (
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
	"testing"
)

func exampleNewConfigManager() (*configManager, error) {
	manager, err := NewConfigManager("test", "F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		return nil, err
	}
	return manager, nil
}

func TestNewConfigManager(t *testing.T) {
	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	config := manager.GetConfiguration()
	fmt.Println(string(config.StatisticsByJson()))
	//fmt.Println(string(config.View()))
	//for path, bytes := range config.Dump() {
	//	fmt.Printf("%s:\n%s\n", path, bytes)
	//}
	//fmt.Println(string(config.Json()))
}

func TestConfiguration_View(t *testing.T) {
	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	config := manager.GetConfiguration()
	fmt.Println(string(config.View()))
}

func TestConfiguration_Dump(t *testing.T) {
	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	config := manager.GetConfiguration()
	for path, bytes := range config.Dump() {
		fmt.Printf("%s:\n%s\n", path, bytes)
	}
}

func TestConfiguration_Json(t *testing.T) {
	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	config := manager.GetConfiguration()
	fmt.Println(string(config.Json()))
}

func TestConfiguration_InsertByQueryer(t *testing.T) {

	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	config := manager.GetConfiguration()
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
	fmt.Println(string(config.View()))
}

func TestConfiguration_InsertLoopConfig(t *testing.T) {

	manager, err := exampleNewConfigManager()
	if err != nil {
		t.Fatal(err)
	}
	config := manager.GetConfiguration()
	location, err := config.Query("location:sep: /test2")
	if err != nil {
		t.Fatal(err)
	}
	locationConfig, err := config.Query("config:sep: F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\conf.d\\location.conf")
	if err != nil {
		t.Fatal(err)
	}
	svr := location.fatherContext()
	loopInclude := parser.NewContext("./conf.d/location.conf", parser_type.TypeInclude, svr.GetIndention().NextIndention())
	err = config.InsertByQueryer(loopInclude, location)
	if err != nil {
		t.Fatal(err)
	}
	q := queryer{
		Parser:    nil,
		fatherCtx: loopInclude,
		selfIndex: 0,
	}
	//err = loopInclude.Insert(locationConfig.Self(), 0)
	err = config.InsertByQueryer(locationConfig.Self(), q)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(config.View()))
}
