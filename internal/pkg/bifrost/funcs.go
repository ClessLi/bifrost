package bifrost

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"io/ioutil"
	"os"
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

// PathExists, 判断文件路径是否存在函数
// 参数:
//     path: 待判断的文件路径字符串
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
func Log(level logger.LogLevel, message string, a ...interface{}) {
	myLogger.Log(level, fmt.Sprintf(message, a...))
}
