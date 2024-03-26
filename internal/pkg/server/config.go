//go:build !viper_yaml2
// +build !viper_yaml2

package server

import (
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/marmotedu/component-base/pkg/util/homedir"
	"github.com/spf13/viper"
)

const (
	// RecommendedHomeDir defines the default directory used to place all bifrost service configurations.
	RecommendedHomeDir = ".bifrost"

	// RecommendedEnvPrefix defines the ENV prefix used by all bifrost service.
	RecommendedEnvPrefix = "BIFROST"
)

// Config is a structure used to configure a GenericGRPCServer.
// Its members are sorted roughly in order of importance for composers.
type Config struct {
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	RA              *RAInfo
	GRPCServing     *GRPCServingInfo
	Middlewares     []string
	Healthz         bool
	EnableProfiling bool
	EnableMetrics   bool
}

// CertKey contains configuration items related to certificate.
type CertKey struct {
	// CertFile is a file containing a PEM-encoded certificate, and possibly the complete certificate chain
	CertFile string
	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
	KeyFile string
}

// SecureServingInfo holds configuration of the TLS server.
type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

// Address join host IP address and host port number into an address string, like: 0.0.0.0:8443.
func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

// InsecureServingInfo holds configuration of the insecure grpc server.
type InsecureServingInfo struct {
	BindAddress string
	BindPort    int
}

func (i *InsecureServingInfo) Address() string {
	return net.JoinHostPort(i.BindAddress, strconv.Itoa(i.BindPort))
}

// RAInfo holds configuration of the Registration Authority server.
type RAInfo struct {
	Host string
	Port int
}

type GRPCServingInfo struct {
	ChunkSize    int
	RecvTimeoutM int
}

func NewGRPCSeringInfo() *GRPCServingInfo {
	return &GRPCServingInfo{
		ChunkSize:    1024,
		RecvTimeoutM: 1,
	}
}

// NewConfig returns a Config struct with the default values.
func NewConfig() *Config {
	return &Config{
		Healthz:         true,
		Middlewares:     []string{},
		EnableProfiling: true,
		EnableMetrics:   true,
	}
}

// CompletedConfig is the completed configuration for GenericServer.
type CompletedConfig struct {
	*Config
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

// NewGRPCServer returns a new instance of GenericGRPCServer from the given config.
func (c CompletedConfig) NewGRPCServer() (*GenericGRPCServer, error) {
	if c.GRPCServing == nil {
		c.GRPCServing = NewGRPCSeringInfo()
	}

	s := &GenericGRPCServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		RAInfo:              c.RA,
		ChunkSize:           c.GRPCServing.ChunkSize,
		ReceiveTimeout:      time.Duration(c.GRPCServing.RecvTimeoutM) * time.Minute,
		healthz:             c.Healthz,
		enableMetrics:       c.EnableMetrics,
		enableProfiling:     c.EnableProfiling,
		middlewares:         c.Middlewares,
	}

	initGenericGRPCServer(s)

	return s, nil
}

// LoadConfig reads in config file and ENV variables if set.
func LoadConfig(cfg string, defaultName string) {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(homedir.HomeDir(), RecommendedHomeDir))
		viper.AddConfigPath("/etc/bifrost")
		viper.SetConfigName(defaultName)
	}

	// Use config file from the flag.
	viper.SetConfigType("yaml")              // set the type of the configuration to yaml.
	viper.AutomaticEnv()                     // read in environment variables that match.
	viper.SetEnvPrefix(RecommendedEnvPrefix) // set ENVIRONMENT variables prefix to IAM.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		logV1.Warnf("WARNING: viper failed to discover and load the configuration file: %s", err.Error())
	}
}
