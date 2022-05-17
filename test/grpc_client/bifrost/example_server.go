package bifrost

import (
	"fmt"
	"time"

	"github.com/yongPhone/bifrost/internal/bifrost"
	"github.com/yongPhone/bifrost/internal/bifrost/config"
	bifrost_opts "github.com/yongPhone/bifrost/internal/bifrost/options"
	"github.com/yongPhone/bifrost/internal/pkg/options"
	log "github.com/yongPhone/bifrost/pkg/log/v1"
)

const (
	localhostIP = "192.168.220.1"
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

	opts.WebServerConfigsOptions.WebServerConfigs = make([]*options.WebServerConfigOptions, 0)
	opts.WebServerConfigsOptions.WebServerConfigs = append(
		opts.WebServerConfigsOptions.WebServerConfigs,
		&options.WebServerConfigOptions{
			ServerName:     "example test",
			ServerType:     "nginx",
			ConfigPath:     "../../nginx/conf/nginx.conf",
			VerifyExecPath: "../../nginx/sbin/nginx.sh",
			LogsDirPath:    "../../nginx/logs",
			BackupDir:      "",
			BackupCycle:    1,
			BackupSaveTime: 1,
		},
	)

	opts.WebServerLogWatcherOptions.WatchTimeout = time.Second * 50
	opts.WebServerLogWatcherOptions.MaxConnections = 1024

	opts.Log.Level = "debug"
	log.Init(opts.Log)
	return bifrost.Run(&config.Config{Options: opts})
}

func serverAddress() string {
	return fmt.Sprintf("%s:%d", localhostIP, testPort)
}
