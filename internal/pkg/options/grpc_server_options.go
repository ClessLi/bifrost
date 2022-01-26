package options

import (
	"github.com/ClessLi/bifrost/internal/pkg/server"
	"github.com/marmotedu/errors"
	"github.com/spf13/pflag"
)

type GRPCServerOptions struct {
	ChunkSize             int `json:"chunksize" mapstructure:"chunksize"`
	ReceiveTimeoutMinutes int `json:"receive-timeout-minutes" mapstructrue:"receive-timeout-minutes"`
}

func NewGRPCServerOptions() *GRPCServerOptions {
	defaults := server.NewGRPCSeringInfo()
	return &GRPCServerOptions{
		ChunkSize:             defaults.ChunkSize,
		ReceiveTimeoutMinutes: defaults.RecvTimeoutM,
	}
}

// AddFlags adds flags for a specific gRPC Server to the specified FlagSet.
func (g *GRPCServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&g.ChunkSize, "grpc.chunksize", g.ChunkSize, ""+
		"Set the max message size in bytes the server can send. Can not less than 100 bytes.")

	fs.IntVar(&g.ReceiveTimeoutMinutes, "grpc.receive-timeout", g.ReceiveTimeoutMinutes, ""+
		"Set the timeout for receiving data. The unit is per minute.")

}

// Validate checks validation of GRPCServerOptions.
func (g *GRPCServerOptions) Validate() []error {
	var errs []error

	if g.ChunkSize < 100 {
		errs = append(errs, errors.Errorf("--grpc.chunksize %d must great equal 100 bytes.", g.ChunkSize))
	}

	if g.ReceiveTimeoutMinutes <= 0 {
		errs = append(errs, errors.Errorf("--grpc.receive-timeout %d must great than 0.", g.ReceiveTimeoutMinutes))
	}

	return errs
}

// ApplyTo applies the run options to the method receiver and returns self.
func (g *GRPCServerOptions) ApplyTo(c *server.Config) error {
	// GRPCServing is required to serve grpc
	c.GRPCServing = &server.GRPCServingInfo{
		ChunkSize:    g.ChunkSize,
		RecvTimeoutM: g.ReceiveTimeoutMinutes,
	}
	return nil
}
