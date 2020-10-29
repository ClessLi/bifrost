// bifrost 包, 该包包含了bifrost项目共用的一些变量，封装了http进程相关接口函数、一些公共函数和守护进程函数
// 创建人：ClessLi
// 创建时间：2020/06/10

package bifrost

import (
	"flag"
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
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
	BifrostConf  = &Config{}
	authDBConfig *AuthDBConfig
	authConfig   *AuthConfig

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
	si           = &systemInfo{}
	svrPidFileKW = nginx.NewKeyWords(nginx.TypeKey, "pid", "*", false, false)

	// bifrost健康状态
	isHealthy = true

	// 初始化信号量
	signalChan = make(chan int)
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
	VERSION = "v1.0.0-alpha.8"

	// Web服务类型
	NGINX WebServerType = "nginx"
	HTTPD WebServerType = "httpd"
)

// Config, bifrost配置文件结构体，定义bifrost配置信息
type Config struct {
	WebServerInfo WebServerInfo `yaml:"WebServerInfo"`
	*AuthDBConfig `yaml:"AuthDBConfig,omitempty"`
	*AuthConfig   `yaml:"AuthConfig,omitempty"`
	LogConfig     `yaml:"LogConfig"`
}

// WebServerInfo, bifrost配置文件对象中web服务器信息结构体，定义管控的web服务器配置文件相关信息
type WebServerInfo struct {
	ListenPort int          `yaml:"listenPort"`
	Servers    []ServerInfo `yaml:"servers,flow"`
}

// WebServerType, web服务器类型对象，定义web服务器所属类型
type WebServerType string

// AuthDBConfig, mysql数据库信息结构体，该库用于存放用户认证信息（可选）
type AuthDBConfig struct {
	DBName   string `yaml:"DBName"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// AuthConfig, 认证信息结构体，记录用户认证信息（可选）
type AuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// LogConfig, bifrost日志信息结构体，定义日志目录、日志级别
type LogConfig struct {
	LogDir string          `yaml:"logDir"`
	Level  logger.LogLevel `yaml:"level"`
}

func (c LogConfig) IsDebugLvl() bool {
	return c.Level >= DEBUG
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
	/* 调测用
	confPath := "./configs/bifrost.yml"
	isExistConfig, pathErr := PathExists(confPath)
	*/
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

	// 加载bifrost配置
	yamlErr := yaml.Unmarshal(confData, BifrostConf)
	if yamlErr != nil {
		fmt.Println(yamlErr)
		flag.Usage()
		os.Exit(1)
	}

	// 初始化认证数据库或认证配置信息
	if BifrostConf.AuthDBConfig != nil {
		authDBConfig = BifrostConf.AuthDBConfig
	} else {
		if BifrostConf.AuthConfig != nil {
			authConfig = BifrostConf.AuthConfig
		} else { // 使用默认认证信息
			authConfig = &AuthConfig{"heimdall", "Bultgang"}
		}
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
	si.OS = fmt.Sprintf("%s %s", platform, release)
	si.BifrostVersion = VERSION

}
