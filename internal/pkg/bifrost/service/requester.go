package service

import (
	"bytes"
	"golang.org/x/net/context"
)

type ViewRequester interface {
	contextQuerier
	serverNameQuerier
	requestTypeQuerier
	tokenQuerier
}

type viewRequester struct {
	context     context.Context
	serverName  string
	requestType RequestType
	token       string
}

func NewViewRequester(ctx context.Context, reqType, serverName, token string) ViewRequester {
	requestType := Unknown
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
	return &viewRequester{
		context:     ctx,
		serverName:  serverName,
		requestType: requestType,
		token:       token,
	}
}

func (v viewRequester) Context() context.Context {
	return v.context
}

func (v viewRequester) GetServerName() string {
	return v.serverName
}

func (v viewRequester) GetRequestType() RequestType {
	return v.requestType
}

func (v viewRequester) GetToken() string {
	return v.token
}

type UpdateRequester interface {
	contextQuerier
	serverNameQuerier
	requestTypeQuerier
	tokenQuerier
	dataQuerier
}

type updateRequester struct {
	context     context.Context
	serverName  string
	requestType RequestType
	token       string
	data        *bytes.Buffer
}

func NewUpdateRequester(ctx context.Context, reqType, serverName, token string, data []byte) UpdateRequester {
	requestType := Unknown
	switch reqType {
	case "UpdateConfig":
		requestType = UpdateConfig
	}
	return &updateRequester{
		context:     ctx,
		serverName:  serverName,
		requestType: requestType,
		token:       token,
		data:        bytes.NewBuffer(data),
	}
}

func (u updateRequester) Context() context.Context {
	return u.context
}

func (u updateRequester) GetServerName() string {
	return u.serverName
}

func (u updateRequester) GetRequestType() RequestType {
	return u.requestType
}

func (u updateRequester) GetToken() string {
	return u.token
}

func (u updateRequester) GetData() []byte {
	return u.data.Bytes()
}

type WatchRequester interface {
	contextQuerier
	serverNameQuerier
	requestTypeQuerier
	tokenQuerier
	objectQuerier
}

type watchRequester struct {
	context     context.Context
	serverName  string
	requestType RequestType
	token       string
	objectName  string
}

func NewWatchRequester(ctx context.Context, reqType, serverName, token, objectName string) WatchRequester {
	requestType := Unknown
	switch reqType {
	case "WatchLog":
		requestType = WatchLog
	}
	return &watchRequester{
		context:     ctx,
		serverName:  serverName,
		requestType: requestType,
		token:       token,
		objectName:  objectName,
	}
}

func (w watchRequester) Context() context.Context {
	return w.context
}

func (w watchRequester) GetServerName() string {
	return w.serverName
}

func (w watchRequester) GetRequestType() RequestType {
	return w.requestType
}

func (w watchRequester) GetToken() string {
	return w.token
}

func (w watchRequester) GetObjectName() string {
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
