package monitor

import (
	"context"
	"fmt"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"github.com/marmotedu/errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"sync"
	"time"
)

const (
	MaxFrequencyPerCycle = 100
	MinCycle             = time.Minute
	MaxCycle             = time.Hour
)

type SystemInfo struct {
	CpuUsePct  string
	MemUsePct  string
	DiskUsePct string
}

type Monitor interface {
	Start() error
	Stop() error
	Report() SystemInfo
}

var _ Monitor = &monitor{}

type monitor struct {
	MonitoringCycle             time.Duration
	MonitoringFrequencyPerCycle int

	signal      chan struct{}
	startLocker sync.Locker
	isStart     bool

	cache       *SystemInfo
	cachemu     *sync.RWMutex
	current     *SystemInfo
	watchLocker sync.Locker
	cannotSync  bool
}

func (m *monitor) Start() error {
	m.startLocker.Lock()
	defer m.startLocker.Unlock()

	// check monitor is already start or not.
	if m.isStart {
		return errors.New("monitor is already start")
	}

	m.isStart = true
	defer func() { m.isStart = false }()

	// init cycle and frequency interval
	cycle := m.MonitoringCycle
	switch {
	case cycle < MinCycle:
		cycle = MinCycle
	case cycle > MaxCycle:
		cycle = MaxCycle
	}

	frequency := m.MonitoringFrequencyPerCycle
	switch {
	case frequency > MaxFrequencyPerCycle:
		frequency = MaxFrequencyPerCycle
	case frequency <= 0:
		frequency = 1
	}

	interval := cycle / time.Duration(frequency)

	todo, done := context.WithCancel(context.TODO())
	defer done()
	go m.watch(todo, cycle, interval)

	<-m.signal
	return gracefulClose(todo, done, time.Second*10)

}

func gracefulClose(ctx context.Context, close context.CancelFunc, timeout time.Duration) error {
	//timeoutC, cancel := context.WithTimeout(ctx, timeout)
	//defer cancel()
	//
	//go close()
	//
	//select {
	//case <-ctx.Done():
	//	log.Info("monitoring graceful closed")
	//	return nil
	//case <-timeoutC.Done():
	//	log.Warn("monitoring graceful close timeout")
	//	return errors.New("graceful close timeout")
	//}
	// TODO: graceful close function
	panic("invalid function")
}

func (m *monitor) Stop() error {
	timeoutC, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	select {
	case m.signal <- struct{}{}:
		log.Info("monitoring stopped")
		return nil
	case <-timeoutC.Done():
		log.Errorf("monitoring stop timeout")
		return errors.New("monitoring stop timeout")
	}
}

func (m *monitor) Report() SystemInfo {
	m.cachemu.RLock()
	defer m.cachemu.RUnlock()
	return *m.cache
}

func (m *monitor) infoSync() {
	if m.cannotSync {
		log.Error("infoSync() call blocked!!")
		return
	}
	m.cannotSync = true
	defer func() {
		m.cannotSync = false
	}()
	m.watchLocker.Lock()
	sysinfo := *m.current
	m.watchLocker.Unlock()
	todo, done := context.WithCancel(context.TODO())
	defer done()
	timeout, tc := context.WithTimeout(todo, time.Second*5)
	defer tc()
	go func() {
		m.cachemu.Lock()
		defer m.cachemu.Unlock()
		*m.cache = sysinfo
		log.Infof("system info sync to cache succeeded.")
		done()
	}()
	select {
	case <-timeout.Done():
		log.Warn("system info sync to cache timeout!")
		return
	case <-todo.Done():
	}
}

func (m *monitor) watch(ctx context.Context, cycle, interval time.Duration) {
	var err error
	defer func() {
		if err != nil {
			log.Warn(err.Error())
			err = m.Stop()
			if err != nil {
				log.Error(err.Error())
			}
		}
	}()
	cpupcts := make([]float64, 0)
	mempcts := make([]float64, 0)
	diskpcts := make([]float64, 0)
	cycleTicker := time.NewTicker(cycle)
	intervalTicker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			log.Info("watching completed.")
			return
		case <-intervalTicker.C:
			var (
				cpupct   []float64
				vmem     *mem.VirtualMemoryStat
				diskstat *disk.UsageStat
			)

			cpupct, err = cpu.Percent(0, false)
			if err != nil {
				return
			}
			cpupcts = append(cpupcts, cpupct[0])
			vmem, err = mem.VirtualMemory()
			if err != nil {
				return
			}
			mempcts = append(mempcts, vmem.UsedPercent)
			diskstat, err = disk.Usage("/")
			if err != nil {
				return
			}
			diskpcts = append(diskpcts, diskstat.UsedPercent)

		case <-cycleTicker.C:
			m.watchLocker.Lock()
			m.current.CpuUsePct = fmt.Sprintf("%.2f", average(cpupcts...))
			cpupcts = cpupcts[:0]
			m.current.MemUsePct = fmt.Sprintf("%.2f", average(mempcts...))
			mempcts = mempcts[:0]
			m.current.DiskUsePct = fmt.Sprintf("%.2f", average(diskpcts...))
			diskpcts = diskpcts[:0]
			m.watchLocker.Unlock()
		}
	}
}

func average(items ...float64) float64 {
	sum := 0.0
	n := 0
	for i, item := range items {
		sum += item
		n = i
	}
	if sum == 0.0 {
		return sum
	}

	return sum / float64(n+1)

}
