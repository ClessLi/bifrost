package grpc_client

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
)

func TestGRPCLoginAndVerify(t *testing.T) {
	authSvcAddr := "192.168.220.11:12320"
	authConn, acErr := grpc.Dial(authSvcAddr, grpc.WithInsecure())
	if acErr != nil {
		panic("connect error")
	}
	defer authConn.Close()
	authClient := authpb.NewAuthServiceClient(authConn)
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
	// token
	t.Log(authResp.Token, authResp.Err)
	// msg: (Token)'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9....' (Err)''
	verifyResp, err := authClient.Verify(context.Background(), &authpb.VerifyRequest{Token: authResp.Token})
	t.Log(verifyResp.Passed, verifyResp.Err)
	// msg: (Passed)'true' (Err)''
}
