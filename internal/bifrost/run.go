package bifrost

import "github.com/ClessLi/bifrost/internal/bifrost/config"

func Run(cfg *config.Config) error {
	l, err := createLogger(cfg)
	if err != nil {
		return err
	}
	err = l.Init()
	if err != nil {
		return err
	}
	defer l.Flush()

	server, err := createBifrostServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
