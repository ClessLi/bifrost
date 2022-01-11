package web_server_config

import (
	"bytes"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/marmotedu/errors"
	"io"
	"time"
)

func (w *webServerConfigServer) Update(stream pbv1.WebServerConfig_UpdateServer) error {
	buffer := bytes.NewBuffer(make([]byte, 0, w.options.ChunkSize))
	defer buffer.Reset()
	var (
		in            *pbv1.ServerConfig
		err           error
		req           *pbv1.ServerConfig
		isTimeout     = false
		recvStartTime = time.Now()
	)

	for !isTimeout {
		isTimeout = time.Since(recvStartTime) >= w.options.RecvTimeoutMinutes*time.Minute
		in, err = stream.Recv()
		// 判断是否传入完毕
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return err
		}

		if req == nil {
			req = &pbv1.ServerConfig{
				ServerName: in.GetServerName(),
			}
		} else {
			if req.ServerName != in.GetServerName() {
				return errors.WithCode(code.ErrValidation, "need server name: '%s', not '%s'", req.ServerName, in.GetServerName())
			}
		}

		buffer.Write(in.GetJsonData())
	}

	if isTimeout {
		return errors.WithCode(code.ErrRequestTimeout, "receive timeout during data wrap")
	}

	if req == nil {
		return errors.WithCode(code.ErrValidation, "update web server config is nil")
	}

	req.JsonData = buffer.Bytes()
	_, resp, err := w.handler.HandlerUpdate().ServeGRPC(stream.Context(), req)
	if err != nil {
		return errors.Wrapf(err, "failed to handle the update operation of the web server config(json-data) - %s", string(req.GetJsonData()))
	}
	return stream.SendAndClose(resp.(*pbv1.Response))
}
