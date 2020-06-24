package bifrost

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func health(c *gin.Context) {
	status := "unkown"
	message := "null"
	h := gin.H{
		"status":  &status,
		"message": &message,
	}

	// 获取接口传参
	strToken, hasToken := c.GetQuery("token")
	if !hasToken {
		status = "failed"
		message = "Token cannot be empty"
		Log(NOTICE, fmt.Sprintf("[%s] token verify failed, message is: '%s'", c.ClientIP(), message))
		c.JSON(http.StatusBadRequest, &h)
		return
	}

	// 校验token
	_, err := verifyAction(strToken)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		status = "failed"
		message = err.Error()
		Log(NOTICE, fmt.Sprintf("[%s] Verified failed", c.ClientIP()))
		c.JSON(http.StatusNotFound, &h)
		return
	}

	status = "success"
	if isHealthy {
		message = "healthy"
	} else {
		message = "unhealthy"
	}
	c.JSON(http.StatusOK, &h)
}
