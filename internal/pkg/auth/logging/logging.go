package logging

import (
	"context"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/auth/service"
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

func (mw loggingMiddleware) Login(ctx context.Context, username, password string, unexpired bool) (ret string, err error) {
	ip, cipErr := getClientIP(ctx)
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Login",
			"clientIP", ip,
			"username", username,
			"password", password,
			"unexpired", unexpired,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())
	if cipErr != nil {
		return ret, cipErr
	}

	ret, err = mw.Service.Login(ctx, username, password, unexpired)
	return ret, err
}

func (mw loggingMiddleware) Verify(ctx context.Context, token string) (ret bool, err error) {
	ip, cipErr := getClientIP(ctx)
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Verify",
			"clientIP", ip,
			"token", token,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now().Local())
	if cipErr != nil {
		return ret, cipErr
	}

	ret, err = mw.Service.Verify(ctx, token)
	return ret, err
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
