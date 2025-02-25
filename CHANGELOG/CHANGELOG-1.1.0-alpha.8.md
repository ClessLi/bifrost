 
<a name="v1.1.0-alpha.8"></a>
## [v1.1.0-alpha.8] - 2025-02-25
### Features
- **bifrost:** add updating check feature in the `bifrostpb`.`WebServerConfig` gRPC service

### BREAKING CHANGE

adjusted the gRPC protocol and added the `OriginalFingerprints` data field to the `ServerConfig` message structure; Optimized the `WebServerConfigService` service interface methods for the gRPC service client.

The `ServerConfig` message structure field has been added as follows:

Addition of message structure fields:

```protobuf
message ServerConfig {
  ...
  bytes OriginalFingerprints = 3;
}
```

The adjustment of the `WebServerConfigService` service interface methods for the gRPC service client are as follows:

Before:

```go
type WebServerConfigService interface {
    GetServerNames() (servernames []string, err error)
    Get(servername string) ([]byte, error)
    Update(servername string, config []byte) error
}
```
After:

```go
import (
    ...
    "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
    utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
    ...
)

type WebServerConfigService interface {
    GetServerNames() (servernames []string, err error)
    Get(servername string) (config configuration.NginxConfig, originalFingerprinter utilsV3.ConfigFingerprinter, err error)
    Update(servername string, config configuration.NginxConfig, originalFingerprints utilsV3.ConfigFingerprints) error
}
```

[v1.1.0-alpha.8]: https://github.com/ClessLi/bifrost/compare/v1.1.0-alpha.7...v1.1.0-alpha.8
