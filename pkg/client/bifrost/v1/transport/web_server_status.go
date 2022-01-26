package transport

import (
	"bytes"
	"context"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"io"
)

type WebServerStatusTransport interface {
	Get() Client
}

type webServerStatusTransport struct {
	getClient Client
}

func (w *webServerStatusTransport) Get() Client {
	return w.getClient
}

func newWebServerStatusClient(conn *grpc.ClientConn, requestFunc grpctransport.EncodeRequestFunc, responseFunc grpctransport.DecodeResponseFunc) Client {
	cli := pbv1.NewWebServerStatusClient(conn)
	return newClient(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, err := requestFunc(ctx, request)
		if err != nil {
			return nil, err
		}

		stream, err := cli.Get(ctx, req.(*pbv1.Null))
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

			buf.Write(d.GetJsonData())
		}

		return responseFunc(ctx, &pbv1.Metrics{JsonData: buf.Bytes()})

	})
}

func newWebServerStatusTransport(transport *transport) WebServerStatusTransport {
	return &webServerStatusTransport{
		getClient: newWebServerStatusClient(
			transport.conn,
			transport.encoderFactory.WebServerStatus().EncodeRequest,
			transport.decoderFactory.WebServerStatus().DecodeResponse,
		),
	}
}
