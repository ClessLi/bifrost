package grpc_health_v1

import (
	"context"
)

const (
	UNKNOWN HealthStatus = iota
	SERVING
	NOT_SERVING
	SERVICE_UNKNOWN
)

type HealthStatus int32

type Service interface {
	Check(ctx context.Context, service string) (HealthStatus, error)
	Watch(ctx context.Context, service string) (<-chan HealthStatus, error)
}

var statusStrings = map[HealthStatus]string{
	UNKNOWN:         "unknown",
	SERVING:         "serving",
	NOT_SERVING:     "not serving",
	SERVICE_UNKNOWN: "service unknown",
}

func StatusString(status HealthStatus) string {
	return statusStrings[status]
}
