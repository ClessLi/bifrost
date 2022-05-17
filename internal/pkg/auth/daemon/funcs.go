package daemon

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/apsdehal/go-logger"

	"github.com/yongPhone/bifrost/internal/pkg/auth/service"
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

func getProc(path string) (*os.Process, error) {
	pid, pidErr := getPid(path)
	if pidErr != nil {
		return nil, pidErr
	}
	return os.FindProcess(pid)
}

func rmPidFile(path string) {
	rmPidFileErr := os.Remove(path)
	if rmPidFileErr != nil {
		Log(ERROR, rmPidFileErr.Error())
	}
	Log(NOTICE, "bifrost.pid has been removed.")
}

// getPid, 查询pid文件并返回pid
// 返回值:
//     pid
//     错误
func getPid(path string) (int, error) {
	// 判断pid文件是否存在
	if _, err := os.Stat(path); err == nil || os.IsExist(err) { // 存在
		// 读取pid文件
		pidBytes, readPidErr := readFile(path)
		if readPidErr != nil {
			Log(ERROR, readPidErr.Error())
			return -1, readPidErr
		}

		// 去除pid后边的换行符
		pidBytes = bytes.TrimRight(pidBytes, "\n")

		// 转码pid
		pid, toIntErr := strconv.Atoi(string(pidBytes))
		if toIntErr != nil {
			Log(ERROR, toIntErr.Error())
			return -1, toIntErr
		}

		return pid, nil
	} else { // 不存在
		return -1, procStatusNotRunning
	}
}

// configCheck, 检查bifrost配置项是否完整
// 返回值:
//     错误
func configCheck() error {
	// 初始化认证数据库或认证配置信息
	if AuthConf.AuthService == nil {
		AuthConf.AuthService = new(service.AuthService)
	}
	if AuthConf.AuthService.Port == 0 {
		AuthConf.AuthService.Port = 12320
	}
	if AuthConf.AuthService.AuthDBConfig == nil && AuthConf.AuthService.AuthConfig == nil {
		AuthConf.AuthService.AuthConfig = &service.AuthConfig{Username: "heimdall", Password: "Bultgang"}
	}
	return nil
}
