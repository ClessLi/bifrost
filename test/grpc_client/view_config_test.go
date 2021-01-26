package grpc_client

import (
	"encoding/json"
	"fmt"
	"github.com/ClessLi/bifrost/pkg/client/auth"
	"github.com/ClessLi/bifrost/pkg/client/bifrost"
	ngJson "github.com/ClessLi/bifrost/pkg/json/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"golang.org/x/net/context"
	"os"
	"testing"
	"time"
)

var (
	authClient     *auth.Client
	bifrostClient  *bifrost.Client
	SvrName        = "bifrost-test"
	initErr        error
	token          string
	bifrostSvrAddr = "192.168.220.11:12321"
)

func init() {
	authSvrAddr := "192.168.220.11:12320"
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
	}
	t.Log(string(data))
}

func TestClientGC(t *testing.T) {
	defer bifrostClient.Close()
	jdata, err := bifrostClient.GetConfig(context.Background(), token, SvrName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jdata))
}

func TestClientUC(t *testing.T) {
	defer bifrostClient.Close()
	jdata, err := bifrostClient.GetConfig(context.Background(), token, SvrName)
	if err != nil {
		t.Fatal(err)
	}
	config, err := ngJson.Unmarshal(jdata)
	if err != nil {
		t.Fatal(err)
	}

	err = config.Insert(config.Children[0], nginx.TypeComment, fmt.Sprintf("#test for client.UpdateConfig at %s", time.Now()))
	if err != nil {
		t.Fatal(err)
	}
	// marshal to json data
	confJson, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := bifrostClient.UpdateConfig(context.Background(), token, SvrName, confJson)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(msg))
}

func TestClientUpdateConfigV2(t *testing.T) {

	defer bifrostClient.Close()
	jdata, err := bifrostClient.GetConfig(context.Background(), token, SvrName)
	if err != nil {
		t.Fatal(err)
	}
	l := loader.NewLoader()
	config, preventer, err := l.LoadFromJsonBytes(jdata)
	if err != nil {
		t.Fatal(err)
	}
	conf := configuration.NewConfiguration(config.(*parser.Config), preventer)
	q, err := conf.Query("comment:sep: :reg: pid.*")
	if err != nil {
		t.Fatal(err)
	}
	err = conf.InsertByQueryer(parser.NewComment(fmt.Sprintf("test for client.UpdateConfig with resolv.V2 at %s", time.Now()), false, q.Self().GetIndention()), q)
	if err != nil {
		t.Fatal(err)
	}

	httpQ, err := conf.Query("http")
	if err != nil {
		t.Fatal(err)
	}
	svrsQ, err := httpQ.QueryAll("server")
	if err != nil {
		t.Fatal(err)
	}

	for _, queryer := range svrsQ {
		portQ, err := queryer.Query("key:sep: :reg: listen\\s+80\\s*")
		if err != nil {
			continue
		}
		_, err = queryer.Query("key:sep: :reg: server_name\\s+test1.com\\s*")
		if err != nil {
			continue
		}
		err = conf.InsertByQueryer(parser.NewComment(fmt.Sprintf("test inline comment for client.UpdateConfig with resolv.V2 at %s", time.Now()), true, portQ.Self().GetIndention()), portQ)
		if err != nil {
			t.Fatal(err)
		}
		break
	}
	//for s, bytes := range conf.Dump() {
	//	t.Log(s, ":", string(bytes))
	//}
	//
	//t.Log(string(conf.View()))

	msg, err := bifrostClient.UpdateConfig(context.Background(), token, SvrName, conf.Json())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(msg))
}

func TestClientVS(t *testing.T) {
	defer bifrostClient.Close()
	jdata, err := bifrostClient.ViewStatistics(context.Background(), token, SvrName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jdata))
}

func TestClientStatus(t *testing.T) {
	defer bifrostClient.Close()
	jdata, err := bifrostClient.Status(context.Background(), token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jdata))
}

func TestClientWatchLog(t *testing.T) {
	defer bifrostClient.Close()
	timeout := time.After(time.Second * 20)
	logWatcher, err := bifrostClient.WatchLog(context.Background(), token, SvrName, "access.log")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = logWatcher.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
	}()
	for {
		select {
		case data := <-logWatcher.GetDataChan():
			//t.Logf(string(data))
			fmt.Println(string(data))
		case err := <-logWatcher.GetErrChan():
			t.Fatal(err.Error())
		case <-timeout:
			t.Log("test end")
			return
		}
	}
}
