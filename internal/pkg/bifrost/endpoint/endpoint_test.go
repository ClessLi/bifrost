package endpoint

import (
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
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

func TestMakeHealthCheckEndpoint(t *testing.T) {
	type args struct {
		in0 service.Service
	}
	type result struct {
		response interface{}
		err      error
	}
	newResult := func(ep endpoint.Endpoint, ctx context.Context, request interface{}) result {
		resp, err := ep(ctx, request)
		return result{
			response: resp,
			err:      err,
		}
	}
	testService := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	tests := []struct {
		name string
		args args
		want result
	}{
		{name: "health check endpoint", args: args{in0: testService}, want: newResult(MakeHealthCheckEndpoint(testService), context.Background(), nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newResult(MakeHealthCheckEndpoint(tt.args.in0), context.Background(), nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeHealthCheckEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeUpdaterEndpoint(t *testing.T) {
	type args struct {
		updater service.Updater
	}
	type result struct {
		response interface{}
		err      error
	}
	newResult := func(ep endpoint.Endpoint, ctx context.Context, request interface{}) result {
		resp, err := ep(ctx, request)
		return result{
			response: resp,
			err:      err,
		}
	}
	testService := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	tests := []struct {
		name string
		args args
		want result
	}{
		{name: "updater endpoint", args: args{updater: testService.Updater()}, want: newResult(MakeUpdaterEndpoint(testService.Updater()), context.Background(), nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newResult(MakeUpdaterEndpoint(tt.args.updater), context.Background(), nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeUpdaterEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeViewerEndpoint(t *testing.T) {
	type args struct {
		viewer service.Viewer
	}
	type result struct {
		response interface{}
		err      error
	}
	newResult := func(ep endpoint.Endpoint, ctx context.Context, request interface{}) result {
		resp, err := ep(ctx, request)
		return result{
			response: resp,
			err:      err,
		}
	}
	testService := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	tests := []struct {
		name string
		args args
		want result
	}{
		{name: "viewer endpoint", args: args{viewer: testService.Viewer()}, want: newResult(MakeViewerEndpoint(testService.Viewer()), context.Background(), nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newResult(MakeViewerEndpoint(tt.args.viewer), context.Background(), nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeViewerEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeWatcherEndpoint(t *testing.T) {
	type args struct {
		watcher service.Watcher
	}
	type result struct {
		response interface{}
		err      error
	}
	newResult := func(ep endpoint.Endpoint, ctx context.Context, request interface{}) result {
		resp, err := ep(ctx, request)
		return result{
			response: resp,
			err:      err,
		}
	}
	testService := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	tests := []struct {
		name string
		args args
		want result
	}{
		{name: "watcher endpoint", args: args{watcher: testService.Watcher()}, want: newResult(MakeWatcherEndpoint(testService.Watcher()), context.Background(), nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newResult(MakeWatcherEndpoint(tt.args.watcher), context.Background(), nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeWatcherEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBifrostEndpoints(t *testing.T) {
	type args struct {
		svc service.Service
	}
	type result struct {
		response interface{}
		err      error
	}
	type results struct {
		viewResult        result
		updateResult      result
		watchResult       result
		healthCheckResult result
	}
	newResult := func(ep endpoint.Endpoint, ctx context.Context, request interface{}) result {
		resp, err := ep(ctx, request)
		return result{
			response: resp,
			err:      err,
		}
	}
	newResults := func(endpoints BifrostEndpoints, ctx context.Context, viewRequest interface{}, updateRequest interface{}, watchRequest interface{}, healthCheckRequest interface{}) results {
		return results{
			viewResult:        newResult(endpoints.ViewerEndpoint, ctx, viewRequest),
			updateResult:      newResult(endpoints.UpdaterEndpoint, ctx, updateRequest),
			watchResult:       newResult(endpoints.WatcherEndpoint, ctx, watchRequest),
			healthCheckResult: newResult(endpoints.HealthCheckEndpoint, ctx, healthCheckRequest),
		}
	}
	testService := service.NewService(new(testViewer), new(testUpdater), newTestWatcher(make(chan []byte), make(chan error)))
	tests := []struct {
		name string
		args args
		want results
	}{
		{name: "new bifrost endpoints", args: args{svc: testService}, want: newResults(NewBifrostEndpoints(testService), context.Background(), nil, nil, nil, nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newResults(NewBifrostEndpoints(tt.args.svc), context.Background(), nil, nil, nil, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBifrostEndpoints() = %v, want %v", got, tt.want)
			}
		})
	}
}
