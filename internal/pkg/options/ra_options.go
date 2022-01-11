package options

import (
	"github.com/ClessLi/bifrost/internal/pkg/server"
	"github.com/marmotedu/errors"
	"github.com/spf13/pflag"
)

// RAOptions contains the options used to connect to the RA server.
type RAOptions struct {
	Host string `json:"host" mapstructure:"host"`
	Port int    `json:"port" mapstructure:"port"`
}

func NewRAOptions() *RAOptions {
	return &RAOptions{}
}

func (ra *RAOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&ra.Host, "ra.host", ra.Host, ""+
		"Specifies the bind address of the Registration Authority server. Set empty to disable.")

	fs.IntVar(&ra.Port, "ra.port", ra.Port, ""+
		"Specifies the bind port of the Registration Authority server.")
}

func (ra *RAOptions) Validate() []error {
	var errs []error

	if len(ra.Host) > 0 && (ra.Port < 1 || ra.Port > 65535) {
		errs = append(
			errs,
			errors.Errorf(
				"--ra.port %v must be between 1 and 65535, inclusive. It cannot be turned off with 0",
				ra.Port,
			),
		)
	}

	return errs
}

func (ra *RAOptions) ApplyTo(c *server.Config) error {
	if ra == nil || ra.Port == 0 {
		c.RA = nil
		return nil
	}

	if c.RA == nil {
		c.RA = &server.RAInfo{}
	}

	c.RA.Host = ra.Host
	c.RA.Port = ra.Port
	return nil
}
