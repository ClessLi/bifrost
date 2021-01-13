// bifrost 包, 该包包含了bifrost项目共用的一些变量，封装了http进程相关接口函数、一些公共函数和守护进程函数
// 创建人：ClessLi
// 创建时间：2020/06/10

package bifrost

import (
	"flag"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service/web_server_manager"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/ClessLi/skirnir/pkg/discover"
	"github.com/apsdehal/go-logger"
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

	// 日志文件
	Logf    *os.File
	Stdoutf *os.File

	// 程序工作目录
	workspace string

	// 进程文件
	pidFilename = "bifrost.pid"
	pidFile     string

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

// Config, bifrost配置文件结构体，定义bifrost配置信息
type Config struct {
	ServiceConfig ServiceConfig `yaml:"Service"`
	//AuthService *AuthService `yaml:"AuthService"`
	*RAConfig `yaml:"RAConfig"`
	LogConfig `yaml:"LogConfig"`
}

type ServiceConfig struct {
	Port                     uint16                                   `yaml:"Port"`
	ChunckSize               int                                      `yaml:"ChunkSize"`
	AuthServerAddr           string                                   `yaml:"AuthServerAddr"`
	WebServerConfigInfos     []web_server_manager.WebServerConfigInfo `yaml:"Infos,flow"`
	BifrostServiceController *service.BifrostServiceController
}

//type Info struct {
//	Name           string                `yaml:"name"`
//	Type           web_server_manager.WebServerType `yaml:"type"`
//	BackupCycle    int                   `yaml:"backupCycle"`
//	BackupSaveTime int                   `yaml:"backupSaveTime"`
//	BackupDir      string                `yaml:"backupDir,omitempty"`
//	ConfPath       string                `yaml:"confPath"`
//	VerifyExecPath string                `yaml:"verifyExecPath"`
//}

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
	return l.Level >= logger.DebugLevel
}

// usage, 重新定义flag.Usage 函数，为bifrost帮助信息提供版本信息及命令行工具传参信息
func usage() {
	_, _ = fmt.Fprintf(os.Stdout, `bifrost version: %s
Usage: %s [-hv] [-f filename] [-s signal]

Options:
`, utils.Version(), os.Args[0])
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
	flag.Parse()
	if *confPath == "" {
		*confPath = "./configs/bifrost.yml"
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("bifrost version: %s\n", utils.Version())
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
	// 初始化bifrost配置
	confData, err := utils.ReadFile(*confPath)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	// 加载bifrost配置
	err = yaml.Unmarshal(confData, BifrostConf)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	// 配置必填项检查
	err = configCheck()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	// 初始化日志
	logDir, err := filepath.Abs(BifrostConf.LogDir)
	if err != nil {
		panic(err)
	}

	logPath := filepath.Join(logDir, "bifrost.log")
	Logf, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	utils.InitLogger(Logf, BifrostConf.LogConfig.Level)

	// 初始化应用运行日志输出
	stdoutPath := filepath.Join(logDir, "bifrost.out")
	Stdoutf, err = os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	os.Stdout = Stdoutf
	os.Stderr = Stdoutf
}

func initConfig() {
	// 初始化web服务器配置服务控制器
	BifrostConf.ServiceConfig.BifrostServiceController = service.NewBifrostServiceController(BifrostConf.ServiceConfig.WebServerConfigInfos...)

	isInit = true
}
