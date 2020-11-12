package bifrost

import (
	"bytes"
	"fmt"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/apsdehal/go-logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

// readFile, 读取文件函数
// 参数:
//     path: 文件路径字符串
// 返回值:
//     文件数据
//     错误
func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

// PathExists, 判断文件路径是否存在函数
// 参数:
//     path: 待判断的文件路径字符串
// 返回值:
//     true: 存在; false: 不存在
//     错误
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	} else {
		return false, nil
	}
}

// Log, 日志记录函数
// 参数:
//     level: 日志级别对象
//     message: 需记录的日志信息字符串
func Log(level logger.LogLevel, message string, a ...interface{}) {
	myLogger.Log(level, fmt.Sprintf(message, a...))
}

func getProc(path string) (*os.Process, error) {
	pid, pidErr := getPid(path)
	if pidErr != nil {
		return nil, pidErr
	}
	return os.FindProcess(pid)
}

func rmPidFile(path string) {
	rmPidFileErr := os.Remove(path)
	if rmPidFileErr != nil {
		Log(ERROR, rmPidFileErr.Error())
	}
	Log(NOTICE, "bifrost.pid has been removed.")
}

// getPid, 查询pid文件并返回pid
// 返回值:
//     pid
//     错误
func getPid(path string) (int, error) {
	// 判断pid文件是否存在
	if _, err := os.Stat(path); err == nil || os.IsExist(err) { // 存在
		// 读取pid文件
		pidBytes, readPidErr := readFile(path)
		if readPidErr != nil {
			Log(ERROR, readPidErr.Error())
			return -1, readPidErr
		}

		// 去除pid后边的换行符
		pidBytes = bytes.TrimRight(pidBytes, "\n")

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

// bifrostConfigCheck, 检查bifrost配置项是否完整
// 返回值:
//     错误
func bifrostConfigCheck() error {
	if BifrostConf == nil {
		return fmt.Errorf("bifrost config load error")
	}
	if len(BifrostConf.Service.ServiceInfos) == 0 {
		return fmt.Errorf("bifrost services config load error")
	}
	if BifrostConf.LogDir == "" {
		return fmt.Errorf("bifrost log config load error")
	}
	// 初始化服务信息配置
	if BifrostConf.Service.ListenPort == 0 {
		BifrostConf.Service.ListenPort = 12321
	}
	if BifrostConf.Service.ChunckSize == 0 {
		BifrostConf.Service.ChunckSize = 4194304
	}
	// 初始化认证数据库或认证配置信息
	if BifrostConf.AuthService == nil {
		BifrostConf.AuthService = new(AuthService)
	}
	if BifrostConf.AuthService.AuthDBConfig != nil {
		authDBConfig = BifrostConf.AuthService.AuthDBConfig
	} else {
		if BifrostConf.AuthService.AuthConfig != nil {
			authConfig = BifrostConf.AuthService.AuthConfig
		} else { // 使用默认认证信息
			authConfig = &AuthConfig{"heimdall", "Bultgang"}
		}
	}
	return nil
}

func serviceRun() {
	fmt.Println("Entering bifrost...")
	Log(DEBUG, "Listening system call signal")
	go ListenSignal(signalChan)
	Log(DEBUG, "Listened system call signal")
	// 初始化各web服务配置信息管控接口
	BifrostConf.Service.Run()

	// 打印启动logo
	fmt.Println(logoStr)

	Log(DEBUG, "Listening ports for bifrost")
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", BifrostConf.Service.ListenPort))
		// 启动bifrost各接口服务
		if err != nil {
			BifrostConf.Service.killCoroutines()
			// 输出子任务运行错误
			Log(CRITICAL, "bifrost Run coroutine has been stoped. Cased by '%s'", err)
			return
		}
		defer lis.Close()
		// 初始化gRPC接口服务端
		// 初始化注册参数
		var opts []grpc.ServerOption
		// Register interceptor
		var interceptor grpc.UnaryServerInterceptor
		interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			ip, err := getClientIP(ctx)
			if err != nil {
				Log(DEBUG, "Failed to get client address")
			}
			Log(DEBUG, "Client address is %s", ip)
			return handler(ctx, req)
		}
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
		grpcSvr := grpc.NewServer(opts...)
		defer grpcSvr.Stop()
		// 注册gRPC接口
		bifrostpb.RegisterAuthServiceServer(grpcSvr, BifrostConf.AuthService)
		bifrostpb.RegisterOperationServiceServer(grpcSvr, BifrostConf.Service)
		// 启动gRPC接口服务
		err = grpcSvr.Serve(lis)
		if err != nil {
			Log(ERROR, err.Error())
		}
		Log(DEBUG, "Stopping listen ports for bifrost")
	}()
	Log(DEBUG, "Listened ports for bifrost")

	select {
	case s := <-signalChan:
		if s == 9 {
			Log(DEBUG, "stopping...")
			break
		}
		Log(DEBUG, "stopping signal error")
	}
	BifrostConf.Service.killCoroutines()
	Log(NOTICE, "bifrost Run coroutine has been stoped.")
	fmt.Println("Exit bifrost")
	return
}
