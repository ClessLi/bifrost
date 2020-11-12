# [Bifrost](https://github.com/ClessLi/bifrost)

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/ClessLi/bifrost?label=bifrost)](https://github.com/ClessLi/bifrost/releases/latest)
![GitHub Releases](https://img.shields.io/github/downloads/ClessLi/bifrost/latest/total)
[![GitHub](https://img.shields.io/github/license/ClessLi/bifrost)](LICENSE)

# 项目介绍

**Bifrost** 是基于golang语言开发的项目，它目前还处于测试阶段，用于对Nginx配置文件解析并提供配置文件展示和修改的接口，支持json、字符串格式与golang结构相互转换。该项目持续更新中。最新可用版本为[v1.0.0-alpha.7](https://github.com/ClessLi/bifrost/tree/v1.0.0-alpha.7) 。

# 项目特点

支持将配置文件、json数据、字符串与配置结构体相互转换

配置结构体支持增加、删除、查询

实现了在加载配置或反序列化json时，防止循环读取配置的功能；实现了nginx配置文件后台更新后，自动热加载的功能

提供配置文件展示和修改及配置信息统计查询的接口

# 合作项目

## [Heimdallr](https://github.com/tanganyu1114/Heimdallr)

基于SDRMS创建的Nginx管理平台，目前已完成nginx配置文件信息，操作系统基础信息的展示

# 使用方法

## 下载地址

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

`configs/bifrost.yml`

配置示例

```yaml
Service:
  listenPort: 12321
  chunkSize: 4194304
  infos:
    -
      name: "bifrost-test"
      type: nginx
      baseURI: "/ng_conf1"
      backupCycle: 1
      backupSaveTime: 7
      backupDir: # 可选，空或未选用时默认以web应用主配置文件所在目录为准
      confPath: "/usr/local/openresty/nginx/conf/nginx.conf"
      verifyExecPath: "/usr/local/openresty/nginx/sbin/nginx"
    -
      name: "bifrost-test2"
      type: nginx
      baseURI: "/ng_conf2"
      backupCycle: 1
      backupSaveTime: 7
      confPath: "/GO_Project/src/bifrost/test/config_test/nginx.conf"
      verifyExecPath: "xxxxxxxxxxxx/nginx"
AuthService:
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

## 命令帮助

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

详见[bifrost_gRPC接口定义](api/protobuf-spec/bifrostpb/bifrost.proto)

注：gRPC服务端口侦听为bifrost配置中Service.listenPort值。

### 接口示例

详见[bifrost_gRPC接口示例](test/grpc_client)