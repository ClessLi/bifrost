package bifrost

import (
	"fmt"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"github.com/apsdehal/go-logger"
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

// Bak, nginx配置文件备份管理函数
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

// bak, nginx配置文件备份函数
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
