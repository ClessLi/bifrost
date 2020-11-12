package grpc_client

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"testing"
)

func TestClientVC(t *testing.T) {
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
