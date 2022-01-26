package options

import (
	"github.com/ClessLi/bifrost/internal/pkg/file_watcher"
	"github.com/marmotedu/errors"
	"github.com/spf13/pflag"
	"time"
)

type WebServerLogWatcherOptions struct {
	MaxConnections int           `json:"max-connections" mapstructure:"max-connections"`
	WatchTimeout   time.Duration `json:"watch-timeout" mapstructure:"watch-timeout"`
}

func NewWebServerLogWatcherOptions() *WebServerLogWatcherOptions {
	defaults := file_watcher.NewConfig()
	return &WebServerLogWatcherOptions{
		MaxConnections: defaults.MaxConnections,
		WatchTimeout:   defaults.OutputTimeout,
	}
}

func (w *WebServerLogWatcherOptions) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&w.MaxConnections, "web-server-log-watcher.max-connections", w.MaxConnections, ""+
		"")

	fs.DurationVar(&w.WatchTimeout, "web-server-log-watcher.watch-timeout", w.WatchTimeout, ""+
		"")
}

func (w *WebServerLogWatcherOptions) Validate() []error {
	var errs []error

	if w.MaxConnections < 1 {
		errs = append(errs, errors.Errorf("--web-server-log-watcher.max-connections %d must great than 0", w.MaxConnections))
	}

	return errs
}

func (w *WebServerLogWatcherOptions) ApplyTo(c *file_watcher.Config) error {
	c.MaxConnections = w.MaxConnections
	c.OutputTimeout = w.WatchTimeout
	return nil
}
