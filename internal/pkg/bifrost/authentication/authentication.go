package authentication

import (
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/ClessLi/bifrost/pkg/client/auth"
	"golang.org/x/net/context"
)

type viewerWithAuthentication struct {
	viewer     service.Viewer
	authSvcCli *auth.Client
}

func (v viewerWithAuthentication) View(requester service.ViewRequester) service.ViewResponder {
	ctx := requester.Context()
	token := requester.GetToken()
	err := checkToken(ctx, token, v.authSvcCli)
	if err != nil {
		return service.NewViewResponder(requester.GetServerName(), []byte(""), err)
	}

	return v.viewer.View(requester)

}

type updaterWithAuthentication struct {
	updater    service.Updater
	authSvcCli *auth.Client
}

func (u updaterWithAuthentication) Update(requester service.UpdateRequester) service.UpdateResponder {
	ctx := requester.Context()
	token := requester.GetToken()
	err := checkToken(ctx, token, u.authSvcCli)
	if err != nil {
		return service.NewUpdateResponder(requester.GetServerName(), err)
	}

	return u.updater.Update(requester)
}

type watcherWithAuthentication struct {
	watcher    service.Watcher
	authSvcCli *auth.Client
}

func (w watcherWithAuthentication) Watch(requester service.WatchRequester) service.WatchResponder {
	ctx := requester.Context()
	token := requester.GetToken()
	err := checkToken(ctx, token, w.authSvcCli)
	if err != nil {
		return service.NewWatchResponder(requester.GetServerName(), nil, nil, nil, err)
	}
	return w.watcher.Watch(requester)
}

type authenticationMiddleware struct {
	viewer  service.Viewer
	updater service.Updater
	watcher service.Watcher
}

func (a authenticationMiddleware) Viewer() service.Viewer {
	return a.viewer
}

func (a authenticationMiddleware) Updater() service.Updater {
	return a.updater
}

func (a authenticationMiddleware) Watcher() service.Watcher {
	return a.watcher
}

func authenticationViewerMiddleware(authSvcCli *auth.Client) service.ViewerMiddleware {
	return func(next service.Viewer) service.Viewer {
		return viewerWithAuthentication{
			viewer:     next,
			authSvcCli: authSvcCli,
		}
	}
}

func authenticationUpdaterMiddleware(authSvcCli *auth.Client) service.UpdaterMiddleware {
	return func(next service.Updater) service.Updater {
		return updaterWithAuthentication{
			updater:    next,
			authSvcCli: authSvcCli,
		}
	}
}

func authenticationWatcherMiddleware(authSvcCli *auth.Client) service.WatcherMiddleware {
	return func(next service.Watcher) service.Watcher {
		return watcherWithAuthentication{
			watcher:    next,
			authSvcCli: authSvcCli,
		}
	}
}

func AuthenticationMiddleware(authServerAddr string) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		authSvcCli, err := auth.NewClient(authServerAddr)
		if err != nil {
			// log
			return nil
		}
		return authenticationMiddleware{
			viewer:  authenticationViewerMiddleware(authSvcCli)(next.Viewer()),
			updater: authenticationUpdaterMiddleware(authSvcCli)(next.Updater()),
			watcher: authenticationWatcherMiddleware(authSvcCli)(next.Watcher()),
		}
	}
}

func checkToken(ctx context.Context, token string, authSvcCli *auth.Client) error {
	if authSvcCli == nil {
		return service.ErrConnToAuthSvr
	}
	pass, err := authSvcCli.Verify(ctx, token)
	if err != nil {
		return err
	}
	if !pass {
		return service.UnknownErrCheckToken
	}
	return nil
}
