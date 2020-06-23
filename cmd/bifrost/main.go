package main

// DONE: 1.权限管理
// DONE: 2.nginx配置定期备份机制
// DONE: 3.日志规范化输出
// DONE: 4.优化守护进程输出，调整标准输出到守护进程日志中

import (
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost"
	"os"
)

func main() {
	defer bifrost.Logf.Close()
	defer bifrost.Stdoutf.Close()

	err := errors.New("unkown signal")
	switch *bifrost.Signal {
	case "":
		err = bifrost.Start()
		if err == nil {
			fmt.Println("bifrost is started")
			os.Exit(0)
		}
	case "stop":
		err = bifrost.Stop()
		if err == nil {
			fmt.Println("bifrost is finished")
			os.Exit(0)
		}
	case "restart":
		err = bifrost.Restart()
		if err == nil {
			fmt.Println("bifrost is restarted")
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
