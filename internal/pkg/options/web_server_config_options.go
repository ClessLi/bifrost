package options

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"

	"github.com/marmotedu/errors"
	"github.com/spf13/pflag"
)

type WebServerConfigOptions struct {
	ServerName               string `json:"server-name"      mapstructure:"server-name"`
	ServerType               string `json:"server-type"      mapstructure:"server-type"`
	ConfigPath               string `json:"config-path"      mapstructure:"config-path"`
	VerifyExecPath           string `json:"verify-exec-path" mapstructure:"verify-exec-path"`
	LogsDirPath              string `json:"logs-dir-path"    mapstructure:"logs-dir-path"`
	BackupDir                string `json:"backup-dir"       mapstructure:"backup-dir"`
	BackupCycle              int    `json:"backup-cycle"     mapstructure:"backup-cycle"`
	BackupsRetentionDuration int    `json:"backups-retention-duration" mapstructure:"backups-retention-duration"`
}

func NewWebServerConfigOptions() *WebServerConfigOptions {
	return &WebServerConfigOptions{
		ServerType:               "nginx",
		BackupCycle:              1,
		BackupsRetentionDuration: 7,
	}
}

func (c *WebServerConfigOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.ServerName, "web-server-config.server-name", c.ServerName, ""+
		"Set the server name of the web server config. It cannot be empty.")

	fs.StringVar(&c.ServerType, "web-server-config.server-type", c.ServerType, ""+
		"Specify the server type of the web server. Currently, only `nginx` is supported.")

	fs.StringVar(&c.ConfigPath, "web-server-config.config-path", c.ConfigPath, ""+
		"Set the path of the web server configuration file."+
		" It cannot be empty and the file exists.")

	fs.StringVar(&c.VerifyExecPath, "web-server-config.verify-exec-path", c.VerifyExecPath, ""+
		"Set the path of the web server configuration verification binary file."+
		" It cannot be empty and the file exists.")

	fs.StringVar(
		&c.LogsDirPath,
		"web-server-config.logs-dir-path",
		filepath.Join(filepath.Dir(filepath.Dir(c.VerifyExecPath)), "logs"),
		""+
			"Set the path of the web server logs dir path.",
	)

	fs.StringVar(&c.BackupDir, "web-server-config.backup-dir", c.BackupDir, ""+
		"Set the special path of the web server configuration file backup directory."+
		" Set empty path, it will backup at the directory of the web server config.")

	fs.IntVar(&c.BackupCycle, "web-server-config.backup-cycle", c.BackupCycle, ""+
		"Set the web server configuration backup cycle. The unit is daily."+
		" Set zero to disable backup.")

	fs.IntVar(&c.BackupsRetentionDuration, "web-server-config.backups-retention-duration", c.BackupsRetentionDuration, ""+
		"Set the retention duration of the web server configuration backup files."+
		" The unit is daily."+
		" Set zero to disable backup.")
}

func (c *WebServerConfigOptions) Validate() []error {
	var errs []error

	// validate server-name
	if len(strings.TrimSpace(c.ServerName)) == 0 {
		errs = append(errs, errors.New("--web-server-config.server-name cannot be empty."))
	}

	// validate server-type
	if c.ServerType != "nginx" {
		errs = append(errs, errors.Errorf("--web-server-config.server-type %s can only be `nginx`", c.ServerType))
	}

	// validate config-path
	if len(strings.TrimSpace(c.ConfigPath)) != 0 { //nolint:nestif
		conff, err := os.Stat(c.ConfigPath)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "--web-server-config.config-path %s check failed.", c.ConfigPath))
		} else {
			if conff.IsDir() {
				errs = append(errs, errors.Errorf("--web-server-config.config-path %s cannot be a directory.", c.ConfigPath))
			}
		}
	} else {
		errs = append(errs, errors.New("--web-server-config.config-path cannot be empty."))
	}

	// validate verify-exec-path
	if len(strings.TrimSpace(c.VerifyExecPath)) != 0 { //nolint:nestif
		vexecf, err := os.Stat(c.VerifyExecPath)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "--web-server-config.verify-exec-path %s check failed.", c.VerifyExecPath))
		} else {
			if vexecf.IsDir() {
				errs = append(errs, errors.Errorf("--web-server-config.verify-exec-path %s cannot be a directory.", c.VerifyExecPath))
			}
		}
	} else {
		errs = append(errs, errors.New("--web-server-config.verify-exec-path cannot be empty."))
	}

	// validate backup-dir
	if len(strings.TrimSpace(c.BackupDir)) > 0 {
		dirf, err := os.Stat(c.BackupDir)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "--web-server-config.backup-dir %s check failed.", c.BackupDir))
		} else {
			if !dirf.IsDir() {
				errs = append(errs, errors.Errorf("--web-server-config.backup-dir %s can only be a directory.", c.BackupDir))
			}
		}
	}

	return errs
}

func (c *WebServerConfigOptions) ApplyToNginx(config *nginx.Config) {
	config.ManagersConfig[c.ServerName] = &configuration.ManagerConfig{
		NginxMainConfigAbsPath: c.ConfigPath,
		//  NginxHome:               c.ServerHome,
		NginxBinFilePath:    c.VerifyExecPath,
		BackupCycleDays:     c.BackupCycle,
		BackupRetentionDays: c.BackupsRetentionDuration,
		BackupDir:           c.BackupDir,
		//  BackupPrefix:            c.BackupPrefix,
	}
}

type WebServerConfigsOptions struct {
	DomainNameServerIPv4 string                    `json:"dns-ipv4" mapstructure:"dns-ipv4"`
	WebServerConfigs     []*WebServerConfigOptions `json:"items" mapstructure:"items"`
}

func NewWebServerConfigsOptions() *WebServerConfigsOptions {
	return &WebServerConfigsOptions{WebServerConfigs: make([]*WebServerConfigOptions, 0)}
}

func (cs *WebServerConfigsOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cs.DomainNameServerIPv4, "web-server-configs.dns-ipv4", cs.DomainNameServerIPv4, ""+
		"Set Domain Name resolver for resolving the proxy service domain name")
}

func (cs *WebServerConfigsOptions) Validate() []error {
	var errs []error

	if net.ParseIP(cs.DomainNameServerIPv4).String() == "<nil>" {
		errs = append(errs, errors.Errorf("--web-server-configs.dns-ipv4 %s is not valid", cs.DomainNameServerIPv4))
	}

	if len(cs.WebServerConfigs) == 0 {
		errs = append(errs, errors.New("web server configs options is null."))
	}

	for i, c := range cs.WebServerConfigs {
		suberrs := c.Validate()
		if len(suberrs) > 0 {
			aErr := errors.NewAggregate(suberrs)
			logV1.Errorf("failed to validate the %dst web server config options, cased by %v", i+1, aErr)
			errs = append(errs, errors.Wrapf(aErr, "failed to validate the %dst web server config options.", i+1))
		}
	}

	return errs
}

func (cs *WebServerConfigsOptions) ApplyGenericOptsTo(config *nginx.Config) {
	config.DomainNameServerIPv4 = cs.DomainNameServerIPv4
}
