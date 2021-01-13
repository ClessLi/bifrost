package web_server_manager

import (
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"sync"
	"time"
)

type nginxConfigServiceController struct {
	service        *nginxConfigServiceWithState
	expectedState  State
	waitGroup      *sync.WaitGroup
	autoBackupChan chan int
	autoReloadChan chan int
}

func (n nginxConfigServiceController) serverName() string {
	return n.service.serverName()
}

func (n nginxConfigServiceController) Status() State {
	return n.service.Status()
}

func (n *nginxConfigServiceController) SetState(state State) {
	n.expectedState = state
}

func (n *nginxConfigServiceController) statusControl() {
	for {
		var err error
		if n.Status() != n.expectedState {
			switch n.expectedState {
			case Normal:
				err = n.toNormal()
			case Disabled:
				err = n.toDisable()
			default:
				n.expectedState = n.Status()
			}
			if err != nil {
				utils.Logger.WarningF("[%s] %s", n.serverName(), err)
			}
			time.Sleep(time.Millisecond * 100)
		} else {
			time.Sleep(time.Second * 10)
		}
	}
}

func (n *nginxConfigServiceController) toNormal() (err error) {
	//Log(DEBUG, "[%s] 初始化bifrost服务相关接口。。。", b.ServiceInfos[i].Name)
	defer func() {
		if err != nil {
			n.service.SetState(Disabled)
		}
		n.service.SetState(Normal)
	}()

	if n.Status() > Initializing {
		err = n.toDisable()
		if err != nil {
			return err
		}
	} else if n.Status() == Initializing {
		return nil
	}

	n.service.SetState(Initializing)
	err = n.service.configLoad()
	if err != nil {
		//fmt.Printf("[%s] load config error: %s\n", n.name, err)
		utils.Logger.ErrorF("[%s] load config error: %s", n.serverName(), err)
		return err
	}

	// DONE: 执行备份与自动加载
	n.autoBackup()
	//Log(DEBUG, "[%s] 载入备份协程", b.ServiceInfos[i].Name)
	n.autoReload()
	//Log(DEBUG, "[%s] 载入自动更新配置协程", b.ServiceInfos[i].Name)
	return nil
}

func (n *nginxConfigServiceController) toDisable() error {
	defer func() {
		n.service.SetState(Disabled)
	}()
	defer n.waitGroup.Wait()
	//utils.Logger.DebugF("[%s] stop backup proc", n.name)
	if n.autoBackupChan != nil {
		//utils.Logger.DebugF("[%s] stop backup proc", n.serverName())
		n.autoBackupChan <- 9
	}
	//utils.Logger.DebugF("[%s] stop config auto reload proc", n.name)
	if n.autoReloadChan != nil {
		//utils.Logger.DebugF("[%s] stop config auto reload proc", n.serverName())
		n.autoReloadChan <- 9
	}
	return nil
}

// autoBackup, ServerInfo的nginx配置文件备份方法
// 参数:
//     c: 整型管道，用于停止备份
func (n *nginxConfigServiceController) autoBackup() {
	n.autoBackupChan = make(chan int)
	go func() {
		n.waitGroup.Add(1)
		defer n.waitGroup.Done()
		for n.Status() > Disabled {
			select {
			case <-time.NewTicker(5 * time.Minute).C: // 每5分钟定时执行备份操作
				utils.Logger.DebugF("[%s] Nginx Config check and backup", n.serverName())
				_ = n.service.configBackup()
			case sig := <-n.autoBackupChan: // 获取管道传入信号
				if sig == 9 { // 为9时，停止备份
					//utils.Logger.InfoF("[%s] Nginx Config autoBackup method stopping...", n.serverName())
					return
				}
			}
		}
	}()
}

// autoReload, ServerInfo的web服务器配置文件自动热加载方法
// 参数:
//     c: 整型管道，用于停止备份
func (n *nginxConfigServiceController) autoReload() {
	n.autoReloadChan = make(chan int)
	go func() {
		n.waitGroup.Add(1)
		defer n.waitGroup.Done()
		for n.Status() > Disabled {
			select {
			case <-time.NewTicker(30 * time.Second).C: // 每30秒检查一次nginx配置文件是否已在后台更新
				utils.Logger.DebugF("[%s] Nginx Config check and reloading", n.serverName())
				reloadErr := n.configReload()
				if reloadErr != nil && reloadErr != nginx.NoReloadRequired {
					utils.Logger.WarningF("[%s] Nginx Config reload failed, cased by '%s'", n.serverName(), reloadErr)
					if n.Status() == Normal {
						n.service.SetState(Abnormal)
					}
				} else if reloadErr == nil {
					utils.Logger.InfoF("[%s] Nginx Config reload successfully", n.serverName())
				}
			case sig := <-n.autoReloadChan: // 获取管道传入信号
				if sig == 9 { // 为9时，停止备份
					//utils.Logger.InfoF("[%s] Nginx Config autoReload method stopping...", n.serverName())
					return
				}
			}
		}
	}()
}

// autoReload, ServerInfo的web服务器配置文件自动热加载子方法
func (n *nginxConfigServiceController) configReload() error {
	// 校验配置文件是否更新
	isSame, checkErr := n.service.checkConfigsHash()
	if checkErr != nil {
		return checkErr
	}

	// 如果有差别，则重新读取配置
	if !isSame {
		//Log(DEBUG, "[%s] reloading nginx config", n.Name)
		return n.service.configLoad()
	}
	return nginx.NoReloadRequired
}

func (n nginxConfigServiceController) GetService() WebServerConfigService {
	return n.service
}

func NewNginxConfigServiceController(info WebServerConfigInfo) WebServerConfigServiceController {
	return &nginxConfigServiceController{
		service:       newNginxConfigServiceWithState(info),
		expectedState: Unknown,
		waitGroup:     new(sync.WaitGroup),
	}
}
