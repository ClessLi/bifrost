package options

import (
	"github.com/ClessLi/bifrost/internal/pkg/server"

	"github.com/spf13/pflag"
)

// ServerRunOptions contains the options while running a generic api server.
type ServerRunOptions struct {
	Healthz     bool     `json:"healthz"     mapstructure:"healthz"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters.
func NewServerRunOptions() *ServerRunOptions {
	defaults := server.NewConfig()

	return &ServerRunOptions{
		Healthz:     defaults.Healthz,
		Middlewares: defaults.Middlewares,
	}
}

// AddFlags adds flags for a specific APIServer to the specified FlagSet.
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.Healthz, "server.healthz", s.Healthz, ""+
		"Add self readiness check and install health check service.")

	fs.StringSliceVar(&s.Middlewares, "server.middlewares", s.Middlewares, ""+
		"List of allowed middlewares for server, comma separated. If this is empty default middlewares will be used.")
}

// Validate checks validation of ServerRunOptions.
func (s *ServerRunOptions) Validate() []error {
	return nil
}

// ApplyTo applies the run options to the method receiver and returns self.
func (s *ServerRunOptions) ApplyTo(c *server.Config) error {
	c.Healthz = s.Healthz
	c.Middlewares = s.Middlewares

	return nil
}
