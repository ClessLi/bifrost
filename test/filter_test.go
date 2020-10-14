package test

import (
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"testing"
)

func TestFilter(t *testing.T) {
	conf, _, err := nginx.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	//keykw := nginx.NewKeyWords("key", "server_name", `^.*com.*`, true, false)
	//svrkw := nginx.NewKeyWords("server", "", "", true, false, keykw)
	//httpServers := conf.QueryAllByKeywords(svrkw)

	//http := nginx.GetHTTP(conf)
	//httpServers := http.QueryAll(nginx.TypeServer, false, "server_name:com")
	//t.Log("Http Servers")
	//for _, server := range httpServers {
	//	t.Log(server.String())
	//}
	t.Log("All Servers")
	allServers := conf.QueryAll(nginx.TypeServer, true, "server_name:^open.*$", "listen:^80$")
	for _, server := range allServers {
		t.Log(server.String())
	}
}

func TestParams(t *testing.T) {
	//conf, err := nginx.Load("test_circle_load/nginx.conf")
	conf, _, err := nginx.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	//for _, param := range conf.Params() {
	//	t.Log(param.String())
	//}
	servers := conf.QueryAll(nginx.TypeServer, true, "server_name:^.*com")
	for _, server := range servers {
		fmt.Printf("server: server_name %s\n", server.Query(nginx.TypeKey, false, "server_name").(*nginx.Key).Value)
		for _, parser := range server.(*nginx.Server).Params() {
			fmt.Println(parser.String())
		}
	}

	//kw := nginx.NewKeyWords(nginx.TypeEvents, "", "", false, true)
	//event := conf.QueryByKeywords(kw).(*nginx.Events)
	//for _, param := range event.Params() {
	//	t.Log(param.String())
	//}

	//http := nginx.GetHTTP(conf)
	//for _, param := range http.Params() {
	//	t.Log(param.String())
	//}
}
