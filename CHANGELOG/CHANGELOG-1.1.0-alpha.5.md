 
<a name="v1.1.0-alpha.5"></a>
## [v1.1.0-alpha.5] - 2025-01-17
### Features
- **web_server_bin_cmd:** add `Web Service Binary Command Execution` gRPC Server

### BREAKING CHANGE

the `bifrost` gRPC server APIs have added `WebServerBinCMD`.`Exec` server

The `protobuf` of the `bifrost` gRPC server APIs has been added as follows:

Protocols of addition:

```protobuf
service WebServerBinCMD {
  rpc Exec(ExecuteRequest) returns (ExecuteResponse) {}
}

message ExecuteRequest {
  string ServerName = 1;
  repeated string Args = 2;
}

message ExecuteResponse {
  bool Successful = 1;
  bytes Stdout = 2;
  bytes Stderr = 3;
}
```

The `bifrost` gRPC server APIs client SDK has been added as follows:

Methods of addition to Client Service Factory:

```go
type Factory interface {
    ...
    WebServerBinCMD() WebServerBinCMDService
}
```

Interface of addition to Client Service:

```go
type WebServerBinCMDService interface {
    Exec(servername string, arg ...string) (isSuccessful bool, stdout, stderr string, err error)
}
```

[v1.1.0-alpha.5]: https://github.com/ClessLi/bifrost/compare/v1.1.0-alpha.4...v1.1.0-alpha.5
