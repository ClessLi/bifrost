package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func RemoveFiles(files []string) error {
	for _, path := range files {
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

// CheckBackups, 检查归档目录下归档文件是否需要清理及是否可以进行归档操作的函数
//
// 参数:
//     name: 归档文件前缀名
//     dir: 归档文件目录路径
//     saveTime: 归档文件保存时间，单位天
//     cycle: 归档操作周期，单位天
//     now: 当前检查时间
// 返回值:
//     true: 需要归档操作; false: 不需要归档
//     错误
func CheckBackups(name, backupDir string, retentionTime, backupCycleTime int, now time.Time) (bool, error) {
	needBackup := true
	saveDate := now.Add(-24 * time.Hour * time.Duration(retentionTime))
	cycleDate := now.Add(-24 * time.Hour * time.Duration(backupCycleTime))
	bakFilePattern := fmt.Sprintf(`^%s\.(\d{8})\.tgz$`, name)
	bakFileReg := regexp.MustCompile(bakFilePattern)

	baks, gErr := filepath.Glob(filepath.Join(backupDir, fmt.Sprintf("%s.*.tgz", name)))
	if gErr != nil {
		return false, gErr
	}

	for i := 0; i < len(baks) && needBackup; i++ {
		bakName := filepath.Base(baks[i])
		if isBak := bakFileReg.MatchString(bakName); isBak {
			bakDate, tpErr := time.ParseInLocation("20060102", bakFileReg.FindStringSubmatch(bakName)[1], now.Location())
			if tpErr != nil {
				return false, tpErr
			}

			if bakDate.Unix() < saveDate.Unix() {
				rmErr := os.Remove(baks[i])
				if rmErr != nil {
					return false, rmErr
				}
			}

			if bakDate.Unix() > cycleDate.Unix() || bakDate.Format("20060102") == now.Format("20060102") {
				needBackup = false
			}

		}
	}

	return needBackup, nil
}
