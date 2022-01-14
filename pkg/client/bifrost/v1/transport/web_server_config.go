package transport

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

const (
	webServerConfigService = "bifrostpb.WebServerConfig"
)

type WebServerConfigTransport interface {
	GetServerNames() *grpctransport.Client
	Get() *grpctransport.Client
	Update() *grpctransport.Client
}

type webServerConfigTransport struct {
	getServerNamesClient *grpctransport.Client
	getClient            *grpctransport.Client
	updateClient         *grpctransport.Client
}

func (w *webServerConfigTransport) GetServerNames() *grpctransport.Client {
	return w.getServerNamesClient
}

func (w *webServerConfigTransport) Get() *grpctransport.Client {
	return w.getClient
}

func (w *webServerConfigTransport) Update() *grpctransport.Client {
	return w.updateClient
}

func newWebServerConfigTransport(transport *transport) WebServerConfigTransport {
	return &webServerConfigTransport{
		getServerNamesClient: grpctransport.NewClient(
			transport.conn,
			webServerConfigService,
			"GetServerNames",
			transport.encoderFactory.WebServerConfig().EncodeRequest,
			transport.decoderFactory.WebServerConfig().DecodeResponse,
			new(pbv1.ServerNames),
		),
		getClient: grpctransport.NewClient(
			transport.conn,
			webServerConfigService,
			"Get",
			transport.encoderFactory.WebServerConfig().EncodeRequest,
			transport.decoderFactory.WebServerConfig().DecodeResponse,
			new(pbv1.ServerConfig),
		),
		updateClient: grpctransport.NewClient(
			transport.conn,
			webServerConfigService,
			"Update",
			transport.encoderFactory.WebServerConfig().EncodeRequest,
			transport.decoderFactory.WebServerConfig().DecodeResponse,
			new(pbv1.Response),
		),
	}
}
