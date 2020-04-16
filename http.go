package main

import (
	"fmt"
	ngJson "github.com/ClessLi/go-nginx-conf-parser/pkg/json"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func GET(config *resolv.Config, c *gin.Context) (h gin.H) {
	t, ok := c.GetQuery("type")
	if !ok {
		t = "string"
	}

	h = gin.H{
		"status":  "unkown",
		"message": "null",
	}
	switch t {
	case "string":
		h["status"] = "success"
		h["message"] = config.String()
	case "json":
		h["status"] = "success"
		h["message"] = config
	default:
		h["status"] = "failed"
		h["message"] = fmt.Sprintf("GET message type <%s> invalid", t)
	}
	return
}

func POST(ngBin string, config *resolv.Config, c *gin.Context) (h gin.H) {
	status := "unkown"
	message := "null"
	h = gin.H{
		"status":  &status,
		"message": &message,
	}

	//confBytes := c.DefaultPostForm("data", "null")
	file, formFileErr := c.FormFile("data")
	if formFileErr != nil {
		message := fmt.Sprintf("FormFile option invalid: <%s>.", formFileErr)
		log(WARN, message)
		h["massage"] = message
		h["status"] = "failed"
		return
	}

	f, fErr := file.Open()
	if fErr != nil {
		message = fmt.Sprintf("Open file failed: <%s>.", fErr)
		log(ERROR, message)
		status = "failed"
		return
	}

	defer f.Close()
	confBytes, rErr := ioutil.ReadAll(f)
	if rErr != nil {
		message := fmt.Sprintf("Read file failed: <%s>.", rErr)
		log(ERROR, message)
		h["status"] = "failed"
		h["message"] = message
		return
	}

	if len(confBytes) > 0 {

		log(DEBUG, fmt.Sprintf("Unmarshal nginx config."))
		newConfig, ujErr := ngJson.Unmarshal(confBytes, &ngJson.Config{})
		if ujErr != nil || len(newConfig.(*resolv.Config).Children) <= 0 || newConfig.(*resolv.Config).Value == "" {
			message = fmt.Sprintf("UnmarshalJson failed. <%s>, data: <%s>", ujErr, confBytes)
			//log(ERROR, fmt.Sprintf("UnmarshalJson failed. <%s>, data: <%s>", ujErr, confBytes))
			log(ERROR, message)
			status = "failed"
			//errChan <- ujErr
			return
		}

		bakPath, bErr := resolv.Backup(config, "")
		if bErr != nil {
			//log(WARN, fmt.Sprintf("Nginx Config backup to %s, but failed. <%s>\n", bakPath, bErr))
			message = fmt.Sprintf("Nginx Config backup to %s, but failed. <%s>\n", bakPath, bErr)
			log(WARN, message)
			//errChan <- bErr
			status = "failed"
			return
		}

		log(INFO, fmt.Sprintf("Nginx Config backup to %s\n", bakPath))
		delErr := resolv.Delete(config)
		if delErr != nil {
			//log(WARN, fmt.Sprintf("Delete nginx config failed. <%s>", delErr))
			message = fmt.Sprintf("Delete nginx config failed. <%s>", delErr)
			log(WARN, message)
			status = "failed"
			return
		}

		log(DEBUG, fmt.Sprintf("Deleted old nginx config."))
		log(DEBUG, fmt.Sprintf("Verify new nginx config."))
		checkErr := resolv.Check(newConfig.(*resolv.Config), ngBin)
		if checkErr != nil {
			//log(WARN, fmt.Sprintf("Nginx config verify failed. <%s>", checkErr))
			message = fmt.Sprintf("Nginx config verify failed. <%s>", checkErr)
			status = "failed"

			log(DEBUG, fmt.Sprintf("Delete new nginx config."))
			delErr := resolv.Delete(newConfig.(*resolv.Config))
			//fmt.Println(newConfig.String())
			if delErr != nil {
				log(WARN, fmt.Sprintf("Delete new nginx config failed. <%s>", delErr))
				//message = fmt.Sprintf("Delete new nginx config failed. <%s>", delErr)
				//log(WARN, message)
				status = "failed"
				message = "New nginx config verify failed. And delete new nginx config failed."
				return
			}

			log(DEBUG, fmt.Sprintf("Rollback nginx config."))
			rollbackErr := resolv.Save(config)
			if rollbackErr != nil {
				log(ERROR, fmt.Sprintf("Nginx config rollback failed. <%s>", rollbackErr))
				status = "failed"
				message = "New nginx config verify failed. And nginx config rollback failed."
				return
			}

			return
		}

		config.Value = newConfig.(*resolv.Config).Value
		config.Children = newConfig.(*resolv.Config).Children
		//config = newConfig.(*resolv.Config)
		log(INFO, "Nginx Config saved successfully")
	} else {
		status = "failed"
		message = fmt.Sprintf("Wrong data: <%s>", confBytes)
		log(ERROR, message)
		return
	}

	status = "success"
	message = "Nginx config update."
	return
}