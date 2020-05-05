package test

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/filter"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"testing"
)

func TestSortInsertInt(t *testing.T) {
	test := []int{345, 2, 2345, 35, 756, 309, 3, 6, 2, 35, 345, 798, 8734, 2}
	list := make([]int, 0, 2)
	list = filter.SortInsertInt(list, test...)
	t.Log(list)
}

func TestSortInsertUniqInt(t *testing.T) {
	test := []int{345, 2, 2345, 35, 756, 309, 3, 6, 2, 35, 345, 798, 8734, 2}
	list := make([]int, 0, 2)
	list = filter.SortInsertUniqInt(list, test...)
	t.Log(list)
}

func TestGetSortServers(t *testing.T) {
	conf, err := resolv.Load("config_test/nginx.conf")

	if err != nil {
		t.Log(err)
	}

	servers := filter.GetHTTPServers(conf, filter.ServersTaggerByPort)
	for _, server := range servers {
		t.Log(server.String())
	}
}
