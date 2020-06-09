package ng_conf_admin

import (
	"encoding/json"
	"fmt"
	ngJson "github.com/ClessLi/go-nginx-conf-parser/pkg/json"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func Run(appConfig *NGConfig, ngConfig *resolv.Config, errChan chan error) {
	_, jerr := json.Marshal(ngConfig)
	//confBytes, jerr := json.Marshal(ngConfig)
	//confBytes, jerr := json.MarshalIndent(ngConfig, "", "    ")
	if jerr != nil {
		errChan <- jerr
	}

	loginURI := "/login/:username/:password"
	verifyURI := "/verify/:token"
	refreshURI := "/refresh/:token"
	apiURI := fmt.Sprintf("%s/:token", appConfig.RelativePath)

	ngBin, absErr := filepath.Abs(appConfig.NginxBin)
	if absErr != nil {
		errChan <- absErr
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
		h := view(appConfig.Name, ngConfig, c)
		c.JSON(http.StatusOK, &h)
	})
	// update
	router.POST(apiURI, func(c *gin.Context) {
		h := update(appConfig.Name, ngBin, ngConfig, c)
		c.JSON(http.StatusOK, &h)
	})

	rErr := router.Run(fmt.Sprintf(":%d", appConfig.Port))
	if rErr != nil {
		// 关闭备份
		bakChan <- 9
		// 输出子任务运行错误
		errChan <- rErr
	}

	// 关闭备份
	bakChan <- 9
	errChan <- nil

}

func view(appName string, config *resolv.Config, c *gin.Context) (h gin.H) {
	status := "unkown"
	var message interface{} = "null"
	h = gin.H{
		"status":  &status,
		"message": &message,
	}
	//token, tokenOK := c.GetQuery("token")
	//if tokenOK {
	//	_, verifyErr := verifyAction(token)
	//	if verifyErr != nil {
	//		status = "failed"
	//		message = verifyErr
	//		return
	//	}
	//} else {
	//	status = "failed"
	//	message = ErrorReasonNoneToken
	//	return
	//}
	token := c.Param("token")
	_, verifyErr := verifyAction(token)
	if verifyErr != nil {
		status = "failed"
		message = verifyErr
		Log(WARN, fmt.Sprintf("[%s] %s", appName, message))
		return
	}

	t, typeOK := c.GetQuery("type")
	if !typeOK {
		t = "string"
	}

	switch t {
	case "string":
		status = "success"
		message = config.String()
	case "json":
		status = "success"
		message = config
	default:
		status = "failed"
		message = fmt.Sprintf("view message type <%s> invalid", t)
	}
	Log(INFO, fmt.Sprintf("[%s] %s", appName, message))
	return
}

func update(appName, ngBin string, ng *resolv.Config, c *gin.Context) (h gin.H) {
	defer resolv.ReleaseConfigsCache()
	status := "unkown"
	message := "null"
	h = gin.H{
		"status":  &status,
		"message": &message,
	}
	//token, tokenOK := c.GetQuery("token")
	//if tokenOK {
	//	_, verifyErr := verifyAction(token)
	//	if verifyErr != nil {
	//		status = "failed"
	//		message = verifyErr.Error()
	//		return
	//	}
	//} else {
	//	status = "failed"
	//	message = ErrorReasonNoneToken
	//	return
	//}
	token := c.Param("token")
	_, verifyErr := verifyAction(token)
	if verifyErr != nil {
		status = "failed"
		message = verifyErr.Error()
		Log(WARN, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
		return
	}

	//confBytes := c.DefaultPostForm("data", "null")
	file, formFileErr := c.FormFile("data")
	if formFileErr != nil {
		message = fmt.Sprintf("FormFile option invalid: <%s>.", formFileErr)
		Log(WARN, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
		status = "failed"
		return
	}

	f, fErr := file.Open()
	if fErr != nil {
		message = fmt.Sprintf("Open file failed: <%s>.", fErr)
		Log(CRITICAL, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
		status = "failed"
		return
	}

	defer f.Close()
	confBytes, rErr := ioutil.ReadAll(f)
	if rErr != nil {
		message := fmt.Sprintf("Read file failed: <%s>.", rErr)
		Log(CRITICAL, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
		h["status"] = "failed"
		h["message"] = message
		return
	}

	if len(confBytes) > 0 {

		Log(NOTICE, fmt.Sprintf("[%s] [%s] Unmarshal nginx ng.", c.ClientIP(), appName))
		newConfig, ujErr := ngJson.Unmarshal(confBytes, &ngJson.Config{})
		if ujErr != nil || len(newConfig.(*resolv.Config).Children) <= 0 || newConfig.(*resolv.Config).Value == "" {
			message = fmt.Sprintf("UnmarshalJson failed. <%s>, data: <%s>", ujErr, confBytes)
			Log(WARN, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
			h["status"] = "failed"
			status = "failed"
			//errChan <- ujErr
			return
		}

		delErr := resolv.Delete(ng)
		if delErr != nil {
			message = fmt.Sprintf("Delete nginx ng failed. <%s>", delErr)
			Log(ERROR, fmt.Sprintf("[%s] [%s] %s", appName, c.ClientIP(), message))
			status = "failed"
			return
		}

		Log(INFO, fmt.Sprintf("[%s] Deleted old nginx ng.", appName))
		Log(INFO, fmt.Sprintf("[%s] Verify new nginx ng.", appName))
		checkErr := resolv.Check(newConfig.(*resolv.Config), ngBin)
		if checkErr != nil {
			message = fmt.Sprintf("Nginx ng verify failed. <%s>", checkErr)
			Log(WARN, fmt.Sprintf("[%s] %s", appName, message))
			status = "failed"

			Log(INFO, fmt.Sprintf("[%s] Delete new nginx ng.", appName))
			delErr := resolv.Delete(newConfig.(*resolv.Config))
			if delErr != nil {
				Log(ERROR, fmt.Sprintf("[%s] Delete new nginx ng failed. <%s>", appName, delErr))
				status = "failed"
				message = "New nginx ng verify failed. And delete new nginx ng failed."
				return
			}

			Log(INFO, fmt.Sprintf("[%s] Rollback nginx ng.", appName))
			rollbackErr := resolv.Save(ng)
			if rollbackErr != nil {
				Log(CRITICAL, fmt.Sprintf("[%s] Nginx ng rollback failed. <%s>", appName, rollbackErr))
				status = "failed"
				message = "New nginx ng verify failed. And nginx ng rollback failed."
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
		return
	}

	status = "success"
	message = "Nginx ng update."
	return
}