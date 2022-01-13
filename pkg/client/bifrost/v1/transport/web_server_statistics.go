package transport

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

const (
	webServerStatisticsService = "bifrostpb.WebServerStatistics"
)

type WebServerStatisticsTransport interface {
	Get() *grpctransport.Client
}

type webServerStatisticsTransport struct {
	getClient *grpctransport.Client
}

func (w *webServerStatisticsTransport) Get() *grpctransport.Client {
	return w.getClient
}

func newWebServerStatisticsTransport(transport *transport) WebServerStatisticsTransport {
	return &webServerStatisticsTransport{
		getClient: grpctransport.NewClient(
			transport.conn,
			webServerStatisticsService,
			"Get",
			transport.encoderFactory.WebServerStatistics().EncodeRequest,
			transport.decoderFactory.WebServerStatistics().DecodeResponse,
			new(pbv1.Statistics),
		),
	}
}
