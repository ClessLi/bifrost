package monitor

import (
	"context"
	"fmt"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"sync"
	"time"

	"github.com/marmotedu/errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

const (
	MaxFrequencyPerCycle = 100
	MinCycle             = time.Minute
	MaxCycle             = time.Hour
	MinSyncInterval      = time.Second
	MaxSyncInterval      = time.Minute * 10
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
	MonitoringSyncInterval      time.Duration
	MonitoringCycle             time.Duration
	MonitoringFrequencyPerCycle int

	ctx         context.Context
	cancel      context.CancelFunc
	procLocker  sync.Locker
	procStarted bool

	cache       *SystemInfo
	cachemu     *sync.RWMutex
	current     *SystemInfo
	watchLocker sync.Locker
	cannotSync  bool
}

func (m *monitor) Start() error {
	m.procLocker.Lock()
	// check monitor is already start or not.
	if m.ctx != nil {
		if m.ctx.Err() == nil || m.procStarted {
			return errors.New("monitor is already start")
		}
	}
	m.procStarted = true
	defer func() { m.procStarted = false }()

	// init process context
	m.ctx, m.cancel = context.WithCancel(context.Background())
	defer m.cancel()

	m.procLocker.Unlock()

	// ini sync duration
	syncDuration := m.MonitoringSyncInterval
	switch {
	case syncDuration < MinSyncInterval:
		syncDuration = MinSyncInterval
	case syncDuration > MaxSyncInterval:
		syncDuration = MaxSyncInterval
	}

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

	// sync system information
	go func() {
		syncWork, syncCancel := context.WithCancel(m.ctx)
		defer syncCancel()
		logV1.Info("start to sync system information")
		syncTicker := time.NewTicker(syncDuration)
		for {
			select {
			case <-syncWork.Done():
				logV1.Info("sync system information stopped")

				return
			case <-syncTicker.C:
				m.infoSync(syncWork)
			}
		}
	}()

	return m.watch(cycle, interval)
}

func (m *monitor) Stop() error {
	if m.ctx == nil || m.cancel == nil {
		return errors.New("monitor is not started or initialization.")
	}
	timeoutC, cancel := context.WithTimeout(m.ctx, time.Second*10)
	defer cancel()
	go func() {
		logV1.Info("monitoring stopping...")
		m.cancel()
		logV1.Debugf("monitoring stop complete.")
	}()
	select {
	case <-m.ctx.Done():
		logV1.Info("monitoring stopped")

		return m.ctx.Err()
	case <-timeoutC.Done():
		logV1.Errorf("monitoring stop timeout")

		return timeoutC.Err()
	}
}

func (m *monitor) Report() SystemInfo {
	m.cachemu.RLock()
	defer m.cachemu.RUnlock()

	return *m.cache
}

func (m *monitor) infoSync(ctx context.Context) {
	if m.cannotSync {
		logV1.Error("infoSync() call blocked!!")

		return
	}
	m.cannotSync = true
	defer func() {
		m.cannotSync = false
	}()
	m.watchLocker.Lock()
	sysinfo := *m.current
	m.watchLocker.Unlock()
	work, done := context.WithCancel(ctx)
	defer done()
	timeout, tc := context.WithTimeout(work, time.Second*5)
	defer tc()
	go func() {
		m.cachemu.Lock()
		defer m.cachemu.Unlock()
		*m.cache = sysinfo
		logV1.Infof("system info sync to cache succeeded.")
		done()
	}()
	select {
	case <-timeout.Done():
		logV1.Warn("system info sync to cache timeout!")

		return
	case <-work.Done():
	}
}

func (m *monitor) watch(cycle, interval time.Duration) error {
	var err error
	defer func() {
		if err != nil {
			logV1.Warn(err.Error())
			err = m.Stop()
			if err != nil {
				logV1.Error(err.Error())
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
		case <-m.ctx.Done():
			logV1.Info("watching completed.")

			return m.ctx.Err()
		case <-intervalTicker.C:
			var (
				cpupct   []float64
				vmem     *mem.VirtualMemoryStat
				diskstat *disk.UsageStat
			)

			cpupct, err = cpu.Percent(0, false)
			if err != nil {
				return err
			}
			cpupcts = append(cpupcts, cpupct[0])
			vmem, err = mem.VirtualMemory()
			if err != nil {
				return err
			}
			mempcts = append(mempcts, vmem.UsedPercent)
			diskstat, err = disk.Usage("/")
			if err != nil {
				return err
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
