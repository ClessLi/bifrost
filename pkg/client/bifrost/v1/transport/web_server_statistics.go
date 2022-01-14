package transport

import (
	"bytes"
	"context"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"io"
)

//const (
//	webServerStatisticsService = "bifrostpb.WebServerStatistics"
//)

type WebServerStatisticsTransport interface {
	Get() Client
}

type webServerStatisticsTransport struct {
	getClient Client
}

func (w *webServerStatisticsTransport) Get() Client {
	return w.getClient
}

func newWebServerStatisticsClient(conn *grpc.ClientConn, requestFunc grpctransport.EncodeRequestFunc, responseFunc grpctransport.DecodeResponseFunc) Client {
	cli := pbv1.NewWebServerStatisticsClient(conn)
	return newClient(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, err := requestFunc(ctx, request)
		if err != nil {
			return nil, err
		}

		stream, err := cli.Get(ctx, req.(*pbv1.ServerName))
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

		return responseFunc(ctx, &pbv1.Statistics{JsonData: buf.Bytes()})

	})
}

func newWebServerStatisticsTransport(transport *transport) WebServerStatisticsTransport {
	return &webServerStatisticsTransport{
		getClient: newWebServerStatisticsClient(
			transport.conn,
			transport.encoderFactory.WebServerStatistics().EncodeRequest,
			transport.decoderFactory.WebServerStatistics().DecodeResponse,
		),
	}
}
