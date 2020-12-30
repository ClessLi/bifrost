package service

import (
	"github.com/ClessLi/bifrost/pkg/client/auth"
	"golang.org/x/net/context"
)

type WithAuthentication struct {
	Service
	authSvcCli *auth.Client
}

func AuthenticationMiddleware(authServerAddr string) ServiceMiddleware {
	return func(next Service) Service {
		authSvcCli, err := auth.NewClient(authServerAddr)
		if err != nil {
			// log
			return nil
		}
		return WithAuthentication{
			Service:    next,
			authSvcCli: authSvcCli,
		}
	}
}

func (a WithAuthentication) Deal(requester Requester) (responder Responder, err error) {
	ctx := requester.GetContext()
	token := requester.GetToken()
	err = a.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return a.Service.Deal(requester)
}

func (a WithAuthentication) Stop() error {
	a.authSvcCli.Close()
	return a.Service.Stop()
}

func (a WithAuthentication) checkToken(ctx context.Context, token string) error {
	if a.authSvcCli == nil {
		return ErrConnToAuthSvr
	}
	pass, err := a.authSvcCli.Verify(ctx, token)
	if err != nil {
		return err
	}
	if !pass {
		return UnknownErrCheckToken
	}
	return nil
}
