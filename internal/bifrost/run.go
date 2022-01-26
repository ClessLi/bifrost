package bifrost

import "github.com/ClessLi/bifrost/internal/bifrost/config"

func Run(cfg *config.Config) error {
	server, err := createBifrostServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
