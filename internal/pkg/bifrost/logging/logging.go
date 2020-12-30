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

func (lmw loggingMiddleware) Deal(requester service.Requester) (responder service.Responder, err error) {
	ip, err := getClientIP(requester.GetContext())

	defer func(begin time.Time) {
		n := 100
		if responder != nil {
			if data, err := responder.Bytes(); data != nil {
				if len(data) < n {
					n = len(data)
				}
				lmw.logger.Log(
					"functions", "Deal",
					"requestType", requester.GetRequestType(),
					"clientIp", ip,
					"token", requester.GetToken(),
					"webServerName", requester.GetServerName(),
					"param", requester.GetParam(),
					"result", string(data[:n])+"...",
					"error", err,
					"took", time.Since(begin),
				)
				return
			} else if watcher, err := responder.GetWatcher(); watcher != nil {
				lmw.logger.Log(
					"functions", "Deal",
					"requestType", requester.GetRequestType(),
					"clientIp", ip,
					"token", requester.GetToken(),
					"webServerName", requester.GetServerName(),
					"param", requester.GetParam(),
					"error", err,
					"took", time.Since(begin),
				)
				return
			}
		}
		lmw.logger.Log(
			"functions", "Deal",
			"requestType", requester.GetRequestType(),
			"clientIp", ip,
			"token", requester.GetToken(),
			"webServerName", requester.GetServerName(),
			"param", requester.GetParam(),
			"responder", responder,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now().Local())
	if err != nil {
		return nil, err
	}
	responder, err = lmw.Service.Deal(requester)
	return responder, err
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
