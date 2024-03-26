package configuration

import (
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ManagerConfig struct {
	NginxMainConfigAbsPath  string
	NginxHome               string
	NginxBinFilePath        string
	RegularlyTaskCycleDelay time.Duration
	BackupCycleDays         int
	BackupRetentionDays     int
	BackupDir               string
	BackupPrefix            string
	BackupTimeZone          *time.Location
}

func (c *ManagerConfig) Complete() (*CompletedManagerConfig, error) {
	// completing check nginx main config file path
	absPath, err := filepath.Abs(c.NginxMainConfigAbsPath)
	if err != nil {
		return nil, err
	}

	// completing nginx home and nginx bin file path
	if len(strings.TrimSpace(c.NginxHome)) == 0 && len(strings.TrimSpace(c.NginxBinFilePath)) == 0 {
		c.NginxHome = filepath.Dir(filepath.Dir(absPath))
		c.NginxBinFilePath = filepath.Join(c.NginxHome, "sbin", "nginx")
	} else if len(strings.TrimSpace(c.NginxHome)) > 0 {
		c.NginxBinFilePath = filepath.Join(c.NginxHome, "sbin", "nginx")
	} else if len(strings.TrimSpace(c.NginxBinFilePath)) > 0 {
		c.NginxHome = filepath.Dir(filepath.Dir(c.NginxBinFilePath))
	}

	// completing regularly task cycle delay time
	if c.RegularlyTaskCycleDelay < time.Second*10 {
		c.RegularlyTaskCycleDelay = time.Second * 10
	}

	// completing nginx config backup options
	if c.BackupCycleDays <= 0 {
		c.BackupCycleDays = 1
	}

	if c.BackupRetentionDays <= 0 {
		c.BackupRetentionDays = 7
	}

	if len(strings.TrimSpace(c.BackupPrefix)) == 0 {
		c.BackupPrefix = "nginx.conf"
	}

	if c.BackupTimeZone == nil {
		c.BackupTimeZone = time.Local
	}

	return &CompletedManagerConfig{c}, err
}

type CompletedManagerConfig struct {
	*ManagerConfig
}

func (cc *CompletedManagerConfig) NewNginxConfigManager() (NginxConfigManager, error) {
	c, err := NewNginxConfigFromFS(cc.NginxMainConfigAbsPath)
	if err != nil {
		return nil, err
	}
	return &nginxConfigManager{
		configuration:           c,
		nginxHome:               cc.NginxHome,
		nginxBinFilePath:        cc.NginxBinFilePath,
		regularlyTaskCycleDelay: cc.RegularlyTaskCycleDelay,
		regularlyTaskSignalChan: make(chan int),
		backupOpts: backupOption{
			backupCycleDays:     cc.BackupCycleDays,
			backupRetentionDays: cc.BackupRetentionDays,
			backupDir:           cc.BackupDir,
			backupPrefix:        cc.BackupPrefix,
			backupTimeZone:      cc.BackupTimeZone,
		},
		wg:        new(sync.WaitGroup),
		isRunning: false,
	}, nil
}
