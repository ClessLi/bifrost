package bifrost

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// TODO: 编写biforst守护进程
func Start() error {
	if os.Getppid() != 1 {
		execPath, pathErr := filepath.Abs(os.Args[0])
		if pathErr != nil {
			return pathErr
		}

		args := append([]string{execPath}, os.Args[1:]...)
		pidFile := filepath.Join(filepath.Dir(os.Args[0]), "biforst.pid")
		if _, err := os.Stat(pidFile); err == nil || os.IsExist(err) {
			pidBytes, readPidErr := readFile(pidFile)
			if readPidErr != nil {
				return readPidErr
			}

			pid, toIntErr := strconv.Atoi(string(pidBytes))
			if toIntErr != nil {
				return toIntErr
			}

			process, procErr := os.FindProcess(pid)
			if procErr != nil {
				return procErr
			}

			return fmt.Errorf("biforst is running, Pid is %d", process.Pid)
		}

		process, procErr := os.StartProcess(execPath, args, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})
		if procErr != nil {
			return procErr
		}

		pidErr := ioutil.WriteFile(filepath.Join(filepath.Dir(os.Args[0]), "biforst.pid"), []byte(fmt.Sprintf("%d", process.Pid)), 644)
		if pidErr != nil {
			return pidErr
		}

		return nil
	}
	return fmt.Errorf("unkonw error")
}
