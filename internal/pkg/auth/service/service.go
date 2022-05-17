package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/net/context"

	"github.com/yongPhone/bifrost/pkg/resolv/nginx"
)

var (
	// 认证接口错误返回.
	ErrorReasonServerBusy    = errors.New("服务器繁忙")
	ErrorReasonRelogin       = errors.New("请重新登陆")
	ErrorReasonWrongPassword = errors.New("用户或密码错误")
	// ErrorReasonNoneToken     = "请通过认证".
)

type Service interface {
	Login(ctx context.Context, username, password string, unexpired bool) (string, error)
	Verify(ctx context.Context, token string) (bool, error)
	// GetPort() int
}

// AuthDBConfig, mysql数据库信息结构体，该库用于存放用户认证信息（可选）.
type AuthDBConfig struct {
	DBName   string `yaml:"DBName"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// AuthConfig, 认证信息结构体，记录用户认证信息（可选）.
type AuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// AuthService, bifrost认证服务结构体，用于用户认证.
type AuthService struct {
	Port          int `yaml:"Port"`
	*AuthDBConfig `yaml:"AuthDBConfig,omitempty"`
	*AuthConfig   `yaml:"AuthConfig,omitempty"`
}

func (s *AuthService) Login(ctx context.Context, username, password string, unexpired bool) (string, error) {
	//ip, err := getClientIP(ctx)
	//if err != nil {
	//	return "", err
	//}

	// 初始化jwt断言对象
	claims := &JWTClaims{
		UserID:      1,
		Username:    username,
		Password:    password,
		FullName:    username,
		Permissions: []string{},
	}
	claims.IssuedAt = &jwt.NumericDate{Time: time.Now().In(nginx.TZ)}
	if unexpired {
		claims.ExpiresAt = &jwt.NumericDate{Time: time.Time{}}
	} else {
		claims.ExpiresAt = &jwt.NumericDate{Time: time.Now().In(nginx.TZ).Add(time.Second * time.Duration(ExpireTime))}
	}

	// 认证用户信息
	if !s.validUser(claims) {
		// Log(WARN, "[%s] Invalid user '%s' or password '%s'.", ip, claims.Username, claims.Password)
		return "", ErrorReasonWrongPassword
	}

	// 生成用户token
	signedToken, err := claims.getToken()
	if err != nil {
		// Log(NOTICE, "[%s] user '%s' login failed, message is: '%s'", ip, username, err.Error())
		return "", err
	}
	// Log(NOTICE, "[%s] user '%s' is login, token is: %s", ip, username, signedToken)

	return signedToken, err
}

func (s *AuthService) Verify(ctx context.Context, token string) (bool, error) {
	//ip, err := getClientIP(ctx)
	//if err != nil {
	//	return false, err
	//}
	_, err := s.verifyAction(token)
	if err != nil {
		// err = fmt.Errorf("[%s] Verified failed: %s", ip, err)
		err = fmt.Errorf("verified failed: %s", err)
		return false, err
	}
	return true, nil
}

// ServiceMiddleware define service middleware.
type ServiceMiddleware func(Service) Service
