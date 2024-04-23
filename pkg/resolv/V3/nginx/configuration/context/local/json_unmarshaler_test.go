package local

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"reflect"
	"testing"
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
    "context-type": "location",
    "value": "~ /target",
    "params":
    [
        {
            "context-type": "inline_comment",
            "value": "target location"
        },
        {
            "context-type": "include",
            "value": "conf.d\\*conf",
            "params":
            [
                "conf.d\\proxy.conf"
            ]
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
            "context-type": "location",
            "value": "~ /test_proxy",
            "params": [
                {
                    "context-type": "directive",
                    "value": "proxy_pass https://baidu.com"
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
            "context-type": "location",
            "value": "~ /test_proxy",
            "params": [
                {
                    "context-type": "directive",
                    "value": "proxy_pass https://baidu.com"
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
			u := &jsonUnmarshaler{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mainUnmarshaler{
				completedMain: tt.fields.completedMain,
			}
			if err := m.UnmarshalJSON(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got, err := m.completedMain.ConfigLines(false); err != nil {
				t.Errorf("view unmarshaled main failed, error = %v", err)
			} else if !reflect.DeepEqual(got, tt.wantCheckConfigLines) {
				t.Errorf("unmarshaled main ConfigLines() = %v, want %v", got, tt.wantCheckConfigLines)
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
