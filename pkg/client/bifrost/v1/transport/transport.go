package transport

import (
	"github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport/decoder"
	"github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport/encoder"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
	"sync"
)

type Client interface {
	Endpoint() endpoint.Endpoint
}

var _ Client = &client{}

type client struct {
	ep endpoint.Endpoint
}

func (c *client) Endpoint() endpoint.Endpoint {
	return c.ep
}

func newClient(ep endpoint.Endpoint) Client {
	return &client{ep: ep}
}

type Factory interface {
	WebServerConfig() WebServerConfigTransport
	WebServerStatistics() WebServerStatisticsTransport
	WebServerStatus() WebServerStatusTransport
	WebServerLogWatcher() WebServerLogWatcherTransport
}

var _ Factory = &transport{}

var (
	onceWebServerConfig     = sync.Once{}
	onceWebServerStatistics = sync.Once{}
	onceWebServerStatus     = sync.Once{}
	onceWebServerLogWatcher = sync.Once{}
	singletonWSCTXP         WebServerConfigTransport
	singletonWSSTXP         WebServerStatisticsTransport
	singletonWSStatusTXP    WebServerStatusTransport
	singletonWSLWTXP        WebServerLogWatcherTransport
)

type transport struct {
	conn           *grpc.ClientConn
	decoderFactory decoder.Factory
	encoderFactory encoder.Factory
}

func (t *transport) WebServerConfig() WebServerConfigTransport {
	onceWebServerConfig.Do(func() {
		if singletonWSCTXP == nil {
			singletonWSCTXP = newWebServerConfigTransport(t)
		}
	})
	if singletonWSCTXP == nil {
		log.Fatal("web server config transport client is nil")

		return nil
	}
	return singletonWSCTXP
}

func (t *transport) WebServerStatistics() WebServerStatisticsTransport {
	onceWebServerStatistics.Do(func() {
		if singletonWSSTXP == nil {
			singletonWSSTXP = newWebServerStatisticsTransport(t)
		}
	})
	if singletonWSSTXP == nil {
		log.Fatal("web server statistics transport client is nil")

		return nil
	}
	return singletonWSSTXP
}

func (t *transport) WebServerStatus() WebServerStatusTransport {
	onceWebServerStatus.Do(func() {
		if singletonWSStatusTXP == nil {
			singletonWSStatusTXP = newWebServerStatusTransport(t)
		}
	})
	if singletonWSStatusTXP == nil {
		log.Fatal("web server status transport client is nil")

		return nil
	}
	return singletonWSStatusTXP
}

func (t *transport) WebServerLogWatcher() WebServerLogWatcherTransport {
	onceWebServerLogWatcher.Do(func() {
		if singletonWSLWTXP == nil {
			singletonWSLWTXP = newWebServerLogWatcherTransport(t)
		}
	})
	if singletonWSLWTXP == nil {
		log.Fatal("web server log watcher transport client is nil")

		return nil
	}
	return singletonWSLWTXP
}

func New(conn *grpc.ClientConn) Factory {
	return &transport{
		conn:           conn,
		decoderFactory: decoder.New(),
		encoderFactory: encoder.New(),
	}
}
