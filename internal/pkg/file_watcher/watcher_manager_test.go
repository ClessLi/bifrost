package file_watcher

import (
	"context"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"sync"
	"testing"
	"time"
)

func TestWatcherManager_Watch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	output1 := make(chan []byte)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case data := <-output1:
				if data == nil {
					log.Infof("watch stopped")
					return
				}
				t.Logf("%s", data)
			case <-ctx.Done():
				return
			}
		}
	}()
	type fields struct {
		config   *Config
		watchers map[string]*FileWatcher
	}
	type args struct {
		file    string
		outputC chan []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "new watcher manager test",
			fields: fields{
				config: &Config{
					MaxConnections: 10,
					OutputTimeout:  time.Second * 20,
				},
				watchers: make(map[string]*FileWatcher),
			},
			args: args{
				file:    "F:\\GO_Project\\src\\bifrost\\test\\nginx\\logs\\access.log",
				outputC: output1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := &WatcherManager{
				config:   tt.fields.config,
				watchers: tt.fields.watchers,
			}
			if err := wm.Watch(tt.args.file, tt.args.outputC); (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}
