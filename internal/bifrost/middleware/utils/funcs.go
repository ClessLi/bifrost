package utils

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"google.golang.org/grpc/peer"
	"net"
	"strings"
)

func GetClientIP(ctx context.Context) string {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "unknown"
	}
	if pr.Addr == net.Addr(nil) {
		return "unknown"
	}
	return pr.Addr.String()
}

func GetAuthnInfo(ctx context.Context) string {
	var info []string
	basicAuthn, ok := ctx.Value(v1.BasicAuthnKey).(string)
	if ok && basicAuthn != "" {
		pair := strings.SplitN(basicAuthn, ":", 2)

		if len(pair) == 2 {
			info = append(info, "Basic Authn: username:"+pair[0]+", password:"+pair[1])
		}

	}

	token, ok := ctx.Value(v1.BearerAuthnTokenKey).(string)
	if ok && token != "" {
		info = append(info, "Bearer Authn Token: "+token)
	}

	if len(info) == 0 {
		return "None"
	}

	return strings.Join(info, " | ")
}
