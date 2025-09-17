package transport //nolint:dupl

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

const (
	webServerBinCMDService = "bifrostpb.WebServerBinCMD"
)

type WebServerBinCMDTransport interface {
	Exec() Client
}

type webServerBinCMDTransport struct {
	execClient Client
}

func (w *webServerBinCMDTransport) Exec() Client {
	return w.execClient
}

func newWebServerBinCMDTransport(transport *transport) WebServerBinCMDTransport {
	return &webServerBinCMDTransport{
		execClient: grpctransport.NewClient(
			transport.conn,
			webServerBinCMDService,
			"Exec",
			transport.encoderFactory.WebServerBinCMD().EncodeRequest,
			transport.decoderFactory.WebServerBinCMD().DecodeResponse,
			new(pbv1.ExecuteResponse),
		),
	}
}
