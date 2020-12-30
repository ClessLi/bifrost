package service

import (
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service/web_server_manager"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"golang.org/x/net/context"
	"sync"
)

type Service interface {
	Deal(requester Requester) (Responder, error)
	Stop() error
	// TODO: WatchLog暂用锁机制，一个日志文件仅允许一个终端访问
}

// BifrostService, bifrost配置文件对象中web服务器信息结构体，定义管控的web服务器配置文件相关信息
type BifrostService struct {
	managers  map[string]web_server_manager.WebServerManager
	monitor   Monitor
	waitGroup *sync.WaitGroup
}

func NewService(managers map[string]web_server_manager.WebServerManager) Service {
	if managers == nil {
		return nil
	}
	var err error
	defer func() {
		if err != nil {
			utils.Logger.FatalF("failed to initialize bifrost service, cased by %s", err)
		}
	}()
	svc := BifrostService{
		managers:  managers,
		monitor:   NewSysInfo(managers),
		waitGroup: new(sync.WaitGroup),
	}

	for name, manager := range svc.managers {
		managerErr := manager.Start()
		if managerErr != nil {
			utils.Logger.Warning(managerErr.Error())
			delete(svc.managers, name)
			break
		}
		svc.waitGroup.Add(1)
	}

	// 监控系统信息
	err = svc.monitor.Start()
	if err != nil {
		return nil
	}
	svc.waitGroup.Add(1)
	return &svc
}

func (b BifrostService) Deal(requester Requester) (responder Responder, err error) {
	switch requester.GetRequestType() {
	case DisplayConfig, GetConfig, ShowStatistics:
		resp, err := b.get(requester.GetContext(), requester.GetServerName(), requester.GetRequestType())
		return NewQueryOrUpdateResponse(resp, err), nil
	case DisplayStatus:
		resp, err := b.status(requester.GetContext())
		return NewQueryOrUpdateResponse(resp, err), nil
	case UpdateConfig:
		err := b.update(requester.GetContext(), requester.GetServerName(), requester.GetRequestType(), requester.GetRequestData(), requester.GetParam())
		if err != nil {
			return nil, err
		}
		return NewQueryOrUpdateResponse([]byte("Update succeeded"), err), nil
	case WatchLog:
		watcher, err := b.watch(requester.GetContext(), requester.GetServerName(), requester.GetRequestType(), requester.GetParam())
		if err != nil {
			return nil, err
		}
		return NewWatcherResponse(watcher), nil
	default:
		return nil, UnknownRequestType
	}
}

// Stop, ServiceInfo关闭协程任务的方法
func (b *BifrostService) Stop() error {
	defer b.waitGroup.Wait()
	err := b.monitor.Stop()
	if err != nil {
		// log
	}
	b.waitGroup.Done()
	for _, manager := range b.managers {
		err := manager.Stop()
		if err != nil {
			// log
			continue
		}
		b.waitGroup.Done()
	}
	return err
}

func (b BifrostService) get(ctx context.Context, svrName string, reqType RequestType) (resp []byte, err error) {
	if reqType == DisplayStatus {
		return b.monitor.DisplayStatus()
	}
	manager, err := b.getWebServerConfigManager(svrName)
	if err != nil {
		return nil, err
	}
	switch reqType {
	case DisplayConfig:
		return manager.DisplayConfig()
	case GetConfig:
		return manager.GetConfig()
	case ShowStatistics:
		return manager.ShowStatistics()
	default:
		return nil, UnknownRequestType
	}
}

func (b *BifrostService) update(ctx context.Context, svrName string, reqType RequestType, uData []byte, params ...interface{}) (err error) {
	manager, err := b.getWebServerConfigManager(svrName)
	if err != nil {
		return err
	}

	switch reqType {
	case UpdateConfig:
		if params != nil && len(params) > 0 {
			param, ok := params[0].(string)
			if ok {
				return manager.UpdateConfig(uData, param)
			}
			return web_server_manager.ErrWrongParamPassedIn
		}
		return ErrParamNotPassedIn
	default:
		return UnknownRequestType
	}
}

func (b BifrostService) watch(ctx context.Context, svrName string, reqType RequestType, params ...interface{}) (watcher web_server_manager.Watcher, err error) {
	manager, err := b.getWebServerConfigManager(svrName)
	if err != nil {
		return nil, err
	}
	switch reqType {
	case WatchLog:
		if params != nil && len(params) > 0 {
			logName, ok := params[0].(string)
			if ok {
				return manager.WatchLog(logName)
			}
		}
		return nil, ErrParamNotPassedIn
	default:
		return nil, UnknownRequestType
	}
}

func (b *BifrostService) status(ctx context.Context) (response []byte, err error) {
	// TODO: SysInfo lock mechanism
	return b.monitor.DisplayStatus()
}

// getWebServerConfigManager, 查询获取WebServerManager对象的方法
// 参数:
//     name: web服务名
func (b BifrostService) getWebServerConfigManager(svrName string) (manager web_server_manager.WebServerManager, err error) {
	var ok bool
	manager, ok = b.managers[svrName]
	if ok {
		return manager, nil
	}
	return nil, ErrUnknownSvrName
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
