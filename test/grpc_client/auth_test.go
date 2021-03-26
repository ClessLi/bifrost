package grpc_client

import (
	"github.com/ClessLi/bifrost/pkg/client/auth"
	"golang.org/x/net/context"
	"testing"
)

func TestGRPCLoginAndVerify(t *testing.T) {
	authSvcAddr := "192.168.220.11:12320"
	client, err := auth.NewClientFromGRPCServerAddress(authSvcAddr)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer client.Close()
	token, loginErr := client.Login(context.Background(), "heimdall", "Bultgang", false)
	if loginErr != nil {
		t.Fatal(loginErr)
		return
	}
	// token
	t.Log(token)
	// msg: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9....'
	pass, verifyErr := client.Verify(context.Background(), token)
	if verifyErr != nil {
		t.Fatal(verifyErr)
		return
	}
	t.Log(pass)
	// msg: 'true'
}
