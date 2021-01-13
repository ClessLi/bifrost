package service

import (
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service/web_server_manager"
)

type BifrostServiceController struct {
	webServerConfigServicesController *web_server_manager.WebServerConfigServicesController
	systemInfo                        *systemInfo
}

func (b *BifrostServiceController) Start() error {
	err := b.webServerConfigServicesController.Start()
	if err != nil {
		return err
	}
	return b.systemInfo.Start()
}

func (b *BifrostServiceController) Stop() error {
	err := b.webServerConfigServicesController.Stop()
	if err != nil {
		return err
	}
	return b.systemInfo.Stop()
}

func (b *BifrostServiceController) GetService() *BifrostService {
	return &BifrostService{
		webServerConfigServicesHandler: b.webServerConfigServicesController.GetServicesHandler(),
		monitor:                        b.systemInfo,
	}
}

func NewBifrostServiceController(infos ...web_server_manager.WebServerConfigInfo) *BifrostServiceController {
	webServerConfigServicesController := web_server_manager.NewWebServerConfigServicesController(infos...)
	systemInfo := NewSysInfo(webServerConfigServicesController.GetServicesHandler())
	return &BifrostServiceController{
		webServerConfigServicesController: webServerConfigServicesController,
		systemInfo:                        systemInfo,
	}
}
