package file_watcher

import (
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"sync"
	"testing"
	"time"
)

func TestWatcherManager_Watch(t *testing.T) {
	type fields struct {
		config   *Config
		mu       sync.RWMutex
		watchers map[string]*FileWatcher
	}
	type args struct {
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
				mu:       sync.RWMutex{},
				watchers: make(map[string]*FileWatcher),
			},
			args: args{
				file: "F:\\GO_Project\\src\\bifrost\\test\\nginx\\logs\\access.log",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := &WatcherManager{
				config:   tt.fields.config,
				mu:       tt.fields.mu,
				watchers: tt.fields.watchers,
			}
			got, err := wm.Watch(tt.args.file)
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
							log.Infof("watch stopped")
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
