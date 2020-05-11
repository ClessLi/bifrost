package ng_conf_admin

import (
	"fmt"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"github.com/apsdehal/go-logger"
	"io/ioutil"
	"os"
	"time"
)

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

func Bak(appConfig *NGConfig, ngConfig *resolv.Config, c chan int) {
	for {
		select {
		case <-time.NewTicker(5 * time.Minute).C:
			bak(appConfig, ngConfig)
		case signal := <-c:
			if signal == 9 {
				Log(NOTICE, fmt.Sprintf("[%s] Nginx Config backup is stop.", appConfig.Name))
				break
			}

		}
	}
}

func bak(appConfig *NGConfig, ngConfig *resolv.Config) {
	bakDate := time.Now().Format("20060102")
	bakName := fmt.Sprintf("nginx.conf.%s.tgz", bakDate)

	bakPath, bErr := resolv.Backup(ngConfig, bakName)
	if bErr != nil && !os.IsExist(bErr) {
		message := fmt.Sprintf("[%s] Nginx Config backup to %s, but failed. <%s>", appConfig.Name, bakPath, bErr)
		Log(CRITICAL, message)
		Log(NOTICE, fmt.Sprintf("[%s] Nginx Config backup is stop.", appConfig.Name))
	} else if bErr == nil {
		Log(NOTICE, fmt.Sprintf("[%s] Nginx Config backup to %s", appConfig.Name, bakPath))
	}

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	} else {
		return false, nil
	}
}

func Log(level logger.LogLevel, message string) {

	myLogger.Log(level, message)
	//fmt.Printf("[%s] [%s] %s\n", level, time.Now().Format(timeFormat), message)

}
