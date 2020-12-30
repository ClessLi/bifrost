package service

import (
	"bytes"
	"golang.org/x/net/context"
)

type Requester interface {
	GetContext() context.Context
	GetRequestType() RequestType
	GetToken() string
	GetServerName() string
	GetParam() string
	GetRequestData() (data []byte)
}

type request struct {
	ctx         context.Context
	requestType RequestType
	token       string
	serverName  string
	param       string
	data        *bytes.Buffer
}

func NewRequest(ctx context.Context, reqType, token, svrName, param string, data []byte) Requester {
	requestType := Unknown
	switch reqType {
	case "DisplayConfig":
		requestType = DisplayConfig
	case "GetConfig":
		requestType = GetConfig
	case "UpdateConfig":
		requestType = UpdateConfig
	case "ShowStatistics":
		requestType = ShowStatistics
	case "DisplayStatus":
		requestType = DisplayStatus
	case "WatchLog":
		requestType = WatchLog
	}
	dataR := bytes.NewBuffer(data)
	return &request{
		ctx:         ctx,
		requestType: requestType,
		token:       token,
		serverName:  svrName,
		param:       param,
		data:        dataR,
	}
}

func (r request) GetContext() context.Context {
	return r.ctx
}

func (r request) GetRequestType() RequestType {
	return r.requestType
}

func (r request) GetToken() string {
	return r.token
}

func (r request) GetServerName() string {
	return r.serverName
}

func (r request) GetParam() string {
	return r.param
}

func (r request) GetRequestData() (data []byte) {
	return r.data.Bytes()
}
