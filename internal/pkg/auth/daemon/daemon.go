package daemon

import (
	"fmt"
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
//
//	错误
func Start() (err error) {
	// 判断当前进程是子进程还是主进程
	if isMain() { // 主进程时
		// 执行主进程
		Log(DEBUG, "Running Main Process")

		// 判断是否已存在子进程
		if pid, pidErr := getPid(pidFile); pidErr == nil {
			process, procErr := os.FindProcess(pid)
			if procErr != nil {
				return procErr
			}

			return fmt.Errorf("bifrost <PID %d> is running", process.Pid)
		} else if pidErr != procStatusNotRunning {
			return pidErr
		}

		// 启动子进程
		Log(NOTICE, "Starting bifrost...")
		os.Stdout = Stdoutf
		os.Stderr = Stdoutf
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
		Log(DEBUG, "Running Sub Process")
		if AuthConf.IsDebugLvl() {
			go func() {
				err := http.ListenAndServe("0.0.0.0:12377", nil)
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
				Log(CRITICAL, err.Error())
			}
			// Log(NOTICE, "bifrost is finished")
		}()

		// 执行bifrost进程

		// 记录pid
		pid := os.Getpid()
		pidErr := ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644)
		if pidErr != nil {
			Log(ERROR, "failed to start bifrost, cased by '%s'", pidErr)
			return pidErr
		}

		// 启动bifrost进程
		Log(NOTICE, "bifrost <PID %d> is running", pid)
		// Run()
		runErr := ServerRun()
		Log(NOTICE, "bifrost <PID %d> is finished", pid)
		return runErr
	}
}

// Stop, 守护进程 stop 方法函数
// 返回值:
//
//	错误
func Stop() error {
	// 判断bifrost进程是否存在
	process, procErr := getProc(pidFile)
	if procErr != nil {
		Log(ERROR, procErr.Error())
		return procErr
	}

	// 存在则关闭进程
	Log(NOTICE, "Stopping bifrost...")
	killErr := process.Signal(syscall.SIGQUIT)
	if killErr != nil {
		if sysErr, ok := killErr.(*os.SyscallError); !ok || sysErr.Syscall != "TerminateProcess" {
			Log(ERROR, killErr.Error())
			return killErr
		} else if ok && sysErr.Syscall == "TerminateProcess" {
			Log(NOTICE, "bifrost is stopping or stopped")
		}
	}

	return nil
}

// Restart, 守护进程 restart 方法函数
// 返回值:
//
//	错误
func Restart() error {
	// 判断当前进程是主进程还是子进程
	if isMain() { // 主进程时
		if err := Stop(); err != nil {
			Log(ERROR, "stop bifrost failed cased by: '%s'", err.Error())
			return err
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

		return Start()
	} else { // 子进程时
		// 传参给子进程重启时，不重启
		return Start()
	}
}

// Status, 守护进程 status 方法函数
// 返回值:
//
//	错误
func Status() (int, error) {
	pid, pidErr := getPid(pidFile)
	if pidErr != nil {
		return -1, pidErr
	}
	_, procErr := os.FindProcess(pid)
	return pid, procErr
}

// isMain, 判断当前进程是否为主进程
// 返回值:
//
//	true: 是主进程; false: 是子进程
func isMain() bool {
	return os.Getppid() != 1
}
