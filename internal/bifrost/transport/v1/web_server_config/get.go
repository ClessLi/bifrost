package web_server_config

import (
	"context"

	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/utils"
)

func (w *webServerConfigServer) GetServerNames(ctx context.Context, null *pbv1.Null) (*pbv1.ServerNames, error) {
	_, resp, err := w.handler.HandlerGetServerNames().ServeGRPC(ctx, null)
	if err != nil {
		return nil, err
	}
	response := resp.(*pbv1.ServerNames)

	return response, nil
}

func (w *webServerConfigServer) Get(r *pbv1.ServerName, stream pbv1.WebServerConfig_GetServer) error {
	_, resp, err := w.handler.HandlerGet().ServeGRPC(stream.Context(), r)
	if err != nil {
		return err
	}

	response := resp.(*pbv1.ServerConfig)
	// send config data
	err = utils.StreamSendMsg(stream, response.GetJsonData(), w.options.ChunkSize, func(msg []byte) interface{} {
		return &pbv1.ServerConfig{JsonData: msg}
	})
	if err != nil {
		return err
	}
	// send config fingerprints
	err = utils.StreamSendMsg(stream, response.GetOriginalFingerprints(), w.options.ChunkSize, func(msg []byte) interface{} {
		return &pbv1.ServerConfig{OriginalFingerprints: msg}
	})

	return stream.Send(&pbv1.ServerConfig{ServerName: response.GetServerName()})
	// return nil
}

func (w *webServerConfigServer) ConnectivityCheckOfProxiedServers(ctx context.Context, pos *pbv1.ServerConfigContextPos) (*pbv1.ContextData, error) {
	_, resp, err := w.handler.HandlerConnectivityCheckOfProxiedServers().ServeGRPC(ctx, pos)
	if err != nil {
		return nil, err
	}
	response := resp.(*pbv1.ContextData)

	return response, nil
}
