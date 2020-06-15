// bifrost 包, 该包包含了bifrost项目共用的一些变量，封装了http进程相关接口函数、一些公共函数和守护进程函数
// 创建人：ClessLi
// 创建时间：2020/06/10

package bifrost

import (
	"flag"
	"fmt"
	"github.com/apsdehal/go-logger"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"sync"
)

var (
	// 传入参数
	confPath = flag.String("f", "./configs/bifrost.yml", "the bifrost `config`uration file path.")
	Signal   = flag.String("s", "", "send `signal` to a master process: stop, restart, status")
	help     = flag.Bool("h", false, "this `help`")
	//confBackupDelay = flag.Duration("b", 10, "how many minutes `delay` for backup nginx config")

	// bifrost配置
	Configs  *ParserConfigs
	dbConfig DBConfig

	// 日志变量
	myLogger *logger.Logger
	// 日志文件
	Logf    *os.File
	Stdoutf *os.File

	// 程序工作目录
	workspace string

	// 进程文件
	pidFilename = "bifrost.pid"
	pidFile     string

	// 协程等待组
	wg sync.WaitGroup

	// 错误变量
	procStatusNotRunning = fmt.Errorf("bifrost is not running")
)

const (
	// 日志级别
	CRITICAL = logger.CriticalLevel
	ERROR    = logger.ErrorLevel
	WARN     = logger.WarningLevel
	NOTICE   = logger.NoticeLevel
	INFO     = logger.InfoLevel
	DEBUG    = logger.DebugLevel

	// 版本号
	VERSION = "v0.0.3-beta.1"
)

// NGConfig, nginx配置文件信息结构体，定义配置文件路径、nginx可执行文件路径和bifrost为其提供接口的路由及侦听端口
type NGConfig struct {
	//Name         string `json:"name"`
	//RelativePath string `json:"relative_path"`
	//Port         int    `json:"port"`
	//ConfPath     string `json:"conf_path"`
	Name         string `yaml:"name"`
	RelativePath string `yaml:"relativePath"`
	Port         int    `yaml:"port"`
	ConfPath     string `yaml:"confPath"`
	NginxBin     string `yaml:"nginxBin"`
}

// DBConfig, mysql数据库信息结构体，用于读写bifrost信息
type DBConfig struct {
	DBName   string `yaml:"DBName"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// LogConfig, bifrost日志信息结构体，定义日志目录、日志级别
type LogConfig struct {
	LogDir string          `yaml:"logDir"`
	Level  logger.LogLevel `yaml:"level"`
}

// ParserConfigs, bifrost配置文件结构体，定义bifrost配置信息
type ParserConfigs struct {
	//NGConfigs []NGConfig `json:"NGConfigs"`
	//NGConfigs  []NGConfig `yaml:"NGConfigs"`
	NGConfigs []NGConfig `yaml:"NGConfigs"`
	DBConfig  `yaml:"DBConfig"`
	LogConfig `yaml:"logConfig"`
}

// usage, 重新定义flag.Usage 函数，为bifrost帮助信息提供版本信息及命令行工具传参信息
func usage() {
	_, _ = fmt.Fprintf(os.Stdout, `bifrost version: %s
Usage: %s [-hv] [-f filename] [-s signal]

Options:
`, VERSION, os.Args[0])
	flag.PrintDefaults()
}

// init, bifrost包初始化函数
func init() {
	// 初始化工作目录
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	workspace = filepath.Dir(ex)
	cdErr := os.Chdir(workspace)
	if cdErr != nil {
		panic(cdErr)
	}

	// 初始化pid文件路径
	pidFile = filepath.Join(workspace, pidFilename)

	// 初始化应用传参
	flag.Usage = usage
	flag.Parse()
	if *confPath == "" {
		*confPath = "./configs/bifrost.yml"
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// 判断传入配置文件目录
	//confPath := "./configs/ng-conf-info.json"
	//confPath := "./configs/bifrost.yml"
	isExistConfig, pathErr := PathExists(*confPath)
	//isExistConfig, pathErr := PathExists(confPath)
	if !isExistConfig {
		if pathErr != nil {
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

	// 初始化bifrost配置
	confData, readErr := readFile(*confPath)
	//confData, readErr := readFile(confPath)
	if readErr != nil {
		fmt.Println(readErr)
		flag.Usage()
		os.Exit(1)
	}

	Configs = &ParserConfigs{}
	//jsonErr := json.Unmarshal(confData, configs)
	jsonErr := yaml.Unmarshal(confData, Configs)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		flag.Usage()
		os.Exit(1)
	}

	// 初始化数据库信息
	dbConfig = Configs.DBConfig

	// 初始化日志
	logDir, absErr := filepath.Abs(Configs.LogDir)
	if absErr != nil {
		panic(absErr)
	}

	logPath := filepath.Join(logDir, "bifrost.log")
	Logf, openErr := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if openErr != nil {
		panic(openErr)
	}

	myLogger, err = logger.New("Bifrost", Configs.Level, Logf)
	if err != nil {
		panic(err)
	}
	myLogger.SetFormat("[%{module}] %{time:2006-01-02 15:04:05.000} [%{level}] %{message}\n")

	// 初始化应用运行日志输出
	stdoutPath := filepath.Join(logDir, "bifrost.out")
	Stdoutf, openErr = os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if openErr != nil {
		panic(openErr)
	}
}