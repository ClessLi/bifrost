 
<a name="v1.1.0-alpha.7"></a>
## [v1.1.0-alpha.7] - 2025-02-19
### Code Refactoring
- delete mistakenly added structure

### Features
- **resolv:** make the `Context` more supportive of functional-style operations on streams of elements

### BREAKING CHANGE

removed the `Query` related methods of the `Context` interface and added the `PosSet` interface.

The `Query` related methods of the `Context` interface has been removed as follows:

Methods of deletion to Context:

```go
type Context interface {
    ...
    QueryByKeyWords(kw KeyWords) Pos
    QueryAllByKeyWords(kw KeyWords) []Pos
    ...
}
```

Methods of addition to Context Interface:

```go
type Context interface {
    ...
    ChildrenPosSet() PosSet
    ...
}
```

Added the `PosSet` interface as follows:

Interface of addition:

```go
type PosSet interface {
    Filter(fn func(pos Pos) bool) PosSet
    Map(fn func(pos Pos) (Pos, error)) PosSet
    MapToPosSet(fn func(pos Pos) PosSet) PosSet
    QueryOne(kw KeyWords) Pos
    QueryAll(kw KeyWords) PosSet
    List() []Pos
    Targets() []Context
    Append(pos ...Pos) PosSet
    AppendWithPosSet(posSet PosSet) PosSet
    Error() error
}
```

Methods of addition to Pos Interface:

```go
type Pos interface {
    ...
    QueryOne(kw KeyWords) Pos
    QueryAll(kw KeyWords) PosSet
}
```

To migrate the code for the `Context` operations, follow the example below:

Before:

```go
package main

import (
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
)

func main() {
    conf, err := configuration.NewNginxConfigFromJsonBytes(jsondata)
    if err != nil {
        panic(err)
    }
    for _, pos := range conf.Main().QueryByKeyWords(context.NewKeyWords(context_type.TypeHttp).
	        SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).Target().
	    QueryByKeyWords(context.NewKeyWords(context_type.TypeServer).
            SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).Target().
        QueryAllByKeyWords(context.NewKeyWords(context_type.TypeDirective).
            SetCascaded(false).
            SetStringMatchingValue("server_name test1.com").
            SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)) { // query `server` context, its server name is "test1.com"
        server, _ := pos.Position()
        if server.QueryByKeyWords(context.NewKeyWords(context_type.TypeDirective).
                SetCascaded(false).
                SetRegexpMatchingValue("^listen 80$").
            SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).Target().Error() != nil { // query `server` context, its listen port is 80
            continue
        }
        // query the "proxy_pass" `directive` context, which is in `if` context(value: "($http_api_name != '')") and `location` context(value: "/test1-location")
        ctx, idx := server.QueryByKeyWords(context.NewKeyWords(context_type.TypeLocation).
                SetRegexpMatchingValue(`^/test1-location$`).
                SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).Target().
            QueryByKeyWords(context.NewKeyWords(context_type.TypeIf).
                SetRegexpMatchingValue(`^\(\$http_api_name != ''\)$`).
                SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).Target().
            QueryByKeyWords(context.NewKeyWords(context_type.TypeDirective).
                SetStringMatchingValue("proxy_pass").
                SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).Position()
        // insert an inline comment after the "proxy_pass" `directive` context
        err = ctx.Insert(local.NewContext(context_type.TypeInlineComment, fmt.Sprintf("[%s]test comments", time.Now().String())), idx+1).Error()
        if err != nil {
			panic(err)
        }
    }
}
```

After:

```go
package main

import (
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
)

func main() {
    conf, err := configuration.NewNginxConfigFromJsonBytes(jsondata)
    if err != nil {
        panic(err)
	}
    ctx, idx := conf.Main().ChildrenPosSet().
        QueryOne(nginx_ctx.NewKeyWords(context_type.TypeHttp).
            SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
        QueryAll(nginx_ctx.NewKeyWords(context_type.TypeServer).
            SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
        Filter( // filter out `server` context positions, theirs server name is "test1.com"
            func(pos nginx_ctx.Pos) bool {
            return pos.QueryOne(context.NewKeyWords(context_type.TypeDirective).
                    SetCascaded(false).
                    SetStringMatchingValue("server_name test1.com").
                    SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
                Target().Error() == nil
            },
        ).
        Filter( // filter out `server` context positions, theirs listen port is 80
            func(pos nginx_ctx.Pos) bool {
                return pos.QueryOne(context.NewKeyWords(context_type.TypeDirective).
                        SetCascaded(false).
                        SetRegexpMatchingValue("^listen 80$").
                        SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
                    Target().Error() == nil
            },
        ).
        // query the "proxy_pass" `directive` context position, which is in `if` context(value: "($http_api_name != '')") and `location` context(value: "/test1-location")
        QueryOne(nginx_ctx.NewKeyWords(context_type.TypeLocation).
            SetRegexpMatchingValue(`^/test1-location$`).
            SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
        QueryOne(nginx_ctx.NewKeyWords(context_type.TypeIf).
            SetRegexpMatchingValue(`^\(\$http_api_name != ''\)$`).
            SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
        QueryOne(nginx_ctx.NewKeyWords(context_type.TypeDirective).
            SetStringMatchingValue("proxy_pass").
            SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
        Position()
    // insert an inline comment after the "proxy_pass" `directive` context
    err = ctx.Insert(local.NewContext(context_type.TypeInlineComment, fmt.Sprintf("[%s]test comments", time.Now().String())), idx+1).Error()
    if err != nil {
        panic(err)
    }
}
```

[v1.1.0-alpha.7]: https://github.com/ClessLi/bifrost/compare/v1.1.0-alpha.6...v1.1.0-alpha.7
