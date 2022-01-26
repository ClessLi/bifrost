package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"reflect"
	"testing"
)

type mockUnknownRequestObject struct {
}

func TestDecodeHealthCheckRequest(t *testing.T) {
	grpcHealthCheckReq := &grpc_health_v1.HealthCheckRequest{}

	type args struct {
		ctx context.Context
		r   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "normal request",
			args: args{
				ctx: context.Background(),
				r:   grpcHealthCheckReq,
			},
			want:    endpoint.HealthRequestInfo{},
			wantErr: false,
		},
		{
			name: "unknown request object",
			args: args{
				ctx: context.Background(),
				r:   new(mockUnknownRequestObject),
			},
			wantErr: true,
		},
		{
			name: "nil request",
			args: args{
				ctx: context.Background(),
				r:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeHealthCheckRequest(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeHealthCheckRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeHealthCheckRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeUpdateRequest(t *testing.T) {
	serverName := "testWebServer"
	token := "testToken"
	data := []byte("testData")

	// want endpoint UpdateRequestInfo
	epUpdateRequestInfo := &endpoint.UpdateRequestInfo{
		UpdateType: "UpdateConfig",
		ServerName: serverName,
		Token:      token,
		Data:       data,
	}
	epUnknownReqTypeRequestInfo := &endpoint.UpdateRequestInfo{
		UpdateType: "?unknown",
		ServerName: serverName,
		Token:      token,
		Data:       data,
	}

	// test grpc request
	grpcUpdateRequest := &bifrostpb.UpdateRequest{
		UpdateType: "UpdateConfig",
		ServerName: serverName,
		Token:      token,
		Data:       data,
	}
	grpcUnknownReqTypeRequest := &bifrostpb.UpdateRequest{
		UpdateType: "?unknown",
		ServerName: serverName,
		Token:      token,
		Data:       data,
	}

	type args struct {
		ctx context.Context
		r   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "normal request",
			args: args{
				ctx: context.Background(),
				r:   grpcUpdateRequest,
			},
			want:    epUpdateRequestInfo,
			wantErr: false,
		},
		{
			name: "unknown req type request",
			args: args{
				ctx: context.Background(),
				r:   grpcUnknownReqTypeRequest,
			},
			want:    epUnknownReqTypeRequestInfo,
			wantErr: false,
		},
		{
			name: "unknown request object",
			args: args{
				ctx: context.Background(),
				r:   new(mockUnknownRequestObject),
			},
			wantErr: true,
		},
		{
			name: "nil request",
			args: args{
				ctx: context.Background(),
				r:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeUpdateRequest(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeUpdateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeUpdateRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeViewRequest(t *testing.T) {
	serverName := "testWebServer"
	token := "testToken"
	displayConfigReqTypeStr := "DisplayConfig"
	getConfigReqTypeStr := "GetConfig"
	showStatisticsReqTypeStr := "ShowStatistics"
	displayServersStatusReqTypeStr := "DisplayServersStatus"
	unknownReqTypeStr := "?unknown"

	// want endpoint ViewRequestInfo
	epDisplayConfigRequestInfo := &endpoint.ViewRequestInfo{
		ViewType:   displayConfigReqTypeStr,
		ServerName: serverName,
		Token:      token,
	}
	epGetConfigRequestInfo := &endpoint.ViewRequestInfo{
		ViewType:   getConfigReqTypeStr,
		ServerName: serverName,
		Token:      token,
	}
	epShowStatisticsRequestInfo := &endpoint.ViewRequestInfo{
		ViewType:   showStatisticsReqTypeStr,
		ServerName: serverName,
		Token:      token,
	}
	epDisplaySvrsStatusRequestInfo := &endpoint.ViewRequestInfo{
		ViewType: displayServersStatusReqTypeStr,
		Token:    token,
	}
	epUnknownReqTypeRequestInfo := &endpoint.ViewRequestInfo{
		ViewType:   unknownReqTypeStr,
		ServerName: serverName,
		Token:      token,
	}

	// test grpc request
	grpcDisplayConfigRequest := &bifrostpb.ViewRequest{
		ViewType:   displayConfigReqTypeStr,
		ServerName: serverName,
		Token:      token,
	}
	grpcGetConfigRequest := &bifrostpb.ViewRequest{
		ViewType:   getConfigReqTypeStr,
		ServerName: serverName,
		Token:      token,
	}
	grpcShowStatisticsRequest := &bifrostpb.ViewRequest{
		ViewType:   showStatisticsReqTypeStr,
		ServerName: serverName,
		Token:      token,
	}
	grpcDisplaySvrsStatusRequest := &bifrostpb.ViewRequest{
		ViewType: displayServersStatusReqTypeStr,
		Token:    token,
	}
	grpcUnknownReqTypeRequest := &bifrostpb.ViewRequest{
		ViewType:   unknownReqTypeStr,
		ServerName: serverName,
		Token:      token,
	}

	type args struct {
		ctx context.Context
		r   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "display config request",
			args: args{
				ctx: context.Background(),
				r:   grpcDisplayConfigRequest,
			},
			want: epDisplayConfigRequestInfo,
		},
		{
			name: "get config request",
			args: args{
				ctx: context.Background(),
				r:   grpcGetConfigRequest,
			},
			want: epGetConfigRequestInfo,
		},
		{
			name: "show statistics request",
			args: args{
				ctx: context.Background(),
				r:   grpcShowStatisticsRequest,
			},
			want: epShowStatisticsRequestInfo,
		},
		{
			name: "display servers status request",
			args: args{
				ctx: context.Background(),
				r:   grpcDisplaySvrsStatusRequest,
			},
			want: epDisplaySvrsStatusRequestInfo,
		},
		{
			name: "unknown req type request",
			args: args{
				ctx: context.Background(),
				r:   grpcUnknownReqTypeRequest,
			},
			want: epUnknownReqTypeRequestInfo,
		},
		{
			name: "unknown request object",
			args: args{
				ctx: context.Background(),
				r:   new(mockUnknownRequestObject),
			},
			wantErr: true,
		},
		{
			name: "nil request",
			args: args{
				ctx: context.Background(),
				r:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeViewRequest(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeViewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeViewRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeWatchRequest(t *testing.T) {
	serverName := "testWebServer"
	token := "testToken"
	logName := "access.log"
	watchLogReqTypeStr := "WatchLog"
	unknownReqTypeStr := "?unknown"

	// want endpoint WatchRequestInfo
	epWatchLogRequestInfo := &endpoint.WatchRequestInfo{
		WatchType:   watchLogReqTypeStr,
		ServerName:  serverName,
		Token:       token,
		WatchObject: logName,
	}
	epUnknownReqTypeRequestInfo := &endpoint.WatchRequestInfo{
		WatchType:   unknownReqTypeStr,
		ServerName:  serverName,
		Token:       token,
		WatchObject: logName,
	}

	// test grpc request
	grpcWatchLogRequest := &bifrostpb.WatchRequest{
		WatchType:   watchLogReqTypeStr,
		ServerName:  serverName,
		Token:       token,
		WatchObject: logName,
	}
	grpcUnknownReqTypeRequst := &bifrostpb.WatchRequest{
		WatchType:   unknownReqTypeStr,
		ServerName:  serverName,
		Token:       token,
		WatchObject: logName,
	}

	type args struct {
		ctx context.Context
		r   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "watch log request",
			args: args{
				ctx: context.Background(),
				r:   grpcWatchLogRequest,
			},
			want: epWatchLogRequestInfo,
		},
		{
			name: "unknown req type request",
			args: args{
				ctx: context.Background(),
				r:   grpcUnknownReqTypeRequst,
			},
			want: epUnknownReqTypeRequestInfo,
		},
		{
			name: "unknown request object",
			args: args{
				ctx: context.Background(),
				r:   new(mockUnknownRequestObject),
			},
			wantErr: true,
		},
		{
			name: "nil request",
			args: args{
				ctx: context.Background(),
				r:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeWatchRequest(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeWatchRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeWatchRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
