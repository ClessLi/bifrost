package utils

import (
	"bytes"
	"fmt"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/marmotedu/errors"
)

const backupDateLayout = `20060102`

func RemoveFiles(files []string) error {
	for _, path := range files {
		// 判断是否为单元测试
		if len(os.Args) > 3 && os.Args[1] == "-test.v" && os.Args[2] == "-test.run" {
			fmt.Printf("remove: %s\n", path)
			continue
		}
		err := os.Remove(path)
		if err != nil {
			return err
		}
		/*// debug test
		fmt.Printf("remove: %s\n", path)
		// debug test end*/
	}
	return nil
}

func getBackupFileRegexp(backupPrefix string) *regexp.Regexp {
	bakFilePattern := `^` + backupPrefix + `\.(\d{8})\.tgz$`
	return regexp.MustCompile(bakFilePattern)
}

func GetBackupFileName(backupPrefix string, now time.Time) string {
	dt := now.Format(backupDateLayout)
	return backupPrefix + "." + dt + ".tgz"
}

// CheckAndCleanBackups 检查归档目录下归档文件是否需要清理及是否可以进行归档操作的函数
//
// 参数:
//
//	backupPrefix: 归档文件前缀名
//	backupDir: 归档文件目录路径
//	retentionTime: 归档文件保存时间，单位天
//	backupCycleTime: 归档操作周期，单位天
//	now: 当前检查时间
//
// 返回值:
//
//	true: 需要归档操作; false: 不需要归档
//	错误
func CheckAndCleanBackups(
	backupPrefix, backupDir string,
	retentionTime, backupCycleTime int,
	now time.Time,
) (bool, error) {
	needBackup := true
	saveDate := now.Add(-24 * time.Hour * time.Duration(retentionTime))
	cycleDate := now.Add(-24 * time.Hour * time.Duration(backupCycleTime))
	if len(strings.TrimSpace(backupPrefix)) == 0 {
		backupPrefix = "nginx.conf"
	} else {
		backupPrefix = strings.TrimSpace(backupPrefix)
	}
	bakFileReg := getBackupFileRegexp(backupPrefix)

	baks, gErr := filepath.Glob(filepath.Join(backupDir, backupPrefix+".*.tgz"))
	if gErr != nil {
		return false, gErr
	}

	for i := 0; i < len(baks) && needBackup; i++ {
		bakName := filepath.Base(baks[i])
		if isBak := bakFileReg.MatchString(bakName); isBak {
			bakDate, tpErr := time.ParseInLocation(
				backupDateLayout,
				bakFileReg.FindStringSubmatch(bakName)[1],
				now.Location(),
			)
			if tpErr != nil {
				return false, errors.Wrapf(tpErr, "failed to resolve archive name '%s'", baks[i])
			}

			// 判断是否需要清理，并清理过期归档
			if bakDate.Unix() < saveDate.Unix() {
				logV1.Infof("cleaning up expired archive '%s'", baks[i])
				rmErr := os.Remove(baks[i])
				if rmErr != nil {
					return false, errors.Wrapf(rmErr, "failed to clean up expired archive '%s'", baks[i])
				}
				logV1.Infof("successfully cleaned up expired archive '%s'", baks[i])
			}

			// 判断该归档是否是最新归档，是反馈不需归档，并退出循环
			if bakDate.Unix() > cycleDate.Unix() || bakDate.Format(backupDateLayout) == now.Format(backupDateLayout) {
				needBackup = false
			}
		}
	}

	return needBackup, nil
}

// GetPid, 查询pid文件并返回pid
// 返回值:
//
//	pid
//	错误
func GetPid(path string) (int, error) {
	// 判断pid文件是否存在
	if _, err := os.Stat(path); err == nil || os.IsExist(err) { // 存在
		// 读取pid文件
		pidBytes, readPidErr := ReadFile(path)
		if readPidErr != nil {
			// Log(ERROR, readPidErr.Error())
			return -1, readPidErr
		}

		// 去除pid后边的换行符
		pidBytes = bytes.TrimRight(pidBytes, "\n")

		// 转码pid
		pid, toIntErr := strconv.Atoi(string(pidBytes))
		if toIntErr != nil {
			// Log(ERROR, toIntErr.Error())
			return -1, toIntErr
		}

		return pid, nil
	} else { // 不存在
		return -1, errors.New("process is not running")
	}
}

// ReadFile, 读取文件函数
// 参数:
//
//	path: 文件路径字符串
//
// 返回值:
//
//	文件数据
//	错误
func ReadFile(path string) ([]byte, error) {
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
