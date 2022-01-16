package monitor

import (
	"context"
	"sync"
	"testing"
	"time"
)

func Test_monitor_infoSync(t *testing.T) {
	type fields struct {
		cache      *SystemInfo
		cachemu    *sync.RWMutex
		current    *SystemInfo
		l          sync.Locker
		cannotSync bool
		timewait   time.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "timeout test",
			fields: fields{
				cache:   &SystemInfo{},
				cachemu: new(sync.RWMutex),
				current: &SystemInfo{
					CpuUsePct:  "10.1%",
					MemUsePct:  "20.3%",
					DiskUsePct: "30.3%",
				},
				l:          new(sync.Mutex),
				cannotSync: false,
				timewait:   time.Second * 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &monitor{
				cache:       tt.fields.cache,
				cachemu:     tt.fields.cachemu,
				current:     tt.fields.current,
				watchLocker: tt.fields.l,
				cannotSync:  tt.fields.cannotSync,
			}
			go func() {
				m.cachemu.RLock()
				defer m.cachemu.RUnlock()
				time.Sleep(tt.fields.timewait)
			}()
			time.Sleep(time.Second)
			m.infoSync()
			t.Log(*m.cache)
		})
	}
}

func Test_gracefulClose(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	//cancelDelay := time.Second*5
	type args struct {
		ctx     context.Context
		close   context.CancelFunc
		timeout time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "graceful close test",
			args: args{
				ctx:     ctx,
				close:   cancel,
				timeout: time.Second * 4,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := gracefulClose(tt.args.ctx, tt.args.close, tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("gracefulClose() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
