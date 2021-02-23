package bifrost

import (
	"fmt"
	"github.com/apsdehal/go-logger"
)

// Config, bifrost配置文件结构体，定义bifrost配置信息
type Config struct {
	ServiceConfig ServiceConfig `yaml:"Service"`
	*RAConfig     `yaml:"RAConfig"`
	LogConfig     `yaml:"LogConfig"`
}

func (c *Config) check() error {
	if c == nil {
		return fmt.Errorf("bifrost config load error")
	}
	if len(c.ServiceConfig.WebServerConfigInfos) == 0 {
		return fmt.Errorf("bifrost services config load error")
	}
	if c.LogDir == "" {
		return fmt.Errorf("bifrost log config load error")
	}
	// 初始化服务信息配置
	if c.ServiceConfig.Port == 0 {
		c.ServiceConfig.Port = 12321
	}
	if c.ServiceConfig.ChunckSize == 0 {
		c.ServiceConfig.ChunckSize = 4194304
	}

	if c.RAConfig != nil {
		if c.RAConfig.Host == "" || c.RAConfig.Port == 0 {
			c.RAConfig = nil
		}
	}
	return nil
}

type ServiceConfig struct {
	Port                 uint16                `yaml:"Port"`
	ChunckSize           int                   `yaml:"ChunkSize"`
	AuthServerAddr       string                `yaml:"AuthServerAddr"`
	WebServerConfigInfos []WebServerConfigInfo `yaml:"Infos,flow"`
}

// WebServerType, web服务器类型对象，定义web服务器所属类型
type WebServerType string

const (
	// Web服务类型
	NGINX WebServerType = "nginx"
	HTTPD WebServerType = "httpd"
)

type WebServerConfigInfo struct {
	Name           string        `yaml:"name"`
	Type           WebServerType `yaml:"type"`
	BackupCycle    int           `yaml:"backupCycle"`
	BackupSaveTime int           `yaml:"backupSaveTime"`
	BackupDir      string        `yaml:"backupDir,omitempty"`
	ConfPath       string        `yaml:"confPath"`
	VerifyExecPath string        `yaml:"verifyExecPath"`
}

type RAConfig struct {
	Host string `yaml:"Host"`
	Port uint16 `yaml:"Port"`
}

// LogConfig, bifrost日志信息结构体，定义日志目录、日志级别
type LogConfig struct {
	LogDir string          `yaml:"logDir"`
	Level  logger.LogLevel `yaml:"level"`
}

func (l LogConfig) IsDebugLvl() bool {
	return l.Level >= logger.DebugLevel
}
