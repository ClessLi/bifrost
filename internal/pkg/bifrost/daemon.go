package bifrost

import (
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// DONE: 编写biforst守护进程
// DONE: 修复Restart函数启动bifrost失败的bug

// Start, 守护进程 start 方法函数
// 返回值:
//     错误
func Start() (err error) {

	// 判断当前进程是子进程还是主进程
	if isMain() { // 主进程时
		// 执行主进程

		// 判断是否已存在子进程
		if pid, pidErr := utils.GetPid(pidFile); pidErr == nil {

			process, procErr := os.FindProcess(pid)
			if procErr != nil {
				return procErr
			}

			return fmt.Errorf("bifrost <PID %d> is running", process.Pid)
		} else if pidErr != utils.ErrProcessNotRunning {
			return pidErr
		}

		// 启动子进程
		fmt.Println("Starting bifrost...")
		exec, pathErr := filepath.Abs(os.Args[0])
		if pathErr != nil {
			return pathErr
		}

		args := append([]string{exec}, os.Args[1:]...)
		_, procErr := os.StartProcess(exec, args, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})
		return procErr

	} else { // 子进程时

		initConfig()
		utils.Logger.Debug("Running Sub Process")
		if BifrostConf.IsDebugLvl() {
			go func() {
				err := http.ListenAndServe("0.0.0.0:12378", nil)
				fmt.Println(err)
			}()
		}

		// 关闭进程后清理pid文件
		defer rmPidFile(pidFile)
		// 进程结束前操作
		defer func() {
			// 捕获panic
			if r := recover(); r != nil {
				err = fmt.Errorf("%s", r)
				utils.Logger.FatalF(err.Error())

			}
			//Log(NOTICE, "bifrost is finished")
		}()

		// 执行bifrost进程

		// 记录pid
		pid := os.Getpid()
		pidErr := ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 644)
		if pidErr != nil {
			utils.Logger.ErrorF("failed to start bifrost, cased by '%s'", pidErr)
			return pidErr
		}

		// 启动bifrost进程
		utils.Logger.NoticeF("bifrost <PID %d> is running", pid)
		//Run()
		runErr := ServerRun()
		utils.Logger.NoticeF("bifrost <PID %d> is finished", pid)
		return runErr
	}
}

// Stop, 守护进程 stop 方法函数
// 返回值:
//     错误
func Stop() error {

	// 判断bifrost进程是否存在
	process, procErr := getProc(pidFile)
	if procErr != nil {
		utils.Logger.Error(procErr.Error())
		return procErr
	}

	// 存在则关闭进程
	utils.Logger.Notice("Stopping bifrost...")
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
		_, procErr := getProc(pidFile)
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

// Restart, 守护进程 restart 方法函数
// 返回值:
//     错误
func Restart() error {
	// 判断当前进程是主进程还是子进程
	if isMain() { // 主进程时
		if err := Stop(); err != nil {
			utils.Logger.ErrorF("stop bifrost failed cased by: '%s'", err.Error())
			return err
		}

		return Start()
	} else { // 子进程时
		// 传参给子进程重启时，不重启
		return Start()
	}
}

// Status, 守护进程 status 方法函数
// 返回值:
//     错误
func Status() (int, error) {
	pid, pidErr := utils.GetPid(pidFile)
	if pidErr != nil {
		return -1, pidErr
	}
	_, procErr := os.FindProcess(pid)
	return pid, procErr
}

// isMain, 判断当前进程是否为主进程
// 返回值:
//     true: 是主进程; false: 是子进程
func isMain() bool {
	return os.Getppid() != 1
}
