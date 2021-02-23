package bifrost

import (
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/ClessLi/skirnir/pkg/discover"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// DONE: 编写biforst守护进程
// DONE: 修复Restart函数启动bifrost失败的bug
// DONE: 新增Daemon接口对象

type Daemon interface {
	Start() error
	Stop() error
	Restart() error
	Status() error
}

type mainDaemon struct {
	subDaemonPid int
}

func (m mainDaemon) Start() error {
	// 执行主进程

	// 判断是否已存在子进程
	if m.subDaemonPid > 0 {
		process, procErr := os.FindProcess(m.subDaemonPid)
		if procErr != nil {
			return procErr
		}
		return fmt.Errorf("bifrost <PID %d> is running", process.Pid)
	}

	// 启动子进程
	fmt.Println("Starting bifrost...")
	exec, pathErr := filepath.Abs(os.Args[0])
	if pathErr != nil {
		return pathErr
	}

	args := append([]string{exec}, os.Args[1:]...)
	_, procErr := os.StartProcess(exec, args, &os.ProcAttr{
		Files: []*os.File{os.Stdin, utils.Stdoutf, utils.Stdoutf},
	})
	return procErr
}

func (m mainDaemon) Stop() error {
	// 判断bifrost进程是否存在
	if m.subDaemonPid <= 0 {
		return utils.ErrProcessNotRunning
	}

	process, procErr := os.FindProcess(m.subDaemonPid)
	if procErr != nil {
		return procErr
	}

	// 存在则关闭进程
	fmt.Println("Stopping bifrost...")
	killErr := process.Signal(syscall.SIGQUIT)
	if killErr != nil {
		if sysErr, ok := killErr.(*os.SyscallError); !ok || sysErr.Syscall != "TerminateProcess" {
			utils.Logger.Error(killErr.Error())
			return killErr
		} else if ok && sysErr.Syscall == "TerminateProcess" {
			utils.Logger.Notice("bifrost is stopping or stopped")
		}
	}

	for i := 0; i < 300; i++ {
		_, procErr = getProc(pidFile)
		if procErr != nil {
			break
		}
		if i == 299 {
			return fmt.Errorf("an unknown error occurred in terminating bifrost")
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (m *mainDaemon) Restart() (err error) {
	err = m.Stop()
	if err != nil {
		return err
	}

	m.subDaemonPid = -1
	return m.Start()
}

func (m mainDaemon) Status() error {
	if m.subDaemonPid <= 0 {
		return utils.ErrProcessNotRunning
	}
	_, procErr := os.FindProcess(m.subDaemonPid)
	return procErr
}

func newMainDaemon(pid int) Daemon {
	return &mainDaemon{subDaemonPid: pid}
}

type subDaemon struct {
	offstageManager service.OffstageManager
	server          Server
	pidFile         string
	signalChan      chan int
	isDebugLvl      bool
}

func (s subDaemon) Start() error {
	utils.Logger.Debug("Running Sub Process")
	if s.isDebugLvl {
		go func() {
			err := http.ListenAndServe("0.0.0.0:12378", nil)
			fmt.Println(err)
		}()
	}

	// 关闭进程后清理pid文件
	defer s.rmPidFile()
	// 进程结束前操作
	defer func() {
		// 捕获panic
		if r := recover(); r != nil {
			err := fmt.Errorf("%s", r)
			utils.Logger.FatalF("panic: %s", err.Error())

		}
	}()

	// 记录pid
	pid := os.Getpid()
	pidErr := ioutil.WriteFile(s.pidFile, []byte(fmt.Sprintf("%d", pid)), 644)
	if pidErr != nil {
		utils.Logger.ErrorF("failed to start bifrost, cased by '%s'", pidErr)
		return pidErr
	}

	// 启动bifrost服务进程
	utils.Logger.NoticeF("bifrost <PID %d> is running", pid)
	runErr := s.serverRun()
	utils.Logger.NoticeF("bifrost <PID %d> is finished", pid)
	return runErr
}

func (s subDaemon) Stop() error {
	panic("subDaemon.Stop method not enabled")
}

func (s subDaemon) Restart() error {
	return s.Start()
}

func (s subDaemon) Status() error {
	panic("subDaemon.Status method not enabled")
}

func (s subDaemon) rmPidFile() {
	rmPidFileErr := os.Remove(s.pidFile)
	if rmPidFileErr != nil {
		utils.Logger.Error(rmPidFileErr.Error())
	}
	utils.Logger.NoticeF("%s has been removed.", s.pidFile)
}

func (s subDaemon) listenSignal() {
	procSigs := make(chan os.Signal, 1)
	signal.Notify(procSigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case sig := <-procSigs:
		utils.Logger.NoticeF("Get system signal: %s", sig.String())
		s.signalChan <- 9
		utils.Logger.Debug("Stop listen system signal")
	}
}

func (s subDaemon) serverRun() error {
	err := s.offstageManager.Start()
	if err != nil {
		return err
	}
	defer func() {
		err := s.offstageManager.Stop()
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}()

	// 开始侦听系统关闭信号
	utils.Logger.Debug("Listening system call signal")
	go s.listenSignal()
	utils.Logger.Debug("Listened system call signal")

	// 启动gRPC服务
	svrErrChan := make(chan error)
	s.server.Start(svrErrChan)

	utils.Logger.Info(logoStr)
	var stopErr error
	select {
	case sig := <-s.signalChan:
		if sig == 9 {
			utils.Logger.Debug("bifrost service is stopping...")
			s.server.Stop()
		}
		utils.Logger.Debug("stop signal error")
	case stopErr = <-svrErrChan:
		s.server.Stop()
		break
	}
	return stopErr
}

func newSubDaemon(manager service.OffstageManager, server Server, pidFile string, signalChan chan int, isDebugLvl bool) Daemon {
	if manager == nil {
		panic("offstage manager is nil")
	}

	return &subDaemon{
		offstageManager: manager,
		server:          server,
		pidFile:         pidFile,
		signalChan:      signalChan,
		isDebugLvl:      isDebugLvl,
	}
}

func NewDaemon() Daemon {
	if isMain() {
		pid, pidErr := utils.GetPid(pidFile)
		if pidErr != nil && pidErr != utils.ErrProcessNotRunning {
			fmt.Println(pidErr)
			os.Exit(1)
		}
		return newMainDaemon(pid)
	}
	// 初始化bifrost配置
	confData, err := utils.ReadFile(*confPath)
	if err != nil {
		panic(err)
	}
	bifrostConf := new(Config)
	// 加载bifrost配置
	err = yaml.Unmarshal(confData, bifrostConf)
	if err != nil {
		panic(err)
	}

	// 配置必填项检查
	err = bifrostConf.check()
	if err != nil {
		panic(err)
	}

	// 初始化日志
	logDir, err := filepath.Abs(bifrostConf.LogDir)
	if err != nil {
		panic(err)
	}

	logPath := filepath.Join(logDir, "bifrost.log")
	utils.Logf, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	utils.InitLogger(utils.Logf, bifrostConf.LogConfig.Level)

	// 初始化应用运行日志输出
	stdoutPath := filepath.Join(logDir, "bifrost.out")
	utils.Stdoutf, err = os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	os.Stdout = utils.Stdoutf
	os.Stderr = utils.Stdoutf

	// 初始化bifrost服务
	errChan := make(chan error)
	svc, managers := newService(bifrostConf, errChan)
	gRPCServer := newGRPCServer(bifrostConf.ServiceConfig.ChunckSize, svc)
	serviceName := "com.github.ClessLi.api.bifrost"
	instanceId := serviceName + "-" + uuid.NewV4().String()
	var registryClient discover.RegistryClient
	if bifrostConf.RAConfig != nil {
		registryClient, err = discover.NewKitConsulRegistryClient(bifrostConf.RAConfig.Host, bifrostConf.RAConfig.Port)
		if err != nil {
			utils.Logger.WarningF("Get Consul Client failed. Cased by: %s", err)
			panic(err)
		}
	}
	listenIP, err := externalIP()
	if err != nil {
		panic(err)
	}
	return newSubDaemon(managers, NewServer(gRPCServer, listenIP, bifrostConf.ServiceConfig.Port, serviceName, instanceId, registryClient), pidFile, make(chan int), bifrostConf.IsDebugLvl())
}

// isMain, 判断当前进程是否为主进程
// 返回值:
//     true: 是主进程; false: 是子进程
func isMain() bool {
	return os.Getppid() != 1
}
