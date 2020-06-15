package bifrost

import (
	"encoding/json"
	"fmt"
	ngJson "github.com/ClessLi/go-nginx-conf-parser/pkg/json"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/statistics"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

//func Run(appConfig *NGConfig, ngConfig *resolv.Config, errChan chan error) {
// Run, bifrost启动函数
// 参数:
//     appConfig: nginx配置文件信息对象
//     ngConfig: nginx配置对象指针
func Run(appConfig NGConfig, ngConfig *resolv.Config) {
	defer wg.Done() // 结束进程前关闭协程等待组
	// 检查nginx配置是否能被正常解析为json
	_, jerr := json.Marshal(ngConfig)

	//confBytes, jerr := json.Marshal(ngConfig)
	//confBytes, jerr := json.MarshalIndent(ngConfig, "", "    ")
	if jerr != nil {
		//errChan <- jerr
		Log(CRITICAL, fmt.Sprintf("%s's coroutine has been stoped. Cased by '%s'", ngConfig.Name, jerr))
		return
	}

	// 初始化接口
	loginURI := "/login"
	verifyURI := "/verify"
	refreshURI := "/refresh"
	//apiURI := fmt.Sprintf("%s/:token", appConfig.RelativePath)
	apiURI := appConfig.RelativePath
	statisticsURI := fmt.Sprintf("%s/statistics", apiURI)

	ngBin, absErr := filepath.Abs(appConfig.NginxBin)
	if absErr != nil {
		//errChan <- absErr
		Log(CRITICAL, fmt.Sprintf("%s's coroutine has been stoped. Cased by '%s'", ngConfig.Name, absErr))
		return
	}

	// 创建备份协程管道及启动备份协程
	bakChan := make(chan int)
	go Bak(appConfig, ngConfig, bakChan)

	router := gin.Default()
	// login
	router.GET(loginURI, login)
	// verify
	router.GET(verifyURI, verify)
	// refresh
	router.GET(refreshURI, refresh)
	// view
	router.GET(apiURI, func(c *gin.Context) {
		h, status := view(appConfig.Name, ngConfig, c)
		c.JSON(status, &h)
	})
	// update
	router.POST(apiURI, func(c *gin.Context) {
		h, status := update(appConfig.Name, ngBin, ngConfig, c)
		c.JSON(status, &h)
	})

	// statistics
	router.GET(statisticsURI, func(c *gin.Context) {
		h, status := statisticsView(appConfig.Name, ngConfig, c)
		c.JSON(status, &h)
	})

	rErr := router.Run(fmt.Sprintf(":%d", appConfig.Port))
	if rErr != nil {
		// 关闭备份
		bakChan <- 9
		// 输出子任务运行错误
		//errChan <- rErr
		Log(CRITICAL, fmt.Sprintf("%s's coroutine has been stoped. Cased by '%s'", ngConfig.Name, rErr))
		return
	}

	// 关闭备份
	bakChan <- 9
	//errChan <- nil
	Log(NOTICE, fmt.Sprintf("%s's coroutine has been stoped", ngConfig.Name))
	return
}

// statisticsView, nginx配置统计信息查询接口函数
// 参数:
//     appName: 子进程标题
//     config: 子进程nginx配置对象指针
//     c: gin上下文对象指针
// 返回值:
//     h: gin.H
//     s: http状态码
func statisticsView(appName string, config *resolv.Config, c *gin.Context) (h gin.H, s int) {
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
		Log(NOTICE, fmt.Sprintf("[%s] request without token", c.ClientIP()))
		return
	}

	// 校验token
	_, verifyErr := verifyAction(token)
	if verifyErr != nil {
		status = "failed"
		h["message"] = verifyErr.Error()
		s = http.StatusBadRequest
		Log(WARN, fmt.Sprintf("[%s] %s", appName, message))
		return
	}

	// 检查接口访问统计查询过滤参数
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

		// 统计查询未过滤数据
		if noHttpSvrsNum == "false" {
			message["httpSvrsNum"] = statistics.HTTPServersNum(config)
		}
		if noHttpSvrNames == "false" {
			message["httpSvrNames"] = statistics.HTTPServerNames(config)
		}
		if noHttpPorts == "false" {
			message["httpPorts"] = statistics.HTTPPorts(config)
		}
		if noLocNum == "false" {
			message["locNum"] = statistics.HTTPLocationsNum(config)
		}
		if noStreamSvrsNum == "false" {
			message["streamSvrsNum"] = statistics.StreamServersNum(config)
		}
		if noStreamPorts == "false" {
			message["streamPorts"] = statistics.StreamPorts(config)
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
	}

	status = "success"
	s = http.StatusOK
	fmt.Println(h)
	return
}

// view, nginx配置查询接口函数
// 参数:
//     appName: 子进程标题
//     config: 子进程nginx配置对象指针
//     c: gin上下文对象指针
// 返回值:
//     h: gin.H
//     s: http返回码
func view(appName string, config *resolv.Config, c *gin.Context) (h gin.H, s int) {
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
		Log(NOTICE, fmt.Sprintf("[%s] request without token", c.ClientIP()))
		return
	}

	// 校验token
	_, verifyErr := verifyAction(token)
	if verifyErr != nil {
		status = "failed"
		message = verifyErr.Error()
		s = http.StatusBadRequest
		Log(WARN, fmt.Sprintf("[%s] %s", appName, message))
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
		status = "success"
		// 修改string切片为string
		var str string
		for _, v := range config.String() {
			str += v
		}
		message = str
		Log(DEBUG, fmt.Sprintf("[%s] %s", appName, message))
	case "json":
		status = "success"
		message = config
		Log(DEBUG, fmt.Sprintf("[%s] %s", appName, message))
	default:
		status = "failed"
		message = fmt.Sprintf("view message type '%s' invalid", t)
		Log(INFO, fmt.Sprintf("[%s] %s", appName, message))
		s = http.StatusBadRequest
	}

	return
}

// update, nginx配置更新接口函数
// 参数:
//     appName: 子进程标题
//     ngBin: 子进程nginx可执行文件路径
//     ng: 子进程nginx配置对象指针
//     c: gin上下文指针
// 返回值:
//     h: gin.H
//     s: http返回码
func update(appName, ngBin string, ng *resolv.Config, c *gin.Context) (h gin.H, s int) {
	// 函数执行完毕前清理nginx配置缓存
	defer resolv.ReleaseConfigsCache()
	// 初始化h
	status := "unkown"
	message := "null"
	h = gin.H{
		"status":  &status,
		"message": &message,
	}

	// 获取接口访问token参数
	token, hasToken := c.GetQuery("token")
	if !hasToken {
		status = "failed"
		message = "Token cannot be empty"
		Log(NOTICE, fmt.Sprintf("[%s] request without token", c.ClientIP()))
		s = http.StatusBadRequest
		return
	}

	// 校验token
	_, verifyErr := verifyAction(token)
	if verifyErr != nil {
		status = "failed"
		message = verifyErr.Error()
		s = http.StatusBadRequest
		Log(WARN, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
		return
	}

	// 获取接口推送数据文件表单
	//confBytes := c.DefaultPostForm("data", "null")
	file, formFileErr := c.FormFile("data")
	if formFileErr != nil {
		message = fmt.Sprintf("FormFile option invalid: <%s>.", formFileErr)
		Log(WARN, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
		status = "failed"
		s = http.StatusBadRequest
		return
	}

	f, fErr := file.Open()
	if fErr != nil {
		message = fmt.Sprintf("Open file failed: <%s>.", fErr)
		Log(CRITICAL, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
		status = "failed"
		s = http.StatusInternalServerError
		return
	}

	defer f.Close() // 关闭文件对象
	confBytes, rErr := ioutil.ReadAll(f)
	if rErr != nil {
		message := fmt.Sprintf("Read file failed: <%s>.", rErr)
		Log(CRITICAL, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
		h["status"] = "failed"
		h["message"] = message
		s = http.StatusInternalServerError
		return
	}

	// 解析接口传入的nginx配置json数据，并完成备份更新
	if len(confBytes) > 0 {

		Log(NOTICE, fmt.Sprintf("[%s] [%s] Unmarshal nginx ng.", c.ClientIP(), appName))
		newConfig, ujErr := ngJson.Unmarshal(confBytes, &ngJson.Config{})
		if ujErr != nil || len(newConfig.(*resolv.Config).Children) <= 0 || newConfig.(*resolv.Config).Value == "" {
			message = fmt.Sprintf("UnmarshalJson failed. <%s>, data: <%s>", ujErr, confBytes)
			Log(WARN, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
			h["status"] = "failed"
			status = "failed"
			//errChan <- ujErr
			s = http.StatusBadRequest
			return
		}

		delErr := resolv.Delete(ng)
		if delErr != nil {
			message = fmt.Sprintf("Delete nginx ng failed. <%s>", delErr)
			Log(ERROR, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
			status = "failed"
			s = http.StatusInternalServerError
			return
		}

		Log(INFO, fmt.Sprintf("[%s] Deleted old nginx ng.", appName))
		Log(INFO, fmt.Sprintf("[%s] Verify new nginx ng.", appName))
		checkErr := resolv.Check(newConfig.(*resolv.Config), ngBin)
		if checkErr != nil {
			s = http.StatusInternalServerError
			message = fmt.Sprintf("Nginx ng verify failed. <%s>", checkErr)
			Log(WARN, fmt.Sprintf("[%s] %s", appName, message))
			status = "failed"

			Log(INFO, fmt.Sprintf("[%s] Delete new nginx ng.", appName))
			delErr := resolv.Delete(newConfig.(*resolv.Config))
			if delErr != nil {
				Log(ERROR, fmt.Sprintf("[%s] Delete new nginx ng failed. <%s>", appName, delErr))
				status = "failed"
				message = "New nginx config verify failed. And delete new nginx config failed."
				return
			}

			Log(INFO, fmt.Sprintf("[%s] Rollback nginx ng.", appName))
			rollbackErr := resolv.Save(ng)
			if rollbackErr != nil {
				Log(CRITICAL, fmt.Sprintf("[%s] Nginx ng rollback failed. <%s>", appName, rollbackErr))
				status = "failed"
				message = "New nginx config verify failed. And nginx config rollback failed."
				return
			}

			return
		}

		ng.Value = newConfig.(*resolv.Config).Value
		ng.Children = newConfig.(*resolv.Config).Children
		//ng = newConfig.(*resolv.Config)
		Log(NOTICE, fmt.Sprintf("[%s] [%s] Nginx Config saved successfully", appName, c.ClientIP()))
	} else {
		status = "failed"
		message = fmt.Sprintf("Wrong data: <%s>", confBytes)
		Log(WARN, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
		s = http.StatusBadRequest
		return
	}

	status = "success"
	message = "Nginx config updated."
	s = http.StatusOK
	return
}