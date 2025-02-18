package nginx

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Backup, 文件归档函数
//
// 参数:
//
//	config: nginx配置文件对象指针
//	name: 归档文件前缀名
//	saveTime: 归档文件保存时间，单位天
//	bakCycle: 归档操作周期，单位天
//	backupDir: 可选参数，备份目录
//
// 返回值:
//
//	bakPath: 归档文件路径
//	err: 错误
func Backup(config *Config, name string, saveTime, bakCycle int, backupDir ...string) (bakPath string, err error) {
	// 归档文件名修饰
	if name != "" {
		name = strings.TrimSpace(name)
		name = strings.TrimRight(name, ".")
	} else {
		name = "nginx.conf"
	}

	// 归档日期初始化
	now := time.Now().In(TZ)
	dt := now.Format("20060102")

	// 归档目录
	isSpecBakDir := false
	if len(backupDir) > 0 && backupDir[0] != "" {
		d, statErr := os.Stat(backupDir[0])
		if statErr == nil {
			isSpecBakDir = d.IsDir()
		}
	}
	confDir := filepath.Dir(config.Value)
	bakName := fmt.Sprintf("%s.%s.tgz", name, dt)
	bakDir := confDir
	bakPath = filepath.Join(confDir, bakName)
	bakPath, err = filepath.Abs(bakPath)
	if err != nil {
		return "", err
	}
	specBakPath := bakPath
	if isSpecBakDir {
		bakDir = backupDir[0]
		specBakPath, err = filepath.Abs(filepath.Join(backupDir[0], bakName))
		if err != nil {
			return "", err
		}
	}

	// 清理过期归档，检查已归档文件
	needBak, err := checkBackups(name, bakDir, saveTime, bakCycle, now)
	if err != nil {
		return "", err
	}

	if needBak {
		caches, err := config.List()
		if err != nil {
			return "", err
		}
		err = tgz(bakPath, caches)
		if err == nil && isSpecBakDir {
			err = os.Rename(bakPath, specBakPath)
			bakPath = specBakPath
		}
	} else {
		err = NoBackupRequired
	}

	return
}

// checkBackups, 检查归档目录下归档文件是否需要清理及是否可以进行归档操作的函数
//
// 参数:
//
//	name: 归档文件前缀名
//	dir: 归档文件目录路径
//	saveTime: 归档文件保存时间，单位天
//	cycle: 归档操作周期，单位天
//	now: 当前检查时间
//
// 返回值:
//
//	true: 需要归档操作; false: 不需要归档
//	错误
func checkBackups(name, dir string, saveTime, cycle int, now time.Time) (bool, error) {
	needBackup := true
	saveDate := now.Add(-24 * time.Hour * time.Duration(saveTime))
	cycleDate := now.Add(-24 * time.Hour * time.Duration(cycle))
	// fmt.Printf("save date now is %s, time is '%s'\n", saveDate.Format("20060102"), saveDate)
	// fmt.Printf("cycle date now is %s, time is '%s'\n", cycleDate.Format("20060102"), cycleDate)
	bakFilePattern := fmt.Sprintf(`^%s\.(\d{8})\.tgz$`, name)
	bakFileReg := regexp.MustCompile(bakFilePattern)

	baks, gErr := filepath.Glob(filepath.Join(dir, fmt.Sprintf("%s.*.tgz", name)))
	if gErr != nil {
		return false, gErr
	}

	// fmt.Printf("date now is %s, time is '%s'\n\n", now.Format("20060102"), now)
	for i := 0; i < len(baks) && needBackup; i++ {
		bakName := filepath.Base(baks[i])
		if isBak := bakFileReg.MatchString(bakName); isBak {
			bakDate, tpErr := time.ParseInLocation("20060102", bakFileReg.FindStringSubmatch(bakName)[1], TZ)
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

			// fmt.Printf("bakDate is %s, time is '%s'\n", bakDate.Format("20060102"), bakDate)
		}
	}

	return needBackup, nil
}

// tgz, 归档操作函数
//
// 参数:
//
//	dest: 归档文件路径
//	caches: 配置文件对象缓存
//	list: 被归档（仅）文件（非目录）切片
//
// 返回值:
//
//	错误
func tgz(dest string, caches Caches) (err error) {
	// 归档目录
	destDir := filepath.Dir(dest)

	// 归档文件初始化
	d, _ := os.Create(dest) // 创建归档文件
	defer d.Close()
	gw := gzip.NewWriter(d) // 转换归档文件为gzip压缩文件
	defer gw.Close()
	tgzw := tar.NewWriter(gw) // 转换归档文件为tar.gz/tgz归档文件
	defer tgzw.Close()

	// 开始归档
	for f := range caches {
		// 初始化被归档文件信息
		relDir, relErr := filepath.Rel(destDir, filepath.Dir(f))
		if relErr != nil {
			return relErr
		}
		fd, opErr := os.Open(f)
		if opErr != nil {
			return opErr
		}

		// 开始归档
		compErr := compress(fd, relDir, tgzw)
		if compErr != nil {
			return compErr
		}
	}
	return
}

// compress, 归档压缩子函数
//
// 参数:
//
//	fd: 被归档文件的系统文件对象指针
//	prefix: 被归档文件的目录路径
//	tgzw: tar文件对象指针
//
// 返回值:
//
//	错误
func compress(fd *os.File, prefix string, tgzw *tar.Writer) error {
	defer fd.Close()

	// 加载被归档文件信息
	info, infoErr := fd.Stat()
	if infoErr != nil {
		return infoErr
	}

	// 文件归档初始化
	header, hErr := tar.FileInfoHeader(info, "")
	if hErr != nil {
		return hErr
	}
	// 调整被归档文件目录信息
	header.Name = filepath.Join(prefix, header.Name)
	// 完成该文件归档
	twErr := tgzw.WriteHeader(header)
	if twErr != nil {
		return twErr
	}
	_, ioErr := io.Copy(tgzw, fd)
	fd.Close()
	if ioErr != nil {
		return ioErr
	}

	return nil
}
