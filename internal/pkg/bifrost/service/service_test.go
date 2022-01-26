package service

import (
	"reflect"
	"testing"
)

type mockUpdater struct {
}

func (mu mockUpdater) Update(_ UpdateRequestInfo) UpdateResponseInfo {
	return NewUpdateResponseInfo("testWebServer", nil)
}

type mockViewer struct {
}

func (mv mockViewer) View(_ ViewRequestInfo) ViewResponseInfo {
	return NewViewResponseInfo("testWebServer", []byte("testData"), nil)
}

type mockWatcher struct {
}

func (mv mockWatcher) Watch(_ WatchRequestInfo) WatchResponseInfo {
	return NewWatchResponseInfo("testWebServer", nil, nil, nil, nil)
}

func TestNewService(t *testing.T) {
	svc := &service{
		viewer:  new(mockViewer),
		updater: new(mockUpdater),
		watcher: new(mockWatcher),
	}

	type args struct {
		viewer  Viewer
		updater Updater
		watcher Watcher
	}
	tests := []struct {
		name      string
		args      args
		want      Service
		wantPanic string
	}{
		{
			name: "new Service",
			args: args{
				viewer:  new(mockViewer),
				updater: new(mockUpdater),
				watcher: new(mockWatcher),
			},
			want: svc,
		},
		{
			name: "new Service with nil Viewer",
			args: args{
				viewer:  nil,
				updater: new(mockUpdater),
				watcher: new(mockWatcher),
			},
			wantPanic: "viewer is nil",
		},
		{
			name: "new Service with nil Updater",
			args: args{
				viewer:  new(mockViewer),
				updater: nil,
				watcher: new(mockWatcher),
			},
			wantPanic: "updater is nil",
		},
		{
			name: "new Service with nil Watcher",
			args: args{
				viewer:  new(mockViewer),
				updater: new(mockUpdater),
				watcher: nil,
			},
			wantPanic: "watcher is nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if gotPanic := recover(); gotPanic != nil && !reflect.DeepEqual(gotPanic, tt.wantPanic) {
					t.Errorf("NewService() panic reason = %v, want %v", gotPanic, tt.wantPanic)
				}
			}()
			if got := NewService(tt.args.viewer, tt.args.updater, tt.args.watcher); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Updater(t *testing.T) {
	type fields struct {
		viewer  Viewer
		updater Updater
		watcher Watcher
	}
	tests := []struct {
		name      string
		fields    fields
		want      Updater
		wantPanic string
	}{
		{
			name: "service.Updater()",
			fields: fields{
				viewer:  new(mockViewer),
				updater: new(mockUpdater),
				watcher: new(mockWatcher),
			},
			want: new(mockUpdater),
		},
		{
			name: "service.Updater() with nil updater member",
			fields: fields{
				viewer:  nil,
				updater: nil,
				watcher: nil,
			},
			wantPanic: "updater is nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				viewer:  tt.fields.viewer,
				updater: tt.fields.updater,
				watcher: tt.fields.watcher,
			}
			defer func() {
				if gotPanic := recover(); gotPanic != nil && !reflect.DeepEqual(gotPanic, tt.wantPanic) {
					t.Errorf("Updater() panic reason = %v, want = %v", gotPanic, tt.wantPanic)
				}
			}()
			if got := s.Updater(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Updater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Viewer(t *testing.T) {
	type fields struct {
		viewer  Viewer
		updater Updater
		watcher Watcher
	}
	tests := []struct {
		name      string
		fields    fields
		want      Viewer
		wantPanic string
	}{
		{
			name: "service.Viewer()",
			fields: fields{
				viewer:  new(mockViewer),
				updater: new(mockUpdater),
				watcher: new(mockWatcher),
			},
			want: new(mockViewer),
		},
		{
			name: "service.Viewer() with nil viewer member",
			fields: fields{
				viewer:  nil,
				updater: nil,
				watcher: nil,
			},
			wantPanic: "viewer is nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				viewer:  tt.fields.viewer,
				updater: tt.fields.updater,
				watcher: tt.fields.watcher,
			}
			defer func() {
				if gotPanic := recover(); gotPanic != nil && !reflect.DeepEqual(gotPanic, tt.wantPanic) {
					t.Errorf("Viewer() panic reason = %v, want = %v", gotPanic, tt.wantPanic)
				}
			}()
			if got := s.Viewer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Viewer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Watcher(t *testing.T) {
	type fields struct {
		viewer  Viewer
		updater Updater
		watcher Watcher
	}
	tests := []struct {
		name      string
		fields    fields
		want      Watcher
		wantPanic string
	}{
		{
			name: "service.Watcher()",
			fields: fields{
				viewer:  new(mockViewer),
				updater: new(mockUpdater),
				watcher: new(mockWatcher),
			},
			want: new(mockWatcher),
		},
		{
			name: "service.Watcher() with nil watcher member",
			fields: fields{
				viewer:  nil,
				updater: nil,
				watcher: nil,
			},
			wantPanic: "watcher is nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				viewer:  tt.fields.viewer,
				updater: tt.fields.updater,
				watcher: tt.fields.watcher,
			}
			defer func() {
				if gotPanic := recover(); gotPanic != nil && !reflect.DeepEqual(gotPanic, tt.wantPanic) {
					t.Errorf("Watcher() panic reason = %v, want = %v", gotPanic, tt.wantPanic)
				}
			}()
			if got := s.Watcher(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Watcher() = %v, want %v", got, tt.want)
			}
		})
	}
}
