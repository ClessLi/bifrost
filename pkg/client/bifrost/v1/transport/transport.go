package transport

import (
	"sync"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"

	"github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport/decoder"
	"github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport/encoder"
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

type transport struct {
	conn           *grpc.ClientConn
	decoderFactory decoder.Factory
	encoderFactory encoder.Factory

	onceWebServerConfig     sync.Once
	onceWebServerStatistics sync.Once
	onceWebServerStatus     sync.Once
	onceWebServerLogWatcher sync.Once
	singletonWSCTXP         WebServerConfigTransport
	singletonWSSTXP         WebServerStatisticsTransport
	singletonWSStatusTXP    WebServerStatusTransport
	singletonWSLWTXP        WebServerLogWatcherTransport
}

func (t *transport) WebServerConfig() WebServerConfigTransport {
	t.onceWebServerConfig.Do(func() {
		if t.singletonWSCTXP == nil {
			t.singletonWSCTXP = newWebServerConfigTransport(t)
		}
	})
	if t.singletonWSCTXP == nil {
		logV1.Fatal("web server config transport client is nil")

		return nil
	}

	return t.singletonWSCTXP
}

func (t *transport) WebServerStatistics() WebServerStatisticsTransport {
	t.onceWebServerStatistics.Do(func() {
		if t.singletonWSSTXP == nil {
			t.singletonWSSTXP = newWebServerStatisticsTransport(t)
		}
	})
	if t.singletonWSSTXP == nil {
		logV1.Fatal("web server statistics transport client is nil")

		return nil
	}

	return t.singletonWSSTXP
}

func (t *transport) WebServerStatus() WebServerStatusTransport {
	t.onceWebServerStatus.Do(func() {
		if t.singletonWSStatusTXP == nil {
			t.singletonWSStatusTXP = newWebServerStatusTransport(t)
		}
	})
	if t.singletonWSStatusTXP == nil {
		logV1.Fatal("web server status transport client is nil")

		return nil
	}

	return t.singletonWSStatusTXP
}

func (t *transport) WebServerLogWatcher() WebServerLogWatcherTransport {
	t.onceWebServerLogWatcher.Do(func() {
		if t.singletonWSLWTXP == nil {
			t.singletonWSLWTXP = newWebServerLogWatcherTransport(t)
		}
	})
	if t.singletonWSLWTXP == nil {
		logV1.Fatal("web server log watcher transport client is nil")

		return nil
	}

	return t.singletonWSLWTXP
}

func New(conn *grpc.ClientConn) Factory {
	return &transport{
		conn:                    conn,
		decoderFactory:          decoder.New(),
		encoderFactory:          encoder.New(),
		onceWebServerConfig:     sync.Once{},
		onceWebServerStatistics: sync.Once{},
		onceWebServerStatus:     sync.Once{},
		onceWebServerLogWatcher: sync.Once{},
	}
}
