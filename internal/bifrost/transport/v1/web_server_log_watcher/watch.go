package web_server_log_watcher

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"io"
)

func (w *webServerLogWatcherServer) Watch(request *pbv1.LogWatchRequest, stream pbv1.WebServerLogWatcher_WatchServer) error {
	reqCtx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	respCtx, resp, err := w.handler.HandlerWatch().ServeGRPC(reqCtx, request) // resp is a *v1.WebServerLog
	if err != nil {
		return err
	}
	respWatcher := resp.(*v1.WebServerLog)

	for {
		select {
		case <-reqCtx.Done():
			return reqCtx.Err()
		case <-respCtx.Done():
			return respCtx.Err()
		case line := <-respWatcher.Lines:
			if line == nil {
				return nil
			}
			line = append(line, '\n')
			err := sendWithChunk(stream, w.options.ChunkSize-2, line)
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF {
				return nil
			}
		}
	}
}

func sendWithChunk(stream pbv1.WebServerLogWatcher_WatchServer, chunksize int, response []byte) error {
	var err error
	n := len(response)
	for i := 0; i < n; i += chunksize {
		if n <= i+chunksize {
			err = stream.Send(&pbv1.Response{Msg: response[i:]})
		} else {
			err = stream.Send(&pbv1.Response{Msg: response[i : i+chunksize]})
		}
		if err != nil {
			return err
		}
	}
	return nil
}
