package service

import (
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

type mockOffstageViewer struct {
}

func (m mockOffstageViewer) DisplayConfig(_ string) ([]byte, error) {
	return []byte("testData"), nil
}

func (m mockOffstageViewer) GetConfig(_ string) ([]byte, error) {
	return []byte("testData"), nil
}

func (m mockOffstageViewer) ShowStatistics(_ string) ([]byte, error) {
	return []byte("testData"), nil
}

func (m mockOffstageViewer) DisplayServersStatus() ([]byte, error) {
	return []byte("testData"), nil
}

func TestNewViewer(t *testing.T) {
	v := &viewer{offstage: new(mockOffstageViewer)}

	type args struct {
		offstage offstageViewer
	}
	tests := []struct {
		name      string
		args      args
		want      Viewer
		wantPanic string
	}{
		{
			name: "normal viewer",
			args: args{offstage: new(mockOffstageViewer)},
			want: v,
		},
		{
			name:      "input nil offstageViewer",
			args:      args{nil},
			wantPanic: "offstage is nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if gotPanic := recover(); gotPanic != nil && !reflect.DeepEqual(gotPanic, tt.wantPanic) {
					t.Errorf("NewViewer() panic reason = %v, want %v", gotPanic, tt.wantPanic)
				}
			}()
			if got := NewViewer(tt.args.offstage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewViewer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_viewer_View(t *testing.T) {
	serverName := "testWebServer"
	token := "testToken"
	displayConfigReqTypeStr := "DisplayConfig"
	getConfigReqTypeStr := "GetConfig"
	showStatisticsReqTypeStr := "ShowStatistics"
	displayStatusReqTypeStr := "DisplayServersStatus"
	unknownReqTypeStr := "?unknown"

	// want data and error
	disConfData, disConfErr := new(mockOffstageViewer).DisplayConfig(serverName)
	getConfData, getConfErr := new(mockOffstageViewer).GetConfig(serverName)
	showStatisticsData, showStatisticsErr := new(mockOffstageViewer).ShowStatistics(serverName)
	disSvrsStatusData, disSvrsStatusErr := new(mockOffstageViewer).DisplayServersStatus()

	// want responseInfo
	disConfRespInfo := NewViewResponseInfo(serverName, disConfData, disConfErr)
	getConfRespInfo := NewViewResponseInfo(serverName, getConfData, getConfErr)
	showStatisticsRespInfo := NewViewResponseInfo(serverName, showStatisticsData, showStatisticsErr)
	disSvrsStatusRespInfo := NewViewResponseInfo("", disSvrsStatusData, disSvrsStatusErr)
	unknownReqTypeErrRespInfo := NewViewResponseInfo(serverName, nil, ErrUnknownRequestType)
	nilReqInfoErrRespInfo := NewViewResponseInfo("", nil, ErrNilRequestInfo)

	// test requestInfo
	displayConfigReqInfo := NewViewRequestInfo(context.Background(), displayConfigReqTypeStr, serverName, token)
	getConfigReqInfo := NewViewRequestInfo(context.Background(), getConfigReqTypeStr, serverName, token)
	showStatisticsReqInfo := NewViewRequestInfo(context.Background(), showStatisticsReqTypeStr, serverName, token)
	displayServersStatusReqInfo := NewViewRequestInfo(context.Background(), displayStatusReqTypeStr, serverName, token)
	unknownReqInfo := NewViewRequestInfo(context.Background(), unknownReqTypeStr, serverName, token)

	type fields struct {
		offstage offstageViewer
	}
	type args struct {
		req ViewRequestInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ViewResponseInfo
	}{
		{
			name:   "display config request",
			fields: fields{offstage: new(mockOffstageViewer)},
			args:   args{req: displayConfigReqInfo},
			want:   disConfRespInfo,
		},
		{
			name:   "get config request",
			fields: fields{offstage: new(mockOffstageViewer)},
			args:   args{req: getConfigReqInfo},
			want:   getConfRespInfo,
		},
		{
			name:   "show statistics request",
			fields: fields{offstage: new(mockOffstageViewer)},
			args:   args{req: showStatisticsReqInfo},
			want:   showStatisticsRespInfo,
		},
		{
			name:   "display servers status request",
			fields: fields{offstage: new(mockOffstageViewer)},
			args:   args{req: displayServersStatusReqInfo},
			want:   disSvrsStatusRespInfo,
		},
		{
			name:   "unknown request",
			fields: fields{offstage: new(mockOffstageViewer)},
			args:   args{req: unknownReqInfo},
			want:   unknownReqTypeErrRespInfo,
		},
		{
			name:   "nil request",
			fields: fields{offstage: new(mockOffstageViewer)},
			args:   args{req: nil},
			want:   nilReqInfoErrRespInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viewer{
				offstage: tt.fields.offstage,
			}
			if got := v.View(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("View() = %v, want %v", got, tt.want)
			}
		})
	}
}
