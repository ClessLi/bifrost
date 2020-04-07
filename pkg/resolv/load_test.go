package resolv

import (
	"encoding/json"
	"testing"
)

func TestLoad(t *testing.T) {
	conf, err := Load("../../test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	t.Log(conf.String())
}

func TestLoadServers(t *testing.T) {
	conf, err := Load("../../test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	t.Log(conf.Servers())

	for _, server := range conf.Servers() {
		t.Log(server.String())
	}
}

func TestLoadServer(t *testing.T) {
	conf, err := Load("../../test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	t.Log(conf.Server().String())

}

func TestMarshalJSON(t *testing.T) {
	conf, err := Load("../../test/nginx.conf")

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
	//l.Add(k)
	//l.Add(k2)
	//s1.Add(l)
	//h.Add(i)
	//h.Add(s1)
	//conf.Add(h)

	//j, jerr := json.MarshalIndent(conf, "", "    ")
	j, jerr := json.Marshal(conf)

	if jerr != nil {
		t.Log(jerr)
	}

	t.Log(string(j))
}

func TestConfig_UnmarshalJSON(t *testing.T) {
	jdata := ``
	conf := NewConf(nil, "")
	//err := conf.UnmarshalJSON([]byte(jdata))
	err := json.Unmarshal([]byte(jdata), &conf)
	if err != nil {
		t.Log(err)
	}

	t.Log(conf)
}
