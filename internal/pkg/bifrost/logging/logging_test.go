package logging

import (
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"os"
	"reflect"
	"testing"
)

type testViewer struct {
}

func (t testViewer) View(info service.ViewRequestInfo) service.ViewResponseInfo {
	return service.NewViewResponseInfo(info.GetServerName(), []byte("view success"), nil)
}

type testUpdater struct {
}

func (t testUpdater) Update(info service.UpdateRequestInfo) service.UpdateResponseInfo {
	return service.NewUpdateResponseInfo(info.GetServerName(), nil)
}

type testWatcher struct {
	dataChan chan []byte
	errChan  chan error
}

func (t testWatcher) Watch(info service.WatchRequestInfo) service.WatchResponseInfo {
	return service.NewWatchResponseInfo(info.GetServerName(), func() error {
		return nil
	}, t.dataChan, t.errChan, nil)
}

func newTestWatcher(dataChan chan []byte, errChan chan error) service.Watcher {
	return &testWatcher{
		dataChan: dataChan,
		errChan:  errChan,
	}
}

func TestLoggingMiddleware(t *testing.T) {
	type args struct {
		logger log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	tests := []struct {
		name string
		args args
		want service.Service
	}{
		{
			name: "test LoggingMiddleware function",
			args: args{logger: logger},
			want: LoggingMiddleware(logger)(svc),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoggingMiddleware(tt.args.logger)(svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoggingMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClientIP(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	wantIp, wantErr := getClientIP(context.Background())
	tests := []struct {
		name    string
		args    args
		wantIp  string
		wantErr bool
	}{
		{
			name:    "test getClientIP function",
			args:    args{ctx: context.Background()},
			wantIp:  wantIp,
			wantErr: wantErr != nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIp, err := getClientIP(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getClientIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIp != tt.wantIp {
				t.Errorf("getClientIP() gotIp = %v, want %v", gotIp, tt.wantIp)
			}
		})
	}
}

func Test_loggingMiddleware_HealthCheck(t *testing.T) {
	type fields struct {
		viewer  service.Viewer
		updater service.Updater
		watcher service.Watcher
		logger  log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	loggingSvc := LoggingMiddleware(logger)(svc)
	tests := []struct {
		name       string
		fields     fields
		wantResult bool
	}{
		{
			name: "test HealthCheck function",
			fields: fields{
				viewer:  loggingSvc.Viewer(),
				updater: loggingSvc.Updater(),
				watcher: loggingSvc.Watcher(),
				logger:  logger,
			},
			wantResult: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lmw := loggingMiddleware{
				viewer:  tt.fields.viewer,
				updater: tt.fields.updater,
				watcher: tt.fields.watcher,
				logger:  tt.fields.logger,
			}
			if gotResult := lmw.HealthCheck(); gotResult != tt.wantResult {
				t.Errorf("HealthCheck() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_loggingMiddleware_Updater(t *testing.T) {
	type fields struct {
		viewer  service.Viewer
		updater service.Updater
		watcher service.Watcher
		logger  log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	loggingSvc := LoggingMiddleware(logger)(svc)
	tests := []struct {
		name   string
		fields fields
		want   service.Updater
	}{
		{
			name: "test loggingMiddleware Updater method",
			fields: fields{
				viewer:  loggingSvc.Viewer(),
				updater: loggingSvc.Updater(),
				watcher: loggingSvc.Watcher(),
				logger:  logger,
			},
			want: LoggingMiddleware(logger)(svc).Updater(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lmw := loggingMiddleware{
				viewer:  tt.fields.viewer,
				updater: tt.fields.updater,
				watcher: tt.fields.watcher,
				logger:  tt.fields.logger,
			}
			if got := lmw.Updater(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Updater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggingMiddleware_Viewer(t *testing.T) {
	type fields struct {
		viewer  service.Viewer
		updater service.Updater
		watcher service.Watcher
		logger  log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	loggingSvc := LoggingMiddleware(logger)(svc)
	tests := []struct {
		name   string
		fields fields
		want   service.Viewer
	}{
		{
			name: "test loggingMiddleware Viewer method",
			fields: fields{
				viewer:  loggingSvc.Viewer(),
				updater: loggingSvc.Updater(),
				watcher: loggingSvc.Watcher(),
				logger:  logger,
			},
			want: LoggingMiddleware(logger)(svc).Viewer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lmw := loggingMiddleware{
				viewer:  tt.fields.viewer,
				updater: tt.fields.updater,
				watcher: tt.fields.watcher,
				logger:  tt.fields.logger,
			}
			if got := lmw.Viewer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Viewer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggingMiddleware_Watcher(t *testing.T) {
	type fields struct {
		viewer  service.Viewer
		updater service.Updater
		watcher service.Watcher
		logger  log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	loggingSvc := LoggingMiddleware(logger)(svc)
	tests := []struct {
		name   string
		fields fields
		want   service.Watcher
	}{
		{
			name: "test loggingMiddleware Viewer method",
			fields: fields{
				viewer:  loggingSvc.Viewer(),
				updater: loggingSvc.Updater(),
				watcher: loggingSvc.Watcher(),
				logger:  logger,
			},
			want: LoggingMiddleware(logger)(svc).Watcher(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lmw := loggingMiddleware{
				viewer:  tt.fields.viewer,
				updater: tt.fields.updater,
				watcher: tt.fields.watcher,
				logger:  tt.fields.logger,
			}
			if got := lmw.Watcher(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Watcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggingUpdaterMiddleware(t *testing.T) {
	type args struct {
		logger log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	tests := []struct {
		name string
		args args
		want service.Updater
	}{
		{
			name: "test loggingUpdaterMiddleware function",
			args: args{logger: logger},
			want: loggingUpdaterMiddleware(logger)(svc.Updater()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loggingUpdaterMiddleware(tt.args.logger)(svc.Updater()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loggingUpdaterMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggingUpdater_Update(t *testing.T) {
	type fields struct {
		updater service.Updater
		logger  log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	type args struct {
		requestInfo service.UpdateRequestInfo
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantResponseInfo service.UpdateResponseInfo
	}{
		{
			name: "test loggingUpdater.Update method",
			fields: fields{
				updater: LoggingMiddleware(logger)(svc).Updater(),
				logger:  logger,
			},
			args:             args{requestInfo: service.NewUpdateRequestInfo(context.Background(), "UpdateConfig", "test", "UNabcde", []byte("test update config"))},
			wantResponseInfo: LoggingMiddleware(logger)(svc).Updater().Update(service.NewUpdateRequestInfo(context.Background(), "UpdateConfig", "test", "UNabcde", []byte("test update config"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := loggingUpdater{
				updater: tt.fields.updater,
				logger:  tt.fields.logger,
			}
			if gotResponseInfo := u.Update(tt.args.requestInfo); !reflect.DeepEqual(gotResponseInfo, tt.wantResponseInfo) {
				t.Errorf("Update() = %v, want %v", gotResponseInfo, tt.wantResponseInfo)
			}
		})
	}
}

func Test_loggingViewerMiddleware(t *testing.T) {
	type args struct {
		logger log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	tests := []struct {
		name string
		args args
		want service.Viewer
	}{
		{
			name: "test loggingViewerMiddleware function",
			args: args{logger: logger},
			want: loggingViewerMiddleware(logger)(svc.Viewer()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loggingViewerMiddleware(tt.args.logger)(svc.Viewer()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loggingViewerMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggingViewer_View(t *testing.T) {
	type fields struct {
		viewer service.Viewer
		logger log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	type args struct {
		requestInfo service.ViewRequestInfo
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantResponseInfo service.ViewResponseInfo
	}{
		{
			name: "test loggingViewer.Viewer method",
			fields: fields{
				viewer: LoggingMiddleware(logger)(svc).Viewer(),
				logger: logger,
			},
			args:             args{requestInfo: service.NewViewRequestInfo(context.Background(), "DisplayConfig", "test", "UNabcde")},
			wantResponseInfo: LoggingMiddleware(logger)(svc).Viewer().View(service.NewViewRequestInfo(context.Background(), "DisplayConfig", "test", "UNabcde")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := loggingViewer{
				viewer: tt.fields.viewer,
				logger: tt.fields.logger,
			}
			if gotResponseInfo := v.View(tt.args.requestInfo); !reflect.DeepEqual(gotResponseInfo, tt.wantResponseInfo) {
				t.Errorf("View() = %v, want %v", gotResponseInfo, tt.wantResponseInfo)
			}
		})
	}
}

func Test_loggingWatcherMiddleware(t *testing.T) {
	type args struct {
		logger log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	tests := []struct {
		name string
		args args
		want service.Watcher
	}{
		{
			name: "test loggingViewerMiddleware function",
			args: args{logger: logger},
			want: loggingWatcherMiddleware(logger)(svc.Watcher()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loggingWatcherMiddleware(tt.args.logger)(svc.Watcher()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loggingWatcherMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggingWatcher_Watch(t *testing.T) {
	type fields struct {
		watcher service.Watcher
		logger  log.Logger
	}
	svc := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	logger := config.KitLogger(os.Stdout)
	type args struct {
		requestInfo service.WatchRequestInfo
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantResponseInfo service.WatchResponseInfo
	}{
		{
			name: "test loggingWatcher.Watcher method",
			fields: fields{
				watcher: LoggingMiddleware(logger)(svc).Watcher(),
				logger:  logger,
			},
			args:             args{requestInfo: service.NewWatchRequestInfo(context.Background(), "WatchLog", "test", "UNabcde", "access.log")},
			wantResponseInfo: LoggingMiddleware(logger)(svc).Watcher().Watch(service.NewWatchRequestInfo(context.Background(), "WatchLog", "test", "UNabcde", "access.log")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := loggingWatcher{
				watcher: tt.fields.watcher,
				logger:  tt.fields.logger,
			}
			if gotResponseInfo := w.Watch(tt.args.requestInfo); !reflect.DeepEqual(gotResponseInfo, tt.wantResponseInfo) {
				t.Errorf("Watch() = %v, want %v", gotResponseInfo, tt.wantResponseInfo)
			}
		})
	}
}
