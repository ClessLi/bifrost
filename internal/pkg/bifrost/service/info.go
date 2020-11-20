package service

import (
	"errors"
	"fmt"
	ngJson "github.com/ClessLi/bifrost/pkg/json/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	UpdateSuccess = []byte("update config success")
	//ErrMainConfigPath      = errors.New("main config path error")
	//ErrLocation            = errors.New("location inconsistency, plz check")
	ErrValidationNotExist = errors.New("the validation process does not exist or is configured incorrectly")
	ErrEmptyConfig        = errors.New("empty configuration error")
	//ErrConfigUnmarshal     = errors.New("config unmarshal error")
	ErrServiceNotAvailable = fmt.Errorf("bifrost service not available")
)

// WebServerType, web服务器类型对象，定义web服务器所属类型
type WebServerType string

// Service, nginx配置文件信息结构体，定义配置文件路径、nginx可执行文件路径和bifrost为其提供接口的路由及侦听端口
type Info struct {
	Name           string        `yaml:"name"`
	Type           WebServerType `yaml:"type"`
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

func (i *Info) update(jsonData []byte) (data []byte, err error) {
	//fmt.Println("获取web服务配置校验二进制文件路径")
	verifyBin, err := filepath.Abs(i.VerifyExecPath)
	if err != nil {
		//Log(CRITICAL, "[%s] %s detailed error: %s", info.Name, bifrostpb.ValidationNotExist, err)
		err = ErrValidationNotExist
		return nil, err
	}

	// check config
	if len(jsonData) > 0 {
		newConfig, err := ngJson.Unmarshal(jsonData)
		if err != nil {
			//Log(DEBUG, "[%s] [%s] buffer detail: %s", info.Name, ip, buffer.String())
			//Log(WARN, "[%s] [%s] %s detailed error: %s", info.Name, ip, bifrostpb.ConfigUnmarshalErr, err)
			return nil, err
		}

		// delete old config

		err = nginx.Delete(i.nginxConfig)
		//message := ""
		if err != nil {
			//message = fmt.Sprintf("Delete nginx ng failed. <%s>", err)
			//Log(ERROR, "[%s] [%s] %s", info.Name, ip, message)
			return nil, err
		}

		//Log(INFO, "[%s] Deleted old nginx config.", info.Name)
		//Log(INFO, "[%s] Verify new nginx config.", info.Name)
		newCaches, err := nginx.SaveWithCheck(newConfig, verifyBin)
		// roll back
		if err != nil {
			//Log(DEBUG, "[%s] Roll back to old nginx config.", info.Name)
			//message = fmt.Sprintf("Nginx ng verify failed. <%s>", err)
			//Log(WARN, "[%s] %s", info.Name, message)

			//Log(INFO, "[%s] Delete new nginx ng.", info.Name)
			var rollErr error
			rollErr = nginx.Delete(newConfig)
			if rollErr != nil {
				//Log(ERROR, "[%s] Delete new nginx ng failed. <%s>", info.Name, err)
				//message = "New nginx config verify failed. And delete new nginx config failed."
				return nil, rollErr
			}

			//Log(INFO, "[%s] Rollback nginx ng.", info.Name)
			_, rollErr = nginx.Save(i.nginxConfig)
			if rollErr != nil {
				//Log(CRITICAL, "[%s] Nginx ng rollback failed. <%s>", info.Name, err)
				//message = "New nginx config verify failed. And nginx config rollback failed."
				return nil, rollErr
			}

			return nil, err
		}
		i.confCaches = newCaches
		i.nginxConfig = newConfig
		i.ConfPath = newConfig.Value

		//Log(NOTICE, "[%s] [%s] Nginx Config saved successfully", info.Name, ip)
		return UpdateSuccess, nil
	} else {
		//Log(WARN, "[%s] [%s] %s", info.Name, ip, bifrostpb.EmptyConfigErr)
		err = ErrEmptyConfig
		return nil, err
	}
}

func (i *Info) NgLoad() error {
	if i.available {
		return i.ngLoad()
	}
	return ErrServiceNotAvailable
}

// ngLoad, ServerInfo的nginx配置加载方法，根据nginx配置文件信息加载nginx配置并记录文件基准信息
// 返回值:
//     nginx配置对象指针
//     错误
func (i *Info) ngLoad() error {
	// 加载nginx配置并获取缓存
	//Log(DEBUG, "[%s] load config...", i.Name)
	path, caches, err := nginx.Load(i.ConfPath)
	if err != nil {
		//Log(DEBUG, "[%s] load config failed: %s", i.Name, err.Error())
		return err
	}

	// 记录缓存
	i.confCaches = caches
	i.ConfPath = path
	i.nginxConfig, err = i.confCaches.GetConfig(i.ConfPath)
	if err != nil {
		//Log(DEBUG, "[%s] load config failed: %s", i.Name, err.Error())
		return err
	}
	//Log(DEBUG, "[%s] load config success", i.Name)

	return nil

}

// Bak, ServerInfo的nginx配置文件备份方法
// 参数:
//     c: 整型管道，用于停止备份
func (i *Info) Bak(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if !i.available {
			//Log(INFO, "[%s] %s, Nginx Config backup is stop.", i.Name, ErrServiceNotAvailable)
		} else {
			//Log(NOTICE, "[%s] Nginx Config backup is stop.", i.Name)
		}
	}()
	i.bakChan = make(chan int, 1)
	for i.available {
		select {
		case <-time.NewTicker(5 * time.Minute).C: // 每5分钟定时执行备份操作
			//Log(DEBUG, "[%s] Nginx Config check and backup", i.Name)
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
func (i *Info) bak() {
	if i.available {
		bakPath, bErr := nginx.Backup(i.nginxConfig, "nginx.conf", i.BackupSaveTime, i.BackupCycle, i.BackupDir)

		if bErr != nil && (!os.IsExist(bErr) && bErr != nginx.NoBackupRequired) { // 备份失败
			//Log(CRITICAL, "[%s] Nginx Config backup to %s, but failed. <%s>", i.Name, bakPath, bErr)
			fmt.Printf("[%s] Nginx Config backup to %s, but failed. <%s>", i.Name, bakPath, bErr)
			//Log(NOTICE, "[%s] Nginx Config backup is stop.", i.Name)
		} else if bErr == nil { // 备份成功
			//Log(INFO, "[%s] Nginx Config backup to %s", i.Name, bakPath)
			fmt.Printf("[%s] Nginx Config backup to %s", i.Name, bakPath)
		}
	}
}

// AutoReload, ServerInfo的web服务器配置文件自动热加载方法
// 参数:
//     c: 整型管道，用于停止备份
func (i *Info) AutoReload(wg *sync.WaitGroup) {
	defer wg.Done()
	//defer func() {
	//	if !i.available {
	//		Log(INFO, "[%s] %s, Nginx Config auto reload is stop.", ErrServiceNotAvailable, i.Name)
	//	} else {
	//		Log(NOTICE, "[%s] Nginx Config auto reload is stop.", i.Name)
	//	}
	//}()
	i.autoReloadChan = make(chan int, 1)
	for i.available {
		select {
		case <-time.NewTicker(30 * time.Second).C: // 每30秒检查一次nginx配置文件是否已在后台更新
			//Log(DEBUG, "[%s] Nginx Config check and reloading", i.Name)
			reloadErr := i.autoReload()
			if reloadErr != nil && reloadErr != nginx.NoReloadRequired {
				//Log(WARN, "[%s] Nginx Config reload failed, cased by '%s'", i.Name, reloadErr)
			} else if reloadErr == nil {
				//Log(INFO, "[%s] Nginx Config reload successfully", i.Name)
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
func (i *Info) autoReload() error {
	if i.available {
		// 校验配置文件是否更新
		isSame, checkErr := i.checkHash()
		if checkErr != nil {
			return checkErr
		}

		// 如果有差别，则重新读取配置
		if !isSame {
			//Log(DEBUG, "[%s] reloading nginx config", i.Name)
			return i.ngLoad()
		}
		return nginx.NoReloadRequired
	}
	return ErrServiceNotAvailable
}

// checkHash, ServerInfo的web服务器配置文件是否已更改校验方法
func (i Info) checkHash() (isSame bool, err error) {
	isSame = true
	for path := range i.confCaches {
		if isSame, err = i.confCaches.CheckHash(path); !isSame {
			return
		}
	}
	return
}

func (i *Info) Enable() {
	i.available = true
	//Log(INFO, "[%s] bifrost service enabled", i.Name)
}

func (i *Info) Disable() {
	i.available = false
	//Log(INFO, "[%s] bifrost service disabled", i.Name)
}
