package bifrost

import (
	"github.com/ClessLi/bifrost/internal/bifrost/config"
	"github.com/ClessLi/bifrost/internal/pkg/logger"
)

func createLogger(cfg *config.Config) (*logger.Logger, error) {
	c := logger.NewConfig()

	err := cfg.Log.ApplyTo(c)
	if err != nil {
		return nil, err
	}

	return c.Complete().NewLogger()
}
