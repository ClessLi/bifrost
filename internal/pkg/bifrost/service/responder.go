package service

import (
	"bytes"
)

type ViewResponder interface {
	serverNameQuerier
	dataResponder
	errorResponder
}

type viewResponder struct {
	serverName string
	data       *bytes.Buffer
	err        error
}

func NewViewResponder(serverName string, data []byte, err error) ViewResponder {
	return &viewResponder{
		serverName: serverName,
		data:       bytes.NewBuffer(data),
		err:        err,
	}
}

func (v viewResponder) GetServerName() string {
	return v.serverName
}

func (v viewResponder) Bytes() []byte {
	return v.data.Bytes()
}

func (v viewResponder) Error() error {
	return v.err
}

type UpdateResponder interface {
	serverNameQuerier
	errorResponder
}

type updateResponder struct {
	serverName string
	err        error
}

func NewUpdateResponder(serverName string, err error) UpdateResponder {
	return &updateResponder{
		serverName: serverName,
		err:        err,
	}
}

func (u updateResponder) GetServerName() string {
	return u.serverName
}

func (u updateResponder) Error() error {
	return u.err
}

type WatchResponder interface {
	serverNameQuerier
	BytesChan() <-chan []byte
	TransferErrorChan() <-chan error
	errorResponder
	Close() error
}

type watchResponder struct {
	serverName      string
	dataChan        <-chan []byte
	transferErrChan <-chan error
	closeFunc       func() error
	err             error
}

func NewWatchResponder(serverName string, closeFunc func() error, dataChan <-chan []byte, transferErrChan <-chan error, err error) WatchResponder {
	return &watchResponder{
		serverName:      serverName,
		dataChan:        dataChan,
		transferErrChan: transferErrChan,
		closeFunc:       closeFunc,
		err:             err,
	}
}

func (w watchResponder) GetServerName() string {
	return w.serverName
}

func (w watchResponder) BytesChan() <-chan []byte {
	return w.dataChan
}

func (w watchResponder) TransferErrorChan() <-chan error {
	return w.transferErrChan
}

func (w watchResponder) Close() error {
	return w.closeFunc()
}

func (w watchResponder) Error() error {
	return w.err
}

type dataResponder interface {
	Bytes() []byte
}

type errorResponder interface {
	Error() error
}
