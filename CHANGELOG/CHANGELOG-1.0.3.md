<a name="v1.0.3"></a>

## [v1.0.3] - 2022-01-26

### Bug Fixes

- update jwt-go
- **web_server_log_watcher:** repair the `WatcherManager` cannot converge quickly, after the `WebServerLogWatcher.Watch`
  client is disconnected
- **web_server_log_watcher:** fix the bug that `FileWatcher` cannot be closed normally
- **web_server_status:** no fatal to get platform information
- **web_server_status:** no fatal to get platform information

### Code Refactoring

- **bifrost:** the `Bifrost` grpc protocol changes
- **middleware:** separate the 'begin' parameter and adjust it to a non required option

### Features

- add web server statistics feature and add bifrost grpc client
- preliminary completion of grpc server service layer logging middleware
- preliminary completion of grpc server service layer logging middleware
- **bifrost:** add web server statistics feature and add bifrost grpc client
- **client:** fix grpc client support grpc streaming
- **client:** fix grpc client support grpc streaming
- **web_server_config:** reconstruct Bifrost grpc service structure
- **web_server_config:** reformulate the transport layer, endpoint layer and service layer
- **web_server_log_watcher:** add file watcher libs, web server log watcher grpc server and options/configs of this
  service
- **web_server_log_watcher:** complete the `WebServerLogWatcher` service test for grpc server and client
- **web_server_log_watcher:** complete the `WebServerLogWatcher` service test for grpc server and client
- **web_server_log_watcher:** add file watcher libs, web server log watcher grpc server and options/configs of this
  service
- **web_server_status:** add feature of web server status
- **web_server_status:** complete the compilation of monitor library, config, options and store
- **web_server_status:** complete the 'WebServerStatus' service store layer
- **web_server_status:** complete the compilation of monitor library, config, options and store
- **web_server_status:** add feature of web server status
- **web_server_status:** complete the preparation of `WebServerStatus` grpc server and client
- **web_server_status:** complete the preparation of `WebServerStatus` grpc server and client
- **web_server_status:** complete the 'WebServerStatus' service store layer

### Pull Requests

- Merge pull request [#6](https://github.com/ClessLi/bifrost/issues/6) from ClessLi/feature/restructure

### BREAKING CHANGE

the grpc protocol of `Bifrost` has been changed, and service authentication is not supported temporarily.

To migrate the code of the `Bifrost` gRPC client follow the example below:

Before:

data, err := bifrostClient.ViewConfig(...)  // view web server config

data, err := bifrostClient.GetConfig(...)  // get web server config

resp, err := bifrostClient.UpdateConfig(..., reqData, ...)  // update web server config

data, err := bifrostClient.Status(...)  // view web server state

data, err := bifrostClient.ViewStatistics(...)  // view web server statistics

watcher, err := bifrostClient.WatchLog(...)  // watching web server log

After:

// web server config

servernames, err := bifrostClient.WebServerConfig().GetServerNames()  // get web server names

data, err := bifrostClient.WebServerConfig().Get(servername)  // view web server config

err := bifrostClient.WebServerConfig().Update(servername, data)  // update web server config

// web server status

metrics, err := bifrostClient.WebServerStatus().Get()  // view web server status

// web server statistics

statistics, err := bifrostClient.WebServerStatistics().Get(servername)  // view web server statistics

// web server log watcher

outputChan, cancel, err := bifrostClient.WebServerLogWatcher().Watch(request)  // watch web server log


[v1.0.3]: https://github.com/ClessLi/bifrost/compare/v1.0.2-alpha.2...v1.0.3
