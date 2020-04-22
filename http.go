package main

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

func run(appConfig *NGConfig, ngConfig *resolv.Config, errChan chan error) {
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
		log(WARN, fmt.Sprintf("[%s] %s", appName, message))
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
	log(INFO, fmt.Sprintf("[%s] %s", appName, message))
	return
}

func update(appName, ngBin string, ng *resolv.Config, c *gin.Context) (h gin.H) {
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
		log(WARN, fmt.Sprintf("[%s] %s", appName, message))
		return
	}

	//confBytes := c.DefaultPostForm("data", "null")
	file, formFileErr := c.FormFile("data")
	if formFileErr != nil {
		message = fmt.Sprintf("FormFile option invalid: <%s>.", formFileErr)
		log(WARN, fmt.Sprintf("[%s] %s", appName, message))
		status = "failed"
		return
	}

	f, fErr := file.Open()
	if fErr != nil {
		message = fmt.Sprintf("Open file failed: <%s>.", fErr)
		log(CRITICAL, fmt.Sprintf("[%s] %s", appName, message))
		status = "failed"
		return
	}

	defer f.Close()
	confBytes, rErr := ioutil.ReadAll(f)
	if rErr != nil {
		message := fmt.Sprintf("Read file failed: <%s>.", rErr)
		log(CRITICAL, fmt.Sprintf("[%s] %s", appName, message))
		h["status"] = "failed"
		h["message"] = message
		return
	}

	if len(confBytes) > 0 {

		log(NOTICE, fmt.Sprintf("[%s] Unmarshal nginx ng.", appName))
		newConfig, ujErr := ngJson.Unmarshal(confBytes, &ngJson.Config{})
		if ujErr != nil || len(newConfig.(*resolv.Config).Children) <= 0 || newConfig.(*resolv.Config).Value == "" {
			message = fmt.Sprintf("UnmarshalJson failed. <%s>, data: <%s>", ujErr, confBytes)
			log(WARN, fmt.Sprintf("[%s] %s", appName, message))
			h["status"] = "failed"
			status = "failed"
			//errChan <- ujErr
			return
		}

		delErr := resolv.Delete(ng)
		if delErr != nil {
			message = fmt.Sprintf("Delete nginx ng failed. <%s>", delErr)
			log(ERROR, fmt.Sprintf("[%s] %s", appName, message))
			status = "failed"
			return
		}

		log(NOTICE, fmt.Sprintf("[%s] Deleted old nginx ng.", appName))
		log(NOTICE, fmt.Sprintf("[%s] Verify new nginx ng.", appName))
		checkErr := resolv.Check(newConfig.(*resolv.Config), ngBin)
		if checkErr != nil {
			message = fmt.Sprintf("Nginx ng verify failed. <%s>", checkErr)
			log(WARN, fmt.Sprintf("[%s] %s", appName, message))
			status = "failed"

			log(NOTICE, fmt.Sprintf("[%s] Delete new nginx ng.", appName))
			delErr := resolv.Delete(newConfig.(*resolv.Config))
			if delErr != nil {
				log(ERROR, fmt.Sprintf("[%s] Delete new nginx ng failed. <%s>", appName, delErr))
				status = "failed"
				message = "New nginx ng verify failed. And delete new nginx ng failed."
				return
			}

			log(NOTICE, fmt.Sprintf("[%s] Rollback nginx ng.", appName))
			rollbackErr := resolv.Save(ng)
			if rollbackErr != nil {
				log(CRITICAL, fmt.Sprintf("[%s] Nginx ng rollback failed. <%s>", appName, rollbackErr))
				status = "failed"
				message = "New nginx ng verify failed. And nginx ng rollback failed."
				return
			}

			return
		}

		ng.Value = newConfig.(*resolv.Config).Value
		ng.Children = newConfig.(*resolv.Config).Children
		//ng = newConfig.(*resolv.Config)
		log(INFO, fmt.Sprintf("[%s] Nginx Config saved successfully", appName))
	} else {
		status = "failed"
		message = fmt.Sprintf("Wrong data: <%s>", confBytes)
		log(WARN, fmt.Sprintf("[%s] %s", appName, message))
		return
	}

	status = "success"
	message = "Nginx ng update."
	return
}
