package main

// DONE: 1.权限管理
// DONE: 2.nginx配置定期备份机制
// DONE: 3.日志规范化输出
// DONE: 4.优化守护进程输出，调整标准输出到守护进程日志中

import (
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"os"
)

func main() {

	defer utils.Logf.Close()
	defer utils.Stdoutf.Close()

	err := errors.New("unkown signal")
	switch *bifrost.Signal {
	case "":
		err = bifrost.Start()
		if err == nil {
			os.Exit(0)
		}
	case "stop":
		err = bifrost.Stop()
		if err == nil {
			if os.Getppid() != 1 {
				fmt.Println("bifrost is stopped")
			}
			os.Exit(0)
		}
	case "restart":
		err = bifrost.Restart()
		if err == nil {
			os.Exit(0)
		}
	case "status":
		pid, statErr := bifrost.Status()
		if statErr != nil {
			fmt.Printf("bifrost is abnormal with error: %s\n", statErr.Error())
			os.Exit(1)
		} else {
			fmt.Printf("bifrost <PID %d> is running\n", pid)
			os.Exit(0)
		}
	}
	fmt.Println(err.Error())
	os.Exit(1)
}
