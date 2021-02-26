package service

import (
	ngLog "github.com/ClessLi/bifrost/pkg/log/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"reflect"
	"sync"
	"testing"
)

func TestNewLogWatcher(t *testing.T) {
	type args struct {
		dataChan  chan []byte
		errChan   chan error
		closeFunc func() error
	}
	type result struct {
		dataChan <-chan []byte
		errChan  <-chan error
		closeErr error
	}
	newResult := func(lw LogWatcher) result {
		return result{
			dataChan: lw.GetDataChan(),
			errChan:  lw.GetTransferErrorChan(),
			closeErr: lw.Close(),
		}
	}
	dataChan := make(chan []byte)
	errChan := make(chan error)
	tests := []struct {
		name string
		args args
		want result
	}{
		{
			name: "test NewLogWatcher function",
			args: args{
				dataChan:  dataChan,
				errChan:   errChan,
				closeFunc: func() error { return nil },
			},
			want: newResult(NewLogWatcher(dataChan, errChan, func() error { return nil })),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newResult(NewLogWatcher(tt.args.dataChan, tt.args.errChan, tt.args.closeFunc)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogWatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestNewMetrics(t *testing.T) {
//	type args struct {
//		webServerInfoFunc func() []WebServerInfo
//		errChan           chan error
//	}
//	tests := []struct {
//		name string
//		args args
//		want Metrics
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewMetrics(tt.args.webServerInfoFunc, tt.args.errChan); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewMetrics() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestNewOffstage(t *testing.T) {
//	type args struct {
//		services       map[string]WebServerConfigService
//		configManagers map[string]configuration.ConfigManager
//		metrics        Metrics
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Offstage
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewOffstage(tt.args.services, tt.args.configManagers, tt.args.metrics); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewOffstage() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestNewWebServerConfigService(t *testing.T) {
	type args struct {
		configuration configuration.Configuration
		serverBinPath string
		logsDir       string
		log           *ngLog.Log
	}
	var testConfiguration configuration.Configuration
	nginxLog := ngLog.NewLog()
	tests := []struct {
		name string
		args args
		want WebServerConfigService
	}{
		{
			name: "test NewWebServerConfigService function",
			args: args{
				configuration: testConfiguration,
				serverBinPath: "/usr/local/nginx/sbin/nginx",
				logsDir:       "/usr/local/nginx/logs",
				log:           nginxLog,
			},
			want: NewWebServerConfigService(testConfiguration, "/usr/local/nginx/sbin/nginx", "/usr/local/nginx/logs", nginxLog),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebServerConfigService(tt.args.configuration, tt.args.serverBinPath, tt.args.logsDir, tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebServerConfigService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWebServerInfo(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want WebServerInfo
	}{
		{
			name: "test NewWebServerInfo function",
			args: args{name: "testWebServer"},
			want: NewWebServerInfo("testWebServer"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebServerInfo(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebServerInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffstage_DisplayConfig(t *testing.T) {
	type fields struct {
		webServerConfigServices map[string]WebServerConfigService
		webServerConfigManagers map[string]configuration.ConfigManager
		metrics                 Metrics
	}
	type args struct {
		serverName string
	}
	l := loader.NewLoader()
	ctx, loopPreventer, err := l.LoadFromFilePath("F:\\GO_Project\\src\\bifrost\\test\\config_test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	serverName := "testWebServer"
	notExistServerName := "testNotExistServer"
	testConfiguration := configuration.NewConfiguration(ctx.(*parser.Config), loopPreventer, new(sync.RWMutex))
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(l, testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
	metrics := NewMetrics(func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}, make(chan error))
	wantOffstage := NewOffstage(wscServices, configManagers, metrics)
	want, wantErr := wantOffstage.DisplayConfig(serverName)
	notExistWant, notExistWantErr := wantOffstage.DisplayConfig(notExistServerName)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test Offstage.DisplayConfig",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args:    args{serverName: serverName},
			want:    want,
			wantErr: wantErr != nil,
		},
		{
			name: "test not exist Web Server DisplayConfig",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args:    args{serverName: notExistServerName},
			want:    notExistWant,
			wantErr: notExistWantErr != nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			got, err := o.DisplayConfig(tt.args.serverName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DisplayConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DisplayConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffstage_DisplayStatus(t *testing.T) {
	type fields struct {
		webServerConfigServices map[string]WebServerConfigService
		webServerConfigManagers map[string]configuration.ConfigManager
		metrics                 Metrics
	}
	l := loader.NewLoader()
	ctx, loopPreventer, err := l.LoadFromFilePath("F:\\GO_Project\\src\\bifrost\\test\\config_test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	serverName := "testWebServer"
	testConfiguration := configuration.NewConfiguration(ctx.(*parser.Config), loopPreventer, new(sync.RWMutex))
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(l, testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
	metrics := NewMetrics(func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}, make(chan error))
	wantOffstage := NewOffstage(wscServices, configManagers, metrics)
	want, wantErr := wantOffstage.DisplayStatus()
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "test Offstage.DisplayStatus method",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			want:    want,
			wantErr: wantErr != nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			got, err := o.DisplayStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("DisplayStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DisplayStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffstage_GetConfig(t *testing.T) {
	type fields struct {
		webServerConfigServices map[string]WebServerConfigService
		webServerConfigManagers map[string]configuration.ConfigManager
		metrics                 Metrics
	}
	type args struct {
		serverName string
	}
	l := loader.NewLoader()
	ctx, loopPreventer, err := l.LoadFromFilePath("F:\\GO_Project\\src\\bifrost\\test\\config_test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	serverName := "testWebServer"
	notExistServerName := "testNotExistServer"
	testConfiguration := configuration.NewConfiguration(ctx.(*parser.Config), loopPreventer, new(sync.RWMutex))
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(l, testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
	metrics := NewMetrics(func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}, make(chan error))
	wantOffstage := NewOffstage(wscServices, configManagers, metrics)
	want, wantErr := wantOffstage.GetConfig(serverName)
	notExistWant, notExistWantErr := wantOffstage.GetConfig(notExistServerName)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test Offstage.DisplayConfig",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args:    args{serverName: serverName},
			want:    want,
			wantErr: wantErr != nil,
		},
		{
			name: "test not exist Web Server DisplayConfig",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args:    args{serverName: notExistServerName},
			want:    notExistWant,
			wantErr: notExistWantErr != nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			got, err := o.GetConfig(tt.args.serverName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffstage_Range(t *testing.T) {
	type fields struct {
		webServerConfigServices map[string]WebServerConfigService
		webServerConfigManagers map[string]configuration.ConfigManager
		metrics                 Metrics
	}
	type args struct {
		rangeFunc func(serverName string, configService WebServerConfigService) bool
	}
	l := loader.NewLoader()
	ctx, loopPreventer, err := l.LoadFromFilePath("F:\\GO_Project\\src\\bifrost\\test\\config_test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	serverName := "testWebServer"
	testConfiguration := configuration.NewConfiguration(ctx.(*parser.Config), loopPreventer, new(sync.RWMutex))
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(l, testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
	metrics := NewMetrics(func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}, make(chan error))
	wantOffstage := NewOffstage(wscServices, configManagers, metrics)
	gotStatus := make(map[string]State)
	wantStatus := make(map[string]State)
	wantOffstage.Range(func(serverName string, configService WebServerConfigService) bool {
		wantStatus[serverName] = configService.ServerStatus()
		return true
	})
	gotVersion := make(map[string]string)
	wantVersion := make(map[string]string)
	wantOffstage.Range(func(serverName string, configService WebServerConfigService) bool {
		wantVersion[serverName] = configService.ServerVersion()
		return true
	})
	gotWebServerInfos := make([]WebServerInfo, 0)
	wantWebServerInfos := make([]WebServerInfo, 0)
	wantOffstage.Range(func(serverName string, configService WebServerConfigService) bool {
		info := NewWebServerInfo(serverName)
		info.Version = configService.ServerVersion()
		info.Status = configService.ServerStatus()
		wantWebServerInfos = append(wantWebServerInfos, info)
		return true
	})
	tests := []struct {
		name   string
		fields fields
		args   args
		got    interface{}
		want   interface{}
	}{
		{
			name: "test Offstage.Range method with ServerStatus",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{rangeFunc: func(serverName string, configService WebServerConfigService) bool {
				gotStatus[serverName] = configService.ServerStatus()
				return true
			}},
			got:  gotStatus,
			want: wantStatus,
		},
		{
			name: "test Offstage.Range method with ServerVersion",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{rangeFunc: func(serverName string, configService WebServerConfigService) bool {
				gotVersion[serverName] = configService.ServerVersion()
				return true
			}},
			got:  gotVersion,
			want: wantVersion,
		},
		{
			name: "test Offstage.Range method with report WebServerInfos",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{rangeFunc: func(serverName string, configService WebServerConfigService) bool {
				info := NewWebServerInfo(serverName)
				info.Version = configService.ServerVersion()
				info.Status = configService.ServerStatus()
				gotWebServerInfos = append(gotWebServerInfos, info)
				return true
			}},
			got:  &gotWebServerInfos,
			want: &wantWebServerInfos,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			o.Range(tt.args.rangeFunc)
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Range() got = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestOffstage_ShowStatistics(t *testing.T) {
	type fields struct {
		webServerConfigServices map[string]WebServerConfigService
		webServerConfigManagers map[string]configuration.ConfigManager
		metrics                 Metrics
	}
	type args struct {
		serverName string
	}
	l := loader.NewLoader()
	ctx, loopPreventer, err := l.LoadFromFilePath("F:\\GO_Project\\src\\bifrost\\test\\config_test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	serverName := "testWebServer"
	notExistServerName := "testNotExistServer"
	testConfiguration := configuration.NewConfiguration(ctx.(*parser.Config), loopPreventer, new(sync.RWMutex))
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(l, testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
	metrics := NewMetrics(func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}, make(chan error))
	wantOffstage := NewOffstage(wscServices, configManagers, metrics)
	want, wantErr := wantOffstage.ShowStatistics(serverName)
	notExistWant, notExistWantErr := wantOffstage.ShowStatistics(notExistServerName)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test Offstage.ShowStatistics",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args:    args{serverName: serverName},
			want:    want,
			wantErr: wantErr != nil,
		},
		{
			name: "test not exist Web Server Statistics",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args:    args{serverName: notExistServerName},
			want:    notExistWant,
			wantErr: notExistWantErr != nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			got, err := o.ShowStatistics(tt.args.serverName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShowStatistics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShowStatistics() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffstage_Start(t *testing.T) {
	type fields struct {
		webServerConfigServices map[string]WebServerConfigService
		webServerConfigManagers map[string]configuration.ConfigManager
		metrics                 Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			if err := o.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOffstage_Stop(t *testing.T) {
	type fields struct {
		webServerConfigServices map[string]WebServerConfigService
		webServerConfigManagers map[string]configuration.ConfigManager
		metrics                 Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			if err := o.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOffstage_UpdateConfig(t *testing.T) {
	type fields struct {
		webServerConfigServices map[string]WebServerConfigService
		webServerConfigManagers map[string]configuration.ConfigManager
		metrics                 Metrics
	}
	type args struct {
		serverName string
		data       []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			if err := o.UpdateConfig(tt.args.serverName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOffstage_WatchLog(t *testing.T) {
	type fields struct {
		webServerConfigServices map[string]WebServerConfigService
		webServerConfigManagers map[string]configuration.ConfigManager
		metrics                 Metrics
	}
	type args struct {
		serverName string
		logName    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    LogWatcher
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			got, err := o.WatchLog(tt.args.serverName, tt.args.logName)
			if (err != nil) != tt.wantErr {
				t.Errorf("WatchLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WatchLog() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServerConfigService_ServerStatus(t *testing.T) {
	type fields struct {
		configuration configuration.Configuration
		serverBinPath string
		logsDir       string
		log           *ngLog.Log
	}
	tests := []struct {
		name   string
		fields fields
		want   State
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WebServerConfigService{
				configuration: tt.fields.configuration,
				serverBinPath: tt.fields.serverBinPath,
				logsDir:       tt.fields.logsDir,
				log:           tt.fields.log,
			}
			if got := w.ServerStatus(); got != tt.want {
				t.Errorf("ServerStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServerConfigService_ServerVersion(t *testing.T) {
	type fields struct {
		configuration configuration.Configuration
		serverBinPath string
		logsDir       string
		log           *ngLog.Log
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WebServerConfigService{
				configuration: tt.fields.configuration,
				serverBinPath: tt.fields.serverBinPath,
				logsDir:       tt.fields.logsDir,
				log:           tt.fields.log,
			}
			if got := w.ServerVersion(); got != tt.want {
				t.Errorf("ServerVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logWatcher_Close(t *testing.T) {
	type fields struct {
		dataChan          chan []byte
		transferErrorChan chan error
		closeFunc         func() error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logWatcher{
				dataChan:          tt.fields.dataChan,
				transferErrorChan: tt.fields.transferErrorChan,
				closeFunc:         tt.fields.closeFunc,
			}
			if err := l.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_logWatcher_GetDataChan(t *testing.T) {
	type fields struct {
		dataChan          chan []byte
		transferErrorChan chan error
		closeFunc         func() error
	}
	tests := []struct {
		name   string
		fields fields
		want   <-chan []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logWatcher{
				dataChan:          tt.fields.dataChan,
				transferErrorChan: tt.fields.transferErrorChan,
				closeFunc:         tt.fields.closeFunc,
			}
			if got := l.GetDataChan(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataChan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logWatcher_GetTransferErrorChan(t *testing.T) {
	type fields struct {
		dataChan          chan []byte
		transferErrorChan chan error
		closeFunc         func() error
	}
	tests := []struct {
		name   string
		fields fields
		want   <-chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logWatcher{
				dataChan:          tt.fields.dataChan,
				transferErrorChan: tt.fields.transferErrorChan,
				closeFunc:         tt.fields.closeFunc,
			}
			if got := l.GetTransferErrorChan(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransferErrorChan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metrics_Report(t *testing.T) {
	type fields struct {
		OS                string
		Time              string
		Cpu               string
		Mem               string
		Disk              string
		StatusList        []WebServerInfo
		BifrostVersion    string
		isStoped          bool
		monitorErrChan    chan error
		webServerInfoFunc func() []WebServerInfo
		locker            *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &metrics{
				OS:                tt.fields.OS,
				Time:              tt.fields.Time,
				Cpu:               tt.fields.Cpu,
				Mem:               tt.fields.Mem,
				Disk:              tt.fields.Disk,
				StatusList:        tt.fields.StatusList,
				BifrostVersion:    tt.fields.BifrostVersion,
				isStoped:          tt.fields.isStoped,
				monitorErrChan:    tt.fields.monitorErrChan,
				webServerInfoFunc: tt.fields.webServerInfoFunc,
				locker:            tt.fields.locker,
			}
			got, err := m.Report()
			if (err != nil) != tt.wantErr {
				t.Errorf("Report() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Report() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metrics_Start(t *testing.T) {
	type fields struct {
		OS                string
		Time              string
		Cpu               string
		Mem               string
		Disk              string
		StatusList        []WebServerInfo
		BifrostVersion    string
		isStoped          bool
		monitorErrChan    chan error
		webServerInfoFunc func() []WebServerInfo
		locker            *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &metrics{
				OS:                tt.fields.OS,
				Time:              tt.fields.Time,
				Cpu:               tt.fields.Cpu,
				Mem:               tt.fields.Mem,
				Disk:              tt.fields.Disk,
				StatusList:        tt.fields.StatusList,
				BifrostVersion:    tt.fields.BifrostVersion,
				isStoped:          tt.fields.isStoped,
				monitorErrChan:    tt.fields.monitorErrChan,
				webServerInfoFunc: tt.fields.webServerInfoFunc,
				locker:            tt.fields.locker,
			}
			if err := m.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_metrics_Stop(t *testing.T) {
	type fields struct {
		OS                string
		Time              string
		Cpu               string
		Mem               string
		Disk              string
		StatusList        []WebServerInfo
		BifrostVersion    string
		isStoped          bool
		monitorErrChan    chan error
		webServerInfoFunc func() []WebServerInfo
		locker            *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &metrics{
				OS:                tt.fields.OS,
				Time:              tt.fields.Time,
				Cpu:               tt.fields.Cpu,
				Mem:               tt.fields.Mem,
				Disk:              tt.fields.Disk,
				StatusList:        tt.fields.StatusList,
				BifrostVersion:    tt.fields.BifrostVersion,
				isStoped:          tt.fields.isStoped,
				monitorErrChan:    tt.fields.monitorErrChan,
				webServerInfoFunc: tt.fields.webServerInfoFunc,
				locker:            tt.fields.locker,
			}
			if err := m.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
