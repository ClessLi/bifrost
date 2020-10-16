package bifrost

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// DONE: 编写biforst守护进程
// DONE: 修复Restart函数启动bifrost失败的bug

// Start, 守护进程 start 方法函数
// 返回值:
//     错误
func Start() (err error) {
	// 判断当前进程是子进程还是主进程
	if isMain() { // 主进程时
		// 执行子进程

		// 判断是否已存在子进程
		if pid, pidErr := getPid(); pidErr == nil {

			process, procErr := os.FindProcess(pid)
			if procErr != nil {
				return procErr
			}

			return fmt.Errorf("bifrost <PID %d> is running", process.Pid)
		} else if pidErr != procStatusNotRunning {
			return pidErr
		}

		// 启动子进程
		Log(NOTICE, "starting bifrost...")
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
		if procErr != nil {
			return procErr
		}

		return nil
	} else { // 子进程时
		// 进程结束前操作
		defer func() {
			// 捕获panic
			if r := recover(); r != nil {
				err = fmt.Errorf("%s", r)
				Log(CRITICAL, err.Error())

			}
			// 进程结束前清理pid文件
			rmPidFileErr := os.Remove(pidFile)
			if rmPidFileErr != nil {
				err = rmPidFileErr
				Log(ERROR, rmPidFileErr.Error())
			}
			Log(NOTICE, "bifrost.pid is removed, bifrost is finished")
		}()

		// 执行bifrost进程

		// 记录pid
		pid := os.Getpid()
		pidErr := ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 644)
		if pidErr != nil {
			Log(ERROR, "failed to start bifrost, cased by '%s'", pidErr)
			return pidErr
		}

		// 启动bifrost进程
		Log(NOTICE, "bifrost <PID %d> is running", pid)
		Run()
		stat := fmt.Sprintf("bifrost <PID %d> is finished", pid)
		Log(NOTICE, stat)
		return fmt.Errorf(stat)
	}
}

// Stop, 守护进程 stop 方法函数
// 返回值:
//     错误
func Stop() error {
	// 判断bifrost进程是否存在
	pid, pidErr := getPid()
	if pidErr != nil {
		return pidErr
	}
	process, procErr := os.FindProcess(pid)
	if procErr != nil {
		Log(ERROR, procErr.Error())
		return procErr
	}

	// 存在则关闭进程
	killErr := process.Kill()
	if killErr != nil {
		if sysErr, ok := killErr.(*os.SyscallError); !ok || sysErr.Syscall != "TerminateProcess" {
			Log(ERROR, killErr.Error())
			return killErr
		} else if ok && sysErr.Syscall == "TerminateProcess" {
			Log(NOTICE, "bifrost is stopping or stopped")
		}
	}

	// 关闭进程后清理pid文件
	rmPidFileErr := os.Remove(pidFile)
	if rmPidFileErr != nil {
		Log(ERROR, rmPidFileErr.Error())
		return rmPidFileErr
	}
	Log(NOTICE, "bifrost.pid is removed, bifrost is finished")

	return nil
}

// getPid, 查询pid文件并返回pid
// 返回值:
//     pid
//     错误
func getPid() (int, error) {
	// 判断pid文件是否存在
	if _, err := os.Stat(pidFile); err == nil || os.IsExist(err) { // 存在
		// 读取pid文件
		pidBytes, readPidErr := readFile(pidFile)
		if readPidErr != nil {
			Log(ERROR, readPidErr.Error())
			return -1, readPidErr
		}

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

// Restart, 守护进程 restart 方法函数
// 返回值:
//     错误
func Restart() error {
	// 判断当前进程是主进程还是子进程
	if isMain() { // 主进程时
		Log(NOTICE, "stopping bifrost...")
		if err := Stop(); err != nil {
			Log(ERROR, "stop bifrost failed cased by: '%s'", err.Error())
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
	pid, pidErr := getPid()
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
