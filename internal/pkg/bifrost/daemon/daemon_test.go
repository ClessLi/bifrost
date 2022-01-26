package daemon

import (
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/apsdehal/go-logger"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewDaemon(t *testing.T) {
	singletonBifrostConf = new(config.Config)
	configData, err := utils.ReadFile("F:\\GO_Project\\src\\bifrost\\test\\configs\\bifrost.yml")
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(configData, singletonBifrostConf)
	if err != nil {
		t.Fatal(err)
	}
	// 初始化日志
	logDir, err := filepath.Abs("F:\\GO_Project\\src\\bifrost\\test\\logs")
	if err != nil {
		panic(err)
	}

	logPath := filepath.Join(logDir, "bifrost.log")
	utils.Logf, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	utils.InitLogger(utils.Logf, singletonBifrostConf.Level)

	// 初始化应用运行日志输出
	stdoutPath := filepath.Join(logDir, "bifrost.out")
	utils.Stdoutf, err = os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	os.Stdout = utils.Stdoutf
	os.Stderr = utils.Stdoutf

	pidFile = filepath.Join(logDir, pidFilename)

	errChan := make(chan error)
	service, managers := newService(errChan)
	ip, err := externalIP()
	if err != nil {
		t.Fatal(err)
	}
	server := NewServer(newGRPCServer(getBifrostConfInstance().ServiceConfig.ChunkSize, service), ip, getBifrostConfInstance().ServiceConfig.Port, "", "", nil)

	daemon := newSubDaemon(managers, server, pidFile, make(chan int), getBifrostConfInstance().IsDebugLvl())
	err = daemon.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMainDaemon_Status(t *testing.T) {
	// 初始化日志
	logDir, err := filepath.Abs("F:\\GO_Project\\src\\bifrost\\test\\logs")
	if err != nil {
		panic(err)
	}

	logPath := filepath.Join(logDir, "bifrost.log")
	utils.Logf, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	utils.InitLogger(utils.Logf, logger.DebugLevel)

	// 初始化应用运行日志输出
	stdoutPath := filepath.Join(logDir, "bifrost.out")
	utils.Stdoutf, err = os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	os.Stdout = utils.Stdoutf
	os.Stderr = utils.Stdoutf
	pidFile = filepath.Join(logDir, pidFilename)
	pid, pidErr := utils.GetPid(pidFile)
	if pidErr != nil {
		t.Fatal(pidErr)
	}
	daemon := newMainDaemon(pid)
	for i := 0; i < 30; i++ {
		fmt.Println(i)
		err = daemon.Status()
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second)
	}

}

func TestMainDaemon_Stop(t *testing.T) {
	// 初始化日志
	logDir, err := filepath.Abs("F:\\GO_Project\\src\\bifrost\\test\\logs")
	if err != nil {
		panic(err)
	}

	logPath := filepath.Join(logDir, "bifrost.log")
	utils.Logf, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	utils.InitLogger(utils.Logf, logger.DebugLevel)

	// 初始化应用运行日志输出
	stdoutPath := filepath.Join(logDir, "bifrost.out")
	utils.Stdoutf, err = os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	os.Stdout = utils.Stdoutf
	os.Stderr = utils.Stdoutf
	pidFile = filepath.Join(logDir, pidFilename)
	pid, pidErr := utils.GetPid(pidFile)
	if pidErr != nil {
		t.Fatal(pidErr)
	}
	daemon := newMainDaemon(pid)
	err = daemon.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
