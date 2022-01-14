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
}

var _ Factory = &transport{}

var (
	onceWebServerConfig     = sync.Once{}
	onceWebServerStatistics = sync.Once{}
	singletonWSCTXP         WebServerConfigTransport
	singletonWSSTXP         WebServerStatisticsTransport
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

func New(conn *grpc.ClientConn) Factory {
	return &transport{
		conn:           conn,
		decoderFactory: decoder.New(),
		encoderFactory: encoder.New(),
	}
}
