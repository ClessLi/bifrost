package local

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
)

func Test_jsonUnmarshaler_UnmarshalJSON(t *testing.T) {
	testTargetFatherCtx := NewContext(context_type.TypeServer, "")
	testMain, err := NewMain("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(NewContext(context_type.TypeInlineComment, "test comment"), 0).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(NewContext(context_type.TypeDirective, "server_name testserver"), 0).
					Insert(
						NewContext(context_type.TypeLocation, "~ /test"),
						1,
					),
				1,
			).
			Insert(testTargetFatherCtx, 2),
		0,
	)
	testIncludedConfig := NewContext(context_type.TypeConfig, "conf.d\\proxy.conf").(*Config)
	testIncludedConfig.ConfigPath, err = newConfigPath(testMain.graph(), testIncludedConfig.Value())
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(testIncludedConfig)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		configGraph      ConfigGraph
		completedContext context.Context
		fatherContext    context.Context
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal test",
			fields: fields{
				configGraph:      testMain.graph(),
				completedContext: context.NullContext(),
				fatherContext:    testTargetFatherCtx,
			},
			args: args{bytes: []byte(
				`{
    "enabled": true,
    "context-type": "location",
    "value": "~ /target",
    "params":
    [
        {
            "context-type": "inline_comment",
            "value": "target location"
        },
        {
            "enabled": true,
            "context-type": "include",
            "value": "conf.d\\*conf"
        }
    ]
}`,
			)},
		},
		{
			name: "unmarshal config",
			fields: fields{
				configGraph:      testMain.graph(),
				completedContext: context.NullContext(),
				fatherContext:    testTargetFatherCtx,
			},
			args: args{bytes: []byte(
				`{
    "context-type": "config",
    "value": "conf.d\\proxy.conf",
    "params": [
        {
            "enabled": true,
            "context-type": "location",
            "value": "~ /test_proxy",
            "params": [
                {
                    "enabled": true,
                    "context-type": "dir_http_proxy_pass",
                    "value": "https://baidu.com",
                    "proxy-pass": {
                        "original-url": "https://baidu.com",
                        "protocol": "https",
                        "addresses": [
                            {
                                "domain-name": "baidu.com",
                                "port": 443,
                                "ipv4-list": [
                                    "10.1.11.111",
                                    "10.1.12.122"
                                ],
                                "resolve-err": null
                            }
                        ],
                        "uri": ""
                    }
                }
            ]
        },
        {
            "context-type": "comment",
            "value": "test proxy end"
        }
    ]
}`,
			)},
		},
		{
			name: "config with unmatched config path",
			fields: fields{
				configGraph:      testMain.graph(),
				completedContext: context.NullContext(),
				fatherContext:    testTargetFatherCtx,
			},
			args: args{bytes: []byte(
				`{
    "context-type": "config",
    "value": "conf.d\\proxy.conf1",
    "params": [
        {
            "enabled": true,
            "context-type": "location",
            "value": "~ /test_proxy",
            "params": [
                {
                    "enabled": true,
                    "context-type": "dir_http_proxy_pass",
                    "value": "https://baidu.com",
                    "proxy-pass": {
                        "original-url": "https://baidu.com",
                        "protocol": "https",
                        "addresses": [
                            {
                                "domain-name": "baidu.com",
                                "port": 443,
                                "ipv4-list": [
                                    "10.1.11.111",
                                    "10.1.12.122"
                                ],
                                "resolve-err": null
                            }
                        ],
                        "uri": ""
                    }
                }
            ]
        },
        {
            "context-type": "comment",
            "value": "test proxy end"
        }
    ]
}`,
			)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &jsonUnmarshaller{
				configGraph:      tt.fields.configGraph,
				completedContext: tt.fields.completedContext,
				fatherContext:    tt.fields.fatherContext,
			}
			if err := u.UnmarshalJSON(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mainUnmarshaler_UnmarshalJSON(t *testing.T) {
	type fields struct {
		completedMain MainContext
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantCheckConfigLines []string
		wantErr              bool
	}{
		{
			name: "normal test",
			args: args{bytes: []byte(
				`{
    "main-config": "C:\\test\\test.conf",
    "configs":
    {
        "C:\\test\\test.conf":
        {
            "enabled": true,
            "context-type": "config",
            "value": "C:\\test\\test.conf",
            "params":
            [
                {
                    "enabled": true,
                    "context-type": "http",
                    "params":
                    [
                        {
                            "context-type": "inline_comment", "value": "test comment"
                        },
                        {
                            "enabled": true,
                            "context-type": "server",
                            "params":
                            [
                                {
                                    "enabled": true,"context-type": "directive", "value": "server_name testserver"
                                },
                                {
                                    "enabled": true,"context-type": "location", "value": "~ /test"
                                },
                                {
                                    "enabled": true,
                                    "context-type": "include",
                                    "value": "conf.d\\include*conf"
                                }
                            ]
                        }
                    ]
                }
            ]
        },
        "conf.d\\include.location1.conf":
        {
            "enabled": true,
            "context-type": "config",
            "value": "conf.d\\include.location1.conf",
            "params":
            [
                {
                    "enabled": true, "context-type": "location", "value": "~ /test1"
                }
            ]
        },
        "conf.d\\include.location2.conf":
        {
            "enabled": true,
            "context-type": "config",
            "value": "conf.d\\include.location2.conf",
            "params":
            [
                {
                    "enabled": true, "context-type": "location", "value": "^~ /test2"
                }
            ]
        }
    }
}`,
			)},
			wantCheckConfigLines: []string{
				"http {    # test comment",
				"    server {",
				"        server_name testserver;",
				"        location ~ /test {",
				"        }",
				"        # include <== conf.d\\include*conf",
				"        location ~ /test1 {",
				"        }",
				"        location ^~ /test2 {",
				"        }",
				"    }",
				"}",
			},
			wantErr: false,
		},
		{
			name: "enabled/disabled children contexts",
			args: args{bytes: []byte(
				`{
    "main-config": "C:\\test\\nginx.conf",
    "configs": {
        "C:\\test\\nginx.conf": {
            "enabled": true,
            "context-type": "config",
            "value": "C:\\test\\nginx.conf",
            "params": [
                {
                    "enabled": true,
                    "context-type": "http",
                    "params": [
                        {
                            "enabled": true,
                            "context-type": "server",
                            "params": [
                                {
                                    "context-type": "inline_comment",
                                    "value": "enabled server with enabled children configs"
                                },
                                {
                                    "enabled": true,
                                    "context-type": "include",
                                    "value": "conf.d/enabled.conf"
                                }
                            ]
                        },
                        {
                            "enabled": true,
                            "context-type": "server",
                            "params": [
                                {
                                    "context-type": "inline_comment",
                                    "value": "enabled server with disabled children configs"
                                },
                                {
                                    "enabled": true,
                                    "context-type": "include",
                                    "value": "conf.d/disabled.conf"
                                }
                            ]
                        },
                        {
                            "enabled": true,
                            "context-type": "server",
                            "params": [
                                {
                                    "context-type": "inline_comment",
                                    "value": "enabled server with disabled include context"
                                },
                                {
                                    "context-type": "include",
                                    "value": "conf.d/enabled.conf"
                                }
                            ]
                        },
                        {
                            "context-type": "server",
                            "params":[
                                {
                                    "context-type": "inline_comment",
                                    "value": "disabled server with enabled children configs"
                                },
                                {
                                    "enabled": true,
                                    "context-type": "include",
                                    "value": "conf.d/enabled.conf"
                                }
                            ]
                        },
                        {
                            "context-type": "server",
                            "params": [
                                {
                                    "context-type": "inline_comment",
                                    "value": "disabled server with disabled children configs"
                                },
                                {
                                    "enabled": true,
                                    "context-type": "include",
                                    "value": "conf.d/disabled.conf"
                                }
                            ]
                        },
                        {
                            "context-type": "server",
                            "params": [
                                {
                                    "context-type": "inline_comment",
                                    "value": "disabled server with disabled include context"
                                },
                                {
                                    "context-type": "include",
                                    "value": "conf.d/enabled.conf"
                                }
                            ]
                        }
                    ]
                }
            ]
        },
        "conf.d/disabled.conf": {
            "enabled": false,
            "context-type": "config",
            "value": "conf.d/disabled.conf",
            "params": [
                {
                    "context-type": "comment",
                    "value": "disabled config"
                },
                {
                    "enabled": true,
                    "context-type": "location",
                    "value": "~ /test",
                    "params": [
                        {
                            "enabled": true,
                            "context-type": "directive",
                            "value": "return 404"
                        }
                    ]
                }
            ]
        },
        "conf.d/enabled.conf": {
            "enabled": true,
            "context-type": "config",
            "value": "conf.d/enabled.conf",
            "params": [
                {
                    "enabled": true,
                    "context-type": "location",
                    "value": "~ /test",
                    "params": [
                        {
                            "enabled": true,
                            "context-type": "directive",
                            "value": "return 200 'test'"
                        }
                    ]
                }
            ]
        }
    }
}`,
			)},
			wantCheckConfigLines: []string{
				"http {",
				"    server {    # enabled server with enabled children configs",
				"        # include <== conf.d/enabled.conf",
				"        location ~ /test {",
				"            return 200 'test';",
				"        }",
				"    }",
				"    server {    # enabled server with disabled children configs",
				"        # include <== conf.d/disabled.conf",
				"        # # disabled config",
				"        # location ~ /test {",
				"        #     return 404;",
				"        # }",
				"    }",
				"    server {    # enabled server with disabled include context",
				"        # # include <== conf.d/enabled.conf",
				"    }",
				"    # server {    # disabled server with enabled children configs",
				"    #     # include <== conf.d/enabled.conf",
				"    #     location ~ /test {",
				"    #         return 200 'test';",
				"    #     }",
				"    # }",
				"    # server {    # disabled server with disabled children configs",
				"    #     # include <== conf.d/disabled.conf",
				"    #     # # disabled config",
				"    #     # location ~ /test {",
				"    #     #     return 404;",
				"    #     # }",
				"    # }",
				"    # server {    # disabled server with disabled include context",
				"    #     # # include <== conf.d/enabled.conf",
				"    # }",
				"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mainUnmarshaller{
				completedMain: tt.fields.completedMain,
			}
			if err := m.UnmarshalJSON(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %++v, wantErr %v", err, tt.wantErr)
			}
			if got, err := m.completedMain.ConfigLines(false); err != nil {
				t.Errorf("view unmarshaled main failed, error = %++v", err)
			} else if !reflect.DeepEqual(got, tt.wantCheckConfigLines) {
				t.Errorf("unmarshaled main ConfigLines() = %v, want %v", strings.Join(got, "\n"), strings.Join(tt.wantCheckConfigLines, "\n"))
			}
		})
	}
}

func Test_jsonMarshal(t *testing.T) {
	e := NewContext(context_type.TypeEvents, "")
	d, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(d))
}
