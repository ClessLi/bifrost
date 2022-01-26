package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"reflect"
	"testing"
	"time"
)

type mockUnknownResponseObject struct {
}

type mockEPErrResponseInfo struct {
	err string
}

func (mee mockEPErrResponseInfo) Error() string {
	return mee.err
}

type mockEPBytesResponseInfo struct {
	ret []byte
	err string
}

func (meb mockEPBytesResponseInfo) Respond() []byte {
	return meb.ret
}

func (meb mockEPBytesResponseInfo) Error() string {
	return meb.err
}

type mockEPWatchResponseInfo struct {
	ret chan endpoint.BytesResponseInfo
}

func (mew mockEPWatchResponseInfo) Respond() <-chan endpoint.BytesResponseInfo {
	return mew.ret
}

func (mew mockEPWatchResponseInfo) Close() error {
	return nil
}

func TestEncodeHealthCheckResponse(t *testing.T) {
	// want grpc response
	grpcHealthyResponse := &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}
	grpcUnhealthyResponse := &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	}

	// test endpoint responseInfo
	epHealthyResponseInfo := endpoint.HealthResponseInfo{Status: true}
	epUnhealthyResponseInfo := endpoint.HealthResponseInfo{Status: false}

	type args struct {
		in0 context.Context
		r   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "healthy response",
			args: args{
				in0: context.Background(),
				r:   epHealthyResponseInfo,
			},
			want: grpcHealthyResponse,
		},
		{
			name: "unhealthy response",
			args: args{
				in0: context.Background(),
				r:   epUnhealthyResponseInfo,
			},
			want: grpcUnhealthyResponse,
		},
		{
			name: "unknown response object",
			args: args{
				in0: context.Background(),
				r:   new(mockUnknownResponseObject),
			},
			wantErr: true,
		},
		{
			name: "nil response",
			args: args{
				in0: context.Background(),
				r:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeHealthCheckResponse(tt.args.in0, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeHealthCheckResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeHealthCheckResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeUpdateResponse(t *testing.T) {
	errMsg := "test error"

	// want grpc response
	grpcErrResponse := &bifrostpb.ErrorResponse{Err: errMsg}

	// test endpoint responseInfo
	epErrResponse := &mockEPErrResponseInfo{err: errMsg}

	type args struct {
		in0 context.Context
		r   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "normal response",
			args: args{
				in0: context.Background(),
				r:   epErrResponse,
			},
			want: grpcErrResponse,
		},
		{
			name: "unknown response object",
			args: args{
				in0: context.Background(),
				r:   new(mockUnknownResponseObject),
			},
			wantErr: true,
		},
		{
			name: "nil response",
			args: args{
				in0: context.Background(),
				r:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeUpdateResponse(tt.args.in0, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeUpdateResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeUpdateResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeViewResponse(t *testing.T) {
	errMsg := "test error"
	resultData := []byte("test data")

	// want grpc response
	grpcBytesResponse := &bifrostpb.BytesResponse{Ret: resultData}
	grpcErrBytesResponse := &bifrostpb.BytesResponse{Err: errMsg}

	// test endpoint responseInfo
	epBytesResponseInfo := &mockEPBytesResponseInfo{ret: resultData}
	epErrBytesResponseInfo := &mockEPBytesResponseInfo{err: errMsg}

	type args struct {
		in0 context.Context
		r   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "normal response",
			args: args{
				in0: context.Background(),
				r:   epBytesResponseInfo,
			},
			want: grpcBytesResponse,
		},
		{
			name: "error response",
			args: args{
				in0: context.Background(),
				r:   epErrBytesResponseInfo,
			},
			want: grpcErrBytesResponse,
		},
		{
			name: "unknown response object",
			args: args{
				in0: context.Background(),
				r:   new(mockUnknownResponseObject),
			},
			wantErr: true,
		},
		{
			name: "nil response",
			args: args{
				in0: context.Background(),
				r:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeViewResponse(tt.args.in0, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeViewResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeViewResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeWatchResponse(t *testing.T) {
	bytesRespInfo := new(mockEPBytesResponseInfo)
	transferErr := "test error"
	respInfoChanForWant := make(chan endpoint.BytesResponseInfo)
	respInfoChanForTest := make(chan endpoint.BytesResponseInfo)

	// want watchResponseInfo
	respInfo := encodeWatchResponse(&mockEPWatchResponseInfo{
		ret: respInfoChanForWant,
	})

	// test endpoint responseInfo
	epWatchResponseInfo := &mockEPWatchResponseInfo{
		ret: respInfoChanForTest,
	}

	type args struct {
		in0 context.Context
		r   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "normal response",
			args: args{
				in0: context.Background(),
				r:   epWatchResponseInfo,
			},
			want: respInfo,
		},
		{
			name: "unknown response object",
			args: args{
				in0: context.Background(),
				r:   new(mockUnknownResponseObject),
			},
			wantErr: true,
		},
		{
			name: "nil response",
			args: args{
				in0: context.Background(),
				r:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeWatchResponse(tt.args.in0, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeWatchResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			gotWatchRespInfo, isRightGot := got.(*watchResponseInfo)
			if !isRightGot {
				t.Errorf("got = %v, is not *watchResponseInfo", got)
				return
			}
			wantWatchRespInfo, isRightWant := tt.want.(*watchResponseInfo)
			if !isRightWant {
				t.Errorf("want = %v, is not *watchResponseInfo", tt.want)
			}

			defer func() {
				gotCloseErr := gotWatchRespInfo.Close()
				wantCloseErr := wantWatchRespInfo.Close()
				if !reflect.DeepEqual(gotCloseErr, wantCloseErr) {
					t.Errorf("*watchResponseInfo.Close() got = %v, want %v", gotCloseErr, wantCloseErr)
				}
			}()

			go func() {
				respInfoChanForWant <- bytesRespInfo
			}()

			go func() {
				respInfoChanForTest <- bytesRespInfo
			}()

			var wantBytesResp *bifrostpb.BytesResponse
			var gotBytesResp *bifrostpb.BytesResponse

			select {
			case wantBytesResp = <-wantWatchRespInfo.Respond():
				break
			case <-time.After(time.Second):
				wantBytesResp = &bifrostpb.BytesResponse{
					Ret: nil,
					Err: transferErr,
				}
			}

			select {
			case gotBytesResp = <-gotWatchRespInfo.Respond():
				break
			case <-time.After(time.Second):
				gotBytesResp = &bifrostpb.BytesResponse{
					Ret: nil,
					Err: transferErr,
				}
			}

			if !reflect.DeepEqual(gotBytesResp, wantBytesResp) {
				t.Errorf("*watchResponseInfo.Respond() got = %v, want %v", gotBytesResp, wantBytesResp)
			}
		})
	}
}

//func Test_newWatchResponseInfo(t *testing.T) {
//	type args struct {
//		epWatchResponseInfo endpoint.WatchResponseInfo
//	}
//	tests := []struct {
//		name string
//		args args
//		want *watchResponseInfo
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := encodeWatchResponse(tt.args.epWatchResponseInfo); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("encodeWatchResponse() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func Test_watchResponseInfo_Close(t *testing.T) {
//	type fields struct {
//		bytesResponseChan chan *bifrostpb.BytesResponse
//		signalChan        chan int
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			wr := watchResponseInfo{
//				bytesResponseChan: tt.fields.bytesResponseChan,
//				signalChan:        tt.fields.signalChan,
//			}
//			if err := wr.Close(); (err != nil) != tt.wantErr {
//				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func Test_watchResponseInfo_Respond(t *testing.T) {
//	type fields struct {
//		bytesResponseChan chan *bifrostpb.BytesResponse
//		signalChan        chan int
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   <-chan *bifrostpb.BytesResponse
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			wr := watchResponseInfo{
//				bytesResponseChan: tt.fields.bytesResponseChan,
//				signalChan:        tt.fields.signalChan,
//			}
//			if got := wr.Respond(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Respond() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
