package bifrost

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"github.com/apsdehal/go-logger"
	"io"
	"io/ioutil"
	"os"
	"time"
)

// readFile, 读取文件函数
// 参数:
//     path: 文件路径字符串
// 返回值:
//     文件数据
//     错误
func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

// ngLoad, nginx配置加载函数，根据nginx配置文件信息加载nginx配置并记录文件基准信息
// 参数:
//     appConfig: nginx配置文件信息对象指针
// 返回值:
//     nginx配置对象指针
//     错误
func ngLoad(appConfig *NGConfig) (*resolv.Config, error) {
	// 加载nginx配置
	ng, err := resolv.Load(appConfig.ConfPath)
	if err != nil {
		return nil, err
	}

	// 获取nginx配置文件列表
	confList, listErr := ng.List()
	if listErr != nil {
		return nil, listErr
	}

	// 计算nginx配置文件哈希基准值
	appConfig.confHash = make(map[string]string, 0)
	for _, path := range confList {
		cHash, hashErr := hash(path)
		if hashErr != nil {
			return nil, hashErr
		}

		appConfig.confHash[path] = cHash
	}

	return ng, nil
}

// hash, 计算文件hash值函数
// 参数:
//     path: 文件路径
// 返回值:
//     文件哈希基准值
//     错误
func hash(path string) (string, error) {
	// 读取文件
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 计算文件数据哈希值
	hash256 := sha256.New()
	_, hashCPErr := io.Copy(hash256, f)
	if hashCPErr != nil {
		return "", hashCPErr
	}
	return hex.EncodeToString(hash256.Sum(nil)), nil
}

// Bak, nginx配置文件备份函数
// 参数:
//     appConfig: nginx配置文件信息对象
//     ngConfig: nginx配置对象指针
//     c: 整型管道，用于停止备份
func Bak(appConfig NGConfig, ngConfig *resolv.Config, c chan int) {
	for {
		select {
		case <-time.NewTicker(5 * time.Minute).C: // 每5分钟定时执行备份操作
			bak(appConfig, ngConfig)
		case signal := <-c: // 获取管道传入信号
			if signal == 9 { // 为9时，停止备份
				Log(NOTICE, fmt.Sprintf("[%s] Nginx Config backup is stop.", appConfig.Name))
				break
			}
		}
	}
}

// bak, nginx配置文件备份子函数
// 参数:
//     appConfig: nginx配置文件信息对象
//     ngConfig: nginx配置对象指针
func bak(appConfig NGConfig, ngConfig *resolv.Config) {
	// 初始化备份文件名
	bakDate := time.Now().Format("20060102")
	bakName := fmt.Sprintf("nginx.conf.%s.tgz", bakDate)

	bakPath, bErr := resolv.Backup(ngConfig, bakName)
	if bErr != nil && !os.IsExist(bErr) { // 备份失败
		Log(CRITICAL, fmt.Sprintf("[%s] Nginx Config backup to %s, but failed. <%s>", appConfig.Name, bakPath, bErr))
		Log(NOTICE, fmt.Sprintf("[%s] Nginx Config backup is stop.", appConfig.Name))
	} else if bErr == nil { // 备份成功
		Log(NOTICE, fmt.Sprintf("[%s] Nginx Config backup to %s", appConfig.Name, bakPath))
	}

}

// AutoReload, nginx配置文件自动热加载函数
// 参数:
//     appConfig: nginx配置文件信息对象指针
//     ngConfig: nginx配置对象指针
//     c: 整型管道，用于停止备份
func AutoReload(appConfig *NGConfig, ngConfig *resolv.Config, c chan int) {
	for {
		select {
		case <-time.NewTicker(30 * time.Second).C: // 每30秒检查一次nginx配置文件是否已在后台更新
			cache, err := autoReload(appConfig)
			if err != nil {
				Log(WARN, fmt.Sprintf("[%s] Nginx Config reload failed, cased by '%s'", appConfig.Name, err))
			} else if cache != nil {
				ngConfig.Value = cache.Value
				ngConfig.Children = cache.Children
				Log(INFO, fmt.Sprintf("[%s] Nginx Config reload successfully", appConfig.Name))
			}
		case signal := <-c: // 获取管道传入信号
			if signal == 9 { // 为9时，停止备份
				Log(NOTICE, fmt.Sprintf("[%s] Nginx Config backup is stop.", appConfig.Name))
				break
			}
		}
	}
}

// autoReload, nginx配置文件自动热加载子函数
// 参数:
//     appConfig: nginx配置文件信息对象指针
//     ngConfig: nginx配置对象指针
func autoReload(appConfig *NGConfig) (*resolv.Config, error) {
	// 校验配置文件是否更新
	isDifferent, checkErr := hashCheck(appConfig)
	if checkErr != nil {
		return nil, checkErr
	}

	// 如果有差别，则重新读取配置
	if isDifferent {
		return ngLoad(appConfig)
	}
	return nil, nil
}

// hashCheck, nginx配置文件是否已更改校验函数
// 参数:
//     appConfig: nginx配置文件信息对象指针
func hashCheck(appConfig *NGConfig) (bool, error) {
	for path, h := range appConfig.confHash {
		// 获取当前配置文件哈希值
		checkHash, checkErr := hash(path)
		if checkErr != nil {
			return false, checkErr
		}

		// 校验哈希，不一致说明需要更新
		if checkHash != h {
			return true, nil
		}

	}
	return false, nil
}

// PathExists, 判断文件路径是否存在函数
// 返回值:
//     true: 存在; false: 不存在
//     错误
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	} else {
		return false, nil
	}
}

// Log, 日志记录函数
// 参数:
//     level: 日志级别对象
//     message: 需记录的日志信息字符串
func Log(level logger.LogLevel, message string) {

	myLogger.Log(level, message)
	//fmt.Printf("[%s] [%s] %s\n", level, time.Now().Format(timeFormat), message)

}
