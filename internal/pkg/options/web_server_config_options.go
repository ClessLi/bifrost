package options

import (
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"github.com/marmotedu/errors"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

type WebServerConfigOptions struct {
	ServerName     string `json:"server-name" mapstructure:"server-name"`
	ServerType     string `json:"server-type" mapstructure:"server-type"`
	ConfigPath     string `json:"config-path" mapstructure:"config-path"`
	VerifyExecPath string `json:"verify-exec-path" mapstructure:"verify-exec-path"`
	BackupDir      string `json:"backup-dir" mapstructure:"backup-dir"`
	BackupCycle    int    `json:"backup-cycle" mapstructure:"backup-cycle"`
	BackupSaveTime int    `json:"backup-save-time" mapstructure:"backup-save-time"`
}

func NewWebServerConfigOptions() *WebServerConfigOptions {
	return &WebServerConfigOptions{
		ServerType:     "nginx",
		BackupCycle:    1,
		BackupSaveTime: 7,
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

	fs.StringVar(&c.BackupDir, "web-server-config.backup-dir", c.BackupDir, ""+
		"Set the special path of the web server configuration file backup directory."+
		" Set empty path, it will backup at the directory of the web server config.")

	fs.IntVar(&c.BackupCycle, "web-server-config.backup-cycle", c.BackupCycle, ""+
		"Set the web server configuration backup cycle. The unit is daily."+
		" Set zero to disable backup.")

	fs.IntVar(&c.BackupSaveTime, "web-server-config.backup-save-time", c.BackupSaveTime, ""+
		"Set the save time of the web server configuration backup file."+
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
	if len(strings.TrimSpace(c.ConfigPath)) == 0 {
		errs = append(errs, errors.New("--web-server-config.config-path cannot be empty."))
	} else {
		conff, err := os.Stat(c.ConfigPath)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "--web-server-config.config-path %s check failed.", c.ConfigPath))
		} else {
			if conff.IsDir() {
				errs = append(errs, errors.Errorf("--web-server-config.config-path %s cannot be a directory.", c.ConfigPath))
			}
		}
	}

	// validate verify-exec-path
	if len(strings.TrimSpace(c.VerifyExecPath)) == 0 {
		errs = append(errs, errors.New("--web-server-config.verify-exec-path cannot be empty."))
	} else {
		vexecf, err := os.Stat(c.VerifyExecPath)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "--web-server-config.verify-exec-path %s check failed.", c.VerifyExecPath))
		} else {
			if vexecf.IsDir() {
				errs = append(errs, errors.Errorf("--web-server-config.verify-exec-path %s cannot be a directory.", c.VerifyExecPath))
			}
		}
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

type WebServerConfigsOptions struct {
	WebServerConfigs []*WebServerConfigOptions `json:"items" mapstructure:"items"`
}

func NewWebServerConfigsOptions() *WebServerConfigsOptions {
	return &WebServerConfigsOptions{WebServerConfigs: make([]*WebServerConfigOptions, 0)}
}

func (cs *WebServerConfigsOptions) Validate() []error {
	var errs []error

	if len(cs.WebServerConfigs) == 0 {
		errs = append(errs, errors.New("web server configs options is null."))
	}

	for i, c := range cs.WebServerConfigs {
		suberrs := c.Validate()
		if len(suberrs) > 0 {
			aErr := errors.NewAggregate(suberrs)
			log.Errorf("failed to validate the %dst web server config options, cased by %v", i+1, aErr)
			errs = append(errs, errors.Wrapf(aErr, "failed to validate the %dst web server config options.", i+1))
		}
	}

	return errs
}
