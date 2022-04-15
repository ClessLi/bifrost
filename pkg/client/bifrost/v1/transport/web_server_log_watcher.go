package transport

import (
	"context"
	"io"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/marmotedu/errors"
	"google.golang.org/grpc"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
)

type WebServerLogWatcherTransport interface {
	Watch() Client
}

type webServerLogWatcherTransport struct {
	watchClient Client
}

func (w *webServerLogWatcherTransport) Watch() Client {
	return w.watchClient
}

func newWebServerLogWatcherClient(
	conn *grpc.ClientConn,
	requestFunc grpctransport.EncodeRequestFunc,
	responseFunc grpctransport.DecodeResponseFunc,
) Client {
	cli := pbv1.NewWebServerLogWatcherClient(conn)

	return newClient(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, err := requestFunc(ctx, request)
		if err != nil {
			return nil, err
		}

		stream, err := cli.Watch(ctx, req.(*pbv1.LogWatchRequest))
		if err != nil {
			return nil, err
		}

		outputC := make(chan []byte)

		go func() {
			defer close(outputC)
			needClose := false
			for !needClose {
				select {
				case outputC <- recvWatcherResponse(stream, &needClose):
				case <-ctx.Done():
					return
				}
			}
		}()

		return responseFunc(ctx, &v1.WebServerLog{Lines: outputC})
	})
}

func recvWatcherResponse(stream pbv1.WebServerLogWatcher_WatchClient, needClose *bool) []byte {
	resp, err := stream.Recv()
	if err != nil && !errors.Is(err, io.EOF) {
		*needClose = true

		return []byte(err.Error())
	}
	if errors.Is(err, io.EOF) {
		*needClose = true

		return nil
	}

	return resp.GetMsg()
}

func newWebServerLogWatcherTransport(transport *transport) WebServerLogWatcherTransport {
	return &webServerLogWatcherTransport{
		watchClient: newWebServerLogWatcherClient(
			transport.conn,
			transport.encoderFactory.WebServerLogWatcher().EncodeRequest,
			transport.decoderFactory.WebServerLogWatcher().DecodeResponse,
		),
	}
}
