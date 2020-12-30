package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//func DecodeViewConfigRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeRequest(ctx, r, "ViewConfig")
//}
//
//func DecodeGetConfigRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeRequest(ctx, r, "GetConfig")
//}
//
//func DecodeUpdateConfigRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeRequest(ctx, r, "UpdateConfig")
//}
//
//func DecodeViewStatisticsRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeRequest(ctx, r, "ViewStatistics")
//}
//
//func DecodeStatusRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeRequest(ctx, r, "Status")
//}
//
//func DecodeWatchLogRequest(ctx context.Context, r interface{}) (request interface{}, err error) {
//	//fmt.Println("decoding watch log")
//	return decodeRequest(ctx, r, "WatchLog")
//}

func DecodeHealthCheckRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return decodeRequest(ctx, r, "HealthCheck")
}

func DecodeBifrostServiceRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return decodeRequest(ctx, r, "BifrostService")
}

func decodeRequest(ctx context.Context, r interface{}, handleType string) (interface{}, error) {
	switch handleType {
	case "BifrostService":
		if req, ok := r.(*bifrostpb.Request); ok {
			return endpoint.Request{
				RequestType: req.RequestType,
				Token:       req.Token,
				ServerName:  req.SvrName,
				Param:       req.Param,
				Data:        req.Data,
			}, nil
		}
		return nil, ErrInvalidBifrostServiceRequest
	case "HealthCheck":
		if _, ok := r.(*grpc_health_v1.HealthCheckRequest); ok {
			return endpoint.HealthRequest{}, nil
		}
		return nil, ErrInvalidHealthCheckRequest
	default:
		//fmt.Printf("decode default request, req class is %T, type is %v\n", r, requestType)
		return nil, ErrUnknownRequest
	}
}
