package options

import (
	"github.com/ClessLi/bifrost/internal/pkg/monitor"
	"github.com/marmotedu/errors"
	"github.com/spf13/pflag"
	"time"
)

type MonitorOptions struct {
	SyncInterval      time.Duration `json:"sync-interval" mapstructure:"sync-interval"`
	CycleTime         time.Duration `json:"cycle-time" mapstructure:"cycle-time"`
	FrequencyPerCycle int           `json:"frequency-per-cycle" mapstructure:"frequency-per-cycle"`
}

func NewMonitorOptions() *MonitorOptions {
	defaults := monitor.NewConfig()
	return &MonitorOptions{
		SyncInterval:      defaults.MonitoringSyncInterval,
		CycleTime:         defaults.MonitoringCycle,
		FrequencyPerCycle: defaults.MonitoringFrequencyPerCycle,
	}
}

func (m *MonitorOptions) AddFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&m.SyncInterval, "monitor.sync-interval", m.SyncInterval, ""+
		"")

	fs.DurationVar(&m.CycleTime, "monitor.cycle-time", m.CycleTime, ""+
		"")

	fs.IntVar(&m.FrequencyPerCycle, "monitor.frequency-per-cycle", m.FrequencyPerCycle, ""+
		"")
}

func (m *MonitorOptions) Validate() []error {
	var errs []error

	if m.FrequencyPerCycle <= 0 {
		errs = append(errs, errors.Errorf("--monitor.frequency-per-cycle %d must great than 0.", m.FrequencyPerCycle))
	}

	return errs
}

func (m *MonitorOptions) ApplyTo(c monitor.Config) error {
	c.MonitoringSyncInterval = m.SyncInterval
	c.MonitoringCycle = m.CycleTime
	c.MonitoringFrequencyPerCycle = m.FrequencyPerCycle
	return nil
}
