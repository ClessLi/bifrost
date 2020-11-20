package grpc_client

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"testing"
)

func TestClientStatus(t *testing.T) {
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
	ret, err := opClient.Status(context.Background(), &opReq)
	if err != nil {
		t.Fatal(err)
		return
	}
	for {
		item, err := ret.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Log(err)
			return
		}

		// status
		t.Log(string(item.Ret))
		// status json data: '{"system":"centos 7.4.1708","time":"2020/11/12 10:15:55","cpu":"0.31","mem":"48.76","disk":"77.21","servers_status":["normal","abnormal"],"servers_version":["nginx version: openresty/1.13.6.2","nginx version: openresty/1.13.6.2"],"bifrost_version":"v1.0.1-alpha.1"}'
	}
}
