package main

// DONE: 1.权限管理
// DONE: 2.nginx配置定期备份机制
// DONE: 3.日志规范化输出
// DONE: 4.优化守护进程输出，调整标准输出到守护进程日志中

import (
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/daemon"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"os"
)

func main() {

	defer utils.Logf.Close()
	defer utils.Stdoutf.Close()

	err := errors.New("unkown signal")
	bifrostDaemon := daemon.NewDaemon()
	switch *daemon.Signal {
	case "":
		err = bifrostDaemon.Start()
	case "stop":
		err = bifrostDaemon.Stop()
	case "restart":
		err = bifrostDaemon.Restart()
	case "status":
		err = bifrostDaemon.Status()
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
