package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"

	"github.com/marmotedu/errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sync/errgroup"
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

	eg *errgroup.Group

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
			return errors.New("the monitor is already start")
		}
	}
	m.procStarted = true
	defer func() { m.procStarted = false }()

	// init process context
	var workCtx context.Context
	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.eg, workCtx = errgroup.WithContext(m.ctx)

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
	m.eg.Go(func() error {
		logV1.Info("start to sync system information")
		syncTicker := time.NewTicker(syncDuration)
		for {
			select {
			case <-workCtx.Done():
				logV1.Info("sync system information stopped")

				return workCtx.Err()
			case <-syncTicker.C:
				m.infoSync(workCtx)
			}
		}
	})

	// system information collection
	m.eg.Go(func() error {
		logV1.Info("start to collect system information")
		err := m.watch(workCtx, cycle, interval)
		if err != nil && !errors.Is(err, context.Canceled) {
			logV1.Warnf("collect system information failed: %v", err)
			err = m.Stop()
			if err != nil {
				logV1.Errorf(err.Error())

				return err
			}
		}

		return nil
	})

	if err := m.ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (m *monitor) Stop() error {
	if m.ctx == nil || m.cancel == nil {
		return errors.New("the monitor is not started or initialization.")
	}
	timeoutC, done := context.WithTimeout(context.TODO(), time.Second*10)
	defer done()

	stopeg, _ := errgroup.WithContext(timeoutC)
	stopeg.Go(func() (err error) {
		defer func() {
			done()
			if errors.Is(err, context.Canceled) {
				logV1.Debugf("the monitor stopped successfully.")
			}
		}()
		logV1.Info("the monitor is stopping...")
		m.cancel()

		return m.eg.Wait()
	})

	if err := stopeg.Wait(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logV1.Errorf("the monitor stopped timeout")

			return errors.New("the monitor stopped timeout")
		} else if !errors.Is(err, context.Canceled) {
			logV1.Errorf("the monitor stopped failed: %v", err)

			return err
		}
		logV1.Info("the monitor has been stopped")

		return nil
	}
	logV1.Errorf("the monitor is already running!")

	return errors.New("the monitor is already running!")
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
	work, done := context.WithTimeout(ctx, time.Second*5)
	defer done()
	go func(done context.CancelFunc) {
		defer done()
		m.cachemu.Lock()
		defer m.cachemu.Unlock()
		*m.cache = sysinfo
	}(done)

	<-work.Done()
	if err := work.Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logV1.Warn("sync the system info to cache timeout!")

			return
		} else if !errors.Is(err, context.Canceled) {
			logV1.Warnf("sync the system info to cache failed: %v", err)

			return
		}
		logV1.Infof("sync the system info to cache successfully!")

		return
	}
	logV1.Errorf("the system information synchronization task has not been stopped!")
}

func (m *monitor) watch(ctx context.Context, cycle, interval time.Duration) error {
	var err error
	cpupcts := make([]float64, 0)
	mempcts := make([]float64, 0)
	diskpcts := make([]float64, 0)
	cycleTicker := time.NewTicker(cycle)
	intervalTicker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			logV1.Info("watching completed.")

			return ctx.Err()
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
