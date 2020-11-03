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
WebServerInfo:
  listenPort: 12321
  servers:
    -
      name: "bifrost-test"
      serverType: nginx
      baseURI: "/ng_conf1"
      backupCycle: 1
      backupSaveTime: 7
      backupDir: # 可选，空或未选用时默认以web应用主配置文件所在目录为准
      confPath: "/usr/local/openresty/nginx/conf/nginx.conf"
      verifyExecPath: "/usr/local/openresty/nginx/sbin/nginx"
    -
      name: "bifrost-test2"
      serverType: nginx
      baseURI: "/ng_conf2"
      backupCycle: 1
      backupSaveTime: 7
      confPath: "/GO_Project/src/bifrost/test/config_test/nginx.conf"
      verifyExecPath: "xxxxxxxxxxxx/nginx"
AuthDBConfig: # 可选，未指定时将考虑AuthConfig
  DBName: "bifrost"
  host: "127.0.0.1"
  port: 3306
  protocol: "tcp"
  user: "heimdall"
  password: "Bultgang"
AuthConfig: # 可选，未指定AuthDBConfig和AuthConfig是，将以"heimdall/Bultgang"作为默认认证信息
  username: "heimdall"
  password: "Bultgang"
LogConfig:
  logDir: "./logs"
  level: 2
```

## 命令帮助

```
> ./bifrost -h
  bifrost version: v1.0.0-alpha.7
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

## 应用接口调用方式

### 访问认证接口

#### 1.登录认证接口

接口地址

`http://<Host>:<Port>/login?username=<username>&password=<password>&unexpired=<true|false>`

返回格式

`json`

请求方式

`http get`

请求示例

`http://127.0.0.1:12321/login?username=heimdall&password=Bultgang`

请求参数

| 名称 | 必填 | 类型 | 说明 |
| :-: | :-: | :-: | :- |
| username | 是 | string | 用户名 |
| password | 是 | string | 用户密码 |
| unexpired | 否 | bool | token是否永不过期，默认为false |

返回参数说明

| 名称 | 类型 | 说明 |
| :-: | :-: | :- |
| json返回示例 | - | - |

json返回示例

```json
{
  "message": "null",
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE2MDcwMzcsImlhdCI6MTU5MTYwMzQzNywidXNlcl9pZCI6MSwicGFzc3dvcmQiOiJuZ2FkbWluIiwidXNlcm5hbWUiOiJuZ2FkbWluIiwiZnVsbF9uYW1lIjoibmdhZG1pbiIsInBlcm1pc3Npb25zIjpbXX0.l5qE1sMBD9VJHspzXlhHNmHhbZiF00YlCafUIsIEJpo"
}
```

#### 2.校验会话token接口

接口地址

`http://<Host>:<Port>/verify?token=<token>`

返回格式

`json`

请求方式

`http get`

请求示例

`http://127.0.0.1:12321/verify?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE2MDcwMzcsImlhdCI6MTU5MTYwMzQzNywidXNlcl9pZCI6MSwicGFzc3dvcmQiOiJuZ2FkbWluIiwidXNlcm5hbWUiOiJuZ2FkbWluIiwiZnVsbF9uYW1lIjoibmdhZG1pbiIsInBlcm1pc3Npb25zIjpbXX0.l5qE1sMBD9VJHspzXlhHNmHhbZiF00YlCafUIsIEJpo`

请求参数

| 名称 | 必填 | 类型 | 说明 |
| :-: | :-: | :-: | :- |
| token | 是 | string | 用户访问认证成功后返回的令牌 |

返回参数说明

| 名称 | 类型 | 说明 |
| :-: | :-: | :- |
| json返回示例 | - | - |

json返回示例

```json
{
  "message": "Certified user 'heimdall'",
  "status": "success"
}
```

#### 3.更新会话token接口

接口地址

`http://<Host>:<Port>/refresh?token=<token>`

返回格式

`json`

请求方式

`http get`

请求示例

`http://127.0.0.1:12321/refresh?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE2MDcwMzcsImlhdCI6MTU5MTYwMzQzNywidXNlcl9pZCI6MSwicGFzc3dvcmQiOiJuZ2FkbWluIiwidXNlcm5hbWUiOiJuZ2FkbWluIiwiZnVsbF9uYW1lIjoibmdhZG1pbiIsInBlcm1pc3Npb25zIjpbXX0.l5qE1sMBD9VJHspzXlhHNmHhbZiF00YlCafUIsIEJpo`

请求参数

| 名称 | 必填 | 类型 | 说明 |
| :-: | :-: | :-: | :- |
| token | 是 | string | 用户访问认证成功后返回的令牌 |

返回参数说明

| 名称 | 类型 | 说明 |
| :-: | :-: | :- |
| json返回示例 | - | - |

json返回示例

```json
{
  "message": "null",
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE2MDc3MDEsImlhdCI6MTU5MTYwMzQzNywidXNlcl9pZCI6MSwicGFzc3dvcmQiOiJuZ2FkbWluIiwidXNlcm5hbWUiOiJuZ2FkbWluIiwiZnVsbF9uYW1lIjoibmdhZG1pbiIsInBlcm1pc3Npb25zIjpbXX0.fDoe4v37XyjmrK4wnfhOUnePwJLdszYXveOfoRXyUj8"
}
```

### Nginx配置信息统计查询接口

接口地址

`http://<Host>:<Port>/<baseURI>/statistics?token=<token>`

注：\<baseURI>为ng管理工具配置中WebServerInfo.servers列表各自元素的baseURI子参数值。

返回格式

`json`

请求方式

`http get`

请求示例

`http://127.0.0.1:12321/ng_conf1/statistics?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE2MDc3MDEsImlhdCI6MTU5MTYwMzQzNywidXNlcl9pZCI6MSwicGFzc3dvcmQiOiJuZ2FkbWluIiwidXNlcm5hbWUiOiJuZ2FkbWluIiwiZnVsbF9uYW1lIjoibmdhZG1pbiIsInBlcm1pc3Npb25zIjpbXX0.fDoe4v37XyjmrK4wnfhOUnePwJLdszYXveOfoRXyUj8`

请求参数

| 名称 | 必填 | 类型 | 说明 |
| :-: | :-: | :-: | :- |
| token | 是 | string | 用户访问认证成功后返回的令牌 |

返回参数说明

| 名称 | 类型 | 说明 |
| :-: | :-: | :- |
| json返回示例 | - | - |

json返回示例

```json
{
  "message": {
    "httpPorts": [80, 990],
    "httpSvrs": {
      "localhost": [80, 990]
    },
    "httpSvrsNum": 2,
    "streamPorts": [5510],
    "streamSvrsNum": 1
  },
  "status": "success"
}
```

### Nginx配置接口

#### 1.查询nginx配置接口

接口地址

`http://<Host>:<Port>/<baseURI>?token=<token>&type=<type>`

注：\<baseURI>为ng管理工具配置中WebServerInfo.servers列表各自元素的baseURI子参数值。

返回格式

`json`

请求方式

`http get`

请求示例

`http://127.0.0.1:12321/ng_conf1?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE2MDc3MDEsImlhdCI6MTU5MTYwMzQzNywidXNlcl9pZCI6MSwicGFzc3dvcmQiOiJuZ2FkbWluIiwidXNlcm5hbWUiOiJuZ2FkbWluIiwiZnVsbF9uYW1lIjoibmdhZG1pbiIsInBlcm1pc3Npb25zIjpbXX0.fDoe4v37XyjmrK4wnfhOUnePwJLdszYXveOfoRXyUj8&type=json`

请求参数

| 名称 | 必填 | 类型 | 说明 |
| :-: | :-: | :-: | :- |
| type | 否 | string | 默认为字符串格式，指定“string”为字符串格式；“json”为json格式 |
| token | 是 | string | 用户访问认证成功后返回的令牌 |

返回参数说明

| 名称 | 类型 | 说明 |
| :-: | :-: | :- |
| json返回示例 | - | - |

json返回示例

> string格式数据返回示例
> ```json
> {
>   "message": "# user  nobody;\nworker_processes 1;\n# error_log  logs/error.log;\n",
>   "status": "success"
> }
> ```

> json格式数据返回示例
> ```json
> {
>   "message": "{\"config\":{\"value\":\"/usr/local/openresty/nginx/conf/nginx.conf\",\"param\":[{\"comments\":\"user  nobody;\",\"inline\":true},{\"name\":\"worker_processes\",\"value\": \"0\"},{\"comments\":\"error_log  logs/error.log;\",\"inline\":false}]}}",
>   "status": "success"
> }
> ```

#### 2.更新nginx配置接口

接口地址

`http://<Host>:<Port>/<baseURI>?token=<token>`

注：\<baseURI>为ng管理工具配置中WebServerInfo.servers列表各自元素的baseURI子参数值。

返回格式

`json`

请求方式

`http post`

请求示例

```shell script
curl 'http://127.0.0.1:12321/ng_conf1?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE2MDc3MDEsImlhdCI6MTU5MTYwMzQzNywidXNlcl9pZCI6MSwicGFzc3dvcmQiOiJuZ2FkbWluIiwidXNlcm5hbWUiOiJuZ2FkbWluIiwiZnVsbF9uYW1lIjoibmdhZG1pbiIsInBlcm1pc3Npb25zIjpbXX0.fDoe4v37XyjmrK4wnfhOUnePwJLdszYXveOfoRXyUj8'\
    -X POST\
    -F "data=@/tmp/ng_conf.json"
```

请求参数

| 名称 | 必填 | 类型 | 说明 |
| :-: | :-: | :-: | :- |
| token | 是 | string | 用户访问认证成功后返回的令牌 |

请求表单参数

| 名称 | 必填 | 类型 | 说明 |
| :-: | :-: | :-: | :- |
| data | 是 | file | 配置文件数据json文件 |

返回参数说明

| 名称 | 类型 | 说明 |
| :-: | :-: | :- |
| json返回示例 | - | - |

json返回示例

```json
{
  "message":"Nginx conf updated.",
  "status":"success"
}
```

### 监控接口

#### 1.bifrost健康状态

接口地址

`http://<Host>:<Port>/health?token=<token>`

返回格式

`json`

请求方式

`http get`

请求示例

`http://127.0.0.1:12321/health?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE2MDcwMzcsImlhdCI6MTU5MTYwMzQzNywidXNlcl9pZCI6MSwicGFzc3dvcmQiOiJuZ2FkbWluIiwidXNlcm5hbWUiOiJuZ2FkbWluIiwiZnVsbF9uYW1lIjoibmdhZG1pbiIsInBlcm1pc3Npb25zIjpbXX0.l5qE1sMBD9VJHspzXlhHNmHhbZiF00YlCafUIsIEJpo`

请求参数

| 名称 | 必填 | 类型 | 说明 |
| :-: | :-: | :-: | :- |
| token | 是 | string | 用户访问认证成功后返回的令牌 |

返回参数说明

| 名称 | 类型 | 说明 |
| :-: | :-: | :- |
| json返回示例 | - | - |

json返回示例

```json
{
  "message": "healthy",
  "status": "success"
}
```

#### 2.系统监控信息接口

接口地址

`http://<Host>:<Port>/sysinfo?token=<token>`

返回格式

`json`

请求方式

`http get`

请求示例

`http://127.0.0.1:12321/sysinfo?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTE2MDcwMzcsImlhdCI6MTU5MTYwMzQzNywidXNlcl9pZCI6MSwicGFzc3dvcmQiOiJuZ2FkbWluIiwidXNlcm5hbWUiOiJuZ2FkbWluIiwiZnVsbF9uYW1lIjoibmdhZG1pbiIsInBlcm1pc3Npb25zIjpbXX0.l5qE1sMBD9VJHspzXlhHNmHhbZiF00YlCafUIsIEJpo`

请求参数

| 名称 | 必填 | 类型 | 说明 |
| :-: | :-: | :-: | :- |
| token | 是 | string | 用户访问认证成功后返回的令牌 |

返回参数说明

| 名称 | 类型 | 说明 |
| :-: | :-: | :- |
| json返回示例 | - | - |

json返回示例

```json
{
  "message": {
    "system": "centos 7.4.1708",
    "time": "2020/10/29 17:32:27",
    "cpu": "2.15",
    "mem": "66.16",
    "disk": "71.97",
    "servers_status": ["normal", "abnormal"],
    "servers_version": ["nginx version: openresty/1.13.6.2", "nginx version: openresty/1.13.6.2"],
    "bifrost_version": "v1.0.0-alpha.8"
  },
  "status": "success"
}
```
