package grpc_client

import (
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/bifrost/pkg/client/bifrost"
	"github.com/ClessLi/skirnir/pkg/discover"
	"golang.org/x/net/context"
	"testing"
)

func TestDiscoverBifrost(t *testing.T) {
	client, err := discover.NewKitConsulDiscoveryClient("192.168.43.151", 8500)
	if err != nil {
		t.Fatal(err)
		return
	}
	instance := client.DiscoverServices("com.github.ClessLi.api.bifrost", config.KitLogger)
	t.Log(instance)
}

func TestClientFromConsul(t *testing.T) {
	bifrostClient.Close()
	client, err := bifrost.NewClientFromConsul("192.168.43.151", 8500)
	if err != nil {
		t.Fatal(err)
		return
	}

	data, err := client.ViewStatistics(context.Background(), token, "bifrost-test")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(data))
}
