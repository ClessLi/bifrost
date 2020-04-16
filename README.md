# 项目介绍
go-nginx-conf-parser 是基于golang语言开发的项目，它目前还处于测试阶段，用于对Nginx配置文件解析并提供配置文件展示和修改的接口，支持json、字符串格式与golang结构相互相互转换。该项目持续更新中。目前可用版本为[v0.0.1]()。

# 项目特点
支持将配置文件、json数据、字符串与配置结构体相互转换
配置结构体支持增加、删除、查询（暂实现查询server上下文结构体）
提供配置文件展示和修改的接口

# 使用方法
## 下载
Windows:  [go-nginx-conf-parser.v0_0_1.win_x64.zip](https://github.com/ClessLi/go-nginx-conf-parser/releases/download/v0.0.1/go-nginx-conf-parser.v0_0_1.win_x64.zip)

Linux: [go-nginx-conf-parser.v0_0_1.linux_x64.zip](https://github.com/ClessLi/go-nginx-conf-parser/releases/download/v0.0.1/go-nginx-conf-parser.v0_0_1.linux_x64.zip)

## 应用配置
> configs/ng-conf-info.yml
```
configs:
  -
    name: "nginx-conf-test"
    relativePath: "/ng_conf"
    port: 18080
    confPath: "/usr/local/openresty/nginx/conf/nginx.conf"
    nginxBin: "/usr/local/openresty/nginx/sbin/nginx"
#  -
#    name: "ng-conf-test2"
#    relativePath: "/ng_conf"
#    port: 28080
#    confPath: "/GO_Project/src/go-nginx-conf-parser/test/config_test/nginx.conf"
#    nginxBin: "xxxxxxxxxxxx/nginx"
```

## 命令帮助
```
> ./go-nginx-conf-parser -h
Usage of ./go-nginx-conf-parser:
  -f conf
    	go-nginx-conf-parser ng-conf-info.y(a)ml path. (default "./configs/ng-conf-info.yml")
  -h help
    	this help
```