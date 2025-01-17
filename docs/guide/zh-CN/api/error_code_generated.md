# 错误码
！！Bifrost 系统(简称Bifrost)错误码列表，由 `codegen -type=int -fullname "Bifrost" -simplename "Bifrost" -doc` 命令生成，不要对此文件做任何更改。
## 功能说明
如果返回结果中存在 `code` 字段，则表示调用 API 接口失败。例如：
```json
{
  "code": 100101,
  "message": "Database error"
}
```
上述返回中 `code` 表示错误码，`message` 表示该错误的具体信息。每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。
## 错误码列表
Bifrost 系统支持的错误码列表如下：

| Identifier                     | Code   | HTTP Code | Description                                                      |
|--------------------------------|--------|-----------|------------------------------------------------------------------|
| ErrSuccess                     | 100001 | 200       | OK                                                               |
| ErrUnknown                     | 100002 | 500       | Internal server error                                            |
| ErrBind                        | 100003 | 400       | Error occurred while binding the request body to the struct      |
| ErrValidation                  | 100004 | 400       | Validation failed                                                |
| ErrTokenInvalid                | 100005 | 401       | Token invalid                                                    |
| ErrPageNotFound                | 100006 | 404       | Page not found                                                   |
| ErrRequestTimeout              | 100007 | 408       | Request timeout                                                  |
| ErrDatabase                    | 100101 | 500       | Database error                                                   |
| ErrDataRepository              | 100201 | 500       | Data Repository error                                            |
| ErrEncrypt                     | 100301 | 401       | Error occurred while encrypting the user password                |
| ErrSignatureInvalid            | 100302 | 401       | Signature is invalid                                             |
| ErrExpired                     | 100303 | 401       | Token expired                                                    |
| ErrInvalidAuthHeader           | 100304 | 401       | Invalid authorization header                                     |
| ErrMissingHeader               | 100305 | 401       | The `Authorization` header was empty                             |
| ErrUserOrPasswordIncorrect     | 100306 | 401       | User or Password was incorrect                                   |
| ErrPermissionDenied            | 100307 | 403       | Permission denied                                                |
| ErrAuthnClientInitFailed       | 100308 | 500       | The `Authentication` client initialization failed                |
| ErrAuthClientNotInit           | 100309 | 500       | The `Authentication` and `Authorization` client not initialized  |
| ErrConnToAuthServerFailed      | 100310 | 500       | Failed to connect to `Authentication` and `Authorization` server |
| ErrEncodingFailed              | 100401 | 500       | Encoding failed due to an error with the data                    |
| ErrDecodingFailed              | 100402 | 500       | Decoding failed due to an error with the data                    |
| ErrInvalidJSON                 | 100403 | 500       | Data is not valid JSON                                           |
| ErrEncodingJSON                | 100404 | 500       | JSON data could not be encoded                                   |
| ErrDecodingJSON                | 100405 | 500       | JSON data could not be decoded                                   |
| ErrInvalidYaml                 | 100406 | 500       | Data is not valid Yaml                                           |
| ErrEncodingYaml                | 100407 | 500       | Yaml data could not be encoded                                   |
| ErrDecodingYaml                | 100408 | 500       | Yaml data could not be decoded                                   |
| ErrConfigurationTypeMismatch   | 110001 | 500       | Configuration type mismatch                                      |
| ErrSameConfigFingerprint       | 110002 | 500       | Same config fingerprint                                          |
| ErrSameConfigFingerprints      | 110003 | 500       | Same config fingerprint between files and configuration          |
| ErrConfigManagerIsRunning      | 110004 | 500       | Config manager is running                                        |
| ErrConfigManagerIsNotRunning   | 110005 | 500       | Config manager is not running                                    |
| ErrWebServerNotFound           | 110006 | 400       | Web server not found                                             |
| ErrConfigurationNotFound       | 110007 | 400       | Web server configuration not found                               |
| ErrParserNotFound              | 110008 | 500       | Parser not found                                                 |
| ErrUnknownKeywordString        | 110009 | 500       | Unknown keyword string                                           |
| ErrInvalidConfig               | 110010 | 500       | Invalid parser.Config                                            |
| ErrParseFailed                 | 110011 | 500       | Config parse failed                                              |
| ErrV3ContextIndexOutOfRange    | 110012 | 500       | Index of the Context's children is out of range                  |
| ErrV3NullContextPosition       | 110013 | 500       | Null Context position                                            |
| ErrV3SetFatherContextFailed    | 110014 | 500       | Set father Context failed                                        |
| ErrV3OperationOnErrorContext   | 110015 | 500       | Performing operations on Error Context                           |
| ErrV3InvalidContext            | 110016 | 500       | Invalid Context                                                  |
| ErrV3InvalidOperation          | 110017 | 500       | Invalid operation                                                |
| ErrV3ContextNotFound           | 110018 | 500       | Queried context not found                                        |
| ErrV3ConversionToContextFailed | 110019 | 500       | Conversion to context failed                                     |
| ErrStopMonitoringTimeout       | 110201 | 500       | Stop monitoring timeout                                          |
| ErrMonitoringServiceSuspension | 110202 | 500       | Monitoring service suspension                                    |
| ErrMonitoringStarted           | 110203 | 500       | Monitoring is already started                                    |
| ErrLogsDirPath                 | 110301 | 500       | Logs dir is not exist or is not a directory                      |
| ErrLogBufferIsNotExist         | 110302 | 500       | Log buffer is not exist                                          |
| ErrLogBufferIsExist            | 110303 | 500       | Log buffer is already exist                                      |
| ErrLogIsLocked                 | 110304 | 500       | Log is locked                                                    |
| ErrLogIsUnlocked               | 110305 | 500       | Log is unlocked                                                  |
| ErrUnknownLockError            | 110306 | 500       | Unknown lock error                                               |

