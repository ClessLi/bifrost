package file_watcher

import (
	"context"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"os"
	"path/filepath"
	"time"

	"github.com/marmotedu/errors"
)

type Config struct {
	MaxConnections int
	OutputTimeout  time.Duration
}

func NewConfig() *Config {
	return &Config{
		MaxConnections: 1000,
		OutputTimeout:  time.Minute * 5,
	}
}

func (c *Config) Complete(path string) (*CompletedConfig, error) {
	abspath, err := absFilePath(path)
	if err != nil {
		return &CompletedConfig{}, err
	}

	return &CompletedConfig{
		filePath: abspath,
		Config:   c,
	}, nil
}

func absFilePath(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", errors.Errorf("%s is a directory", path)
	}

	return filepath.Abs(path)
}

type CompletedConfig struct {
	filePath string
	*Config
}

func (cc *CompletedConfig) NewFileWatcher(firstOutputCtx context.Context) (*FileWatcher, <-chan []byte, error) {
	watcher, err := newFileWatcher(context.Background(), cc)
	if err != nil {
		return nil, nil, err
	}
	output, err := watcher.Output(firstOutputCtx)
	if err != nil {
		return nil, nil, err
	}
	go func() {
		err := watcher.Start()
		if err != nil {
			logV1.Warnf(err.Error())
		}
	}()

	return watcher, output, nil
}
