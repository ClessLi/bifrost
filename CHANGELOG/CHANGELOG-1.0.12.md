 
<a name="v1.0.12"></a>
## [v1.0.12] - 2024-04-24
### Features
- **resolv:** refactored the `Context` `JSON` format serialization structure to simplify serialization/deserialization operations, making the serialization structure more universal and unified.

### BREAKING CHANGE

Changed the serialization structure of the `Context` `JSON` format, as well as the new functions for the `Directive` and `Comment` context objects.

The `Context` `JSON` format serialization structure changes as follows:

Before:

```json
{
    "main-config": "C:\\test\\test.conf",
    "configs":
    {
        "C:\\test\\test.conf":
        [
            {
                "http":
                {
                    "params":
                    [
                        {
                            "inline": true, "comments": "test comment"
                        },
                        {
                            "server":
                            {
                                "params":
                                [
                                    {
                                        "directive": "server_name", "params": "testserver"
                                    },
                                    {
                                        "location": {"value": "~ /test"}
                                    },
                                    {
                                        "include":
                                        {
                                            "value": "conf.d\\include*conf",
                                            "params": ["conf.d\\include.location1.conf", "conf.d\\include.location2.conf"]
                                        }
                                    }
                                ]
                            }
                        }
                    ]
                }
            }
        ],
        "conf.d\\include.location1.conf":
        [
            {
                "location": {"value": "~ /test1"}
            }
        ],
        "conf.d\\include.location2.conf":
        [
            {
                "location": {"value": "^~ /test2"}
            }
        ]
    }
}
```

After:

```json
{
    "main-config": "C:\\test\\test.conf",
    "configs":
    {
        "C:\\test\\test.conf":
        {
            "context-type": "config",
            "value": "C:\\test\\test.conf",
            "params":
            [
                {
                    "context-type": "http",
                    "params":
                    [
                        {
                            "context-type": "inline_comment", "value": "test comment"
                        },
                        {
                            "context-type": "server",
                            "params":
                            [
                                {
                                    "context-type": "directive", "value": "server_name testserver"
                                },
                                {
                                    "context-type": "location", "value": "~ /test"
                                },
                                {
                                    "context-type": "include",
                                    "value": "conf.d\\include*conf",
                                    "params": ["conf.d\\include.location1.conf", "conf.d\\include.location2.conf"]
                                }
                            ]
                        }
                    ]
                }
            ]
        },
        "conf.d\\include.location1.conf":
        {
            "context-type": "config",
            "value": "conf.d\\include.location1.conf",
            "params":
            [
                {
                    "context-type": "location", "value": "~ /test1"
                }
            ]
        },
        "conf.d\\include.location2.conf":
        {
            "context-type": "config",
            "value": "conf.d\\include.location2.conf",
            "params":
            [
                {
                    "context-type": "location", "value": "^~ /test2"
                }
            ]
        }
    }
}
```

The code migration example for creating function calls for `Directive` and `Comment` context objects is as follows:

Before:

```go
import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
)

// new directive context
newDirective := local.NewDirective("some_directive", "some params")
// new comment context
newComment := local.NewComment("some comments", false)
newComment := local.NewComment("some inline comments", true)
```

After:

```go
import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
)

// new directive context
newDirective := local.NewContext(context_type.TypeDirective, "some_directive some params")
// new comment context
newComment := local.NewContext(context_type.TypeComment, "some comments")
newInlineComment := local.NewContext(context_type.TypeInlineComment, "some inline comments")
```

[v1.0.12]: https://github.com/ClessLi/bifrost/compare/v1.0.11...v1.0.12
