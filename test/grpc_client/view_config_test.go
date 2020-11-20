package grpc_client

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"testing"
)

func TestClientVC(t *testing.T) {
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

	ret, err := opClient.ViewConfig(context.Background(), &opReq)
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
		// view config
		t.Log(string(item.Ret))
		// data: '# test for api.UpdateConfig
		//       # user  nobody;
		//       worker_processes 1;
		//       ...'
	}

}
