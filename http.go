package main

import (
	"fmt"
	ngJson "github.com/ClessLi/go-nginx-conf-parser/pkg/json"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

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

func update(appName, ngBin string, config *resolv.Config, c *gin.Context) (h gin.H) {
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
		log(ERROR, fmt.Sprintf("[%s] %s", appName, message))
		status = "failed"
		return
	}

	defer f.Close()
	confBytes, rErr := ioutil.ReadAll(f)
	if rErr != nil {
		message := fmt.Sprintf("Read file failed: <%s>.", rErr)
		log(ERROR, fmt.Sprintf("[%s] %s", appName, message))
		h["status"] = "failed"
		h["message"] = message
		return
	}

	if len(confBytes) > 0 {

		log(DEBUG, fmt.Sprintf("[%s] Unmarshal nginx config.", appName))
		newConfig, ujErr := ngJson.Unmarshal(confBytes, &ngJson.Config{})
		if ujErr != nil || len(newConfig.(*resolv.Config).Children) <= 0 || newConfig.(*resolv.Config).Value == "" {
			message = fmt.Sprintf("UnmarshalJson failed. <%s>, data: <%s>", ujErr, confBytes)
			log(ERROR, fmt.Sprintf("[%s] %s", appName, message))
			h["status"] = "failed"
			status = "failed"
			//errChan <- ujErr
			return
		}

		delErr := resolv.Delete(config)
		if delErr != nil {
			message = fmt.Sprintf("Delete nginx config failed. <%s>", delErr)
			log(WARN, fmt.Sprintf("[%s] %s", appName, message))
			status = "failed"
			return
		}

		log(DEBUG, fmt.Sprintf("[%s] Deleted old nginx config.", appName))
		log(DEBUG, fmt.Sprintf("[%s] Verify new nginx config.", appName))
		checkErr := resolv.Check(newConfig.(*resolv.Config), ngBin)
		if checkErr != nil {
			message = fmt.Sprintf("Nginx config verify failed. <%s>", checkErr)
			log(ERROR, fmt.Sprintf("[%s] %s", appName, message))
			status = "failed"

			log(DEBUG, fmt.Sprintf("[%s] Delete new nginx config.", appName))
			delErr := resolv.Delete(newConfig.(*resolv.Config))
			if delErr != nil {
				log(ERROR, fmt.Sprintf("[%s] Delete new nginx config failed. <%s>", appName, delErr))
				status = "failed"
				message = "New nginx config verify failed. And delete new nginx config failed."
				return
			}

			log(DEBUG, fmt.Sprintf("[%s] Rollback nginx config.", appName))
			rollbackErr := resolv.Save(config)
			if rollbackErr != nil {
				log(ERROR, fmt.Sprintf("[%s] Nginx config rollback failed. <%s>", appName, rollbackErr))
				status = "failed"
				message = "New nginx config verify failed. And nginx config rollback failed."
				return
			}

			return
		}

		config.Value = newConfig.(*resolv.Config).Value
		config.Children = newConfig.(*resolv.Config).Children
		//config = newConfig.(*resolv.Config)
		log(INFO, fmt.Sprintf("[%s] Nginx Config saved successfully", appName))
	} else {
		status = "failed"
		message = fmt.Sprintf("Wrong data: <%s>", confBytes)
		log(ERROR, fmt.Sprintf("[%s] %s", appName, message))
		return
	}

	status = "success"
	message = "Nginx config update."
	return
}
