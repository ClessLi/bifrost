package config

import (
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/apsdehal/go-logger"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

// Config, bifrost配置文件结构体，定义bifrost配置信息
type Config struct {
	ServiceConfig ServiceConfig `yaml:"Service"`
	*RAConfig     `yaml:"RAConfig"`
	LogConfig     `yaml:"LogConfig"`
}

func NewBifrostConfig(configFilePath string) *Config {

	// 初始化bifrost配置
	confData, err := utils.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	BifrostConf := new(Config)
	// 加载bifrost配置
	err = yaml.Unmarshal(confData, BifrostConf)
	if err != nil {
		panic(err)
	}

	// 配置必填项检查
	err = BifrostConf.check()
	if err != nil {
		panic(err)
	}

	// 初始化日志
	logDir, err := filepath.Abs(BifrostConf.LogDir)
	if err != nil {
		panic(err)
	}

	logPath := filepath.Join(logDir, "bifrost.log")
	utils.Logf, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	utils.InitLogger(utils.Logf, BifrostConf.LogConfig.Level)

	// 初始化应用运行日志输出
	stdoutPath := filepath.Join(logDir, "bifrost.out")
	utils.Stdoutf, err = os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	os.Stdout = utils.Stdoutf
	os.Stderr = utils.Stdoutf

	return BifrostConf
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
	if c.ServiceConfig.ChunkSize == 0 {
		c.ServiceConfig.ChunkSize = 4194304
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
	ChunkSize            int                   `yaml:"ChunkSize"`
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
