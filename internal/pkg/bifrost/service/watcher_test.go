package service

import (
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

type mockLogWatcher struct {
}

func (ml mockLogWatcher) GetDataChan() <-chan []byte {
	return nil
}

func (ml mockLogWatcher) GetTransferErrorChan() <-chan error {
	return nil
}

func (ml mockLogWatcher) Close() error {
	return nil
}

type mockOffstageWatcher struct {
}

func (m mockOffstageWatcher) WatchLog(_, _ string) (LogWatcher, error) {
	return new(mockLogWatcher), nil
}

func TestNewWatcher(t *testing.T) {
	w := &watcher{offstage: new(mockOffstageWatcher)}

	type args struct {
		offstage offstageWatcher
	}
	tests := []struct {
		name      string
		args      args
		want      Watcher
		wantPanic string
	}{
		{
			name: "normal watcher",
			args: args{offstage: new(mockOffstageWatcher)},
			want: w,
		},
		{
			name:      "input nil offstageWatcher",
			args:      args{offstage: nil},
			wantPanic: "offstage is nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if gotPanic := recover(); gotPanic != nil && !reflect.DeepEqual(gotPanic, tt.wantPanic) {
					t.Errorf("NewWatcher() panic reason = %v, want %v", gotPanic, tt.wantPanic)
				}
			}()
			if got := NewWatcher(tt.args.offstage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_watcher_Watch(t *testing.T) {
	serverName := "testWebServer"
	token := "testToken"
	logName := "access.log"
	watchLogReqTypeStr := "WatchLog"
	unknownReqTypeStr := "?unknown"

	// want LogWatcher instance and error
	logWatcherIns, logWatcherErr := new(mockOffstageWatcher).WatchLog(serverName, logName)
	if logWatcherIns == nil {
		t.Fatal("logWatcherIns is nil")
	}

	// want responseInfo
	respInfo := NewWatchResponseInfo(serverName, logWatcherIns.Close, logWatcherIns.GetDataChan(), logWatcherIns.GetTransferErrorChan(), logWatcherErr)
	unknownReqTypeErrRespInfo := NewWatchResponseInfo(serverName, func() error { return ErrInvalidResponseInfo }, nil, nil, ErrUnknownRequestType)
	nilReqErrRespInfo := NewWatchResponseInfo("", func() error { return ErrInvalidResponseInfo }, nil, nil, ErrNilRequestInfo)

	// test requestInfo
	reqInfo := NewWatchRequestInfo(context.Background(), watchLogReqTypeStr, serverName, token, logName)
	unknownReqInfo := NewWatchRequestInfo(context.Background(), unknownReqTypeStr, serverName, token, logName)

	type fields struct {
		offstage offstageWatcher
	}
	type args struct {
		req WatchRequestInfo
	}
	type result struct {
		serverName      string
		bytesChan       <-chan []byte
		transferErrChan <-chan error
		closeErr        error
		err             error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   result
	}{
		{
			name:   "normal request",
			fields: fields{offstage: new(mockOffstageWatcher)},
			args:   args{req: reqInfo},
			want: result{
				serverName:      respInfo.GetServerName(),
				bytesChan:       respInfo.BytesChan(),
				transferErrChan: respInfo.TransferErrorChan(),
				closeErr:        respInfo.Close(),
				err:             respInfo.Error(),
			},
		},
		{
			name:   "unknown request",
			fields: fields{offstage: new(mockOffstageWatcher)},
			args:   args{req: unknownReqInfo},
			want: result{
				serverName:      unknownReqTypeErrRespInfo.GetServerName(),
				bytesChan:       unknownReqTypeErrRespInfo.BytesChan(),
				transferErrChan: unknownReqTypeErrRespInfo.TransferErrorChan(),
				closeErr:        unknownReqTypeErrRespInfo.Close(),
				err:             unknownReqTypeErrRespInfo.Error(),
			},
		},
		{
			name:   "nil request",
			fields: fields{offstage: new(mockOffstageWatcher)},
			args:   args{req: nil},
			want: result{
				serverName:      nilReqErrRespInfo.GetServerName(),
				bytesChan:       nilReqErrRespInfo.BytesChan(),
				transferErrChan: nilReqErrRespInfo.TransferErrorChan(),
				closeErr:        nilReqErrRespInfo.Close(),
				err:             nilReqErrRespInfo.Error(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := watcher{
				offstage: tt.fields.offstage,
			}
			gotResp := w.Watch(tt.args.req)
			got := result{
				serverName:      gotResp.GetServerName(),
				bytesChan:       gotResp.BytesChan(),
				transferErrChan: gotResp.TransferErrorChan(),
				closeErr:        gotResp.Close(),
				err:             gotResp.Error(),
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("result of Watch() = %v, want %v", got, tt.want)
			}
		})
	}
}
