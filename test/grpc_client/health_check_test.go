package grpc_client

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	conn, _ := grpc.Dial(bifrostSvrAddr, grpc.WithInsecure())
	defer conn.Close()
	c := grpc_health_v1.NewHealthClient(conn)
	reply, _ := c.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	t.Log(reply.Status)
}
