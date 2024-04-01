 
<a name="v1.0.9"></a>
## [v1.0.9] - 2024-04-01
### Bug Fixes
- fix 1 `Dependabot` alert.
- **resolv:** fix string containing matching logic in the `KeyWords`.`Match()` method.

### Features
- **resolv:** interface the `Main` Context.
- **resolv:** preliminary completion of the development and testing of the `resolve` V3 version, as well as the update of the `bifrost` service to the `resolve` V3 version.
- **resolv:** complete the writing of the `KeyWords` class and conduct preliminary unit testing of related methods for this class
- **resolv:** preliminary completion of V3 `local` library unit testing and repair.
- **resolv:** complete functional unit testing of V3 `local`.`Include` Context.
- **resolv:** add V3 resolv lib.
- **resolv:** add V3 resolv lib.

### BREAKING CHANGE

replacing the nginx configuration resolving library from version `V2` to `V3`.

To migrate the code follow the example below:

Before:

```go
import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
)

nginxConfFromPath, err := configuration.NewConfigurationFromPath(configAbsPath)
nginxConfFromJsonBytes, err := configuration.NewConfigurationFromJsonBytes(configJsonBytes)
```

After:

```go
import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
)

nginxConfFromPath, err := configuration.NewNginxConfigFromFS(configAbsPath)
nginxConfFromJsonBytes, err := configuration.NewNginxConfigFromJsonBytes(configJsonBytes)
```

Example for querying and inserting context:

```go
import (
	"fmt"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
)

conf, err := configuration.NewNginxConfigFromJsonBytes(jsondata)
if err != nil {
	t.Fatal(err)
}
for _, pos := range conf.Main().QueryByKeyWords(context.NewKeyWords(context_type.TypeHttp)).Target().
	QueryByKeyWords(context.NewKeyWords(context_type.TypeServer)).Target().
	QueryAllByKeyWords(context.NewKeyWords(context_type.TypeDirective).SetStringMatchingValue("server_name test1.com")) {  // query `server` context, its server name is "test1.com"
	server, _ := pos.Position()
	if server.QueryByKeyWords(context.NewKeyWords(context_type.TypeDirective).SetRegexpMatchingValue("^listen 80$")).Target().Error() != nil {  // query `server` context, its listen port is 80
		continue
	}
	// query the "proxy_pass" `directive` context, which is in `if` context(value: "($http_api_name != '')") and `location` context(value: "/test1-location")
	ctx, idx := server.QueryByKeyWords(context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`^/test1-location$`)).Target().
		QueryByKeyWords(context.NewKeyWords(context_type.TypeIf).SetRegexpMatchingValue(`^\(\$http_api_name != ''\)$`)).Target().
		QueryByKeyWords(context.NewKeyWords(context_type.TypeDirective).SetStringMatchingValue("proxy_pass")).Position()
	// insert an inline comment after the "proxy_pass" `directive` context
	err = ctx.Insert(local.NewComment(fmt.Sprintf("[%s]test comments", time.Now().String()), true), idx+1).Error()
	if err != nil {
		return err
	}
}
```

Examples for building nginx context object:

```go
import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
)

// new main context
newMainContext, err := local.NewMain("/usr/local/nginx/conf/nginx.conf")
// new directive context
newDirective := local.NewDirective("some_directive", "some params")
// new comment context
newComment := local.NewComment("some comments", false)
// new other context
newConfig := local.NewContext(context_type.TypeConfig, "conf.d/location.conf")
newInclude := local.NewContext(context_type.TypeInclude, "conf.d/*.conf")
newHttp := local.NewContext(context_type.TypeHttp, "")
...
```

[v1.0.9]: https://github.com/ClessLi/bifrost/compare/v1.0.8...v1.0.9
