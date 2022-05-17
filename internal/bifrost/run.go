package bifrost

import "github.com/yongPhone/bifrost/internal/bifrost/config"

func Run(cfg *config.Config) error {
	server, err := createBifrostServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
