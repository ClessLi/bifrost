package options

import (
	"encoding/json"

	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"

	genericoptions "github.com/yongPhone/bifrost/internal/pkg/options"
	"github.com/yongPhone/bifrost/internal/pkg/server"
	log "github.com/yongPhone/bifrost/pkg/log/v1"
)

type Options struct {
	GenericServerRunOptions    *genericoptions.ServerRunOptions           `json:"server"                 mapstructure:"server"`
	SecureServing              *genericoptions.SecureServingOptions       `json:"secure"                 mapstructure:"secure"`
	InsecureServing            *genericoptions.InsecureServingOptions     `json:"insecure"               mapstructure:"insecure"`
	RAOptions                  *genericoptions.RAOptions                  `json:"ra"                     mapstructure:"ra"`
	GRPCServing                *genericoptions.GRPCServerOptions          `json:"grpc"                   mapstructure:"grpc"`
	WebServerConfigsOptions    *genericoptions.WebServerConfigsOptions    `json:"web-server-configs"     mapstructure:"web-server-configs"`
	MonitorOptions             *genericoptions.MonitorOptions             `json:"monitor"                mapstructure:"monitor"`
	WebServerLogWatcherOptions *genericoptions.WebServerLogWatcherOptions `json:"web-server-log-watcher" mapstructure:"web-server-log-watcher"`
	Log                        *log.Options                               `json:"log"                    mapstructure:"log"`
}

func NewOptions() *Options {
	return &Options{
		GenericServerRunOptions:    genericoptions.NewServerRunOptions(),
		SecureServing:              genericoptions.NewSecureServingOptions(),
		InsecureServing:            genericoptions.NewInsecureServingOptions(),
		RAOptions:                  genericoptions.NewRAOptions(),
		GRPCServing:                genericoptions.NewGRPCServerOptions(),
		WebServerConfigsOptions:    genericoptions.NewWebServerConfigsOptions(),
		MonitorOptions:             genericoptions.NewMonitorOptions(),
		WebServerLogWatcherOptions: genericoptions.NewWebServerLogWatcherOptions(),
		Log:                        log.NewOptions(),
	}
}

func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.RAOptions.AddFlags(fss.FlagSet("RA options"))
	o.GRPCServing.AddFlags(fss.FlagSet("gRPC serving"))
	o.MonitorOptions.AddFlags(fss.FlagSet("monitor"))
	o.WebServerLogWatcherOptions.AddFlags(fss.FlagSet("log watcher"))
	o.Log.AddFlags(fss.FlagSet("log"))

	return fss
}

func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

func (o *Options) Strings() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func (o *Options) Complete() error {
	return o.SecureServing.Complete()
}
