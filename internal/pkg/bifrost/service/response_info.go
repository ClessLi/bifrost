package service

import (
	"bytes"
)

type dataResponseInfo interface {
	Bytes() []byte
}

type errorResponseInfo interface {
	Error() error
}

type ViewResponseInfo interface {
	serverNameQuerier
	dataResponseInfo
	errorResponseInfo
}

type viewResponseInfo struct {
	serverName string
	dataBuffer *bytes.Buffer
	err        error
}

func NewViewResponseInfo(serverName string, data []byte, respErr error) ViewResponseInfo {
	return &viewResponseInfo{
		serverName: serverName,
		dataBuffer: bytes.NewBuffer(data),
		err:        respErr,
	}
}

func (v viewResponseInfo) GetServerName() string {
	return v.serverName
}

func (v viewResponseInfo) Bytes() []byte {
	return v.dataBuffer.Bytes()
}

func (v viewResponseInfo) Error() error {
	return v.err
}

type UpdateResponseInfo interface {
	serverNameQuerier
	errorResponseInfo
}

type updateResponseInfo struct {
	serverName string
	err        error
}

func NewUpdateResponseInfo(serverName string, respErr error) UpdateResponseInfo {
	return &updateResponseInfo{
		serverName: serverName,
		err:        respErr,
	}
}

func (u updateResponseInfo) GetServerName() string {
	return u.serverName
}

func (u updateResponseInfo) Error() error {
	return u.err
}

type WatchResponseInfo interface {
	serverNameQuerier
	BytesChan() <-chan []byte
	TransferErrorChan() <-chan error
	errorResponseInfo
	Close() error
}

type watchResponseInfo struct {
	serverName      string
	dataChan        <-chan []byte
	transferErrChan <-chan error
	closeFunc       func() error
	err             error
}

func NewWatchResponseInfo(serverName string, closeFunc func() error, dataChan <-chan []byte, transferErrChan <-chan error, respErr error) WatchResponseInfo {
	if closeFunc == nil {
		closeFunc = func() error {
			return nil
		}
	}
	return &watchResponseInfo{
		serverName:      serverName,
		dataChan:        dataChan,
		transferErrChan: transferErrChan,
		closeFunc:       closeFunc,
		err:             respErr,
	}
}

func (w watchResponseInfo) GetServerName() string {
	return w.serverName
}

func (w watchResponseInfo) BytesChan() <-chan []byte {
	return w.dataChan
}

func (w watchResponseInfo) TransferErrorChan() <-chan error {
	return w.transferErrChan
}

func (w watchResponseInfo) Close() error {
	return w.closeFunc()
}

func (w watchResponseInfo) Error() error {
	return w.err
}
