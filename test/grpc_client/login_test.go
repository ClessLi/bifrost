package grpc_client

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
)

func TestGRPCLogin(t *testing.T) {
	serviceAddr := "192.168.220.11:12321"
	conn, cStreamErr := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if cStreamErr != nil {
		panic("connect error")
	}

	defer conn.Close()
	authClient := bifrostpb.NewAuthServiceClient(conn)
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
	// token
	t.Log(authResp.Token, authResp.Err)
	// msg: (Token)'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9....' (Err)''
}
