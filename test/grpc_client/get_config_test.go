package grpc_client

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	ngJson "github.com/ClessLi/bifrost/pkg/json/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"testing"
)

func TestClientGetConfig(t *testing.T) {
	serviceAddr := "192.168.220.11:12321"
	conn, cStreamErr := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if cStreamErr != nil {
		panic("connect error")
	}

	defer conn.Close()
	authClient := bifrostpb.NewAuthServiceClient(conn)
	opClient := bifrostpb.NewOperationServiceClient(conn)
	authReq := bifrostpb.AuthRequest{
		Username:  "heimdall",
		Password:  "Bultgang",
		Unexpired: false,
	}
	authResp, err := authClient.Login(context.Background(), &authReq)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	opReq := bifrostpb.OperateRequest{
		Token:    authResp.Token,
		Location: "bifrost-test",
	}
	confResp, err := opClient.GetConfig(context.Background(), &opReq)
	if err != nil {
		t.Fatal(err)
		return
	}

	var config *nginx.Config
	for {
		item, err := confResp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Log(err)
			return
		}

		// get config
		config, err = ngJson.Unmarshal(item.Ret.JData)
		if err != nil {
			t.Fatal(err)
			return
		}

		// print path of config, after nginx json unmarshal
		t.Log(config.Value)
		// path: '/usr/local/openresty/nginx/conf/nginx.conf'
		// print config json data from response
		t.Log(string(item.Ret.JData))
		// data: '{"config":{"value":"/usr/local/openresty/nginx/conf/nginx.conf","param":[{"comments":"test for api.UpdateConfig","inline":false},...]}}'
	}
}
