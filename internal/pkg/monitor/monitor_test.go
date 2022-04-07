package monitor

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/marmotedu/errors"
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
			m.infoSync(context.TODO())
			t.Log(*m.cache)
		})
	}
}

func Test_monitor_Start(t *testing.T) {
	type fields struct {
		MonitoringSyncDuration      time.Duration
		MonitoringCycle             time.Duration
		MonitoringFrequencyPerCycle int
		ctx                         context.Context
		cancel                      context.CancelFunc
		procLocker                  sync.Locker
		procStarted                 bool
		cache                       *SystemInfo
		cachemu                     *sync.RWMutex
		current                     *SystemInfo
		watchLocker                 sync.Locker
		cannotSync                  bool
	}
	tests := []struct {
		name   string
		fields fields
		//multi   int
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				MonitoringSyncDuration:      time.Second,
				MonitoringCycle:             time.Second,
				MonitoringFrequencyPerCycle: 30,
				ctx:                         nil,
				cancel:                      nil,
				procLocker:                  new(sync.Mutex),
				procStarted:                 false,
				cache:                       new(SystemInfo),
				cachemu:                     new(sync.RWMutex),
				current:                     new(SystemInfo),
				watchLocker:                 new(sync.Mutex),
				cannotSync:                  false,
			},
			//multi:   10,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &monitor{
				MonitoringSyncInterval:      tt.fields.MonitoringSyncDuration,
				MonitoringCycle:             tt.fields.MonitoringCycle,
				MonitoringFrequencyPerCycle: tt.fields.MonitoringFrequencyPerCycle,
				ctx:                         tt.fields.ctx,
				cancel:                      tt.fields.cancel,
				procLocker:                  tt.fields.procLocker,
				procStarted:                 tt.fields.procStarted,
				cache:                       tt.fields.cache,
				cachemu:                     tt.fields.cachemu,
				current:                     tt.fields.current,
				watchLocker:                 tt.fields.watchLocker,
				cannotSync:                  tt.fields.cannotSync,
			}
			errs := make([]error, 0)
			wg := new(sync.WaitGroup)
			//for i := 0; i < tt.multi; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				err := m.Start()
				if err != nil {
					errs = append(errs, err)
				}
			}()
			go func() {
				defer wg.Done()
				time.Sleep(time.Minute * 2)
				err := m.Stop()
				if err != nil {
					errs = append(errs, err)
				}
			}()
			//}
			wg.Wait()
			if err := errors.NewAggregate(errs); (err != nil) != tt.wantErr {
				//t.Errorf("%d times Start() error = %v, wantErr %v", tt.multi, err, tt.wantErr)
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
