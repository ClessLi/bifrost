package transport

import (
	"bytes"
	"context"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/marmotedu/errors"
	"google.golang.org/grpc"
	"io"
)

const (
	webServerConfigService = "bifrostpb.WebServerConfig"
)

type WebServerConfigTransport interface {
	GetServerNames() Client
	Get() Client
	Update() Client
}

type webServerConfigTransport struct {
	getServerNamesClient Client
	getClient            Client
	updateClient         Client
}

func (w *webServerConfigTransport) GetServerNames() Client {
	return w.getServerNamesClient
}

func (w *webServerConfigTransport) Get() Client {
	return w.getClient
}

func (w *webServerConfigTransport) Update() Client {
	return w.updateClient
}

func newWebServerConfigGetClient(conn *grpc.ClientConn, requestFunc grpctransport.EncodeRequestFunc, responseFunc grpctransport.DecodeResponseFunc) Client {
	cli := pbv1.NewWebServerConfigClient(conn)
	return newClient(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, err := requestFunc(ctx, request)
		if err != nil {
			return nil, err
		}

		stream, err := cli.Get(ctx, req.(*pbv1.ServerName))
		if err != nil {
			return nil, err
		}

		buf := bytes.NewBuffer(nil)
		for {
			d, err := stream.Recv()
			if err != nil && err != io.EOF {
				return nil, err
			}
			if err == io.EOF {
				break
			}
			if d.GetServerName() != "" && d.GetServerName() != req.(*pbv1.ServerName).GetName() {
				return nil, errors.Errorf("the web server config is incorrect: got config of `%s`, want config of `%s`", d.GetServerName(), req.(*pbv1.ServerName).GetName())
			}
			buf.Write(d.GetJsonData())
		}

		return responseFunc(
			ctx,
			&pbv1.ServerConfig{
				ServerName: req.(*pbv1.ServerName).GetName(),
				JsonData:   buf.Bytes(),
			},
		)

	})
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
		getClient: newWebServerConfigGetClient(
			transport.conn,
			transport.encoderFactory.WebServerConfig().EncodeRequest,
			transport.decoderFactory.WebServerConfig().DecodeResponse,
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
