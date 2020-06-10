package main

// DONE:1.权限管理
// DONE:2.nginx配置定期备份机制
// DONE:3.日志规范化输出

import (
	"fmt"
	"github.com/ClessLi/go-nginx-conf-parser/internal/pkg/bifrost"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
)

func main() {
	defer bifrost.Logf.Close()
	for _, ngConfig := range bifrost.Configs.NGConfigs {
		ng, err := resolv.Load(ngConfig.ConfPath)

		if err != nil {
			fmt.Println(err)
			continue
		}

		errChan := make(chan error)

		go bifrost.Run(&ngConfig, ng, errChan)

		err = <-errChan
		if err != nil {
			bifrost.Log(bifrost.CRITICAL, fmt.Sprintf("%s's coroutine has been stoped. Cased by <%s>", ngConfig.Name, err))
		} else {
			bifrost.Log(bifrost.NOTICE, fmt.Sprintf("%s's coroutine has been stoped", ngConfig.Name))
		}
	}

}
