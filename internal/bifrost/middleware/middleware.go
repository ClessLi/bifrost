package middleware

import (
	"github.com/ClessLi/bifrost/internal/bifrost/middleware/logging"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
)

type Middleware func(svcv1.ServiceFactory) svcv1.ServiceFactory

type Middlewares map[string]Middleware

var defaultMiddlewares = Middlewares{
	"logging": func(svc svcv1.ServiceFactory) svcv1.ServiceFactory {
		return logging.New(svc)
	},
}

func GetMiddlewares() Middlewares {
	return defaultMiddlewares
}
