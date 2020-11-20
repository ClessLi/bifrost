package grpc_client

import (
	"encoding/json"
	"fmt"
	"github.com/ClessLi/bifrost/pkg/client/auth"
	"github.com/ClessLi/bifrost/pkg/client/bifrost"
	ngJson "github.com/ClessLi/bifrost/pkg/json/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"golang.org/x/net/context"
	"os"
	"testing"
)

var (
	authClient    *auth.Client
	bifrostClient *bifrost.Client
	SvrName       = "bifrost-test"
	initErr       error
	token         string
)

func init() {
	authSvrAddr := "192.168.220.11:12320"
	bifrostSvrAddr := "192.168.220.11:12321"
	username := "heimdall"
	password := "Bultgang"
	authClient, initErr = auth.NewClient(authSvrAddr)
	if initErr != nil {
		fmt.Println(initErr)
		os.Exit(1)
	}
	defer authClient.Close()
	bifrostClient, initErr = bifrost.NewClient(bifrostSvrAddr)
	if initErr != nil {
		fmt.Println(initErr)
		os.Exit(2)
	}
	token, initErr = authClient.Login(context.Background(), username, password, false)
	if initErr != nil {
		fmt.Println(initErr)
		os.Exit(3)
	}
}

func TestClientVC(t *testing.T) {
	defer bifrostClient.Close()
	data, err := bifrostClient.ViewConfig(context.Background(), token, SvrName)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(data))
}

func TestClientGC(t *testing.T) {
	defer bifrostClient.Close()
	jdata, err := bifrostClient.GetConfig(context.Background(), token, SvrName)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(jdata))
}

func TestClientUC(t *testing.T) {
	defer bifrostClient.Close()
	jdata, err := bifrostClient.GetConfig(context.Background(), token, SvrName)
	if err != nil {
		t.Fatal(err)
		return
	}
	config, err := ngJson.Unmarshal(jdata)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = config.Insert(config.Children[0], nginx.TypeComment, "#test for client.UpdateConfig")
	if err != nil {
		t.Fatal(err)
		return
	}
	// marshal to json data
	confJson, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
		return
	}

	msg, err := bifrostClient.UpdateConfig(context.Background(), token, SvrName, confJson)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(msg))
}

func TestClientVS(t *testing.T) {
	defer bifrostClient.Close()
	jdata, err := bifrostClient.ViewStatistics(context.Background(), token, SvrName)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(jdata))
}

func TestClientStatus(t *testing.T) {
	defer bifrostClient.Close()
	jdata, err := bifrostClient.Status(context.Background(), token)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(jdata))
}
