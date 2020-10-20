package bifrost

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"os"
	"time"
)

// ServerInfo, nginx配置文件信息结构体，定义配置文件路径、nginx可执行文件路径和bifrost为其提供接口的路由及侦听端口
type ServerInfo struct {
	Name           string        `yaml:"name"`
	ServerType     WebServerType `yaml:"serverType"`
	BaseURI        string        `yaml:"baseURI"`
	BackupCycle    int           `yaml:"backupCycle"`
	BackupSaveTime int           `yaml:"backupSaveTime"`
	BackupDir      string        `yaml:"backupDir,omitempty"`
	ConfPath       string        `yaml:"confPath"`
	VerifyExecPath string        `yaml:"verifyExecPath"`
	confCaches     nginx.Caches
	nginxConfig    *nginx.Config
}

// ngLoad, ServerInfo的nginx配置加载方法，根据nginx配置文件信息加载nginx配置并记录文件基准信息
// 返回值:
//     nginx配置对象指针
//     错误
func (si *ServerInfo) ngLoad() error {
	// 加载nginx配置并获取缓存
	path, caches, err := nginx.Load(si.ConfPath)
	if err != nil {
		return err
	}

	// 记录缓存
	si.confCaches = caches
	si.ConfPath = path
	si.nginxConfig, err = si.confCaches.GetConfig(si.ConfPath)
	if err != nil {
		return err
	}

	return nil

}

// Bak, ServerInfo的nginx配置文件备份方法
// 参数:
//     c: 整型管道，用于停止备份
func (si ServerInfo) Bak(signal chan int) {
	for {
		select {
		case <-time.NewTicker(5 * time.Minute).C: // 每5分钟定时执行备份操作
			si.bak()
		case s := <-signal: // 获取管道传入信号
			if s == 9 { // 为9时，停止备份
				Log(NOTICE, "[%s] Nginx Config backup is stop.", si.Name)
				break
			}
		}
	}
}

// bak, ServerInfo的nginx配置文件备份子方法
// 参数:
func (si ServerInfo) bak() {
	config, confErr := si.confCaches.GetConfig(si.ConfPath)
	if confErr != nil {
		Log(CRITICAL, "[%s] Nginx Config backup failed, cased by %s.", si.Name, confErr)
		return
	}

	bakPath, bErr := nginx.Backup(config, "nginx.conf", si.BackupSaveTime, si.BackupCycle, si.BackupDir)

	if bErr != nil && (!os.IsExist(bErr) && bErr != nginx.NoBackupRequired) { // 备份失败
		Log(CRITICAL, "[%s] Nginx Config backup to %s, but failed. <%s>", si.Name, bakPath, bErr)
		Log(NOTICE, "[%s] Nginx Config backup is stop.", si.Name)
	} else if bErr == nil { // 备份成功
		Log(NOTICE, "[%s] Nginx Config backup to %s", si.Name, bakPath)
	}

}

// AutoReload, ServerInfo的web服务器配置文件自动热加载方法
// 参数:
//     c: 整型管道，用于停止备份
func (si *ServerInfo) AutoReload(signal chan int) {
	for {
		select {
		case <-time.NewTicker(30 * time.Second).C: // 每30秒检查一次nginx配置文件是否已在后台更新
			reloadErr := si.autoReload()
			if reloadErr != nil && reloadErr != nginx.NoReloadRequired {
				Log(WARN, "[%s] Nginx Config reload failed, cased by '%s'", si.Name, reloadErr)
			} else if reloadErr == nil {
				Log(INFO, "[%s] Nginx Config reload successfully", si.Name)
			}
		case s := <-signal: // 获取管道传入信号
			if s == 9 { // 为9时，停止备份
				Log(NOTICE, "[%s] Nginx Config backup is stop.", si.Name)
				break
			}
		}
	}
}

// autoReload, ServerInfo的web服务器配置文件自动热加载子方法
func (si *ServerInfo) autoReload() error {
	// 校验配置文件是否更新
	isSame, checkErr := si.checkHash()
	if checkErr != nil {
		return checkErr
	}

	// 如果有差别，则重新读取配置
	if !isSame {
		Log(DEBUG, "[%s] reloading nginx config", si.Name)
		return si.ngLoad()
	}
	return nginx.NoReloadRequired
}

// checkHash, ServerInfo的web服务器配置文件是否已更改校验方法
func (si ServerInfo) checkHash() (isSame bool, err error) {
	isSame = true
	for path := range si.confCaches {
		if isSame, err = si.confCaches.CheckHash(path); !isSame {
			return
		}
	}
	return
}
