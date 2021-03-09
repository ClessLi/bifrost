package service

import (
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

type mockOffstageUpdater struct {
}

func (m mockOffstageUpdater) UpdateConfig(_ string, _ []byte) error {
	return nil
}

func TestNewUpdater(t *testing.T) {
	u := &updater{offstage: new(mockOffstageUpdater)}

	type args struct {
		offstage offstageUpdater
	}
	tests := []struct {
		name      string
		args      args
		want      Updater
		wantPanic string
	}{
		{
			name: "normal updater",
			args: args{offstage: new(mockOffstageUpdater)},
			want: u,
		},
		{
			name:      "input nil offstageUpdater",
			args:      args{offstage: nil},
			wantPanic: "offstage is nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if gotPanic := recover(); gotPanic != nil && !reflect.DeepEqual(gotPanic, tt.wantPanic) {
					t.Errorf("NewUpdater() panic reason = %v, want %v", gotPanic, tt.wantPanic)
				}
			}()
			if got := NewUpdater(tt.args.offstage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUpdater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updater_Update(t *testing.T) {
	serverName := "testWebServer"
	token := "testToken"
	data := []byte("testData")

	// want error
	updateConfigErr := new(mockOffstageUpdater).UpdateConfig(serverName, data)

	// want responseInfo
	respInfo := NewUpdateResponseInfo(serverName, updateConfigErr)
	unknownReqTypeErrRespInfo := NewUpdateResponseInfo(serverName, ErrUnknownRequestType)
	nilReqInfoErrRespInfo := NewUpdateResponseInfo("", ErrNilRequestInfo)

	// test requestInfo
	reqInfo := NewUpdateRequestInfo(context.Background(), "UpdateConfig", serverName, token, data)
	unknownReqInfo := NewUpdateRequestInfo(context.Background(), "?unkown", serverName, token, data)

	type fields struct {
		offstage offstageUpdater
	}
	type args struct {
		req UpdateRequestInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   UpdateResponseInfo
	}{
		{
			name:   "normal request",
			fields: fields{offstage: new(mockOffstageUpdater)},
			args:   args{req: reqInfo},
			want:   respInfo,
		},
		{
			name:   "unknown request",
			fields: fields{offstage: new(mockOffstageUpdater)},
			args:   args{req: unknownReqInfo},
			want:   unknownReqTypeErrRespInfo,
		},
		{
			name:   "nil request",
			fields: fields{offstage: new(mockOffstageUpdater)},
			args:   args{req: nil},
			want:   nilReqInfoErrRespInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &updater{
				offstage: tt.fields.offstage,
			}
			if got := u.Update(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
