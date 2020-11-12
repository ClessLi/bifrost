package bifrost

import (
	"fmt"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"net"
	"time"
)

// AuthService, bifrost认证服务结构体，用于用户认证
type AuthService struct {
	*AuthDBConfig `yaml:"AuthDBConfig,omitempty"`
	*AuthConfig   `yaml:"AuthConfig,omitempty"`
}

func (a AuthService) Login(ctx context.Context, req *bifrostpb.AuthRequest) (resp *bifrostpb.AuthResponse, err error) {
	resp = &bifrostpb.AuthResponse{}
	ip, err := getClientIP(ctx)
	if err != nil {
		resp.Err = err.Error()
		return resp, err
	}

	// 初始化jwt断言对象
	claims := &JWTClaims{
		UserID:      1,
		Username:    req.Username,
		Password:    req.Password,
		FullName:    req.Username,
		Permissions: []string{},
	}
	claims.IssuedAt = time.Now().In(nginx.TZ).Unix()
	if req.Unexpired {
		claims.ExpiresAt = 0
	} else {
		claims.ExpiresAt = time.Now().In(nginx.TZ).Add(time.Second * time.Duration(ExpireTime)).Unix()
	}

	// 认证用户信息
	if !validUser(claims) {
		Log(WARN, "[%s] Invalid user '%s' or password '%s'.", ip, claims.Username, claims.Password)
		err = fmt.Errorf(ErrorReasonWrongPassword)
		resp.Err = ErrorReasonWrongPassword
		return resp, err
	}

	// 生成用户token
	signedToken, err := getToken(claims)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		Log(NOTICE, "[%s] user '%s' login failed, message is: '%s'", ip, req.Username, err.Error())
		resp.Err = err.Error()
		return resp, err
	}
	Log(NOTICE, "[%s] user '%s' is login, token is: %s", ip, req.Username, signedToken)

	resp.Token = signedToken
	return resp, err
}

func Verify(ctx context.Context, token string) (ip string, err error) {
	ip, err = getClientIP(ctx)
	if err != nil {
		return
	}
	_, err = verifyAction(token)
	if err != nil {
		err = fmt.Errorf("[%s] Verified failed: %s", ip, err)
		return
	}
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
