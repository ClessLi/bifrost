package local

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"reflect"
	"regexp"
	"testing"
)

func TestRegisterJsonRegMatcher(t *testing.T) {
	type args struct {
		contextType context_type.ContextType
		regexp      *regexp.Regexp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterJsonRegMatcher(tt.args.contextType, tt.args.regexp); (err != nil) != tt.wantErr {
				t.Errorf("RegisterJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterJsonUnmarshalerBuilder(t *testing.T) {
	type args struct {
		contextType context_type.ContextType
		newFunc     func() JsonUnmarshalContext
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterJsonUnmarshalerBuilder(tt.args.contextType, tt.args.newFunc); (err != nil) != tt.wantErr {
				t.Errorf("RegisterJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_jsonUnmarshalComment_GetChildren(t *testing.T) {
	type fields struct {
		Comments string
		Inline   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []*json.RawMessage
	}{
		{
			name: "return null",
			want: []*json.RawMessage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := jsonUnmarshalComment{
				Comments: tt.fields.Comments,
				Inline:   tt.fields.Inline,
			}
			if got := c.GetChildren(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChildren() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUnmarshalComment_GetValue(t *testing.T) {
	type fields struct {
		Comments string
		Inline   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "null comments",
			fields: fields{Comments: ""},
			want:   "",
		},
		{
			name:   "some comments",
			fields: fields{Comments: "   test xxx   "},
			want:   "   test xxx   ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := jsonUnmarshalComment{
				Comments: tt.fields.Comments,
				Inline:   tt.fields.Inline,
			}
			if got := c.GetValue(); got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUnmarshalComment_Type(t *testing.T) {
	type fields struct {
		Comments string
		Inline   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   context_type.ContextType
	}{
		{
			name:   "inline comment",
			fields: fields{Inline: true},
			want:   context_type.TypeInlineComment,
		},
		{
			name:   "not inline comment",
			fields: fields{Inline: false},
			want:   context_type.TypeComment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := jsonUnmarshalComment{
				Comments: tt.fields.Comments,
				Inline:   tt.fields.Inline,
			}
			if got := c.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUnmarshalContext_GetChildren(t *testing.T) {
	type fields struct {
		Value       string
		Children    []*json.RawMessage
		contextType context_type.ContextType
	}
	tests := []struct {
		name   string
		fields fields
		want   []*json.RawMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := jsonUnmarshalContext{
				Value:       tt.fields.Value,
				Children:    tt.fields.Children,
				contextType: tt.fields.contextType,
			}
			if got := u.GetChildren(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChildren() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUnmarshalContext_GetValue(t *testing.T) {
	type fields struct {
		Value       string
		Children    []*json.RawMessage
		contextType context_type.ContextType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := jsonUnmarshalContext{
				Value:       tt.fields.Value,
				Children:    tt.fields.Children,
				contextType: tt.fields.contextType,
			}
			if got := u.GetValue(); got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUnmarshalContext_Type(t *testing.T) {
	type fields struct {
		Value       string
		Children    []*json.RawMessage
		contextType context_type.ContextType
	}
	tests := []struct {
		name   string
		fields fields
		want   context_type.ContextType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := jsonUnmarshalContext{
				Value:       tt.fields.Value,
				Children:    tt.fields.Children,
				contextType: tt.fields.contextType,
			}
			if got := u.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUnmarshalDirective_GetChildren(t *testing.T) {
	type fields struct {
		Name   string
		Params string
	}
	tests := []struct {
		name   string
		fields fields
		want   []*json.RawMessage
	}{
		{
			name: "return null",
			want: []*json.RawMessage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := jsonUnmarshalDirective{
				Name:   tt.fields.Name,
				Params: tt.fields.Params,
			}
			if got := d.GetChildren(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChildren() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUnmarshalDirective_GetValue(t *testing.T) {
	type fields struct {
		Name   string
		Params string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "directive without params",
			fields: fields{Name: "server_name"},
			want:   "server_name",
		},
		{
			name: "directive with params",
			fields: fields{
				Name:   " server_name",
				Params: "testserver.com ",
			},
			want: "server_name testserver.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := jsonUnmarshalDirective{
				Name:   tt.fields.Name,
				Params: tt.fields.Params,
			}
			if got := d.GetValue(); got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUnmarshalDirective_Type(t *testing.T) {
	type fields struct {
		Name   string
		Params string
	}
	tests := []struct {
		name   string
		fields fields
		want   context_type.ContextType
	}{
		{
			name: "directive type",
			want: context_type.TypeDirective,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := jsonUnmarshalDirective{
				Name:   tt.fields.Name,
				Params: tt.fields.Params,
			}
			if got := d.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUnmarshaler_UnmarshalJSON(t *testing.T) {
	testTargetFatherCtx := NewContext(context_type.TypeServer, "")
	testMain, err := NewMain("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(NewComment("test comment", true), 0).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(NewDirective("server_name", "testserver"), 0).
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
		unmarshalContext JsonUnmarshalContext
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
				unmarshalContext: &jsonUnmarshalLocation{jsonUnmarshalContext{contextType: context_type.TypeLocation}},
				configGraph:      testMain.graph(),
				completedContext: context.NullContext(),
				fatherContext:    testTargetFatherCtx,
			},
			args: args{bytes: []byte(
				`{"location":
    {
        "value": "~ /target",
        "params":
        [
            {
                "inline": true,
                "comments": "target location"
            },
            {
                "include":
                {
                    "value": "conf.d\\*conf",
                    "params":
                    [
                        "conf.d\\proxy.conf"
                    ]
                }
            }
        ]
    }
}`,
			)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &jsonUnmarshaler{
				unmarshalContext: tt.fields.unmarshalContext,
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

func Test_jsonUnmarshaler_nextUnmarshaler(t *testing.T) {
	type fields struct {
		unmarshalContext JsonUnmarshalContext
		configGraph      ConfigGraph
		completedContext context.Context
		fatherContext    context.Context
	}
	type args struct {
		message *json.RawMessage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *jsonUnmarshaler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &jsonUnmarshaler{
				unmarshalContext: tt.fields.unmarshalContext,
				configGraph:      tt.fields.configGraph,
				completedContext: tt.fields.completedContext,
				fatherContext:    tt.fields.fatherContext,
			}
			if got := u.nextUnmarshaler(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nextUnmarshaler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonUnmarshaler_unmarshalInclude(t *testing.T) {
	type fields struct {
		unmarshalContext JsonUnmarshalContext
		configGraph      ConfigGraph
		completedContext context.Context
		fatherContext    context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &jsonUnmarshaler{
				unmarshalContext: tt.fields.unmarshalContext,
				configGraph:      tt.fields.configGraph,
				completedContext: tt.fields.completedContext,
				fatherContext:    tt.fields.fatherContext,
			}
			if err := u.unmarshalInclude(); (err != nil) != tt.wantErr {
				t.Errorf("unmarshalInclude() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mainUnmarshaler_UnmarshalJSON(t *testing.T) {
	type fields struct {
		unmarshalContext *jsonUnmarshalMain
		completedMain    MainContext
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
			name:   "normal test",
			fields: fields{unmarshalContext: new(jsonUnmarshalMain)},
			args: args{bytes: []byte(
				`{
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
}`,
			)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mainUnmarshaler{
				unmarshalContext: tt.fields.unmarshalContext,
				completedMain:    tt.fields.completedMain,
			}
			if err := m.UnmarshalJSON(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerCommentJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerCommentJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerCommentJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerCommentJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerCommentJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerCommentJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerDirectiveJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerDirectiveJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerDirectiveJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerEventsJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerEventsJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerEventsJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerEventsJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerEventsJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerEventsJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerGEOJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerGEOJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerGEOJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerGEOJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerGEOJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerGEOJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerHTTPJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerHTTPJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerHTTPJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerHTTPJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerHTTPJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerHTTPJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerIFJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerIFJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerIFJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerIFJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerIFJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerIFJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerIncludeJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerIncludeJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerIncludeJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerIncludeJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerIncludeJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerIncludeJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerJsonRegMatchers(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerJsonRegMatchers(); (err != nil) != tt.wantErr {
				t.Errorf("registerJsonRegMatchers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerJsonUnmarshalerBuilders(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerJsonUnmarshalerBuilders(); (err != nil) != tt.wantErr {
				t.Errorf("registerJsonUnmarshalerBuilders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerLimitExceptJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerLimitExceptJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerLimitExceptJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerLimitExceptJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerLimitExceptJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerLimitExceptJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerLocationJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerLocationJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerLocationJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerLocationJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerLocationJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerLocationJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerMapJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerMapJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerMapJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerMapJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerMapJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerMapJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerServerJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerServerJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerServerJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerServerJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerServerJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerServerJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerStreamJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerStreamJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerStreamJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerStreamJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerStreamJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerStreamJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerTypesJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerTypesJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerTypesJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerTypesJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerTypesJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerTypesJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerUpstreamJsonRegMatcher(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerUpstreamJsonRegMatcher(); (err != nil) != tt.wantErr {
				t.Errorf("registerUpstreamJsonRegMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerUpstreamJsonUnmarshalerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerUpstreamJsonUnmarshalerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerUpstreamJsonUnmarshalerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
