# [Bifrost](https://github.com/ClessLi/bifrost)

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/ClessLi/bifrost?label=bifrost)](https://github.com/ClessLi/bifrost/releases/latest)
![GitHub Releases](https://img.shields.io/github/downloads/ClessLi/bifrost/latest/total)
[![GitHub](https://img.shields.io/github/license/ClessLi/bifrost)](LICENSE)

# 项目介绍

**Bifrost**
是基于golang语言开发的项目，它目前还处于测试阶段，用于对Nginx配置文件解析并提供配置文件展示和修改的接口，支持json、字符串格式与golang结构相互转换。该项目持续更新中。最新可用版本为[v1.1.0-alpha.11](https://github.com/ClessLi/bifrost/tree/v1.1.0-alpha.11)

# 项目特点

支持将配置文件、json数据、字符串与配置结构体相互转换

配置结构体支持增加、删除、查询

实现了在加载配置或反序列化json时，防止循环读取配置的功能；实现了nginx配置文件后台更新后，自动热加载的功能

提供配置文件展示和修改及配置信息统计查询，及主机系统状况信息查询的gRPC接口

# 合作项目

## [Heimedallr-Reborn](https://github.com/tanganyu1114/heimdallr-reborn)

nginx后管平台

基于 gin-vue-admin 框架制作

目前仅支持配置文件查看和日志信息读取

配置nginx信息功能还在开发中

# 使用方法

## 编译

如果你需要重新编译Bifrost项目，可以执行以下 2 步：

1. 克隆源码

```bash
$ git clone https://github.com/ClessLi/bifrost $GOPATH/src/github.com/clessli/bifrost
```

2. 编译

```bash
$ cd $GOPATH/src/github.com/clessli/bifrost
$ go build cmd/bifrost
```

## 应用配置

配置路径

`bifrost: configs/bifrost.yml`

配置示例

`bifrost`
```yaml
server:
  healthz: true  # 是否开启健康检查，如果开启会安装 healthz gRPC服务，默认 true

# 服务配置
# secure:  # gRPC 安全模式配置，目前暂默认不启用
insecure:
  bind-address: 0.0.0.0  # 绑定的不安全 IP 地址，设置为 0.0.0.0 表示使用全部网络接口，默认为 127.0.0.1，建议设置为提供服务网卡ip或域名，在注册服务到注册中心时会用到，避免服务发现异常

# gRPC服务参数配置
grpc:
  chunksize: 2048  # 传输带宽配置，单位（Byte），范围（100~65535）
  receiv-timeout-minutes: 3

# Web Server Config 相关配置
web-server-configs:
  dns-ipv4: "114.114.114.114"  # dns的ip地址，用于解析被反向代理URL的域名地址
  items:
    - server-name: "bifrost-test"  # WebServer 名称
      server-type: "nginx"  # WebServer 类型，目前暂仅支持 nginx
      config-path: "/usr/local/nginx/conf/nginx.conf"  # WebServer 配置文件路径
      verify-exec-path: "/usr/local/nginx/sbin/nginx"  # WebServer 配置文件校验用可执行文件路径，目前仅支持 nginx 的应用运行二进制文件路径
      logs-dir-path: "/usr/local/nginx/logs"  # WebServer 日志存放路径
      backup-dir: ""  # WebServer 配置文件自动备份路径，为空时将使用`config-path`文件的目录路径作为备份目录路径
      backup-cycle: 1  # WebServer 配置文件自动备份周期时长，单位（天），为0时不启用自动备份
      backups-retention-duration: 7  # WebServer 配置文件自动备份归档保存时长，单位（天），为0时不启用自动备份

# 注册中心配置
# RA:  # 注册中心地址配置
#   Host: "192.168.0.11"
#   Port: 8500

# 日志配置
log:
# 启用开发模式
# development: false

# 禁用日志输出
# disable-caller: false

# 禁用日志追踪
# disable-stracktrace: false

# 启用代色彩的日志记录
# enable-color: false

# 错误日志最低级别, 默认为“warn”
# error-level: warn

# 错误日志输出路径，默认为“logs/biforst-error.log”
# error-output-paths:
# - logs/bifrost-error.log

# 日志输出格式，支持“console”，“json”，默认为“console”
# format: console

# Info日志最低级别，默认为“info”
# info-level: info

# Info日志输出路径，默认为“logs/bifrost.log”
# info-output-paths:
# - logs/bifrost.log
```

## 命令帮助

`bifrost`
```bash
$ ./bifrost -h
The Bifrost is used to parse the nginx configuration file 
and provide an interface for displaying and modifying the configuration file.
It supports the mutual conversion of JSON, string format and golang structure.
The Bifrost services to do the api objects management with gRPC protocol.

Find more Bifrost information at:
    https://github.com/ClessLi/bifrost/blob/master/docs/guide/en-US/cmd/bifrost.md

Usage:
  bifrost [flags]

Generic flags:

      --server.healthz
                Add self readiness check and install health check service. (default true)
      --server.middlewares strings
                List of allowed middlewares for server, comma separated. If this is empty default middlewares will be used.

Secure serving flags:

      --secure.bind-address string
                The IP address on which to listen for the --secure.bind-port port. The associated interface(s) must be reachable by the rest of the engine, and by CLI/web clients. If blank, all interfaces will be used (0.0.0.0 for all
                IPv4s interfaces and :: for all IPv6 interfaces). (default "0.0.0.0")
      --secure.bind-port int
                The port on which to serve 12421 with authentication and authorization. Set to zero to disable.
      --secure.tls.cert-dir string
                The directory where the TLS certs are located. If --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file are provided, this flag will be ignored. (default "/var/run/bifrost")
      --secure.tls.cert-key.cert-file string
                File containing the default x509 Certificate for TLS. (CA cert, if any, concatenated after server cert).
      --secure.tls.cert-key.private-key-file string
                File containing the default x509 private key matching --secure.tls.cert-key.cert-file.
      --secure.tls.pair-name string
                The name which will be used with --secure.tls.cert-dir to make a cert and key filenames. It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key (default "bifrost")

Insecure serving flags:

      --insecure.bind-address string
                The IP address on which to serve the --insecure.bind-port (set to 0.0.0.0 for all IPv4s interfaces and :: for all IPv6 interfaces). (default "0.0.0.0")
      --insecure.bind-port int
                The port on which to serve unsecured, unauthenticated access. It is assumed that firewall rules are set up such that this port is not reachable from outside of the deployed machine and that port 12321 on the bifrost public
                address is proxied to this port. Set to zero to disable. (default 12321)

RA options flags:

      --ra.host string
                Specifies the bind address of the Registration Authority server. Set empty to disable.
      --ra.port int
                Specifies the bind port of the Registration Authority server.

GRPC serving flags:

      --grpc.chunksize int
                Set the max message size in bytes the server can send. Can not less than 100 bytes. (default 1024)
      --grpc.receive-timeout int
                Set the timeout for receiving data. The unit is per minute. (default 1)

Web server configs flags:

      --web-server-configs.dns-ipv4 string
                Set Domain Name resolver for resolving the proxy service domain name

Monitor flags:

      --monitor.cycle-time duration
                 (default 2m0s)
      --monitor.frequency-per-cycle int
                 (default 10)
      --monitor.sync-interval duration
                 (default 1m0s)

Log watcher flags:

      --web-server-log-watcher.max-connections int
                 (default 1000)
      --web-server-log-watcher.watch-timeout duration
                 (default 5m0s)

Log flags:

      --log.development
                Development puts the logger in development mode, which changes the behavior of DPanicLevel and takes stacktraces more liberally.
      --log.disable-caller
                Disable output of caller information in the log.
      --log.disable-stacktrace
                Disable the log to record a stack trace for all messages at or above panic level.
      --log.enable-color
                Enable output ansi colors in plain format logs.
      --log.error-level LEVEL
                Minimum Error log output LEVEL. (default "warn")
      --log.error-output-paths strings
                Output paths of Error log. (default [logs\bifrost_error.log])
      --log.format FORMAT
                Log output FORMAT, support plain or json format. (default "console")
      --log.info-level LEVEL
                Minimum Info log output LEVEL. (default "info")
      --log.info-output-paths strings
                Output paths of Info log. (default [logs\bifrost.log])
      --log.inner-error-output-paths strings
                Inner Error output paths of log. (default [stderr])
      --log.name string
                The name of the logger.

Global flags:

  -c, --config FILE
                Read configuration from specified FILE, support JSON, TOML, YAML, HCL, or Java properties formats.
  -h, --help
                help for bifrost
      --version version[=true]
                Print version information and quit.
```

## 配置解析库

### Nginx配置管理器

Nginx配置管理器提供配置读取、更新、保存、备份及重载，方法详见其接口文档（[NginxConfigManager](pkg/resolv/V3/nginx/configuration/nginx_config_manager.go)）

实例化方法如下：

```go
package main

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
)

func main() {
	nginxConfFromPath, err := configuration.NewNginxConfigFromFS(configAbsPath)
	nginxConfFromJsonBytes, err := configuration.NewNginxConfigFromJsonBytes(configJsonBytes)
	...
}
```

Nginx配置上下文对象检索与插入示例如下：

```go
package main

import (
	"fmt"
	"time"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginxCtx "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
)

func main() {
	conf, err := configuration.NewNginxConfigFromJsonBytes(jsondata)
	if err != nil {
		panic(err)
	}
	ctx, idx := conf.Main().ChildrenPosSet().
		QueryOne(nginxCtx.NewKeyWordsByType(context_type.TypeHttp).
			SetSkipQueryFilter(nginxCtx.SkipDisabledCtxFilterFunc)).
		QueryAll(nginxCtx.NewKeyWordsByType(context_type.TypeServer).
			SetSkipQueryFilter(nginxCtx.SkipDisabledCtxFilterFunc)).
		Filter( // filter out `server` context positions, theirs server name is "test1.com"
			func(pos nginxCtx.Pos) bool {
				return pos.QueryOne(nginxCtx.NewKeyWordsByType(context_type.TypeDirective).
					SetCascaded(false).
					SetStringMatchingValue("server_name test1.com").
					SetSkipQueryFilter(nginxCtx.SkipDisabledCtxFilterFunc)).
					Target().Error() == nil
			},
		).
		Filter( // filter out `server` context positions, theirs listen port is 80
			func(pos nginxCtx.Pos) bool {
				return pos.QueryOne(nginxCtx.NewKeyWordsByType(context_type.TypeDirective).
					SetCascaded(false).
					SetRegexpMatchingValue("^listen 80$").
					SetSkipQueryFilter(nginxCtx.SkipDisabledCtxFilterFunc)).
					Target().Error() == nil
			},
		).
		// query the "proxy_pass" `directive` context position, which is in `if` context(value: "($http_api_name != '')") and `location` context(value: "/test1-location")
		QueryOne(nginxCtx.NewKeyWordsByType(context_type.TypeLocation).
			SetRegexpMatchingValue(`^/test1-location$`).
			SetSkipQueryFilter(nginxCtx.SkipDisabledCtxFilterFunc)).
		QueryOne(nginxCtx.NewKeyWordsByType(context_type.TypeIf).
			SetRegexpMatchingValue(`^\(\$http_api_name != ''\)$`).
			SetSkipQueryFilter(nginxCtx.SkipDisabledCtxFilterFunc)).
		QueryOne(nginxCtx.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass).
			SetSkipQueryFilter(nginxCtx.SkipDisabledCtxFilterFunc)).
		Position()
	// insert an inline comment after the "proxy_pass" `directive` context
	err = ctx.Insert(local.NewContext(context_type.TypeInlineComment, fmt.Sprintf("[%s]test comments", time.Now().String())), idx+1).Error()
	if err != nil {
		panic(err)
	}
}
```

Nginx配置上下文对象新建示例如下：

```go
package main

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
)

func main() {
	// new main context
	newMainContext, err := local.NewMain("/usr/local/nginx/conf/nginx.conf")
	// new directive context
	newDirective := local.NewContext(context_type.TypeDirective, "some_directive some params")
	// new comment context
	newComment := local.NewContext(context_type.TypeComment, "some comments")
	newInlineComment := local.NewContext(context_type.TypeInlineComment, "some inline comments")
	// new http type proxy_pass directive
	newHttpProxyPass := local.NewContext(context_type.TypeDirHTTPProxyPass, "https://example.com/test/")
	// new stream type proxy_pass directive
	newStreamProxyPass := local.NewContext(context_type.TypeDirStreamProxyPass, "example.com:22")
	// new other context
	newConfig := local.NewContext(context_type.TypeConfig, "conf.d/location.conf")
	newInclude := local.NewContext(context_type.TypeInclude, "conf.d/*.conf")
	newHttp := local.NewContext(context_type.TypeHttp, "")
	...
}
```

## 接口文档

支持web服务器（暂仅支持nginx）配置文件查看、序列化导出（json）、配置更新、配置统计信息查看、web服务器状态信息查看，及web服务器（暂仅支持nginx）日志监看功能

详见

[bifrost_gRPC接口定义](api/protobuf-spec/bifrostpb/v1/bifrost.proto)

## 客户端

结合go-kit框架编写的客户端对象

### bifrost客户端

通过"pkg/client/bifrost/client.NewClient"函数可生成bifrost服务客户端

详见[bifrost客户端](pkg/client/bifrost/v1/client.go)

<h3 id="test">客户端使用示例</h3>

详见[客户端测试示例](test/grpc_client)
