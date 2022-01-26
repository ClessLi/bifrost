package service

import (
	"bytes"
	"golang.org/x/net/context"
)

type contextQuerier interface {
	Context() context.Context
}

type tokenQuerier interface {
	GetToken() string
}

type requestTypeQuerier interface {
	GetRequestType() RequestType
}

type serverNameQuerier interface {
	GetServerName() string
}

type watchedObjectQuerier interface {
	GetWatchedObjectName() string
}

type dataQuerier interface {
	GetData() []byte
}

type ViewRequestInfo interface {
	contextQuerier
	serverNameQuerier
	requestTypeQuerier
	tokenQuerier
}

type viewRequestInfo struct {
	context     context.Context
	serverName  string
	requestType RequestType
	token       string
}

func NewViewRequestInfo(ctx context.Context, reqTypeStr, serverName, token string) ViewRequestInfo {
	requestType := UnknownReqType
	switch reqTypeStr {
	case "DisplayConfig":
		requestType = DisplayConfig
	case "GetConfig":
		requestType = GetConfig
	case "ShowStatistics":
		requestType = ShowStatistics
	case "DisplayServersStatus":
		requestType = DisplayServersStatus
	}
	return &viewRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: requestType,
		token:       token,
	}
}

func (v viewRequestInfo) Context() context.Context {
	return v.context
}

func (v viewRequestInfo) GetServerName() string {
	return v.serverName
}

func (v viewRequestInfo) GetRequestType() RequestType {
	return v.requestType
}

func (v viewRequestInfo) GetToken() string {
	return v.token
}

type UpdateRequestInfo interface {
	contextQuerier
	serverNameQuerier
	requestTypeQuerier
	tokenQuerier
	dataQuerier
}

type updateRequestInfo struct {
	context     context.Context
	serverName  string
	requestType RequestType
	token       string
	dataBuffer  *bytes.Buffer
}

func NewUpdateRequestInfo(ctx context.Context, reqTypeStr, serverName, token string, data []byte) UpdateRequestInfo {
	requestType := UnknownReqType
	switch reqTypeStr {
	case "UpdateConfig":
		requestType = UpdateConfig
	}
	return &updateRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: requestType,
		token:       token,
		dataBuffer:  bytes.NewBuffer(data),
	}
}

func (u updateRequestInfo) Context() context.Context {
	return u.context
}

func (u updateRequestInfo) GetServerName() string {
	return u.serverName
}

func (u updateRequestInfo) GetRequestType() RequestType {
	return u.requestType
}

func (u updateRequestInfo) GetToken() string {
	return u.token
}

func (u updateRequestInfo) GetData() []byte {
	return u.dataBuffer.Bytes()
}

type WatchRequestInfo interface {
	contextQuerier
	serverNameQuerier
	requestTypeQuerier
	tokenQuerier
	watchedObjectQuerier
}

type watchRequestInfo struct {
	context           context.Context
	serverName        string
	requestType       RequestType
	token             string
	watchedObjectName string
}

func NewWatchRequestInfo(ctx context.Context, reqTypeStr, serverName, token, watchedObjectName string) WatchRequestInfo {
	requestType := UnknownReqType
	switch reqTypeStr {
	case "WatchLog":
		requestType = WatchLog
	}
	return &watchRequestInfo{
		context:           ctx,
		serverName:        serverName,
		requestType:       requestType,
		token:             token,
		watchedObjectName: watchedObjectName,
	}
}

func (w watchRequestInfo) Context() context.Context {
	return w.context
}

func (w watchRequestInfo) GetServerName() string {
	return w.serverName
}

func (w watchRequestInfo) GetRequestType() RequestType {
	return w.requestType
}

func (w watchRequestInfo) GetToken() string {
	return w.token
}

func (w watchRequestInfo) GetWatchedObjectName() string {
	return w.watchedObjectName
}
