package logging

import (
	"context"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc/peer"
	"net"
	"time"
)

// loggingMiddleware Make a new type
// that contains BifrostService interface and logger instance
type loggingMiddleware struct {
	service.Service
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func LoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return loggingMiddleware{next, logger}
	}
}

func (mw loggingMiddleware) ViewConfig(ctx context.Context, token, svrName string) (data []byte, err error) {
	ip, cipErr := getClientIP(ctx)
	defer func(begin time.Time) {
		n := 100
		if len(data) < n {
			n = len(data)
		}
		mw.logger.Log(
			"functions", "ViewConfig",
			"clientIp", ip,
			"token", token,
			"svrName", svrName,
			"result", string(data[:n])+"...",
			"took", time.Since(begin),
		)
	}(time.Now().Local())
	if cipErr != nil {
		return data, cipErr
	}
	data, err = mw.Service.ViewConfig(ctx, token, svrName)
	return data, err
}

func (mw loggingMiddleware) GetConfig(ctx context.Context, token, svrName string) (jsonData []byte, err error) {
	ip, cipErr := getClientIP(ctx)
	defer func(begin time.Time) {
		n := 100
		if len(jsonData) < n {
			n = len(jsonData)
		}
		mw.logger.Log(
			"functions", "GetConfig",
			"clientIp", ip,
			"token", token,
			"svrName", svrName,
			"result", string(jsonData[:n])+"...",
			"took", time.Since(begin),
		)
	}(time.Now().Local())
	if cipErr != nil {
		return jsonData, cipErr
	}
	jsonData, err = mw.Service.GetConfig(ctx, token, svrName)
	return jsonData, err
}

func (mw loggingMiddleware) UpdateConfig(ctx context.Context, token, svrName string, jsonData []byte) (data []byte, err error) {
	ip, cipErr := getClientIP(ctx)
	defer func(begin time.Time) {
		n := 100
		if len(jsonData) < n {
			n = len(jsonData)
		}
		mw.logger.Log(
			"functions", "UpdateConfig",
			"clientIp", ip,
			"token", token,
			"svrName", svrName,
			"reqJsonData", string(jsonData[:n])+"...",
			"errorMsg", err,
			"took", time.Since(begin),
		)
	}(time.Now().Local())
	if cipErr != nil {
		return nil, cipErr
	}
	data, err = mw.Service.UpdateConfig(ctx, token, svrName, jsonData)
	return data, err
}

func (mw loggingMiddleware) ViewStatistics(ctx context.Context, token, svrName string) (jsonData []byte, err error) {
	ip, cipErr := getClientIP(ctx)
	defer func(begin time.Time) {
		n := 100
		if len(jsonData) < n {
			n = len(jsonData)
		}
		mw.logger.Log(
			"functions", "ViewStatistics",
			"clientIp", ip,
			"token", token,
			"svrName", svrName,
			"result", string(jsonData[:n])+"...",
			"took", time.Since(begin),
		)
	}(time.Now().Local())
	if cipErr != nil {
		return jsonData, cipErr
	}
	jsonData, err = mw.Service.ViewStatistics(ctx, token, svrName)
	return jsonData, err
}

func (mw loggingMiddleware) Status(ctx context.Context, token string) (jsonData []byte, err error) {
	ip, cipErr := getClientIP(ctx)
	defer func(begin time.Time) {
		n := 100
		if len(jsonData) < n {
			n = len(jsonData)
		}
		mw.logger.Log(
			"functions", "Status",
			"clientIp", ip,
			"token", token,
			"result", string(jsonData[:n])+"...",
			"took", time.Since(begin),
		)
	}(time.Now().Local())
	if cipErr != nil {
		return jsonData, cipErr
	}
	jsonData, err = mw.Service.Status(ctx, token)
	return jsonData, err
}

func (mw loggingMiddleware) WatchLog(ctx context.Context, token, svrName, logName string, dataChan chan<- []byte, signal <-chan int) (err error) {
	ip, cipErr := getClientIP(ctx)
	defer func(begin time.Time) {
		mw.logger.Log(
			"functions", "WatchLog",
			"clientIp", ip,
			"token", token,
			"svrName", svrName,
			"logFile", logName,
			"error", err,
			"during", time.Since(begin),
		)
	}(time.Now().Local())
	if cipErr != nil {
		return cipErr
	}
	err = mw.Service.WatchLog(ctx, token, svrName, logName, dataChan, signal)
	return err
}

func (mw loggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "HealthChcek",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now().Local())
	result = true
	return
}

func getClientIP(ctx context.Context) (ip string, err error) {
	//md, ok := metadata.FromIncomingContext(ctx)
	pr, ok := peer.FromContext(ctx)
	if !ok {
		err = fmt.Errorf("getClientIP, invoke FromContext() failed")
		//err = fmt.Errorf("getClientIP, invoke FromIncomingContext() failed")
		return
	}
	if pr.Addr == net.Addr(nil) {
		err = fmt.Errorf("getClientIP, peer.Addr is nil")
		return
	}
	//fmt.Println(md)
	//ips := md.Get("x-real-ip")
	//if len(ips) == 0 {
	//	err = fmt.Errorf("get real ip failed")
	//	return
	//}
	//ip = ips[0]
	ip = pr.Addr.String()
	return ip, nil
}
