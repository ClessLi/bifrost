package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClessLi/go-nginx-conf-parser/internal/pkg/password"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)

const (
	ErrorReasonServerBusy    = "服务器繁忙"
	ErrorReasonRelogin       = "请重新登陆"
	ErrorReasonWrongPassword = "用户或密码错误"
	//ErrorReasonNoneToken     = "请通过认证"
)

type JWTClaims struct { // token里面添加用户信息，验证token后可能会用到用户信息
	jwt.StandardClaims
	UserID      int      `json:"user_id"`
	Password    string   `json:"password"`
	Username    string   `json:"username"`
	FullName    string   `json:"full_name"`
	Permissions []string `json:"permissions"`
}

var (
	ExpireTime = 3600 // token有效期
)

func login(c *gin.Context) {
	username := c.Param("username")
	passwd := c.Param("password")
	claims := &JWTClaims{
		UserID:      1,
		Username:    username,
		Password:    passwd,
		FullName:    username,
		Permissions: []string{},
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()

	status := "unkown"
	var token interface{} = "null"
	var message interface{} = "null"
	h := gin.H{
		"status":  &status,
		"token":   &token,
		"message": &message,
	}

	signedToken, err := getToken(claims)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		status = "failed"
		message = err
		log(NOTICE, fmt.Sprintf("[%s] user <%s> login failed, message is: <%s>", c.ClientIP(), username, message))
		c.JSON(http.StatusNotFound, &h)
		return
	}
	log(NOTICE, fmt.Sprintf("[%s] user <%s> is login, token is: %s", c.ClientIP(), username, signedToken))

	status = "success"
	token = signedToken
	//c.String(http.StatusOK, signedToken)
	c.JSON(http.StatusOK, &h)
}

func verify(c *gin.Context) {
	strToken := c.Param("token")
	status := "unkown"
	message := "null"
	h := gin.H{
		"status":  &status,
		"message": &message,
	}
	claim, err := verifyAction(strToken)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		status = "failed"
		message = err.Error()
		log(NOTICE, fmt.Sprintf("[%s] Verified failed", c.ClientIP()))
		c.JSON(http.StatusNotFound, &h)
		return
	}
	//c.String(http.StatusOK, "Certified user ", claim.Username)
	status = "success"
	message = fmt.Sprintf("Certified user <%s>", claim.Username)
	log(NOTICE, fmt.Sprintf("[%s] %s", c.ClientIP(), message))
	c.JSON(http.StatusOK, &h)
}

func refresh(c *gin.Context) {
	strToken := c.Param("token")
	status := "unkown"
	var token interface{} = "null"
	var message interface{} = "null"
	h := gin.H{
		"status":  &status,
		"token":   &token,
		"message": &message,
	}
	claims, err := verifyAction(strToken)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		status = "failed"
		message = err
		log(NOTICE, fmt.Sprintf("[%s] Verified failed", c.ClientIP()))
		c.JSON(http.StatusNotFound, &h)
		return
	}
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	signedToken, err := getToken(claims)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		status = "failed"
		message = err
		log(NOTICE, fmt.Sprintf("[%s] refresh token failed", c.ClientIP()))
		c.JSON(http.StatusNotFound, &h)
		return
	}
	//c.String(http.StatusOK, signedToken)
	status = "success"
	token = signedToken
	c.JSON(http.StatusOK, &h)
}

func verifyAction(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(password.Secret), nil
	})
	if err != nil {
		log(WARN, err.Error())
		//return nil, errors.New(ErrorReasonServerBusy)
		return nil, errors.New(ErrorReasonRelogin)
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New(ErrorReasonRelogin)
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, errors.New(ErrorReasonRelogin)
	}
	log(INFO, fmt.Sprintf("Verify user <%s>...", claims.Username))
	//fmt.Println("verify")
	return claims, nil
}

func getToken(claims *JWTClaims) (string, error) {
	if !validUser(claims) {
		log(WARN, fmt.Sprintf("invalid user <%s> or password <%s>.", claims.Username, claims.Password))
		return "", errors.New(ErrorReasonWrongPassword)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(password.Secret))
	if err != nil {
		log(WARN, err.Error())
		return "", errors.New(ErrorReasonServerBusy)
	}
	return signedToken, nil
}

func validUser(claims *JWTClaims) bool {
	sqlStr := fmt.Sprintf("SELECT `password` FROM `%s`.`user` WHERE `user_name` = \"%s\" LIMIT 1;", dbConfig.DBName, claims.Username)
	checkPasswd, err := getPasswd(sqlStr)
	if err != nil && err != sql.ErrNoRows {
		log(ERROR, err.Error())
		return false
	} else if err == sql.ErrNoRows {
		log(NOTICE, fmt.Sprintf("user <%s> is not exist in ng_admin", claims.Username))
		return false
	}

	return password.Password(claims.Password) == checkPasswd
}

func getPasswd(sqlStr string) (string, error) {
	mysqlUrl := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8", dbConfig.User, dbConfig.Password, dbConfig.Protocol, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	//fmt.Println(mysqlUrl)
	db, dbConnErr := sql.Open("mysql", mysqlUrl)
	if dbConnErr != nil {
		log(ERROR, dbConnErr.Error())
		return "", dbConnErr
	}

	defer db.Close()

	rows, queryErr := db.Query(sqlStr)
	if queryErr != nil {
		log(WARN, queryErr.Error())
		return "", queryErr
	}

	_, rowErr := rows.Columns()
	if rowErr == sql.ErrNoRows {
		return "", rowErr
	}

	for rows.Next() {
		var passwd string
		scanErr := rows.Scan(&passwd)
		if scanErr != nil {
			log(WARN, scanErr.Error())
			return "", scanErr
		}

		if passwd != "" {
			return passwd, nil
		}
	}

	return "", errors.New("sql: unkown error")
}
