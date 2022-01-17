package bifrost

import (
	"fmt"
	"github.com/ClessLi/bifrost/internal/bifrost"
	"github.com/ClessLi/bifrost/internal/bifrost/config"
	bifrost_opts "github.com/ClessLi/bifrost/internal/bifrost/options"
	"github.com/ClessLi/bifrost/internal/pkg/options"
	"time"
)

const (
	localhostIP = "192.168.200.1"
	testPort    = 12321
)

func exampleServerRun() error {
	opts := bifrost_opts.NewOptions()

	opts.SecureServing = nil

	opts.InsecureServing.BindAddress = localhostIP
	opts.InsecureServing.BindPort = testPort

	opts.RAOptions = nil

	opts.GRPCServing.ChunkSize = 100

	opts.MonitorOptions.SyncInterval = time.Second * 10
	opts.MonitorOptions.CycleTime = time.Second * 30
	opts.MonitorOptions.FrequencyPerCycle = 10

	opts.WebServerConfigsOptions.WebServerConfigs = make([]*options.WebServerConfigOptions, 0)
	opts.WebServerConfigsOptions.WebServerConfigs = append(opts.WebServerConfigsOptions.WebServerConfigs, &options.WebServerConfigOptions{
		ServerName:     "example test",
		ServerType:     "nginx",
		ConfigPath:     "../../nginx/conf/nginx.conf",
		VerifyExecPath: "../../nginx/sbin/nginx.sh",
		BackupDir:      "",
		BackupCycle:    0,
		BackupSaveTime: 0,
	})

	return bifrost.Run(&config.Config{Options: opts})
}

func serverAddress() string {
	return fmt.Sprintf("%s:%d", localhostIP, testPort)
}
