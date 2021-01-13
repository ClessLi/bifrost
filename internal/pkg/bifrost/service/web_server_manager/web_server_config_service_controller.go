package web_server_manager

import (
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"sync"
	"time"
)

type WebServerConfigServiceController interface {
	serverName() string
	Status() State
	SetState(state State)
	GetService() WebServerConfigService
	statusControl()
	autoBackup()
	autoReload()
}

func newWebServerConfigServiceController(info WebServerConfigInfo) WebServerConfigServiceController {
	switch info.Type {
	case NGINX:
		controller := NewNginxConfigServiceController(info)
		go controller.statusControl()
		return controller
	}
	return nil
}

type WebServerConfigServicesController struct {
	controllers map[string]WebServerConfigServiceController
}

func (m *WebServerConfigServicesController) Stop() error {
	return m.offstagesStatusControl(Disabled)
}

func (m *WebServerConfigServicesController) Start() error {
	return m.offstagesStatusControl(Normal)
}

func (m WebServerConfigServicesController) Status(serverName string) (State, error) {
	offstage, err := m.getOffstageController(serverName)
	if err != nil {
		return Unknown, err
	}
	return offstage.Status(), nil
}

func (m WebServerConfigServicesController) GetServicesHandler() *WebServerConfigServicesHandler {
	services := make([]WebServerConfigService, 0)
	for _, controller := range m.controllers {
		services = append(services, controller.GetService())
	}
	return newWebServerConfigServiceHandler(services...)
}

func (m WebServerConfigServicesController) getOffstageController(serverName string) (WebServerConfigServiceController, error) {
	if controller, ok := m.controllers[serverName]; ok {
		return controller, nil
	}
	return nil, ErrOffstageNotExist
}

func (m *WebServerConfigServicesController) offstagesStatusControl(expectedState State) error {
	if expectedState == Unknown || expectedState == Abnormal {
		return ErrWrongStateExpectation
	}
	workQueue := new(utils.StringQueue)
	workTimesLimit := make(map[string]int)
	//workIsWait := make(map[string]bool)
	workIsWait := new(sync.Map)
	for serverName := range m.controllers {
		workQueue.Add(serverName)
		m.controllers[serverName].SetState(expectedState)
		workTimesLimit[serverName] = 100
		workIsWait.Store(serverName, false)
	}
	var workErr error
	for !workQueue.IsEmpty() {
		serverName := workQueue.Poll()
		if workTimesLimit[serverName] < 1 {
			delete(workTimesLimit, serverName)
			if workErr == nil {
				workErr = fmt.Errorf("contral web server config manager error, cased by: <%s: %s>", serverName, "limit of working times")
			}
			workErr = fmt.Errorf("%s; <%s: %s>", workErr, serverName, "limit of working times")
			continue
		}

		currentState, err := m.Status(serverName)
		if err != nil {
			if workErr == nil {
				workErr = fmt.Errorf("web server config manager stop error, cased by: <%s: %s>", serverName, err)
			}
			workErr = fmt.Errorf("%s; <%s: %s>", workErr, serverName, err)
		}
		if currentState != expectedState {
			isWait, ok := workIsWait.Load(serverName)
			if ok && !isWait.(bool) {
				workIsWait.Store(serverName, true)
				go func() {
					workTimesLimit[serverName] -= 1
					time.Sleep(time.Millisecond * 200)
					workIsWait.Store(serverName, false)
				}()
			}
			workQueue.Add(serverName)
			time.Sleep(time.Millisecond)
		}
	}
	return workErr
}

func NewWebServerConfigServicesController(infos ...WebServerConfigInfo) *WebServerConfigServicesController {
	if infos == nil || len(infos) < 1 {
		return nil
	}
	c := &WebServerConfigServicesController{controllers: make(map[string]WebServerConfigServiceController)}
	for _, info := range infos {
		controller := newWebServerConfigServiceController(info)
		if controller == nil {
			continue
		}
		c.controllers[controller.serverName()] = controller
	}
	return c
}
