package service

import (
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	ngLog "github.com/ClessLi/bifrost/pkg/server_log/nginx"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
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
	serverName := "testWebServer"
	notExistServerName := "testNotExistServer"
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(loader.NewLoader(), testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
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
	serverName := "testWebServer"
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(loader.NewLoader(), testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
	metrics := NewMetrics(func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}, make(chan error))
	wantOffstage := NewOffstage(wscServices, configManagers, metrics)
	want, wantErr := wantOffstage.DisplayServersStatus()
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "test Offstage.DisplayServersStatus method",
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
			got, err := o.DisplayServersStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("DisplayServersStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DisplayServersStatus() got = %v, want %v", got, tt.want)
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
	serverName := "testWebServer"
	notExistServerName := "testNotExistServer"
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(loader.NewLoader(), testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
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
	serverName := "testWebServer"
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(loader.NewLoader(), testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
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
	serverName := "testWebServer"
	notExistServerName := "testNotExistServer"
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(loader.NewLoader(), testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
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
	serverName := "testWebServer"
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(loader.NewLoader(), testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
	metrics := NewMetrics(func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}, make(chan error))
	wantOffstage := NewOffstage(wscServices, configManagers, metrics)
	wantErr := wantOffstage.Start()
	_ = wantOffstage.Stop()
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test Offstage.Start method",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			wantErr: wantErr != nil,
		},
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
			defer o.Stop()
		})
	}
}

func TestOffstage_Stop(t *testing.T) {
	type fields struct {
		webServerConfigServices map[string]WebServerConfigService
		webServerConfigManagers map[string]configuration.ConfigManager
		metrics                 Metrics
	}
	serverName := "testWebServer"
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(loader.NewLoader(), testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
	metrics := NewMetrics(func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}, make(chan error))
	wantOffstage := NewOffstage(wscServices, configManagers, metrics)
	_ = wantOffstage.Start()
	wantErr := wantOffstage.Stop()
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test Offstage.Stop method",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			wantErr: wantErr != nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			_ = o.Start()
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
	serverName := "testWebServer"
	notExistServerName := "testNotExistServer"
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "null", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(loader.NewLoader(), testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
	metrics := NewMetrics(func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}, make(chan error))
	wantOffstage := NewOffstage(wscServices, configManagers, metrics)
	enabledData, err := wantOffstage.GetConfig(serverName)
	if err != nil {
		t.Fatal(err)
	}
	disabledData := []byte("disabled config data")
	_ = wantOffstage.Start()
	wantErr := wantOffstage.UpdateConfig(serverName, enabledData)
	_ = wantOffstage.Stop()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		isStart bool
	}{
		{
			name: "1) test Offstage.UpdateConfig method with disabled config data, not started and not exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: notExistServerName,
				data:       disabledData,
			},
			wantErr: wantErr != nil,
			isStart: false,
		},
		{
			name: "2) test Offstage.UpdateConfig method with disabled config data, not started and exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: serverName,
				data:       disabledData,
			},
			wantErr: wantErr != nil,
			isStart: false,
		},
		{
			name: "3) test Offstage.UpdateConfig method with disabled config data, started and not exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: notExistServerName,
				data:       disabledData,
			},
			wantErr: wantErr != nil,
			isStart: true,
		},
		{
			name: "4) test Offstage.UpdateConfig method with enabled config data, not started and not exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: notExistServerName,
				data:       enabledData,
			},
			wantErr: wantErr != nil,
			isStart: false,
		},
		{
			name: "5) test Offstage.UpdateConfig method with disabled config data, started and exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: serverName,
				data:       disabledData,
			},
			wantErr: wantErr != nil,
			isStart: true,
		},
		{
			name: "6) test Offstage.UpdateConfig method with enabled config data, not started and exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: serverName,
				data:       enabledData,
			},
			wantErr: wantErr != nil,
			isStart: false,
		},
		{
			name: "7) test Offstage.UpdateConfig method with enabled config data, started and not exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: notExistServerName,
				data:       enabledData,
			},
			wantErr: wantErr != nil,
			isStart: true,
		},
		{
			name: "8) test Offstage.UpdateConfig method with enabled config data, started and exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: serverName,
				data:       enabledData,
			},
			wantErr: wantErr != nil,
			isStart: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			if tt.isStart {
				_ = o.Start()
			}
			if err := o.UpdateConfig(tt.args.serverName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.isStart {
				_ = o.Stop()
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
	serverName := "testWebServer"
	notExistServerName := "testNotExistServer"
	testLogData := []byte("test access log\n")
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wscServices := map[string]WebServerConfigService{serverName: NewWebServerConfigService(testConfiguration, "null", "F:\\GO_Project\\src\\bifrost\\test\\nginx\\logs", ngLog.NewLog())}
	configManagers := map[string]configuration.ConfigManager{serverName: configuration.NewNginxConfigurationManager(loader.NewLoader(), testConfiguration, "null", "", 1, 7, new(sync.RWMutex))}
	metrics := NewMetrics(func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}, make(chan error))
	wantOffstage := NewOffstage(wscServices, configManagers, metrics)
	_, wantDisabledErr := wantOffstage.WatchLog(serverName, "error.log")
	_, wantNotExistErr := wantOffstage.WatchLog(notExistServerName, "access.log")
	wantLogWatcher, wantErr := wantOffstage.WatchLog(serverName, "access.log")
	if wantLogWatcher == nil {
		t.Fatal(wantErr)
	}
	time.Sleep(time.Second)

	accessLog, err := os.OpenFile("F:\\GO_Project\\src\\bifrost\\test\\nginx\\logs\\access.log", os.O_APPEND, 0755)
	if err != nil {
		t.Fatal(err)
	}
	_, err = accessLog.Write(testLogData)
	if err != nil {
		_ = accessLog.Close()
		t.Fatal(err)
	}
	_ = accessLog.Close()
	var want []byte
	select {
	case <-time.After(time.Second * 5):
		t.Fatal("catch want access.log data timeout")
	case want = <-wantLogWatcher.GetDataChan():
		t.Logf("want data: %s", want)
		break
	}

	wantCloseErr := wantLogWatcher.Close()
	time.Sleep(time.Second)
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         []byte
		wantErr      bool
		wantCloseErr bool
	}{
		{
			name: "test Offstage.WatchLog method with enabled log and exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: serverName,
				logName:    "access.log",
			},
			want:         want,
			wantErr:      wantErr != nil,
			wantCloseErr: wantCloseErr != nil,
		},
		{
			name: "test Offstage.WatchLog method with disabled log and exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: serverName,
				logName:    "error.log",
			},
			want:         nil,
			wantErr:      wantDisabledErr != nil,
			wantCloseErr: false,
		},
		{
			name: "test Offstage.WatchLog method with not exist web server",
			fields: fields{
				webServerConfigServices: wscServices,
				webServerConfigManagers: configManagers,
				metrics:                 metrics,
			},
			args: args{
				serverName: notExistServerName,
				logName:    "access.log",
			},
			want:         nil,
			wantErr:      wantNotExistErr != nil,
			wantCloseErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer time.Sleep(time.Second)
			o := Offstage{
				webServerConfigServices: tt.fields.webServerConfigServices,
				webServerConfigManagers: tt.fields.webServerConfigManagers,
				metrics:                 tt.fields.metrics,
			}
			gotLogWatcher, err := o.WatchLog(tt.args.serverName, tt.args.logName)
			if (err != nil) != tt.wantErr {
				t.Errorf("WatchLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			time.Sleep(time.Second)
			gotAccessLog, err := os.OpenFile("F:\\GO_Project\\src\\bifrost\\test\\nginx\\logs\\access.log", os.O_APPEND, 0755)
			if err != nil {
				t.Fatal(err)
			}
			_, err = gotAccessLog.Write(testLogData)
			if err != nil {
				_ = gotAccessLog.Close()
				t.Fatal(err)
			}

			_ = gotAccessLog.Close()

			var got []byte
			select {
			case <-time.After(time.Second * 5):
				t.Errorf("catch got access.log data timeout")
			case got = <-gotLogWatcher.GetDataChan():
				break
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WatchLog() got = %v, want %v", got, tt.want)
			}
			err = gotLogWatcher.Close()
			if (err != nil) != tt.wantCloseErr {
				t.Errorf("LogWatcher.Close() error = %v, wantCloseErr %v", err, tt.wantCloseErr)
				return
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
	logDir := "F:\\GO_Project\\src\\bifrost\\test\\nginx\\logs"
	nginxLog := ngLog.NewLog()
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wantWebServerConfigService := NewWebServerConfigService(testConfiguration, "null", logDir, nginxLog)
	want := wantWebServerConfigService.ServerStatus()

	tests := []struct {
		name   string
		fields fields
		want   State
	}{
		{
			name: "test WebServerConfigService.ServerStatus method",
			fields: fields{
				configuration: testConfiguration,
				serverBinPath: "null",
				logsDir:       logDir,
				log:           nginxLog,
			},
			want: want,
		},
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
	logDir := "F:\\GO_Project\\src\\bifrost\\test\\nginx\\logs"
	nginxLog := ngLog.NewLog()
	testConfiguration, err := configuration.NewConfigurationFromPath("F:\\GO_Project\\src\\bifrost\\test\\nginx\\conf\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	wantWebServerConfigService := NewWebServerConfigService(testConfiguration, "null", logDir, nginxLog)
	want := wantWebServerConfigService.ServerVersion()
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test WebServerConfigService.ServerVersion method",
			fields: fields{
				configuration: testConfiguration,
				serverBinPath: "null",
				logsDir:       logDir,
				log:           nginxLog,
			},
			want: want,
		},
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
	serverName := "testWebServer"
	wsiFunc := func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}
	wantMetrics := NewMetrics(wsiFunc, make(chan error))
	want, wantErr := wantMetrics.Report()
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "test metrics Report method",
			fields: fields{
				OS:                wantMetrics.(*metrics).OS,
				Time:              wantMetrics.(*metrics).Time,
				Cpu:               wantMetrics.(*metrics).Cpu,
				Mem:               wantMetrics.(*metrics).Mem,
				Disk:              wantMetrics.(*metrics).Disk,
				StatusList:        make([]WebServerInfo, 0),
				BifrostVersion:    config.GetVersion(),
				isStoped:          true,
				monitorErrChan:    make(chan error),
				webServerInfoFunc: wsiFunc,
				locker:            new(sync.RWMutex),
			},
			want:    want,
			wantErr: wantErr != nil,
		},
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
				t.Errorf("Report() got = %s, want %s", got, tt.want)
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
	serverName := "testWebServer"
	wsiFunc := func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}
	wantMetrics := NewMetrics(wsiFunc, make(chan error))
	wantErr := wantMetrics.Start()
	_ = wantMetrics.Stop()
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test metrics Start method",
			fields: fields{
				OS:                wantMetrics.(*metrics).OS,
				Time:              wantMetrics.(*metrics).Time,
				Cpu:               wantMetrics.(*metrics).Cpu,
				Mem:               wantMetrics.(*metrics).Mem,
				Disk:              wantMetrics.(*metrics).Disk,
				StatusList:        make([]WebServerInfo, 0),
				BifrostVersion:    config.GetVersion(),
				isStoped:          true,
				monitorErrChan:    make(chan error),
				webServerInfoFunc: wsiFunc,
				locker:            new(sync.RWMutex),
			},
			wantErr: wantErr != nil,
		},
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
			_ = m.Stop()
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
	serverName := "testWebServer"
	wsiFunc := func() []WebServerInfo {
		return []WebServerInfo{{
			Name:    serverName,
			Status:  Normal,
			Version: "test version",
		}}
	}
	wantMetrics := NewMetrics(wsiFunc, make(chan error))
	_ = wantMetrics.Start()
	wantErr := wantMetrics.Stop()
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test metrics Stop method",
			fields: fields{
				OS:                wantMetrics.(*metrics).OS,
				Time:              wantMetrics.(*metrics).Time,
				Cpu:               wantMetrics.(*metrics).Cpu,
				Mem:               wantMetrics.(*metrics).Mem,
				Disk:              wantMetrics.(*metrics).Disk,
				StatusList:        make([]WebServerInfo, 0),
				BifrostVersion:    config.GetVersion(),
				isStoped:          true,
				monitorErrChan:    make(chan error),
				webServerInfoFunc: wsiFunc,
				locker:            new(sync.RWMutex),
			},
			wantErr: wantErr != nil,
		},
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
			_ = m.Start()
			if err := m.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
