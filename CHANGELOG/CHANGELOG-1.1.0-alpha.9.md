 
<a name="v1.1.0-alpha.9"></a>
## [v1.1.0-alpha.9] - 2025-10-22
### Bug Fixes
- fix known vulnerabilities, CVE IDs as follows: `CVE-2025-30204`, `CVE-2025-22872`, `CVE-2025-22870`

### Code Refactoring
- **bifrost:** adjust the concurrency architecture of the `Bifrost` server
- **bifrost:** combine the `sync`.`errgroup` package and adjust the concurrency architecture of the `monitor` package
- **bifrost:** adjust the concurrency architecture of the `monitor` package
- **bifrost:** adjust the matching logic of the `KeyWords` interface object

### Features
- **bifrost:** the `NginxConfigManager` server has added a new feature that can automatically and periodically parse the domain names of servers being reverse-proxied
- **bifrost:** added proxy information parsing and context validation functions for the `HTTPProxyPass` and `StreamProxyPass` contexts
- **bifrost:** added `HTTPProxyPass` and `StreamProxyPass` contexts, which are sub-objects of the `Directive` context; completed the proxy information parsing mechanism for these two newly added `Directive` contexts.
- **web_server_config:** add `Network Connectivity Check Of The Proxied Server` gRPC Server

### BREAKING CHANGE

the `bifrost` gRPC server APIs have added `WebServerConfig`.`ConnectivityCheckOfProxiedServers` server

The `protobuf` of the `bifrost` gRPC server APIs has been added as follows:

Protocols of addition:

```protobuf
...
service WebServerConfig {
  ...
  rpc ConnectivityCheckOfProxiedServers(ServerConfigContextPos) returns (ContextData) {}
}

...

message ServerConfigContextPos {
  string ServerName = 1;
  ContextPos ContextPos = 2;
  bytes OriginalFingerprints = 3;
}

message ContextPos {
  string ConfigPath = 1;
  repeated int32 Pos = 2;
}

message ContextData {
  bytes JsonData = 1;
}

...
```

The `bifrost` gRPC server APIs client SDK has been added as follows:

Methods of addition to `WebServerConfigService` Client Service:

```go
type WebServerConfigService interface {
    ...
    ConnectivityCheckOfProxiedServers(servname string, proxyPass local.ProxyPass, originalFingerprints utilsV3.ConfigFingerprints) (resp local.ProxyPass, err error)
    ...
}
```

change the method name `SetMatchingFilter` in the `KeyWords` interface to `AppendMatchingFilter`; change the name of the `NewKeyWords` function, which builds the `KeyWords` interface object, to `NewKeyWordsByType`, and add a new `NewKeyWords` function to build a `KeyWords` interface object for custom matching filter function.

The method of the `KeyWords` interface has been changed as follows:

Before:

```go
type KeyWords interface {
    ...
    SetMatchingFilter(filterFunc func(targetCtx Context) bool) KeyWords
}
```

After:

```go
type KeyWords interface {
    ...
    AppendMatchingFilter(filterFunc func(targetCtx Context) bool) KeyWords
}
```

The modification of the `NewKeyWords` function is as follows:

Before:

```go
func NewKeyWords(ctxtype context_type.ContextType) KeyWords {
    ...
}
```

After:

```go
func NewKeyWordsByType(ctxtype context_type.ContextType) KeyWords {
    ...
}
```

Added a constructor function `NewKeyWords` for custom injection matching filter function, as follows:

```go
func NewKeyWords(matchingFilterFunc func(targetCtx Context) bool) KeyWords {
    ...
}
```

the `SetMatchingFilter` method has been added to the `context`.`KeyWords` interface, which is used for the custom matching mechanism of target context keyword retrieval.

The method of the `KeyWords` interface has been added as follows:

Addition of methods:

```go
type KeyWords interface {
    ...
    SetMatchingFilter(filterFunc func(targetCtx Context) bool) KeyWords
```

[v1.1.0-alpha.9]: https://github.com/ClessLi/bifrost/compare/v1.1.0-alpha.8...v1.1.0-alpha.9
