package authentication

import (
	authService "github.com/ClessLi/bifrost/internal/pkg/auth/service"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/ClessLi/bifrost/pkg/client/auth"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"testing"
)

type testViewer struct {
}

func (t testViewer) View(info service.ViewRequestInfo) service.ViewResponseInfo {
	return service.NewViewResponseInfo(info.GetServerName(), []byte("view success"), nil)
}

type testUpdater struct {
}

func (t testUpdater) Update(info service.UpdateRequestInfo) service.UpdateResponseInfo {
	return service.NewUpdateResponseInfo(info.GetServerName(), nil)
}

type testWatcher struct {
	dataChan chan []byte
	errChan  chan error
}

func (t testWatcher) Watch(info service.WatchRequestInfo) service.WatchResponseInfo {
	return service.NewWatchResponseInfo(info.GetServerName(), func() error {
		return nil
	}, t.dataChan, t.errChan, nil)
}

func newTestWatcher(dataChan chan []byte, errChan chan error) service.Watcher {
	return &testWatcher{
		dataChan: dataChan,
		errChan:  errChan,
	}
}

// username: test, password: P@s5w0rd, token unexpired: UNabcde, token: abcde
type testAuthService struct {
}

func (t testAuthService) Login(_ context.Context, username, password string, unexpired bool) (string, error) {
	tokenStr := ""
	if strings.EqualFold(username, "test") && strings.EqualFold(password, "P@s5w0rd") {
		if unexpired {
			tokenStr += "UN"
		}
		tokenStr += "abcde"
	}
	if tokenStr == "" {
		return tokenStr, authService.ErrorReasonWrongPassword
	}
	return tokenStr, nil
}

func (t testAuthService) Verify(_ context.Context, token string) (bool, error) {
	switch token {
	case "UNabcde":
		return true, nil
	case "abcde":
		return false, authService.ErrorReasonRelogin
	default:
		return false, authService.ErrorReasonWrongPassword
	}
}

func TestViewerWithAuthentication_View(t *testing.T) {
	authCli := &auth.Client{Service: testAuthService{}}
	viewer := new(testViewer)
	authViewer := authenticationViewerMiddleware(authCli)(viewer)
	unexpiredReqInfo := service.NewViewRequestInfo(context.Background(), "test", "authTest", "UNabcde")
	expiredReqInfo := service.NewViewRequestInfo(context.Background(), "test", "authTest", "abcde")
	wrongTokenWithBlockReqInfo := service.NewViewRequestInfo(context.Background(), "test", "authTest", " abcde")
	wrongTokenWithReqInfo := service.NewViewRequestInfo(context.Background(), "test", "authTest", "asldfj")
	respInfo := authViewer.View(unexpiredReqInfo)
	t.Logf("view with unexpired token. svcName: %s, response: %s, error: %s", respInfo.GetServerName(), respInfo.Bytes(), respInfo.Error())
	respInfo = authViewer.View(expiredReqInfo)
	t.Logf("view with expired token. svcName: %s, response: %s, error: %s", respInfo.GetServerName(), respInfo.Bytes(), respInfo.Error())
	respInfo = authViewer.View(wrongTokenWithBlockReqInfo)
	t.Logf("view with wrong token (which right token with block). svcName: %s, response: %s, error: %s", respInfo.GetServerName(), respInfo.Bytes(), respInfo.Error())
	respInfo = authViewer.View(wrongTokenWithReqInfo)
	t.Logf("view with wrong token. svcName: %s, response: %s, error: %s", respInfo.GetServerName(), respInfo.Bytes(), respInfo.Error())
}

func TestUpdaterWithAuthentication_Update(t *testing.T) {
	authCli := &auth.Client{Service: testAuthService{}}
	updater := new(testUpdater)
	authUpdater := authenticationUpdaterMiddleware(authCli)(updater)
	unexpiredReqInfo := service.NewUpdateRequestInfo(context.Background(), "test", "authTest", "UNabcde", nil)
	expiredReqInfo := service.NewUpdateRequestInfo(context.Background(), "test", "authTest", "abcde", nil)
	wrongTokenWithBlockReqInfo := service.NewUpdateRequestInfo(context.Background(), "test", "authTest", " abcde", nil)
	wrongTokenWithReqInfo := service.NewUpdateRequestInfo(context.Background(), "test", "authTest", "asldfj", nil)
	respInfo := authUpdater.Update(unexpiredReqInfo)
	t.Logf("view with unexpired token. svcName: %s, error: %s", respInfo.GetServerName(), respInfo.Error())
	respInfo = authUpdater.Update(expiredReqInfo)
	t.Logf("view with expired token. svcName: %s, error: %s", respInfo.GetServerName(), respInfo.Error())
	respInfo = authUpdater.Update(wrongTokenWithBlockReqInfo)
	t.Logf("view with wrong token (which right token with block). svcName: %s, error: %s", respInfo.GetServerName(), respInfo.Error())
	respInfo = authUpdater.Update(wrongTokenWithReqInfo)
	t.Logf("view with wrong token. svcName: %s, error: %s", respInfo.GetServerName(), respInfo.Error())
}

func TestWatcherWithAuthentication_Watch(t *testing.T) {
	authCli := &auth.Client{Service: testAuthService{}}
	dataChan := make(chan []byte)
	errChan := make(chan error)
	watcher := newTestWatcher(dataChan, errChan)
	authWatcher := authenticationWatcherMiddleware(authCli)(watcher)
	unexpiredReqInfo := service.NewWatchRequestInfo(context.Background(), "test", "authTest", "UNabcde", "access.log")
	expiredReqInfo := service.NewWatchRequestInfo(context.Background(), "test", "authTest", "abcde", "access.log")
	wrongTokenWithBlockReqInfo := service.NewWatchRequestInfo(context.Background(), "test", "authTest", " abcde", "access.log")
	wrongTokenWithReqInfo := service.NewWatchRequestInfo(context.Background(), "test", "authTest", "asldfj", "access.log")
	respInfo := authWatcher.Watch(unexpiredReqInfo)
	t.Logf("view with unexpired token. svcName: %s, bytes channel is nil: %s, error channel is nil: %s, error: %s, Close() error: %s", respInfo.GetServerName(), strconv.FormatBool(respInfo.BytesChan() == nil), strconv.FormatBool(respInfo.TransferErrorChan() == nil), respInfo.Error(), respInfo.Close())
	respInfo = authWatcher.Watch(expiredReqInfo)
	t.Logf("view with expired token. svcName: %s, bytes channel is nil: %s, error channel is nil: %s, error: %s, Close() error: %s", respInfo.GetServerName(), strconv.FormatBool(respInfo.BytesChan() == nil), strconv.FormatBool(respInfo.TransferErrorChan() == nil), respInfo.Error(), respInfo.Close())
	respInfo = authWatcher.Watch(wrongTokenWithBlockReqInfo)
	t.Logf("view with wrong token (which right token with block). svcName: %s, bytes channel is nil: %s, error channel is nil: %s, error: %s, Close() error: %s", respInfo.GetServerName(), strconv.FormatBool(respInfo.BytesChan() == nil), strconv.FormatBool(respInfo.TransferErrorChan() == nil), respInfo.Error(), respInfo.Close())
	respInfo = authWatcher.Watch(wrongTokenWithReqInfo)
	t.Logf("view with wrong token. svcName: %s, bytes channel is nil: %s, error channel is nil: %s, error: %s, Close() error: %s", respInfo.GetServerName(), strconv.FormatBool(respInfo.BytesChan() == nil), strconv.FormatBool(respInfo.TransferErrorChan() == nil), respInfo.Error(), respInfo.Close())
}
