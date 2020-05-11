package main

// DONE:1.权限管理
// DONE:2.nginx配置定期备份机制
// DONE:3.日志规范化输出

import (
	"fmt"
	"github.com/ClessLi/go-nginx-conf-parser/internal/pkg/ng_conf_admin"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
)

const (
//ERROR      = "ERROR"
//WARN       = "WARN"
//NOTICE       = "NOTICE"
//DEBUG      = "DEBUG"
//timeFormat = "2006-01-02 15:04:05.013"
)

//func init() {
//	flag.Usage = usage
//}

//func usage() {
//	fmt.Fprintf(os.Stdout, `go-nginx-conf-parser version: v0.0.1`)
//	flag.Usage()
//}

func main() {
	defer ng_conf_admin.Logf.Close()
	for _, ngConfig := range ng_conf_admin.Configs.NGConfigs {
		ng, err := resolv.Load(ngConfig.ConfPath)

		if err != nil {
			fmt.Println(err)
			continue
		}

		errChan := make(chan error)

		go ng_conf_admin.Run(&ngConfig, ng, errChan)

		err = <-errChan
		if err != nil {
			ng_conf_admin.Log(ng_conf_admin.CRITICAL, fmt.Sprintf("%s's coroutine has been stoped. Cased by <%s>", ngConfig.Name, err))
		} else {
			ng_conf_admin.Log(ng_conf_admin.NOTICE, fmt.Sprintf("%s's coroutine has been stoped", ngConfig.Name))
		}
	}

}
