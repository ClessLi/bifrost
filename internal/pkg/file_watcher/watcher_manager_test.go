package file_watcher

import (
	"context"
	"sync"
	"testing"
	"time"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
)

func TestWatcherManager_Watch(t *testing.T) {
	type fields struct {
		config   *Config
		watchers map[string]*FileWatcher
	}
	type args struct {
		ctx  context.Context
		file string
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
				ctx:  context.Background(),
				file: "../../../test/nginx/logs/access.log",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := &WatcherManager{
				config:   tt.fields.config,
				mu:       sync.RWMutex{},
				watchers: tt.fields.watchers,
			}
			got, err := wm.Watch(tt.args.ctx, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			wg := new(sync.WaitGroup)
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case data := <-got:
						if data == nil {
							logV1.Infof("watch stopped")
							return
						}
						t.Logf("%s", data)
					}
				}
			}()
			wg.Wait()
		})
	}
}
