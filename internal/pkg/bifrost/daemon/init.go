// bifrost 包, 该包包含了bifrost项目共用的一些变量，封装了http进程相关接口函数、一些公共函数和守护进程函数
// 创建人：ClessLi
// 创建时间：2020/06/10

package daemon

import (
	"flag"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"os"
	"path/filepath"
)

var (
	// 传入参数
	confPath = flag.String("f", "./configs/bifrost.yml", "the bifrost `config`uration file path.")
	Signal   = flag.String("s", "", "send `signal` to a master process: stop, restart, status")
	help     = flag.Bool("h", false, "this `help`")
	version  = flag.Bool("v", false, "this `version`")
	//confBackupDelay = flag.Duration("b", 10, "how many minutes `delay` for backup nginx config")

	// 程序工作目录
	workspace string

	// 进程文件
	pidFilename = "bifrost.pid"
	pidFile     string
)

// usage, 重新定义flag.Usage 函数，为bifrost帮助信息提供版本信息及命令行工具传参信息
func usage() {
	_, _ = fmt.Fprintf(os.Stdout, `bifrost version: %s
Usage: %s [-hv] [-f filename] [-s signal]

Options:
`, config.GetVersion(), os.Args[0])
	flag.PrintDefaults()
}

// init, bifrost包初始化函数
func init() {
	// DONE: nginx配置文件后台更新后自动热加载功能
	// 初始化工作目录
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	workspace = filepath.Dir(ex)
	err = os.Chdir(workspace)
	if err != nil {
		panic(err)
	}

	// 初始化pid文件路径
	pidFile = filepath.Join(workspace, pidFilename)

	// 初始化应用传参
	flag.Usage = usage
	// 判断是否为单元测试
	if len(os.Args) > 3 && os.Args[1] == "-test.v" && os.Args[2] == "-test.run" {
		fmt.Println(workspace)
		return
	}
	flag.Parse()

	if *confPath == "" {
		*confPath = "./configs/bifrost.yml"
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("bifrost version: %s\n", config.GetVersion())
		os.Exit(0)
	}

	// 判断传入配置文件目录
	isExistConfig, err := utils.PathExists(*confPath)
	if !isExistConfig {
		if err != nil {
			fmt.Println("The bifrost config file", "'"+*confPath+"'", "is not found.")
		} else {
			fmt.Println("Unkown error of the bifrost config file.")
		}
		flag.Usage()
		os.Exit(1)
	}

	// 判断传入信号
	if *Signal != "" && *Signal != "stop" && *Signal != "restart" && *Signal != "status" {
		flag.Usage()
		os.Exit(1)
	}
}
