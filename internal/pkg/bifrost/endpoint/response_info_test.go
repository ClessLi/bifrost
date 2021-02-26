package endpoint

import (
	"bytes"
	"errors"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"reflect"
	"testing"
	"time"
)

func Test_bytesResponseInfo_Error(t *testing.T) {
	type fields struct {
		Result *bytes.Buffer
		Err    error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "test normal response", fields: fields{
			Result: bytes.NewBuffer([]byte("")),
			Err:    nil,
		}, want: ""},
		{name: "test error response", fields: fields{
			Result: bytes.NewBuffer([]byte("")),
			Err:    errors.New("invalid response"),
		}, want: "invalid response"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := bytesResponseInfo{
				Result: tt.fields.Result,
				Err:    tt.fields.Err,
			}
			if got := br.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bytesResponseInfo_Respond(t *testing.T) {
	type fields struct {
		Result *bytes.Buffer
		Err    error
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{name: "test normal response", fields: fields{
			Result: bytes.NewBuffer([]byte("successful response")),
			Err:    nil,
		}, want: []byte("successful response")},
		{name: "test error response", fields: fields{
			Result: bytes.NewBuffer([]byte("")),
			Err:    errors.New("invalid response"),
		}, want: []byte("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := bytesResponseInfo{
				Result: tt.fields.Result,
				Err:    tt.fields.Err,
			}
			if got := br.Respond(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Respond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorResponseInfo_Error(t *testing.T) {
	type fields struct {
		Err error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "test null error", fields: fields{Err: nil}, want: ""},
		{name: "test error", fields: fields{Err: errors.New("error response")}, want: "error response"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			er := errorResponseInfo{
				Err: tt.fields.Err,
			}
			if got := er.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newUpdateResponseInfo(t *testing.T) {
	type args struct {
		svcResponseInfo service.UpdateResponseInfo
	}
	tests := []struct {
		name string
		args args
		want ErrorResponseInfo
	}{
		{name: "test update success response", args: args{svcResponseInfo: service.NewUpdateResponseInfo("test", nil)}, want: newUpdateResponseInfo(service.NewUpdateResponseInfo("test", nil))},
		{name: "test update fail response", args: args{svcResponseInfo: service.NewUpdateResponseInfo("test", errors.New("update failed"))}, want: newUpdateResponseInfo(service.NewUpdateResponseInfo("test", errors.New("update failed")))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUpdateResponseInfo(tt.args.svcResponseInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUpdateResponseInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newViewResponseInfo(t *testing.T) {
	type args struct {
		svcResponseInfo service.ViewResponseInfo
	}
	tests := []struct {
		name string
		args args
		want BytesResponseInfo
	}{
		{name: "test view success response", args: args{svcResponseInfo: service.NewViewResponseInfo("test", []byte("view test"), nil)}, want: newViewResponseInfo(service.NewViewResponseInfo("test", []byte("view test"), nil))},
		{name: "test view fail response", args: args{svcResponseInfo: service.NewViewResponseInfo("test", nil, errors.New("view failed"))}, want: newViewResponseInfo(service.NewViewResponseInfo("test", nil, errors.New("view failed")))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newViewResponseInfo(tt.args.svcResponseInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newViewResponseInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newWatchResponseInfo(t *testing.T) {
	type args struct {
		svcResponseInfo       service.WatchResponseInfo
		signalChan            chan int
		bytesResponseInfoChan chan BytesResponseInfo
	}
	type result struct {
		bytesResponseInfoChan <-chan BytesResponseInfo
		responseErr           error
	}
	newResult := func(info WatchResponseInfo) result {
		return result{
			bytesResponseInfoChan: info.Respond(),
			responseErr:           info.Close(),
		}
	}
	normalCloseFunc := func() error { return nil }
	abnormalCloseFunc := func() error { return service.ErrInvalidResponseInfo }
	dataChan := make(chan []byte)
	errChan := make(chan error)
	bytesRespInfoChan := make(chan BytesResponseInfo)
	signalChan := make(chan int)
	normalSvcResponseInfo := service.NewWatchResponseInfo("test", normalCloseFunc, dataChan, errChan, nil)
	abnormalSvcResponseInfo := service.NewWatchResponseInfo("test", abnormalCloseFunc, nil, nil, errors.New("watch failed"))
	tests := []struct {
		name string
		args args
		want result
	}{
		{
			name: "test watch success response",
			args: args{
				svcResponseInfo:       normalSvcResponseInfo,
				signalChan:            signalChan,
				bytesResponseInfoChan: bytesRespInfoChan,
			},
			want: newResult(newWatchResponseInfo(normalSvcResponseInfo, signalChan, bytesRespInfoChan)),
		},
		{
			name: "test watch fail response",
			args: args{
				svcResponseInfo:       abnormalSvcResponseInfo,
				signalChan:            signalChan,
				bytesResponseInfoChan: bytesRespInfoChan,
			},
			want: newResult(newWatchResponseInfo(abnormalSvcResponseInfo, signalChan, bytesRespInfoChan)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newResult(newWatchResponseInfo(tt.args.svcResponseInfo, tt.args.signalChan, tt.args.bytesResponseInfoChan)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newWatchResponseInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_watchResponseInfo_Close(t *testing.T) {
	type fields struct {
		Result           chan BytesResponseInfo
		signalChan       chan int
		serviceCloseFunc func() error
	}
	normalResult := make(chan BytesResponseInfo)
	abnormalResult := make(chan BytesResponseInfo)
	normalSignalChan := make(chan int)
	go func() { <-normalSignalChan }()
	abnormalSignalChan := make(chan int)
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "test watch response info close success", fields: fields{
			Result:           normalResult,
			signalChan:       normalSignalChan,
			serviceCloseFunc: func() error { return nil },
		}, wantErr: false},
		{name: "test watch response info close timeout", fields: fields{
			Result:           abnormalResult,
			signalChan:       abnormalSignalChan,
			serviceCloseFunc: func() error { return nil },
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := watchResponseInfo{
				Result:           tt.fields.Result,
				signalChan:       tt.fields.signalChan,
				serviceCloseFunc: tt.fields.serviceCloseFunc,
			}
			if err := wr.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_watchResponseInfo_Respond(t *testing.T) {
	type fields struct {
		Result           chan BytesResponseInfo
		signalChan       chan int
		serviceCloseFunc func() error
	}
	normalResult := make(chan BytesResponseInfo)
	abnormalResult := make(chan BytesResponseInfo)
	tests := []struct {
		name   string
		fields fields
		want   <-chan BytesResponseInfo
	}{
		{name: "test watch respond success", fields: fields{
			Result:     normalResult,
			signalChan: make(chan int),
			serviceCloseFunc: func() error {
				return nil
			},
		}, want: normalResult},
		{name: "test watch respond failed", fields: fields{
			Result:     abnormalResult,
			signalChan: make(chan int),
			serviceCloseFunc: func() error {
				select {
				case <-time.After(time.Second * 10):
					return ErrWatchResponseInfoCloseTimeout
				}
			},
		}, want: abnormalResult},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := watchResponseInfo{
				Result: tt.fields.Result,
				//closeFunc: tt.fields.closeFunc,
			}
			if got := wr.Respond(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Respond() = %v, want %v", got, tt.want)
			}
		})
	}
}
