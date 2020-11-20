package service

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
)

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
			//Log(ERROR, readPidErr.Error())
			return -1, readPidErr
		}

		// 去除pid后边的换行符
		pidBytes = bytes.TrimRight(pidBytes, "\n")

		// 转码pid
		pid, toIntErr := strconv.Atoi(string(pidBytes))
		if toIntErr != nil {
			//Log(ERROR, toIntErr.Error())
			return -1, toIntErr
		}

		return pid, nil
	} else { // 不存在
		return -1, ErrProcessNotRunning
	}
}

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
