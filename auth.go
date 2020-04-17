package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClessLi/go-nginx-conf-parser/internal/pkg/password"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	ErrorreasonServerbusy    = "服务器繁忙"
	ErrorreasonRelogin       = "请重新登陆"
	ErrorreasonWrongpassword = "用户或密码错误"
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
	signedToken, err := getToken(claims)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, signedToken)
}

func verify(c *gin.Context) {
	strToken := c.Param("token")
	claim, err := verifyAction(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, "verify,", claim.Username)
}

func refresh(c *gin.Context) {
	strToken := c.Param("token")
	claims, err := verifyAction(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	signedToken, err := getToken(claims)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, signedToken)
}

func verifyAction(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(password.Secret), nil
	})
	if err != nil {
		return nil, errors.New(ErrorreasonServerbusy)
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New(ErrorreasonRelogin)
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, errors.New(ErrorreasonRelogin)
	}
	fmt.Println("verify")
	return claims, nil
}

func getToken(claims *JWTClaims) (string, error) {
	if !validUser(claims) {
		return "", errors.New(ErrorreasonWrongpassword)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(password.Secret))
	if err != nil {
		return "", errors.New(ErrorreasonServerbusy)
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
		log(INFO, fmt.Sprintf("user <%s> of nginx_admin is not exist", claims.Username))
		return false
	}

	return password.Password(claims.Password) == checkPasswd
}

func getPasswd(sqlStr string) (string, error) {
	mysqlUrl := fmt.Sprintf("%s:@%s(%s:%d)/%s", dbConfig.User, dbConfig.Protocol, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	db, dbConnErr := sql.Open("mysql", mysqlUrl)
	if dbConnErr != nil {
		log(ERROR, dbConnErr.Error())
		return "", dbConnErr
	}

	rows, queryErr := db.Query(sqlStr)
	if queryErr != nil {
		log(WARN, queryErr.Error())
		return "", queryErr
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
