package bifrost

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/password"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)

const (
	// 认证接口错误返回
	ErrorReasonServerBusy    = "服务器繁忙"
	ErrorReasonRelogin       = "请重新登陆"
	ErrorReasonWrongPassword = "用户或密码错误"
	//ErrorReasonNoneToken     = "请通过认证"
)

// JWTClaims, jwt断言对象，定义认证接口校验的用户信息
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

// login, 用户登录函数，定义用户登录认证接口函数
// 参数:
//     c: gin.Context 对象指针
func login(c *gin.Context) {
	// 初始化
	status := "unkown"
	var token interface{} = "null"
	var message interface{} = "null"
	h := gin.H{
		"status":  &status,
		"token":   &token,
		"message": &message,
	}

	// 获取接口请求传参
	username, hasusername := c.GetQuery("username")
	passwd, haspasswd := c.GetQuery("password")
	if !hasusername || !haspasswd {
		status = "failed"
		message = "check your username or password"
		Log(NOTICE, fmt.Sprintf("[%s] login failed, message is: '%s'", c.ClientIP(), message))
		c.JSON(http.StatusBadRequest, &h)
		return
	}
	// 判断区分临时、永久令牌

	isUnExp := false
	unexpired, hasunexp := c.GetQuery("unexpired")
	if !hasunexp {
		unexpired = "false"
	}
	switch unexpired {
	case "true":
		isUnExp = true
	case "false":
		isUnExp = false
	default:
		status = "failed"
		message = fmt.Sprintf("invalid param unexpired=%s", unexpired)
		Log(NOTICE, fmt.Sprintf("[%s] login failed, message is: '%s'", c.ClientIP(), message))
		c.JSON(http.StatusBadRequest, &h)
		return
	}

	// 初始化jwt断言对象
	claims := &JWTClaims{
		UserID:      1,
		Username:    username,
		Password:    passwd,
		FullName:    username,
		Permissions: []string{},
	}
	claims.IssuedAt = time.Now().Unix()
	if isUnExp {
		claims.ExpiresAt = 0
	} else {
		claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()
	}

	// 认证用户信息
	if !validUser(claims) {
		Log(WARN, fmt.Sprintf("Invalid user '%s' or password '%s'.", claims.Username, claims.Password))
		status = "failed"
		message = ErrorReasonWrongPassword
		c.JSON(http.StatusOK, &h)
		return
	}

	// 生成用户token
	signedToken, err := getToken(claims)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		status = "failed"
		message = err.Error()
		Log(NOTICE, fmt.Sprintf("[%s] user '%s' login failed, message is: '%s'", c.ClientIP(), username, message))
		c.JSON(http.StatusOK, &h)
		return
	}
	Log(NOTICE, fmt.Sprintf("[%s] user '%s' is login, token is: %s", c.ClientIP(), username, signedToken))

	status = "success"
	token = signedToken
	//c.String(http.StatusOK, signedToken)
	c.JSON(http.StatusOK, &h)
}

// verify, token校验函数，定义token校验认证接口函数
// 参数:
//     c: gin.Context 对象指针
func verify(c *gin.Context) {
	// 初始化
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
	claim, err := verifyAction(strToken)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		status = "failed"
		message = err.Error()
		Log(NOTICE, fmt.Sprintf("[%s] Verified failed", c.ClientIP()))
		c.JSON(http.StatusNotFound, &h)
		return
	}
	status = "success"
	message = fmt.Sprintf("Certified user '%s'", claim.Username)
	Log(NOTICE, fmt.Sprintf("[%s] %s", c.ClientIP(), message))
	c.JSON(http.StatusOK, &h)
}

// refresh, token更新函数，定义token更新认证接口函数
// 参数:
//     c: gin.Context 对象指针
func refresh(c *gin.Context) {
	// 初始化
	status := "unkown"
	var token interface{} = "null"
	var message interface{} = "null"
	h := gin.H{
		"status":  &status,
		"token":   &token,
		"message": &message,
	}

	// 获取接口传参
	strToken, hasToken := c.GetQuery("token")
	if !hasToken {
		status = "failed"
		message = "Token cannot be empty"
		Log(NOTICE, fmt.Sprintf("[%s] token refresh failed, message is: '%s'", c.ClientIP(), message))
		c.JSON(http.StatusBadRequest, &h)
		return
	}

	// 校验token
	claims, err := verifyAction(strToken)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		status = "failed"
		message = err.Error()
		Log(NOTICE, fmt.Sprintf("[%s] Verified failed", c.ClientIP()))
		c.JSON(http.StatusNotFound, &h)
		return
	}

	// 重新生成token
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	signedToken, err := getToken(claims)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		status = "failed"
		message = err
		Log(NOTICE, fmt.Sprintf("[%s] refresh token failed", c.ClientIP()))
		c.JSON(http.StatusNotFound, &h)
		return
	}
	//c.String(http.StatusOK, signedToken)
	status = "success"
	token = signedToken
	c.JSON(http.StatusOK, &h)
}

// verifyAction, 认证token有效性函数
// 参数:
//     strToken: token字符串
// 返回值:
//     用户jwt断言对象指针
//     错误
func verifyAction(strToken string) (*JWTClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(password.Secret), nil
	})
	if err != nil {
		Log(WARN, err.Error())
		//return nil, errors.New(ErrorReasonServerBusy)
		return nil, errors.New(ErrorReasonRelogin)
	}

	// 转换jwt断言对象
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New(ErrorReasonRelogin)
	}
	Log(INFO, fmt.Sprintf("Verify user '%s'...", claims.Username))

	// 认证用户信息
	if !validUser(claims) {
		Log(WARN, fmt.Sprintf("Invalid user '%s' or password '%s'.", claims.Username, claims.Password))
		return nil, errors.New(ErrorReasonWrongPassword)
	}

	if err := token.Claims.Valid(); err != nil {
		return nil, errors.New(ErrorReasonRelogin)
	}
	Log(INFO, fmt.Sprintf("Username '%s' passed verification", claims.Username))

	// 通过返回有效用户jwt断言对象
	return claims, nil
}

// getToken, token生成函数，根据jwt断言对象编码为token
// 参数:
//     claims: 用户jwt断言对象指针
// 返回值:
//     token字符串
//     错误
func getToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(password.Secret))
	if err != nil {
		Log(WARN, err.Error())
		return "", errors.New(ErrorReasonServerBusy)
	}
	return signedToken, nil
}

// validUser, 用户认证函数，判断用户是否有效
// 参数:
//     claims: 用户jwt断言对象指针
// 返回值:
//     用户是否有效
func validUser(claims *JWTClaims) bool {
	if authDBConfig == nil {
		return claims.Username == authConfig.Username && claims.Password == authConfig.Password
	}
	sqlStr := fmt.Sprintf("SELECT `password` FROM `%s`.`user` WHERE `user_name` = \"%s\" LIMIT 1;", authDBConfig.DBName, claims.Username)
	checkPasswd, err := getPasswd(sqlStr)
	if err != nil && err != sql.ErrNoRows {
		Log(ERROR, err.Error())
		return false
	} else if err == sql.ErrNoRows {
		Log(NOTICE, fmt.Sprintf("user '%s' is not exist in bifrost", claims.Username))
		return false
	}

	return password.Password(claims.Password) == checkPasswd
}

// getPasswd, 用户密码查询函数
// 参数:
//     sqlStr: 查询语句
// 返回值:
//     用户加密密码
//     错误
func getPasswd(sqlStr string) (string, error) {
	mysqlUrl := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8", authDBConfig.User, authDBConfig.Password, authDBConfig.Protocol, authDBConfig.Host, authDBConfig.Port, authDBConfig.DBName)
	//fmt.Println(mysqlUrl)
	db, dbConnErr := sql.Open("mysql", mysqlUrl)
	if dbConnErr != nil {
		Log(ERROR, dbConnErr.Error())
		return "", dbConnErr
	}

	defer db.Close()

	rows, queryErr := db.Query(sqlStr)
	if queryErr != nil {
		Log(WARN, queryErr.Error())
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
			Log(WARN, scanErr.Error())
			return "", scanErr
		}

		if passwd != "" {
			return passwd, nil
		}
	}

	return "", errors.New("sql: unkown error")
}
