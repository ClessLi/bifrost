package service

import (
	"bytes"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service/web_server_manager"
)

type Responder interface {
	Bytes() ([]byte, error)
	GetWatcher() (web_server_manager.Watcher, error)
}

type queryOrUpdateResponse struct {
	response *bytes.Buffer
	err      error
}

func (qr queryOrUpdateResponse) Bytes() ([]byte, error) {
	return qr.response.Bytes(), qr.err
}

func (qr queryOrUpdateResponse) GetWatcher() (web_server_manager.Watcher, error) {
	return nil, ErrNotWatcherResponse
}

func NewQueryOrUpdateResponse(response []byte, err error) Responder {
	return &queryOrUpdateResponse{
		response: bytes.NewBuffer(response),
		err:      err,
	}
}

type watcherResponse struct {
	response web_server_manager.Watcher
}

func (wr watcherResponse) Bytes() ([]byte, error) {
	return nil, ErrNotBytesResponse
}

func (wr watcherResponse) GetWatcher() (web_server_manager.Watcher, error) {
	return wr.response, nil
}

func NewWatcherResponse(response web_server_manager.Watcher) Responder {
	return &watcherResponse{response: response}
}
