package bifrost

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/ClessLi/bifrost/internal/bifrost"
	"github.com/ClessLi/bifrost/internal/bifrost/config"
	bifrost_opts "github.com/ClessLi/bifrost/internal/bifrost/options"
	"github.com/ClessLi/bifrost/internal/pkg/options"
)

const (
	localhostIP = "127.0.0.1"
	testPort    = 12321
)

func exampleServerRun() error {
	opts := bifrost_opts.NewOptions()

	opts.SecureServing = nil

	opts.InsecureServing.BindAddress = localhostIP
	opts.InsecureServing.BindPort = testPort

	opts.RAOptions = nil

	opts.GRPCServing.ChunkSize = 100

	opts.MonitorOptions.SyncInterval = time.Second * 1
	opts.MonitorOptions.CycleTime = time.Second * 5
	opts.MonitorOptions.FrequencyPerCycle = 5

	absConfigPath, err := filepath.Abs("../../nginx/conf/nginx.conf")
	if err != nil {
		return err
	}
	opts.WebServerConfigsOptions.DomainNameServerIPv4 = "114.114.114.114"
	opts.WebServerConfigsOptions.WebServerConfigs = make([]*options.WebServerConfigOptions, 0)
	opts.WebServerConfigsOptions.WebServerConfigs = append(
		opts.WebServerConfigsOptions.WebServerConfigs,
		&options.WebServerConfigOptions{
			ServerName:               "example test",
			ServerType:               "nginx",
			ConfigPath:               absConfigPath,
			VerifyExecPath:           verifyExecPath,
			LogsDirPath:              "../../nginx/logs",
			BackupDir:                "",
			BackupCycle:              1,
			BackupsRetentionDuration: 1,
		},
	)

	opts.WebServerLogWatcherOptions.WatchTimeout = time.Second * 50
	opts.WebServerLogWatcherOptions.MaxConnections = 1024

	opts.Log.InfoLevel = "debug"

	return bifrost.Run(&config.Config{Options: opts})
}

func serverAddress() string {
	return fmt.Sprintf("%s:%d", localhostIP, testPort)
}
