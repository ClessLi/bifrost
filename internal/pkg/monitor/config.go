package monitor

import (
	"sync"
	"time"
)

// Config is a structure used to configure a Monitor.
// Its members are sorted roughly in order of importance for composers.
type Config struct {
	MonitoringSyncInterval      time.Duration
	MonitoringCycle             time.Duration
	MonitoringFrequencyPerCycle int
}

// NewConfig returns a Config struct with the default values.
func NewConfig() *Config {
	return &Config{
		MonitoringSyncInterval:      time.Minute * 1,
		MonitoringCycle:             time.Minute * 2,
		MonitoringFrequencyPerCycle: 10,
	}
}

// CompletedConfig is the completed configuration for Monitor.
type CompletedConfig struct {
	*Config
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

func (c CompletedConfig) NewMonitor() (Monitor, error) {
	if c.Config == nil {
		c.Config = NewConfig()
	}

	unknownSysInfo := &SystemInfo{
		CpuUsePct:  "unknown",
		MemUsePct:  "unknown",
		DiskUsePct: "unknown",
	}

	m := &monitor{
		MonitoringSyncInterval:      c.MonitoringSyncInterval,
		MonitoringCycle:             c.MonitoringCycle,
		MonitoringFrequencyPerCycle: c.MonitoringFrequencyPerCycle,
		ctx:                         nil,
		cancel:                      nil,
		procLocker:                  new(sync.Mutex),
		procStarted:                 false,
		cache:                       unknownSysInfo,
		cachemu:                     new(sync.RWMutex),
		current:                     unknownSysInfo,
		watchLocker:                 new(sync.Mutex),
		cannotSync:                  false,
	}

	return Monitor(m), nil
}
