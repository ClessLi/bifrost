package service

import (
	"bytes"
	"golang.org/x/net/context"
)

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

func NewViewRequestInfo(ctx context.Context, reqType, serverName, token string) ViewRequestInfo {
	requestType := UnknownReqType
	switch reqType {
	case "DisplayConfig":
		requestType = DisplayConfig
	case "GetConfig":
		requestType = GetConfig
	case "ShowStatistics":
		requestType = ShowStatistics
	case "DisplayStatus":
		requestType = DisplayStatus
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
	data        *bytes.Buffer
}

func NewUpdateRequestInfo(ctx context.Context, reqType, serverName, token string, data []byte) UpdateRequestInfo {
	requestType := UnknownReqType
	switch reqType {
	case "UpdateConfig":
		requestType = UpdateConfig
	}
	return &updateRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: requestType,
		token:       token,
		data:        bytes.NewBuffer(data),
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
	return u.data.Bytes()
}

type WatchRequestInfo interface {
	contextQuerier
	serverNameQuerier
	requestTypeQuerier
	tokenQuerier
	objectQuerier
}

type watchRequestInfo struct {
	context     context.Context
	serverName  string
	requestType RequestType
	token       string
	objectName  string
}

func NewWatchRequestInfo(ctx context.Context, reqType, serverName, token, objectName string) WatchRequestInfo {
	requestType := UnknownReqType
	switch reqType {
	case "WatchLog":
		requestType = WatchLog
	}
	return &watchRequestInfo{
		context:     ctx,
		serverName:  serverName,
		requestType: requestType,
		token:       token,
		objectName:  objectName,
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

func (w watchRequestInfo) GetObjectName() string {
	return w.objectName
}

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

type objectQuerier interface {
	GetObjectName() string
}

type dataQuerier interface {
	GetData() []byte
}
