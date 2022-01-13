package service

import (
	"context"
	"github.com/ClessLi/bifrost/pkg/client/bifrost/v1/endpoint"
)

var ctxIns context.Context

type Factory interface {
	WebServerConfig() WebServerConfigService
	WebServerStatistics() WebServerStatisticsService
}

type factory struct {
	eps endpoint.Factory
}

func (f *factory) WebServerConfig() WebServerConfigService {
	return newWebServerConfigService(f)
}

func (f *factory) WebServerStatistics() WebServerStatisticsService {
	return newWebServerStatisticsService(f)
}

func New(endpoint endpoint.Factory) Factory {
	return &factory{eps: endpoint}
}

func SetContext(ctx context.Context) {
	ctxIns = ctx
}

func GetContext() context.Context {
	if ctxIns == nil {
		return context.Background()
	}
	return ctxIns
}
