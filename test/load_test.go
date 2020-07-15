package test

import (
	"encoding/json"
	ngJson "github.com/ClessLi/bifrost/pkg/json/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"io/ioutil"
	"testing"
)

func TestLoad(t *testing.T) {
	//conf, err := resolv.Load("config_test/nginx.conf")
	conf, err := nginx.Load("test_circle_load/nginx.conf")

	if err != nil {
		t.Log(err)
	} else {
		t.Log(conf.String())
	}

}

func TestLoadServers(t *testing.T) {
	conf, err := nginx.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	t.Log(conf.Servers())

	for _, server := range conf.Servers() {
		t.Log(server.String())
	}
}

func TestLoadServer(t *testing.T) {
	conf, err := nginx.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	t.Log(conf.Server().String())

}

func TestMarshalJSON(t *testing.T) {
	//conf, err := resolv.Load("config_test/nginx.conf")
	conf, err := nginx.Load("test_circle_load/nginx.conf")

	if err != nil {
		t.Log(err)
	}
	//conf := NewConf(nil, "test.conf")
	//h := NewHttp()
	//s1 := NewServer()
	//l := NewLocation("/")
	//k := NewKey("$test","$testV")
	//k2 := NewKey("$remote_addr - $remote_user [$time_local] \"$request\" ","")
	//i := NewInclude("../../test/mime.types","../../test/mime.types")
	//l.add(k)
	//l.add(k2)
	//s1.add(l)
	//h.add(i)
	//h.add(s1)
	//conf.add(h)

	j, jerr := json.MarshalIndent(conf, "", "    ")
	//j, jerr := json.Marshal(conf)

	if jerr != nil {
		t.Log(jerr)
	}

	t.Log(string(j))
}
func TestConfig_UnmarshalJSON(t *testing.T) {
	jdata, rerr := ioutil.ReadFile("test_circle_load/test.json")
	if rerr != nil {
		t.Log(rerr)
	}
	//conf := NewConf(nil, "")
	//err := conf.UnmarshalToJSON([]byte(jdata))
	//err := json.unmarshal([]byte(jdata), &conf)
	conf, err := ngJson.Unmarshal(jdata)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(conf)
		t.Log(conf.String())
	}

}
