package logger

import (
	v1log "github.com/ClessLi/component-base/pkg/log/v1"
	"path/filepath"
)

type Config struct {
	InfoLogOpts *v1log.Options
	ErrLogOpts  *v1log.Options
}

func NewConfig() *Config {
	return &Config{
		InfoLogOpts: v1log.NewOptions(),
		ErrLogOpts:  v1log.NewOptions(),
	}
}

type CompletedConfig struct {
	*Config
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

func (c CompletedConfig) NewLogger() (*Logger, error) {
	return &Logger{
		initFunc: func() error {
			// 创建日志目录
			for _, infoOutputPath := range c.InfoLogOpts.OutputPaths {
				if infoOutputPath == "stdout" || infoOutputPath == "stderr" {
					continue
				}
				err := createLogDir(filepath.Dir(infoOutputPath))
				if err != nil {
					return err
				}
			}
			for _, errorOutputPath := range c.ErrLogOpts.OutputPaths {
				if errorOutputPath == "stdout" || errorOutputPath == "stderr" {
					continue
				}
				err := createLogDir(filepath.Dir(errorOutputPath))
				if err != nil {
					return err
				}
			}

			v1log.Init(c.InfoLogOpts, c.ErrLogOpts)
			return nil
		},
		flush: func() {
			v1log.Flush()
		},
	}, nil
}
