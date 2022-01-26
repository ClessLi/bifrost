package service

import (
	"bytes"
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

func TestNewUpdateRequestInfo(t *testing.T) {
	ctx := context.Background()
	updateConfigReqTypeStr := "UpdateConfig"
	unknownReqTypeStr := "?unknown"
	nullReqTypeStr := ""
	serverName := "testWebServer"
	token := "testToken"
	data := []byte("testData")
	dataBuff := bytes.NewBuffer(data)

	updateConfigReqInfo := &updateRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: UpdateConfig,
		token:       token,
		dataBuffer:  dataBuff,
	}
	unknownReqInfo := &updateRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: UnknownReqType,
		token:       token,
		dataBuffer:  dataBuff,
	}
	updateConfigWithNullDataReqInfo := &updateRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: UpdateConfig,
		token:       token,
		dataBuffer:  bytes.NewBuffer(nil),
	}

	type args struct {
		ctx        context.Context
		reqTypeStr string
		serverName string
		token      string
		data       []byte
	}

	tests := []struct {
		name string
		args args
		want UpdateRequestInfo
	}{
		{
			name: "normal update config request",
			args: args{
				ctx:        ctx,
				reqTypeStr: updateConfigReqTypeStr,
				serverName: serverName,
				token:      token,
				data:       data,
			},
			want: updateConfigReqInfo,
		},
		{
			name: "input unknown request type str",
			args: args{
				ctx:        ctx,
				reqTypeStr: unknownReqTypeStr,
				serverName: serverName,
				token:      token,
				data:       data,
			},
			want: unknownReqInfo,
		},
		{
			name: "input null request type str",
			args: args{
				ctx:        ctx,
				reqTypeStr: nullReqTypeStr,
				serverName: serverName,
				token:      token,
				data:       data,
			},
			want: unknownReqInfo,
		},
		{
			name: "input null data",
			args: args{
				ctx:        ctx,
				reqTypeStr: updateConfigReqTypeStr,
				serverName: serverName,
				token:      token,
				data:       nil,
			},
			want: updateConfigWithNullDataReqInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUpdateRequestInfo(tt.args.ctx, tt.args.reqTypeStr, tt.args.serverName, tt.args.token, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUpdateRequestInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewViewRequestInfo(t *testing.T) {
	ctx := context.Background()
	displayConfigReqTypeStr := "DisplayConfig"
	getConfigReqTypeStr := "GetConfig"
	showStatisticsReqTypeStr := "ShowStatistics"
	displayStatusReqTypeStr := "DisplayServersStatus"
	unknownReqTypeStr := "?unknown"
	nullReqTypeStr := ""
	serverName := "testWebServer"
	token := "testToken"

	displayConfigReqInfo := &viewRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: DisplayConfig,
		token:       token,
	}
	getConfigReqInfo := &viewRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: GetConfig,
		token:       token,
	}
	showStatisticsReqInfo := &viewRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: ShowStatistics,
		token:       token,
	}
	displayStatusReqInfo := &viewRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: DisplayServersStatus,
		token:       token,
	}
	unknownReqInfo := &viewRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: UnknownReqType,
		token:       token,
	}

	type args struct {
		ctx        context.Context
		reqTypeStr string
		serverName string
		token      string
	}
	tests := []struct {
		name string
		args args
		want ViewRequestInfo
	}{
		{
			name: "display config request",
			args: args{
				ctx:        ctx,
				reqTypeStr: displayConfigReqTypeStr,
				serverName: serverName,
				token:      token,
			},
			want: displayConfigReqInfo,
		},
		{
			name: "get config request",
			args: args{
				ctx:        ctx,
				reqTypeStr: getConfigReqTypeStr,
				serverName: serverName,
				token:      token,
			},
			want: getConfigReqInfo,
		},
		{
			name: "show statistics request",
			args: args{
				ctx:        ctx,
				reqTypeStr: showStatisticsReqTypeStr,
				serverName: serverName,
				token:      token,
			},
			want: showStatisticsReqInfo,
		},
		{
			name: "display status request",
			args: args{
				ctx:        ctx,
				reqTypeStr: displayStatusReqTypeStr,
				serverName: serverName,
				token:      token,
			},
			want: displayStatusReqInfo,
		},
		{
			name: "unknown request",
			args: args{
				ctx:        ctx,
				reqTypeStr: unknownReqTypeStr,
				serverName: serverName,
				token:      token,
			},
			want: unknownReqInfo,
		},
		{
			name: "input null request type str",
			args: args{
				ctx:        ctx,
				reqTypeStr: nullReqTypeStr,
				serverName: serverName,
				token:      token,
			},
			want: unknownReqInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewViewRequestInfo(tt.args.ctx, tt.args.reqTypeStr, tt.args.serverName, tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewViewRequestInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWatchRequestInfo(t *testing.T) {
	ctx := context.Background()
	watchLogReqTypeStr := "WatchLog"
	unknownReqTypeStr := "?unknown"
	nullReqTypeStr := ""
	serverName := "testWebServer"
	token := "testToken"
	logName := "access.log"
	nullWatchedObjectName := ""

	watchLogReqInfo := &watchRequestInfo{
		context:           ctx,
		serverName:        serverName,
		requestType:       WatchLog,
		token:             token,
		watchedObjectName: logName,
	}
	unknownReqInfo := &watchRequestInfo{
		context:           ctx,
		serverName:        serverName,
		requestType:       UnknownReqType,
		token:             token,
		watchedObjectName: logName,
	}
	nullWatchedObjectReqInfo := &watchRequestInfo{
		context:           ctx,
		serverName:        serverName,
		requestType:       WatchLog,
		token:             token,
		watchedObjectName: nullWatchedObjectName,
	}

	type args struct {
		ctx               context.Context
		reqTypeStr        string
		serverName        string
		token             string
		watchedObjectName string
	}
	tests := []struct {
		name string
		args args
		want WatchRequestInfo
	}{
		{
			name: "watch access.log request",
			args: args{
				ctx:               ctx,
				reqTypeStr:        watchLogReqTypeStr,
				serverName:        serverName,
				token:             token,
				watchedObjectName: logName,
			},
			want: watchLogReqInfo,
		},
		{
			name: "unknown request",
			args: args{
				ctx:               ctx,
				reqTypeStr:        unknownReqTypeStr,
				serverName:        serverName,
				token:             token,
				watchedObjectName: logName,
			},
			want: unknownReqInfo,
		},
		{
			name: "null request type str input",
			args: args{
				ctx:               ctx,
				reqTypeStr:        nullReqTypeStr,
				serverName:        serverName,
				token:             token,
				watchedObjectName: logName,
			},
			want: unknownReqInfo,
		},
		{
			name: "null watched object name input",
			args: args{
				ctx:               ctx,
				reqTypeStr:        watchLogReqTypeStr,
				serverName:        serverName,
				token:             token,
				watchedObjectName: nullWatchedObjectName,
			},
			want: nullWatchedObjectReqInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWatchRequestInfo(tt.args.ctx, tt.args.reqTypeStr, tt.args.serverName, tt.args.token, tt.args.watchedObjectName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWatchRequestInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateRequestInfo_GetData(t *testing.T) {
	ctx := context.Background()
	serverName := "testWebServer"
	token := "testToken"
	data := []byte("testData")
	dataBuff := bytes.NewBuffer(data)
	nullDataBuff := bytes.NewBuffer(nil)
	nullData := nullDataBuff.Bytes()

	type fields struct {
		context     context.Context
		serverName  string
		requestType RequestType
		token       string
		dataBuffer  *bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "normal",
			fields: fields{
				context:     ctx,
				serverName:  serverName,
				requestType: UpdateConfig,
				token:       token,
				dataBuffer:  dataBuff,
			},
			want: data,
		},
		{
			name: "null data",
			fields: fields{
				context:     ctx,
				serverName:  serverName,
				requestType: UpdateConfig,
				token:       token,
				dataBuffer:  nullDataBuff,
			},
			want: nullData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := updateRequestInfo{
				context:     tt.fields.context,
				serverName:  tt.fields.serverName,
				requestType: tt.fields.requestType,
				token:       tt.fields.token,
				dataBuffer:  tt.fields.dataBuffer,
			}
			if got := u.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}
