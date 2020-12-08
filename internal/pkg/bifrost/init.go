// bifrost 包, 该包包含了bifrost项目共用的一些变量，封装了http进程相关接口函数、一些公共函数和守护进程函数
// 创建人：ClessLi
// 创建时间：2020/06/10

package bifrost

import (
	"flag"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/ClessLi/skirnir/pkg/discover"
	"github.com/apsdehal/go-logger"
	"github.com/shirou/gopsutil/host"
	"gopkg.in/yaml.v2"
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

	// bifrost配置
	BifrostConf = new(Config)
	//authDBConfig *AuthDBConfig
	//authConfig   *AuthConfig

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

	// 错误变量
	procStatusNotRunning = fmt.Errorf("process is not running")

	// 初始化systemInfo监控对象
	//sysInfo      = &systemInfo{}
	//svrPidFileKW = nginx.NewKeyWords(nginx.TypeKey, "pid", "*", false, false)

	// 初始化信号量
	signalChan = make(chan int)
	isInit     bool

	// 服务实例id
	instanceId      string
	discoveryClient discover.RegistryClient
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
	VERSION = "v1.0.1-alpha.2"
)

// Config, bifrost配置文件结构体，定义bifrost配置信息
type Config struct {
	Service *service.BifrostService `yaml:"Service"`
	//AuthService *AuthService `yaml:"AuthService"`
	*RAConfig `yaml:"RAConfig"`
	LogConfig `yaml:"LogConfig"`
}

type RAConfig struct {
	Host string `yaml:"Host"`
	Port uint16 `yaml:"Port"`
}

// LogConfig, bifrost日志信息结构体，定义日志目录、日志级别
type LogConfig struct {
	LogDir string          `yaml:"logDir"`
	Level  logger.LogLevel `yaml:"level"`
}

func (l LogConfig) IsDebugLvl() bool {
	return l.Level >= DEBUG
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
	// DONE: nginx配置文件后台更新后自动热加载功能
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

	if *version {
		fmt.Printf("bifrost version: %s\n", VERSION)
		os.Exit(0)
	}

	// 判断传入配置文件目录
	isExistConfig, pathErr := PathExists(*confPath)
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
	if readErr != nil {
		fmt.Println(readErr)
		flag.Usage()
		os.Exit(1)
	}

	// 加载bifrost配置
	yamlErr := yaml.Unmarshal(confData, BifrostConf)
	if yamlErr != nil {
		fmt.Println(yamlErr)
		flag.Usage()
		os.Exit(1)
	}

	// 配置必填项检查
	checkErr := configCheck()
	if checkErr != nil {
		fmt.Println(checkErr)
		flag.Usage()
		os.Exit(1)
	}

	// 初始化日志
	logDir, absErr := filepath.Abs(BifrostConf.LogDir)
	if absErr != nil {
		panic(absErr)
	}

	logPath := filepath.Join(logDir, "bifrost.log")
	Logf, openErr := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if openErr != nil {
		panic(openErr)
	}

	myLogger, err = logger.New("Bifrost", BifrostConf.Level, Logf)
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

	platform, _, release, OSErr := host.PlatformInformation()
	if OSErr != nil {
		Log(CRITICAL, "bifrost is stopped, cased by '%s'", OSErr)
		os.Exit(1)
	}
	service.SysInfo.OS = fmt.Sprintf("%s %s", platform, release)
	service.SysInfo.BifrostVersion = VERSION
	isInit = true
}
