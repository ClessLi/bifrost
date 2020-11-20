package grpc_client

import (
	"encoding/json"
	"fmt"
	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	ngJson "github.com/ClessLi/bifrost/pkg/json/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"testing"
)

func TestClientUpdateConfig(t *testing.T) {
	authSvcAddr := "192.168.220.11:12320"
	authConn, acErr := grpc.Dial(authSvcAddr, grpc.WithInsecure())
	if acErr != nil {
		panic("connect error")
	}
	defer authConn.Close()
	serviceAddr := "192.168.220.11:12321"
	conn, cStreamErr := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if cStreamErr != nil {
		panic("connect error")
	}

	defer conn.Close()
	authClient := authpb.NewAuthServiceClient(authConn)
	opClient := bifrostpb.NewBifrostServiceClient(conn)
	authReq := authpb.AuthRequest{
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
		Token:   authResp.Token,
		SvrName: "bifrost-test",
	}
	confResp, err := opClient.GetConfig(context.Background(), &opReq)
	if err != nil {
		t.Fatal(err)
		return
	}

	var config *nginx.Config
	// get config
	for {
		item, err := confResp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Log(err)
			return
		}

		config, err = ngJson.Unmarshal(item.Ret.JData)
		if err != nil {
			t.Fatal(err)
			return
		}

	}

	// modify config
	err = config.Insert(config.Children[0], nginx.TypeComment, "#test for bifrostpb.UpdateConfig")
	if err != nil {
		t.Fatal(err)
		return
	}
	// marshal to json data
	jdata, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
		return
	}
	//t.Log(string(jdata))
	// init request
	confReq := &bifrostpb.ConfigRequest{
		Token:   authResp.Token,
		SvrName: "bifrost-test",
		Req:     &bifrostpb.Config{},
	}
	// init client stream
	upStream, err := opClient.UpdateConfig(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}
	chunckSize := 4 * 1024 * 1024
	// update config with client stream
	for i := 0; i < len(jdata); i += chunckSize {
		if i+chunckSize >= len(jdata) {
			confReq.Req.JData = jdata[i:]
			fmt.Println("end")
		} else {
			confReq.Req.JData = jdata[i : i+chunckSize]
			fmt.Println("idx", i)
		}
		err = upStream.Send(confReq)
	}
	upRet, err := upStream.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
		return
	}

	// print response of update config
	t.Log(string(upRet.Ret), upRet.Err)
	// response: (Ret)'update config success'
}
