package bifrost

import (
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"os"
	"sync"
	"time"
)

var ServiceNotAvailable = fmt.Errorf("bifrost service not available")

// WebServerType, web服务器类型对象，定义web服务器所属类型
type WebServerType string

// Service, nginx配置文件信息结构体，定义配置文件路径、nginx可执行文件路径和bifrost为其提供接口的路由及侦听端口
type ServiceInfo struct {
	Name           string        `yaml:"name"`
	Type           WebServerType `yaml:"type"`
	BaseURI        string        `yaml:"baseURI"`
	BackupCycle    int           `yaml:"backupCycle"`
	BackupSaveTime int           `yaml:"backupSaveTime"`
	BackupDir      string        `yaml:"backupDir,omitempty"`
	ConfPath       string        `yaml:"confPath"`
	VerifyExecPath string        `yaml:"verifyExecPath"`
	confCaches     nginx.Caches
	nginxConfig    *nginx.Config
	available      bool
	bakChan        chan int
	autoReloadChan chan int
}

func (i *ServiceInfo) NgLoad() error {
	if i.available {
		return i.ngLoad()
	}
	return ServiceNotAvailable
}

// ngLoad, ServerInfo的nginx配置加载方法，根据nginx配置文件信息加载nginx配置并记录文件基准信息
// 返回值:
//     nginx配置对象指针
//     错误
func (i *ServiceInfo) ngLoad() error {
	// 加载nginx配置并获取缓存
	Log(DEBUG, "[%s] load config...", i.Name)
	path, caches, err := nginx.Load(i.ConfPath)
	if err != nil {
		Log(DEBUG, "[%s] load config failed: %s", i.Name, err.Error())
		return err
	}

	// 记录缓存
	i.confCaches = caches
	i.ConfPath = path
	i.nginxConfig, err = i.confCaches.GetConfig(i.ConfPath)
	if err != nil {
		Log(DEBUG, "[%s] load config failed: %s", i.Name, err.Error())
		return err
	}
	Log(DEBUG, "[%s] load config success", i.Name)

	return nil

}

// Bak, ServerInfo的nginx配置文件备份方法
// 参数:
//     c: 整型管道，用于停止备份
func (i *ServiceInfo) Bak(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if !i.available {
			Log(INFO, "[%s] %s, Nginx Config backup is stop.", i.Name, ServiceNotAvailable)
		} else {
			Log(NOTICE, "[%s] Nginx Config backup is stop.", i.Name)
		}
	}()
	i.bakChan = make(chan int, 1)
	for i.available {
		select {
		case <-time.NewTicker(5 * time.Minute).C: // 每5分钟定时执行备份操作
			Log(DEBUG, "[%s] Nginx Config check and backup", i.Name)
			i.bak()
		case sig := <-i.bakChan: // 获取管道传入信号
			if sig == 9 { // 为9时，停止备份
				goto stopHere
			}
		}
	}
stopHere:
	return
}

// bak, ServerInfo的nginx配置文件备份子方法
// 参数:
func (i *ServiceInfo) bak() {
	if i.available {
		bakPath, bErr := nginx.Backup(i.nginxConfig, "nginx.conf", i.BackupSaveTime, i.BackupCycle, i.BackupDir)

		if bErr != nil && (!os.IsExist(bErr) && bErr != nginx.NoBackupRequired) { // 备份失败
			Log(CRITICAL, "[%s] Nginx Config backup to %s, but failed. <%s>", i.Name, bakPath, bErr)
			Log(NOTICE, "[%s] Nginx Config backup is stop.", i.Name)
		} else if bErr == nil { // 备份成功
			Log(INFO, "[%s] Nginx Config backup to %s", i.Name, bakPath)
		}
	}
}

// AutoReload, ServerInfo的web服务器配置文件自动热加载方法
// 参数:
//     c: 整型管道，用于停止备份
func (i *ServiceInfo) AutoReload(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if !i.available {
			Log(INFO, "[%s] %s, Nginx Config auto reload is stop.", ServiceNotAvailable, i.Name)
		} else {
			Log(NOTICE, "[%s] Nginx Config auto reload is stop.", i.Name)
		}
	}()
	i.autoReloadChan = make(chan int, 1)
	for i.available {
		select {
		case <-time.NewTicker(30 * time.Second).C: // 每30秒检查一次nginx配置文件是否已在后台更新
			Log(DEBUG, "[%s] Nginx Config check and reloading", i.Name)
			reloadErr := i.autoReload()
			if reloadErr != nil && reloadErr != nginx.NoReloadRequired {
				Log(WARN, "[%s] Nginx Config reload failed, cased by '%s'", i.Name, reloadErr)
			} else if reloadErr == nil {
				Log(INFO, "[%s] Nginx Config reload successfully", i.Name)
			}
		case sig := <-i.autoReloadChan: // 获取管道传入信号
			if sig == 9 { // 为9时，停止备份
				goto stopHere
			}
		}
	}
stopHere:
	return
}

// autoReload, ServerInfo的web服务器配置文件自动热加载子方法
func (i *ServiceInfo) autoReload() error {
	if i.available {
		// 校验配置文件是否更新
		isSame, checkErr := i.checkHash()
		if checkErr != nil {
			return checkErr
		}

		// 如果有差别，则重新读取配置
		if !isSame {
			Log(DEBUG, "[%s] reloading nginx config", i.Name)
			return i.ngLoad()
		}
		return nginx.NoReloadRequired
	}
	return ServiceNotAvailable
}

// checkHash, ServerInfo的web服务器配置文件是否已更改校验方法
func (i ServiceInfo) checkHash() (isSame bool, err error) {
	isSame = true
	for path := range i.confCaches {
		if isSame, err = i.confCaches.CheckHash(path); !isSame {
			return
		}
	}
	return
}

func (i *ServiceInfo) Enable() {
	i.available = true
	Log(INFO, "[%s] bifrost service enabled", i.Name)
}

func (i *ServiceInfo) Disable() {
	i.available = false
	Log(INFO, "[%s] bifrost service disabled", i.Name)
}
