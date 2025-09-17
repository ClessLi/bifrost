package test

import (
	"testing"

	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
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
	path, caches, err := nginx.Load("config_test/nginx.conf")
	if err != nil {
		t.Log(err)
	}

	conf, err := caches.GetConfig(path)
	if err != nil {
		t.Log(err)

		return
	}
	servers := nginx.GetHTTPServers(conf, nginx.ServerName, nginx.ServerPort)
	// servers := resolv.GetHTTPServers(conf, resolv.ServerName)
	// servers := statistics.GetHTTPServers(conf, statistics.ServerName)
	for _, server := range servers {
		// caches := nginx.NewCaches()
		// t.Log(server.string(&caches))
		t.Log(server.String())
	}
}
