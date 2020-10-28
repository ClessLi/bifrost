package bifrost

import (
	"encoding/json"
	"fmt"
	ngJson "github.com/ClessLi/bifrost/pkg/json/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	nginxStatistics "github.com/ClessLi/bifrost/pkg/statistics/nginx"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// Run, bifrost启动函数
func Run() {

	Log(DEBUG, "Listening system call signal")
	go ListenSignal(signalChan)
	Log(DEBUG, "Listened system call signal")

	// 初始化接口
	loginURI := "/login"
	verifyURI := "/verify"
	refreshURI := "/refresh"
	sysinfoURI := "/sysinfo"
	healthURI := "/health"

	router := gin.Default()
	// login
	router.GET(loginURI, login)
	// verify
	router.GET(verifyURI, verify)
	// refresh
	router.GET(refreshURI, refresh)
	// sysinfo
	router.GET(sysinfoURI, sysInfo)
	// health
	router.GET(healthURI, health)

	// 初始化管道指针切片
	// 创建备份协程管道指针切片
	bakChans := make([]*chan int, 0)
	// 创建自动热加载配置协程管道指针切片
	autoReloadChans := make([]*chan int, 0)
	// 管道切片长度，用于指定管道索引和管道长度
	chansLen := 0

	// 初始化各web服务配置信息管控接口
	for i := 0; i < len(BifrostConf.WebServerInfo.Servers); i++ {
		switch BifrostConf.WebServerInfo.Servers[i].ServerType {
		case NGINX:

			// 创建备份协程管道及启动备份协程
			bakChan := make(chan int)
			bakChans = append(bakChans, &bakChan)

			// 创建自动热加载配置协程管道及启动自动热加载配置协程
			autoReloadChan := make(chan int)
			autoReloadChans = append(autoReloadChans, &autoReloadChan)

			nginxConfAPIInit(router, &BifrostConf.WebServerInfo.Servers[i], &chansLen, &bakChans, &autoReloadChans)

		case HTTPD:
			// TODO: apache httpd配置解析器
			continue
		default:
			continue
		}
	}

	// 监控系统信息
	monitorChan := make(chan int)
	go monitoring(monitorChan)

	// 打印启动logo
	fmt.Println(logoStr)

	Log(DEBUG, "Listening ports for bifrost")
	go func() {
		// 启动bifrost各接口服务
		rErr := router.Run(fmt.Sprintf(":%d", BifrostConf.WebServerInfo.ListenPort))
		if rErr != nil {
			killCoroutine(chansLen, bakChans, autoReloadChans)
			monitorChan <- 9
			// 输出子任务运行错误
			Log(CRITICAL, "bifrost Run coroutine has been stoped. Cased by '%s'", rErr)
			return
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
	killCoroutine(chansLen, bakChans, autoReloadChans)
	monitorChan <- 9
	Log(NOTICE, "bifrost Run coroutine has been stoped.")
	return
}

// nginxConfAPIInit, nginx配置管控接口初始化函数
// 参数:
//     router: gin引擎指针
//     serverInfo: bifrost配置对象中web服务器信息配置对象指针
//     chansLen: 管道切片长度值指针
//     bakChans: 备份管道指针切片指针
//     autoReloadChans: 自动热加载管道指针切片指针
func nginxConfAPIInit(router *gin.Engine, serverInfo *ServerInfo, chansLen *int, bakChans *[]*chan int, autoReloadChans *[]*chan int) {

	//fmt.Println("初始化", serverInfo.Name, "web服务相关接口。。。")
	loadErr := serverInfo.ngLoad()
	if loadErr != nil {
		Log(ERROR, "[%s] load config error: %s", serverInfo.Name, loadErr)
		*chansLen++
		return
	}
	//ngConfig, confErr := serverInfo.confCaches.GetConfig(serverInfo.ConfPath)
	//if confErr != nil {
	//	Log(ERROR, "[%s] load config error: %s", serverInfo.Name, confErr)
	//	*chansLen++
	//	return
	//}
	// 检查nginx配置是否能被正常解析为json
	//fmt.Println("校验nginx配置。。。")
	_, jerr := json.Marshal(serverInfo.nginxConfig)
	if jerr != nil {
		Log(CRITICAL, "[%s] coroutine has been stoped. Cased by '%s'", serverInfo.Name, jerr)
		*chansLen++
		return
	}
	//fmt.Println("获取接口URI")
	apiURI := serverInfo.BaseURI
	//fmt.Println("生成统计接口URI")
	statisticsURI := fmt.Sprintf("%s/statistics", apiURI)

	//fmt.Println("载入备份协程")
	go serverInfo.Bak(*(*bakChans)[*chansLen])
	//fmt.Println("载入自动更新配置协程")
	go serverInfo.AutoReload(*(*autoReloadChans)[*chansLen])

	*chansLen++

	// view
	//fmt.Println("载入查询接口")
	router.GET(apiURI, func(c *gin.Context) {
		h, status := view(serverInfo, c)
		c.JSON(status, &h)
	})

	// update
	//fmt.Println("载入更新接口")
	router.POST(apiURI, func(c *gin.Context) {
		h, status := update(serverInfo, c)
		c.JSON(status, &h)
	})

	// statistics
	//fmt.Println("载入统计接口")
	router.GET(statisticsURI, func(c *gin.Context) {
		h, status := statisticsView(serverInfo, c)
		c.JSON(status, &h)
	})
}

// killCoroutine, 关闭协程任务函数
// 参数:
//     chansLen: 管道指针切片元素个数
//     chansArray: 传入的管道指针切片集合
func killCoroutine(chansLen int, chansArray ...[]*chan int) {
	for i := 0; i < chansLen; i++ {
		for _, chans := range chansArray {
			*chans[i] <- 9
		}
	}
}

// statisticsView, nginx配置统计信息查询接口函数
// 参数:
//     si: ServerInfo对象指针
//     c: gin上下文对象指针
// 返回值:
//     h: gin.H
//     s: http状态码
func statisticsView(si *ServerInfo, c *gin.Context) (h gin.H, s int) {
	// 初始化h
	status := "unkown"
	message := make(gin.H, 0)
	h = gin.H{
		"status":  &status,
		"message": &message,
	}

	// 获取接口访问参数
	token, hasToken := c.GetQuery("token")
	if !hasToken {
		status = "failed"
		h["message"] = "Token cannot be empty"
		s = http.StatusBadRequest
		Log(NOTICE, "[%s] request without token", c.ClientIP())
		return
	}

	// 校验token
	_, verifyErr := verifyAction(token)
	if verifyErr != nil {
		status = "failed"
		h["message"] = verifyErr.Error()
		s = http.StatusBadRequest
		Log(WARN, "[%s] %s", si.Name, message)
		return
	}

	/*	// 检查接口访问统计查询过滤参数
		noHttpSvrsNum, notHasHttpSvrsNum := c.GetQuery("NoHttpSvrsNum")
		if !notHasHttpSvrsNum {
			noHttpSvrsNum = "false"
		}

		noHttpSvrNames, notHasHttpSvrNames := c.GetQuery("NoHttpSvrNames")
		if !notHasHttpSvrNames {
			noHttpSvrNames = "false"
		}

		noHttpPorts, notHasHttpPorts := c.GetQuery("NoHttpPorts")
		if !notHasHttpPorts {
			noHttpPorts = "false"
		}

		noLocNum, notHasLocationsNum := c.GetQuery("NoLocsNum")
		if !notHasLocationsNum {
			noLocNum = "false"
		}

		noStreamSvrsNum, notHasStreamSvrsNum := c.GetQuery("NoStreamSvrsNum")
		if !notHasStreamSvrsNum {
			noStreamSvrsNum = "false"
		}

		noStreamPorts, notHasStreamPorts := c.GetQuery("NoStreamPorts")
		if !notHasStreamPorts {
			noStreamPorts = "false"
		}

		// 判断统计查询过滤参数是否为bool型
		if (noHttpSvrsNum == "true" || noHttpSvrsNum == "false") && (noHttpSvrNames == "true" || noHttpSvrNames == "false") && (noHttpPorts == "true" || noHttpPorts == "false") && (noLocNum == "true" || noLocNum == "false") && (noStreamSvrsNum == "true" || noStreamSvrsNum == "false") && (noStreamPorts == "true" || noStreamPorts == "false") {
			// 统计查询过滤参数不可都为真
			if noHttpSvrsNum == "true" && noHttpSvrNames == "true" && noHttpPorts == "true" && noLocNum == "true" && noStreamSvrsNum == "true" && noStreamPorts == "true" {
				status = "failed"
				h["message"] = "invalid params."
				s = http.StatusBadRequest
				return
			}

			// 如果配置解析失败，返回解析失败
			if config == nil {
				status = "failed"
				message := "configuration resolution failed"
				Log(ERROR, "[%s] %s", appName, message)
				s = http.StatusInternalServerError
				return
			}

			// 统计查询未过滤数据
			if noHttpSvrsNum == "false" {
				message["httpSvrsNum"] = nginxStatistics.HTTPServersNum(config)
			}
			if noHttpSvrNames == "false" {
				message["httpSvrNames"] = nginxStatistics.HTTPServerNames(config)
			}
			if noHttpPorts == "false" {
				message["httpPorts"] = nginxStatistics.HTTPPorts(config)
			}
			if noLocNum == "false" {
				message["locNum"] = nginxStatistics.HTTPLocationsNum(config)
			}
			if noStreamSvrsNum == "false" {
				message["streamSvrsNum"] = nginxStatistics.StreamServersNum(config)
			}
			if noStreamPorts == "false" {
				message["streamPorts"] = nginxStatistics.StreamServers(config)
			}

			// 无数据时，返回无数据
			if len(message) == 0 {
				status = "failed"
				h["message"] = "no data"
				s = http.StatusOK
				return
			}
		} else {
			status = "failed"
			h["message"] = "invalid params."
			s = http.StatusBadRequest
			return
		}*/

	httpServersNum, httpServers := nginxStatistics.HTTPServers(si.nginxConfig)
	message["httpSvrsNum"] = httpServersNum
	message["httpSvrs"] = httpServers
	message["httpPorts"] = nginxStatistics.HTTPPorts(si.nginxConfig)
	streamServersNum, streamPorts := nginxStatistics.StreamServers(si.nginxConfig)
	message["streamSvrsNum"] = streamServersNum
	message["streamPorts"] = streamPorts

	status = "success"
	s = http.StatusOK
	return
}

// view, nginx配置查询接口函数
// 参数:
//     si: ServerInfo对象指针
//     c: gin上下文对象指针
// 返回值:
//     h: gin.H
//     s: http返回码
func view(si *ServerInfo, c *gin.Context) (h gin.H, s int) {
	// 初始化h
	status := "unkown"
	var message interface{} = "null"
	h = gin.H{
		"status":  &status,
		"message": &message,
	}

	// 获取查询接口访问token参数
	token, hasToken := c.GetQuery("token")
	if !hasToken {
		status = "failed"
		message = "Token cannot be empty"
		s = http.StatusBadRequest
		Log(NOTICE, "[%s] request without token", c.ClientIP())
		return
	}

	// 校验token
	_, verifyErr := verifyAction(token)
	if verifyErr != nil {
		status = "failed"
		message = verifyErr.Error()
		s = http.StatusBadRequest
		Log(WARN, "[%s] %s", si.Name, message)
		return
	}

	// 获取查询接口访问type参数
	t, typeOK := c.GetQuery("type")
	if !typeOK {
		t = "string"
	}

	s = http.StatusOK
	// 组装查询数据
	switch t {
	case "string":
		// 防止配置对象为空时，接口异常
		if si.nginxConfig == nil {
			status = "failed"
			message = "configuration resolution failed"
			Log(ERROR, "[%s] %s", si.Name, message)
			s = http.StatusInternalServerError
		} else {
			status = "success"
			// 修改string切片为string
			// @ClessLi 多次字符串+= 内存开销过大，决定用strings.Join([]string, string)一次性拼接
			/*var str string
			caches := nginx.NewCaches()
			for _, v := range si.nginxConfig.string(&caches) {
				str += v
			}
			message = str*/
			//caches := nginx.NewCaches()
			//message = strings.Join(si.nginxConfig.string(&caches), "")
			message = strings.Join(si.nginxConfig.String(), "")
			Log(DEBUG, "[%s] %s", si.Name, message)
		}
	case "json":
		// 防止配置对象为空时，接口异常
		if si.nginxConfig == nil {
			status = "failed"
			message = "configuration resolution failed"
			Log(ERROR, "[%s] %s", si.Name, message)
			s = http.StatusInternalServerError
		} else {
			status = "success"
			data, marshalErr := json.Marshal(si.nginxConfig)
			if marshalErr != nil {
				Log(ERROR, "[%s] %s", si.Name, marshalErr)
			}
			message = string(data)
			Log(DEBUG, "[%s] %s", si.Name, message)
		}
	default:
		status = "failed"
		message = fmt.Sprintf("view message type '%s' invalid", t)
		Log(INFO, "[%s] %s", si.Name, message)
		s = http.StatusBadRequest
	}

	return
}

// update, nginx配置更新接口函数
// 参数:
//     si: ServerInfo对象指针
//     c: gin上下文指针
// 返回值:
//     h: gin.H
//     s: http返回码
func update(si *ServerInfo, c *gin.Context) (h gin.H, s int) {
	// 初始化h
	status := "unkown"
	message := "null"
	h = gin.H{
		"status":  &status,
		"message": &message,
	}

	//fmt.Println("获取web服务配置校验二进制文件路径")
	verifyBin, absErr := filepath.Abs(si.VerifyExecPath)
	if absErr != nil {
		s = http.StatusInternalServerError
		message = fmt.Sprintf("The validation process does not exist or is configured incorrectly.")
		status = "failed"
		Log(CRITICAL, "[%s] %s detailed error: %s", si.Name, message, absErr)
		return
	}

	// 获取接口访问token参数
	token, hasToken := c.GetQuery("token")
	if !hasToken {
		status = "failed"
		message = "Token cannot be empty"
		Log(NOTICE, "[%s] request without token", c.ClientIP())
		s = http.StatusBadRequest
		return
	}

	// 校验token
	_, verifyErr := verifyAction(token)
	if verifyErr != nil {
		status = "failed"
		message = verifyErr.Error()
		s = http.StatusBadRequest
		Log(WARN, "[%s] [%s] %s", si.Name, c.ClientIP(), message)
		return
	}

	// 获取接口推送数据文件表单
	//confBytes := c.DefaultPostForm("data", "null")
	file, formFileErr := c.FormFile("data")
	if formFileErr != nil {
		message = fmt.Sprintf("FormFile option invalid: <%s>.", formFileErr)
		Log(WARN, "[%s] [%s] %s", si.Name, c.ClientIP(), message)
		status = "failed"
		s = http.StatusBadRequest
		return
	}

	f, fErr := file.Open()
	if fErr != nil {
		message = fmt.Sprintf("Open file failed: <%s>.", fErr)
		Log(CRITICAL, "[%s] [%s] %s", si.Name, c.ClientIP(), message)
		status = "failed"
		s = http.StatusInternalServerError
		return
	}

	defer f.Close() // 关闭文件对象
	confBytes, rErr := ioutil.ReadAll(f)
	if rErr != nil {
		message := fmt.Sprintf("Read file failed: <%s>.", rErr)
		Log(CRITICAL, "[%s] [%s] %s", si.Name, c.ClientIP(), message)
		h["status"] = "failed"
		h["message"] = message
		s = http.StatusInternalServerError
		return
	}

	// 解析接口传入的nginx配置json数据，并完成备份更新
	if len(confBytes) > 0 {

		Log(NOTICE, "[%s] [%s] unmarshal nginx ng.", c.ClientIP(), si.Name)
		newConfig, ujErr := ngJson.Unmarshal(confBytes)
		if ujErr != nil || len(newConfig.Children) <= 0 || newConfig.Value == "" {
			message = fmt.Sprintf("UnmarshalJson failed. <%s>, data: <%s>", ujErr, confBytes)
			Log(WARN, "[%s] [%s] %s", si.Name, c.ClientIP(), message)
			h["status"] = "failed"
			status = "failed"
			//errChan <- ujErr
			s = http.StatusBadRequest
			return
		}

		delErr := nginx.Delete(si.nginxConfig)
		if delErr != nil {
			message = fmt.Sprintf("Delete nginx ng failed. <%s>", delErr)
			Log(ERROR, "[%s] [%s] %s", si.Name, c.ClientIP(), message)
			status = "failed"
			s = http.StatusInternalServerError
			return
		}

		Log(INFO, "[%s] Deleted old nginx ng.", si.Name)
		Log(INFO, "[%s] Verify new nginx ng.", si.Name)
		newCaches, checkErr := nginx.SaveWithCheck(newConfig, verifyBin)
		if checkErr != nil {
			s = http.StatusInternalServerError
			message = fmt.Sprintf("Nginx ng verify failed. <%s>", checkErr)
			Log(WARN, "[%s] %s", si.Name, message)
			status = "failed"

			Log(INFO, "[%s] Delete new nginx ng.", si.Name)
			delErr := nginx.Delete(newConfig)
			if delErr != nil {
				Log(ERROR, "[%s] Delete new nginx ng failed. <%s>", si.Name, delErr)
				status = "failed"
				message = "New nginx config verify failed. And delete new nginx config failed."
				return
			}

			Log(INFO, "[%s] Rollback nginx ng.", si.Name)
			_, rollbackErr := nginx.Save(si.nginxConfig)
			if rollbackErr != nil {
				Log(CRITICAL, "[%s] Nginx ng rollback failed. <%s>", si.Name, rollbackErr)
				status = "failed"
				message = "New nginx config verify failed. And nginx config rollback failed."
				return
			}

			return
		}
		si.confCaches = newCaches
		si.nginxConfig = newConfig
		si.ConfPath = newConfig.Value

		Log(NOTICE, "[%s] [%s] Nginx Config saved successfully", si.Name, c.ClientIP())
	} else {
		status = "failed"
		message = fmt.Sprintf("Wrong data: <%s>", confBytes)
		Log(WARN, "[%s] [%s] %s", si.Name, c.ClientIP(), message)
		s = http.StatusBadRequest
		return
	}

	status = "success"
	message = "Nginx config updated."
	s = http.StatusOK
	return
}
