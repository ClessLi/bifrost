package bifrost

import (
	"github.com/yongPhone/bifrost/internal/bifrost/config"
	"github.com/yongPhone/bifrost/internal/bifrost/options"
	"github.com/yongPhone/bifrost/pkg/app"
	log "github.com/yongPhone/bifrost/pkg/log/v1"
)

const commandDesc = `The Bifrost is used to parse the nginx configuration file 
and provide an interface for displaying and modifying the configuration file.
It supports the mutual conversion of JSON, string format and golang structure.
The Bifrost services to do the api objects management with gRPC protocol.

Find more Bifrost information at:
    https://github.com/ClessLi/bifrost/blob/master/docs/guide/en-US/cmd/bifrost.md`

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("Bifrost",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		// init auth api client
		// auth.Init(opts.AuthAPIClient)

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
