package test

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"testing"
)

func TestSortInsertInt(t *testing.T) {
	test := []int{345, 2, 2345, 35, 756, 309, 3, 6, 2, 35, 345, 798, 8734, 2}
	list := make([]int, 0, 2)
	list = nginx.SortInsertInt(list, test...)
	t.Log(list)
}

func TestSortInsertUniqInt(t *testing.T) {
	test := []int{345, 2, 2345, 35, 756, 309, 3, 6, 2, 35, 345, 798, 8734, 2}
	list := make([]int, 0, 2)
	list = nginx.SortInsertUniqInt(list, test...)
	t.Log(list)
}

func TestGetSortServers(t *testing.T) {
	conf, _, err := nginx.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	servers := nginx.GetHTTPServers(conf, nginx.ServerName, nginx.ServerPort)
	//servers := resolv.GetHTTPServers(conf, resolv.ServerName)
	//servers := statistics.GetHTTPServers(conf, statistics.ServerName)
	for _, server := range servers {
		t.Log(server.String())
	}
}
