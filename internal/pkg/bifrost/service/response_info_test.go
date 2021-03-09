package service

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestNewViewResponseInfo(t *testing.T) {
	serverName := "testWebServer"
	data := []byte("testData")
	err := errors.New("test error")

	respInfo := &viewResponseInfo{
		serverName: serverName,
		dataBuffer: bytes.NewBuffer(data),
		err:        nil,
	}
	errRespInfo := &viewResponseInfo{
		serverName: serverName,
		dataBuffer: bytes.NewBuffer(nil),
		err:        err,
	}
	nullRespInfo := &viewResponseInfo{
		serverName: serverName,
		dataBuffer: bytes.NewBuffer(nil),
		err:        nil,
	}

	type args struct {
		serverName string
		data       []byte
		err        error
	}
	tests := []struct {
		name string
		args args
		want ViewResponseInfo
	}{
		{
			name: "normal response",
			args: args{
				serverName: serverName,
				data:       data,
				err:        nil,
			},
			want: respInfo,
		},
		{
			name: "response with error",
			args: args{
				serverName: serverName,
				data:       nil,
				err:        err,
			},
			want: errRespInfo,
		},
		{
			name: "input null response data and null error",
			args: args{
				serverName: serverName,
				data:       nil,
				err:        nil,
			},
			want: nullRespInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewViewResponseInfo(tt.args.serverName, tt.args.data, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewViewResponseInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUpdateResponseInfo(t *testing.T) {
	serverName := "testWebServer"
	err := errors.New("test error")

	respInfo := &updateResponseInfo{
		serverName: serverName,
		err:        nil,
	}
	errRespInfo := &updateResponseInfo{
		serverName: serverName,
		err:        err,
	}

	type args struct {
		serverName string
		err        error
	}
	tests := []struct {
		name string
		args args
		want UpdateResponseInfo
	}{
		{
			name: "normal response",
			args: args{
				serverName: serverName,
				err:        nil,
			},
			want: respInfo,
		},
		{
			name: "response with error",
			args: args{
				serverName: serverName,
				err:        err,
			},
			want: errRespInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUpdateResponseInfo(tt.args.serverName, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUpdateResponseInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWatchResponseInfo(t *testing.T) {
	serverName := "testWebServer"
	err := errors.New("test error")
	closeErr := errors.New("test close error")
	dataChan := make(chan []byte)
	errChan := make(chan error)
	closeFunc := func() error { return closeErr }

	respInfo := &watchResponseInfo{
		serverName:      serverName,
		dataChan:        dataChan,
		transferErrChan: errChan,
		closeFunc:       closeFunc,
		err:             nil,
	}
	errRespInfo := &watchResponseInfo{
		serverName:      serverName,
		dataChan:        nil,
		transferErrChan: nil,
		closeFunc:       closeFunc,
		err:             err,
	}
	nilCloseFuncRespInfo := &watchResponseInfo{
		serverName:      serverName,
		dataChan:        dataChan,
		transferErrChan: errChan,
		closeFunc: func() error {
			return nil
		},
		err: nil,
	}

	type args struct {
		serverName      string
		closeFunc       func() error
		dataChan        <-chan []byte
		transferErrChan <-chan error
		err             error
	}
	type result struct {
		serverName      string
		dataChan        <-chan []byte
		transferErrChan <-chan error
		closeErr        error
		err             error
	}
	tests := []struct {
		name string
		args args
		want WatchResponseInfo
	}{
		{
			name: "normal response",
			args: args{
				serverName:      serverName,
				closeFunc:       closeFunc,
				dataChan:        dataChan,
				transferErrChan: errChan,
				err:             nil,
			},
			want: respInfo,
		},
		{
			name: "response with error",
			args: args{
				serverName:      serverName,
				closeFunc:       closeFunc,
				dataChan:        nil,
				transferErrChan: nil,
				err:             err,
			},
			want: errRespInfo,
		},
		{
			name: "input nil close function",
			args: args{
				serverName:      serverName,
				closeFunc:       nil,
				dataChan:        dataChan,
				transferErrChan: errChan,
				err:             nil,
			},
			want: nilCloseFuncRespInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWatchResponseInfo(tt.args.serverName, tt.args.closeFunc, tt.args.dataChan, tt.args.transferErrChan, tt.args.err)
			gotResult := result{
				serverName:      got.GetServerName(),
				dataChan:        got.BytesChan(),
				transferErrChan: got.TransferErrorChan(),
				closeErr:        got.Close(),
				err:             got.Error(),
			}
			wantResult := result{
				serverName:      tt.want.GetServerName(),
				dataChan:        tt.want.BytesChan(),
				transferErrChan: tt.want.TransferErrorChan(),
				closeErr:        tt.want.Close(),
				err:             tt.want.Error(),
			}
			if !reflect.DeepEqual(gotResult, wantResult) {
				t.Errorf("NewWatchResponseInfo() = %v, want %v", gotResult, wantResult)
			}
		})
	}
}

func Test_watchResponseInfo_nilDataChan_and_nilTransferErrChan(t *testing.T) {
	resp := NewWatchResponseInfo("test", func() error {
		return nil
	}, nil, nil, nil)
	dataChan := resp.BytesChan()
	errChan := resp.TransferErrorChan()
	select {
	case data := <-dataChan:
		t.Log(data)
	case err := <-errChan:
		t.Log(err)
	case <-time.After(time.Second):
		t.Logf("channel recv timeout, dataChan type is %T, value is %v; transferErrChan type is %T, value is %v", dataChan, dataChan, errChan, errChan)
	}
}
