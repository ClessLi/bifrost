# [Bifrost](https://github.com/ClessLi/bifrost)

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/ClessLi/bifrost?label=bifrost)](https://github.com/ClessLi/bifrost/releases/latest)
![GitHub Releases](https://img.shields.io/github/downloads/ClessLi/bifrost/latest/total)
[![GitHub](https://img.shields.io/github/license/ClessLi/bifrost)](LICENSE)

# 项目介绍

**Bifrost** 是基于golang语言开发的项目，它目前还处于测试阶段，用于对Nginx配置文件解析并提供配置文件展示和修改的接口，支持json、字符串格式与golang结构相互转换。该项目持续更新中。最新可用测试版本为[v1.0.1-alpha.1](https://github.com/ClessLi/bifrost/tree/v1.0.1-alpha.1) （v1.0.1-alpha.*取消http协议接口，改用gRPC协议接口。该版本仍在开发、测试中）。

# 项目特点

支持将配置文件、json数据、字符串与配置结构体相互转换

配置结构体支持增加、删除、查询

实现了在加载配置或反序列化json时，防止循环读取配置的功能；实现了nginx配置文件后台更新后，自动热加载的功能

提供配置文件展示和修改及配置信息统计查询，及主机系统状况信息查询的gRPC接口

# 合作项目

## [Heimdallr](https://github.com/tanganyu1114/Heimdallr)

基于SDRMS创建的Nginx管理平台，目前已完成nginx配置文件信息，操作系统基础信息的展示

# 使用方法

## 下载地址

bifrost-auth-v1.0.1-alpha.1

> Windows: [bifrost-auth.v1_0_1.alpha_1.win_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.1-alpha.1/bifrost-auth.v1_0_1.alpha_1.win_x64.zip)
> 
> Linux: [bifrost-auth.v1_0_1.alpha_1.linux_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.1-alpha.1/bifrost-auth.v1_0_1.alpha_1.linux_x64.zip)

bifrost-v1.0.1-alpha.1

> Windows: [bifrost.v1_0_1.alpha_1.win_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.1-alpha.1/bifrost.v1_0_1.alpha_1.win_x64.zip)
> 
> Linux: [bifrost.v1_0_1.alpha_1.linux_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.1-alpha.1/bifrost.v1_0_1.alpha_1.linux_x64.zip)

bifrost-v1.0.0-alpha.7

> Windows: [bifrost.v1_0_0.alpha_7.win_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.0-alpha.7/bifrost.v1_0_0.alpha_7.win_x64.zip)
> 
> Linux: [bifrost.v1_0_0.alpha_7.linux_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.0-alpha.7/bifrost.v1_0_0.alpha_7.linux_x64.zip)

bifrost-v1.0.0-alpha.6

> Windows: [bifrost.v1_0_0.alpha_6.win_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.0-alpha.6/bifrost.v1_0_0.alpha_6.win_x64.zip)
> 
> Linux: [bifrost.v1_0_0.alpha_6.linux_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.0-alpha.6/bifrost.v1_0_0.alpha_6.linux_x64.zip)

bifrost-v1.0.0-alpha.5

> Windows: [bifrost.v1_0_0.alpha_5.win_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.0-alpha.5/bifrost.v1_0_0.alpha_5.win_x64.zip)
> 
> Linux: [bifrost.v1_0_0.alpha_5.linux_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.0-alpha.5/bifrost.v1_0_0.alpha_5.linux_x64.zip)

bifrost-v1.0.0-alpha.4

> Windows: [bifrost.v1_0_0.alpha_4.win_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.0-alpha.4/bifrost.v1_0_0.alpha_4.win_x64.zip)
> 
> Linux: [bifrost.v1_0_0.alpha_4.linux_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.0-alpha.4/bifrost.v1_0_0.alpha_4.linux_x64.zip)

bifrost-v1.0.0-alpha.1

> Windows: [bifrost.v1_0_0.alpha_1.win_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.0-alpha.1/bifrost.v1_0_0.alpha_1.win_x64.zip)
> 
> Linux: [bifrost.v1_0_0.alpha_1.linux_x64](https://github.com/ClessLi/bifrost/releases/download/v1.0.0-alpha.1/bifrost.v1_0_0.alpha_1.linux_x64.zip)

bifrost-v0.0.3

> Windows: [bifrost.v0_0_3.win_x64](https://github.com/ClessLi/bifrost/releases/download/v0.0.3/bifrost.v0_0_3.win_x64.zip)
> 
> Linux: [bifrost.v0_0_3.linux_x64](https://github.com/ClessLi/bifrost/releases/download/v0.0.3/bifrost.v0_0_3.linux_x64.zip)

## 应用配置

配置路径

`bifrost: configs/bifrost.yml`

`bifrost-auth: configs/auth.yml`

配置示例

`bifrost-auth`
```yaml
AuthService:
  Port: 12320
  AuthDBConfig: # 可选，未指定时将考虑AuthConfig
    DBName: "bifrost"
    host: "127.0.0.1"
    port: 3306
    protocol: "tcp"
    user: "heimdall"
    password: "Bultgang"
  AuthConfig: # 可选，未指定AuthDBConfig和AuthConfig时，将以"heimdall/Bultgang"作为默认认证信息
    username: "heimdall"
    password: "Bultgang"
LogConfig:
  logDir: "./logs"
  level: 2
```

`bifrost`
```yaml
Service:
  Port: 12321
  ChunkSize: 4194304
  AuthServerAddr: "127.0.0.1:12320"
  Infos:
    -
      name: "bifrost-test"
      type: nginx
      backupCycle: 1
      backupSaveTime: 7
      backupDir: # 可选，空或未选用时默认以web应用主配置文件所在目录为准
      confPath: "/usr/local/openresty/nginx/conf/nginx.conf"
      verifyExecPath: "/usr/local/openresty/nginx/sbin/nginx"
    -
      name: "bifrost-test2"
      type: nginx
      backupCycle: 1
      backupSaveTime: 7
      confPath: "/GO_Project/src/bifrost/test/config_test/nginx.conf"
      verifyExecPath: "xxxxxxxxxxxx/nginx"
LogConfig:
  logDir: "./logs"
  level: 2
```

## 命令帮助

`bifrost-auth`
```
> ./bifrost-auth -h
  bifrost-auth version: v1.0.1-alpha.1
  Usage: ./bifrost-auth [-hv] [-f filename] [-s signal]
  
  Options:
    -f config
      	the bifrost-auth configuration file path. (default "./configs/auth.yml")
    -h help
      	this help
    -s signal
      	send signal to a master process: stop, restart, status
    -v version
      	this version
```

`bifrost`
```
> ./bifrost -h
  bifrost version: v1.0.1-alpha.1
  Usage: ./bifrost [-hv] [-f filename] [-s signal]
  
  Options:
    -f config
      	the bifrost configuration file path. (default "./configs/bifrost.yml")
    -h help
      	this help
    -s signal
      	send signal to a master process: stop, restart, status
    -v version
      	this version 
```

## 接口文档

支持web服务器（暂仅支持nginx）配置文件查看、序列化导出（json）、配置更新、配置统计信息查看、web服务器状态信息查看，及web服务器（暂仅支持nginx）日志监看功能

详见

[bifrost_gRPC接口定义](api/protobuf-spec/bifrostpb/bifrost.proto)

注：gRPC服务端口侦听为bifrost配置中Service.Port值。

[auth_gRPC接口定义](api/protobuf-spec/authpb/auth.proto)

注：gRPC服务端口侦听为auth配置中AuthService.Port值。

### ~~接口示例~~(已改用客户端方式，详见[客户端使用示例](#test))

~~详见~~[~~bifrost_gRPC接口示例~~](test/grpc_client)

## 客户端

结合go-kit框架编写的客户端对象

### 认证客户端

通过"pkg/client/auth/client.NewClient"函数可生成认证服务客户端

详见[认证客户端](pkg/client/auth/client.go)

### bifrost客户端

通过"pkg/client/bifrost/client.NewClient"函数可生成bifrost服务客户端

详见[bifrost客户端](pkg/client/bifrost/client.go)

<h3 id="test">客户端使用示例</h3>

详见[客户端测试示例](test/grpc_client)