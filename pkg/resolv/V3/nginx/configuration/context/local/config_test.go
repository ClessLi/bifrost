package local

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/dominikbraun/graph"
)

func TestConfig_Clone(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
			),
		0,
	)
	testCloneChildren := make([]context.Context, 0)
	for _, child := range testMain.MainConfig().Children {
		testCloneChildren = append(testCloneChildren, child.Clone())
	}
	type fields struct {
		BasicContext BasicContext
		ConfigPath   context.ConfigPath
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name: "clone without graph and ConfigPath",
			fields: fields{
				BasicContext: testMain.MainConfig().BasicContext,
				ConfigPath:   testMain.MainConfig().ConfigPath,
			},
			want: &Config{
				BasicContext: BasicContext{
					Enabled:        true,
					ContextType:    context_type.TypeConfig,
					ContextValue:   testMain.MainConfig().ContextValue,
					Children:       testCloneChildren,
					father:         context.NullContext(),
					self:           context.NullContext(),
					headStringFunc: testMain.MainConfig().headStringFunc,
					tailStringFunc: testMain.MainConfig().tailStringFunc,
				},
				ConfigPath: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				ConfigPath:   tt.fields.ConfigPath,
			}
			got := c.Clone()
			gotlines, err := got.ConfigLines(false)
			if err != nil {
				t.Errorf("got.ConfigLines() return error: %v", err)

				return
			}
			wantlines, err := tt.want.ConfigLines(false)
			if err != nil {
				t.Errorf("want.ConfigLines() return error: %v", err)

				return
			}
			if reflect.DeepEqual(got.(*Config).Children, tt.want.(*Config).Children) ||
				!reflect.DeepEqual(gotlines, wantlines) ||
				got.Father().Type() != tt.want.Father().Type() ||
				got.Type() != tt.want.Type() ||
				got.Value() != tt.want.Value() ||
				got.IsEnabled() != tt.want.IsEnabled() ||
				!reflect.DeepEqual(got.(*Config).ConfigPath, tt.want.(*Config).ConfigPath) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_ConfigLines(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").Insert(
			NewContext(context_type.TypeServer, "").Insert(
				NewContext(context_type.TypeDirective, "listen 80"),
				0,
			).Insert(
				NewContext(context_type.TypeDirective, "server_name example.com"),
				1,
			).Insert(
				NewContext(context_type.TypeInclude, "conf.d/disabled_location.conf"),
				2,
			).Insert(
				NewContext(context_type.TypeInclude, "conf.d/strange_location.conf"),
				3,
			),
			0,
		).Insert(
			NewContext(context_type.TypeComment, "disabled server context"),
			1,
		).Insert(
			NewContext(context_type.TypeServer, "").Disable().Insert(
				NewContext(context_type.TypeInlineComment, "disabled server"),
				0,
			).Insert(
				NewContext(context_type.TypeDirective, "listen 8080"),
				1,
			).Insert(
				NewContext(context_type.TypeDirective, "server_name example.com"),
				2,
			).Insert(
				NewContext(context_type.TypeLocation, "~ /disabled-location").Disable().Insert(
					NewContext(context_type.TypeDirective, "proxy_pass http://disabled-url"),
					0,
				),
				3,
			).Insert(
				NewContext(context_type.TypeInclude, "conf.d/disabled_location.conf"),
				4,
			).Insert(
				NewContext(context_type.TypeInclude, "conf.d/strange_location.conf"),
				5,
			),
			2,
		),
		0,
	)
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "conf.d/disabled_location.conf").Disable().Insert(
			NewContext(context_type.TypeComment, "disabled config"),
			0,
		).Insert(
			NewContext(context_type.TypeLocation, "~ /test").Insert(
				NewContext(context_type.TypeDirective, "return 404"),
				0,
			),
			1,
		).Insert(
			NewContext(context_type.TypeLocation, "~ /has-disabled-ctx").Insert(
				NewContext(context_type.TypeComment, "disabled if ctx"),
				0,
			).Insert(
				NewContext(context_type.TypeIf, "($is_enabled ~* false)").Disable().Insert(
					NewContext(context_type.TypeDirective, "set $is_enabled true").Disable(),
					0,
				).Insert(
					NewContext(context_type.TypeDirective, "return 404"),
					1,
				),
				1,
			),
			2,
		).Insert(
			NewContext(context_type.TypeComment, "}"),
			3,
		).Insert(
			NewContext(context_type.TypeComment, "}"),
			4,
		).(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "conf.d/strange_location.conf").Insert(
			NewContext(context_type.TypeComment, "strange config"),
			0,
		).Insert(
			NewContext(context_type.TypeLocation, "~ /normal-loc").Insert(
				NewContext(context_type.TypeDirective, "return 200"),
				0,
			),
			1,
		).Insert(
			NewContext(context_type.TypeComment, "location ~ /strange-loc {"),
			2,
		).Insert(
			NewContext(context_type.TypeComment, "    if ($strange ~* this_is_a_strange_if_ctx) {"),
			3,
		).Insert(
			NewContext(context_type.TypeComment, "        return 404;"),
			4,
		).Insert(
			NewContext(context_type.TypeComment, "    proxy_pass http://strange_url;"),
			5,
		).Insert(
			NewContext(context_type.TypeComment, "}"),
			6,
		).(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	disabledConfig, err := testMain.GetConfig("C:\\test\\conf.d\\disabled_location.conf")
	if err != nil {
		t.Fatal(err)
	}
	strangeConfig, err := testMain.GetConfig("C:\\test\\conf.d\\strange_location.conf")
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		BasicContext BasicContext
		ConfigPath   context.ConfigPath
	}
	type args struct {
		isDumping bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "normal test for dump",
			fields: fields{
				BasicContext: testMain.MainConfig().BasicContext,
				ConfigPath:   testMain.MainConfig().ConfigPath,
			},
			args: args{isDumping: true},
			want: []string{
				"http {",
				"    server {",
				"        listen 80;",
				"        server_name example.com;",
				"        include conf.d/disabled_location.conf;",
				"        include conf.d/strange_location.conf;",
				"    }",
				"    # disabled server context",
				"    # server {    # disabled server",
				"    #     listen 8080;",
				"    #     server_name example.com;",
				"    #     # location ~ /disabled-location {",
				"    #     #     proxy_pass http://disabled-url;",
				"    #     # }",
				"    #     include conf.d/disabled_location.conf;",
				"    #     include conf.d/strange_location.conf;",
				"    # }",
				"}",
			},
		},
		{
			name: "dumping for the disabled config",
			fields: fields{
				BasicContext: disabledConfig.BasicContext,
				ConfigPath:   disabledConfig.ConfigPath,
			},
			args: args{isDumping: true},
			want: []string{
				"# # disabled config",
				"# location ~ /test {",
				"#     return 404;",
				"# }",
				"# location ~ /has-disabled-ctx {",
				"#     # disabled if ctx",
				"#     # if ($is_enabled ~* false) {",
				"#     #     # set $is_enabled true;",
				"#     #     return 404;",
				"#     # }",
				"# }",
				"# # }",
				"# # }",
			},
		},
		{
			name: "dumping for the strange config",
			fields: fields{
				BasicContext: strangeConfig.BasicContext,
				ConfigPath:   strangeConfig.ConfigPath,
			},
			args: args{isDumping: true},
			want: []string{
				"# strange config",
				"location ~ /normal-loc {",
				"    return 200;",
				"}",
				"# location ~ /strange-loc {",
				"#     if ($strange ~* this_is_a_strange_if_ctx) {",
				"#         return 404;",
				"#     proxy_pass http://strange_url;",
				"# }",
			},
		},
		{
			name: "normal test for view",
			fields: fields{
				BasicContext: testMain.MainConfig().BasicContext,
				ConfigPath:   testMain.MainConfig().ConfigPath,
			},
			args: args{isDumping: false},
			want: []string{
				"http {",
				"    server {",
				"        listen 80;",
				"        server_name example.com;",
				"        # include <== conf.d/disabled_location.conf",
				"        # # disabled config",
				"        # location ~ /test {",
				"        #     return 404;",
				"        # }",
				"        # location ~ /has-disabled-ctx {",
				"        #     # disabled if ctx",
				"        #     # if ($is_enabled ~* false) {",
				"        #     #     # set $is_enabled true;",
				"        #     #     return 404;",
				"        #     # }",
				"        # }",
				"        # # }",
				"        # # }",
				"        # include <== conf.d/strange_location.conf",
				"        # strange config",
				"        location ~ /normal-loc {",
				"            return 200;",
				"        }",
				"        # location ~ /strange-loc {",
				"        #     if ($strange ~* this_is_a_strange_if_ctx) {",
				"        #         return 404;",
				"        #     proxy_pass http://strange_url;",
				"        # }",
				"    }",
				"    # disabled server context",
				"    # server {    # disabled server",
				"    #     listen 8080;",
				"    #     server_name example.com;",
				"    #     # location ~ /disabled-location {",
				"    #     #     proxy_pass http://disabled-url;",
				"    #     # }",
				"    #     # include <== conf.d/disabled_location.conf",
				"    #     # # disabled config",
				"    #     # location ~ /test {",
				"    #     #     return 404;",
				"    #     # }",
				"    #     # location ~ /has-disabled-ctx {",
				"    #     #     # disabled if ctx",
				"    #     #     # if ($is_enabled ~* false) {",
				"    #     #     #     # set $is_enabled true;",
				"    #     #     #     return 404;",
				"    #     #     # }",
				"    #     # }",
				"    #     # # }",
				"    #     # # }",
				"    #     # include <== conf.d/strange_location.conf",
				"    #     # strange config",
				"    #     location ~ /normal-loc {",
				"    #         return 200;",
				"    #     }",
				"    #     # location ~ /strange-loc {",
				"    #     #     if ($strange ~* this_is_a_strange_if_ctx) {",
				"    #     #         return 404;",
				"    #     #     proxy_pass http://strange_url;",
				"    #     # }",
				"    # }",
				"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				ConfigPath:   tt.fields.ConfigPath,
			}
			got, err := c.ConfigLines(tt.args.isDumping)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigLines() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigLines() got = %v, want %v", strings.Join(got, "\n"), strings.Join(tt.want, "\n"))
			}
		})
	}
}

func TestConfig_SetFather(t *testing.T) {
	testMain, err := NewMain("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		BasicContext BasicContext
		ConfigPath   context.ConfigPath
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "return error",
			wantErr: true,
		},
		{
			name:    "invalid father context",
			args:    args{ctx: NewContext(context_type.TypeConfig, "test.conf")},
			wantErr: true,
		},
		{
			name:    "normal test",
			fields:  fields{BasicContext: BasicContext{father: context.NullContext()}},
			args:    args{ctx: testMain},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				ConfigPath:   tt.fields.ConfigPath,
			}
			if err := c.SetFather(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SetFather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_SetValue(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
			),
		0,
	)
	err = testMain.AddConfig(NewContext(context_type.TypeConfig, "a.conf").(*Config))
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		BasicContext BasicContext
		ConfigPath   context.ConfigPath
	}
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "failed to set value",
			fields: fields{
				BasicContext: testMain.MainConfig().BasicContext,
				ConfigPath:   testMain.MainConfig().ConfigPath,
			},
			args:    args{v: "C:\\test\\a.conf"},
			wantErr: true,
		},
		{
			name: "set value",
			fields: fields{
				BasicContext: aConfig.BasicContext,
				ConfigPath:   aConfig.ConfigPath,
			},
			args:    args{v: "b.conf"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				ConfigPath:   tt.fields.ConfigPath,
			}
			cache, err := c.Father().(MainContext).GetConfig(configHash(c))
			if err != nil {
				t.Fatal(err)
			}
			if err := cache.SetValue(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_isInGraph(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
			),
		0,
	)
	nilGraphConfig := NewContext(context_type.TypeConfig, "nilgraph").(*Config)
	nilGraphConfig.ConfigPath, _ = newConfigPath(testMain, nilGraphConfig.Value())

	nilConfigPathConfig := NewContext(context_type.TypeConfig, "nilpath").(*Config)
	nilConfigPathConfig.SetFather(testMain)

	notInGraphConfig := NewContext(context_type.TypeConfig, "notingraph").(*Config)
	notInGraphConfig.ConfigPath, _ = newConfigPath(testMain, notInGraphConfig.Value())
	notInGraphConfig.SetFather(testMain)

	inGraphConfig := NewContext(context_type.TypeConfig, "ingraph").(*Config)
	inGraphConfig.ConfigPath, _ = newConfigPath(testMain, inGraphConfig.Value())
	inGraphConfig.SetFather(testMain)
	err = inGraphConfig.Father().(MainContext).AddConfig(inGraphConfig)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		BasicContext BasicContext
		ConfigPath   context.ConfigPath
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "config has nil graph",
			fields: fields{
				BasicContext: nilGraphConfig.BasicContext,
				ConfigPath:   nilGraphConfig.ConfigPath,
			},
			want: false,
		},
		{
			name: "config has nil config path",
			fields: fields{
				BasicContext: nilConfigPathConfig.BasicContext,
				ConfigPath:   nilConfigPathConfig.ConfigPath,
			},
			want: false,
		},
		{
			name: "config has not been added into a graph",
			fields: fields{
				BasicContext: notInGraphConfig.BasicContext,
				ConfigPath:   notInGraphConfig.ConfigPath,
			},
			want: false,
		},
		{
			name: "config has been added into a graph",
			fields: fields{
				BasicContext: inGraphConfig.BasicContext,
				ConfigPath:   inGraphConfig.ConfigPath,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				ConfigPath:   tt.fields.ConfigPath,
			}
			if got := c.isInGraph(); got != tt.want {
				t.Errorf("isInGraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_mainContext(t *testing.T) {
	type fields struct {
		BasicContext BasicContext
		ConfigPath   context.ConfigPath
	}
	tests := []struct {
		name    string
		fields  fields
		want    MainContext
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				ConfigPath:   tt.fields.ConfigPath,
			}
			got, err := c.mainContext()
			if (err != nil) != tt.wantErr {
				t.Errorf("mainContext() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mainContext() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configGraph_AddConfig(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
			),
		0,
	)
	err = testMain.AddConfig(NewContext(context_type.TypeConfig, "a.conf").(*Config))
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	bConfig := NewContext(context_type.TypeConfig, "b.conf").(*Config)
	bConfig.SetFather(testMain)
	bConfig.ConfigPath, _ = newConfigPath(testMain, bConfig.Value())
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		wantErr bool
	}{
		{
			name:    "add nil config",
			fields:  testMain.graph(),
			args:    args{config: nil},
			wantErr: true,
		},
		{
			name:    "add an already exist config",
			fields:  testMain.graph(),
			args:    args{config: aConfig},
			wantErr: true,
		},
		{
			name:    "normal test",
			fields:  testMain.graph(),
			args:    args{config: bConfig},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.AddConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("AddConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_AddEdge(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
			),
		0,
	)
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "a.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "b.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	bConfig, err := testMain.GetConfig("C:\\test\\b.conf")
	if err != nil {
		t.Fatal(err)
	}

	otherMain, _ := NewMain("C:\\test1\\nginx.conf")
	err = otherMain.AddConfig(
		NewContext(context_type.TypeConfig, "other.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	inOtherGraphConfig, err := otherMain.GetConfig("C:\\test1\\other.conf")
	if err != nil {
		t.Fatal(err)
	}

	nullpathConfig := NewContext(context_type.TypeConfig, "").(*Config)
	nullpathConfig.ConfigPath = &context.AbsConfigPath{}

	excludeConfig := NewContext(context_type.TypeConfig, "exclude.conf").(*Config)
	excludeConfig.ConfigPath, _ = newConfigPath(testMain, excludeConfig.Value())

	type args struct {
		src *Config
		dst *Config
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		wantErr bool
	}{
		{
			name:   "nil source config",
			fields: testMain.graph(),
			args: args{
				src: nil,
				dst: aConfig,
			},
			wantErr: true,
		},
		{
			name:   "nil destination config",
			fields: testMain.graph(),
			args: args{
				src: aConfig,
				dst: nil,
			},
			wantErr: true,
		},
		{
			name:   "source config with null config path",
			fields: testMain.graph(),
			args: args{
				src: nullpathConfig,
				dst: aConfig,
			},
			wantErr: true,
		},
		{
			name:   "destination config with null config path",
			fields: testMain.graph(),
			args: args{
				src: aConfig,
				dst: nullpathConfig,
			},
			wantErr: true,
		},
		{
			name:   "source config is exclude from the graph",
			fields: testMain.graph(),
			args: args{
				src: excludeConfig,
				dst: aConfig,
			},
			wantErr: true,
		},
		{
			name:   "destination config is exclude from the graph",
			fields: testMain.graph(),
			args: args{
				src: aConfig,
				dst: excludeConfig,
			},
			wantErr: true,
		},
		{
			name:   "source config is in the other graph",
			fields: testMain.graph(),
			args: args{
				src: inOtherGraphConfig,
				dst: aConfig,
			},
			wantErr: true,
		},
		{
			name:   "destination config is in the other graph",
			fields: testMain.graph(),
			args: args{
				src: aConfig,
				dst: inOtherGraphConfig,
			},
			wantErr: true,
		},
		{
			name:   "normal test",
			fields: testMain.graph(),
			args: args{
				src: aConfig,
				dst: bConfig,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.AddEdge(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("AddEdge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_GetConfig(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
			),
		0,
	)
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "a.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		fullpath string
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name:    "wrong config path",
			fields:  testMain.graph(),
			args:    args{fullpath: "wrong/config/path.conf"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "normal test",
			fields:  testMain.graph(),
			args:    args{fullpath: "C:\\test\\a.conf"},
			want:    aConfig,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			got, err := c.GetConfig(tt.args.fullpath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configGraph_ListConfigs(t *testing.T) {
	type fields struct {
		graph      graph.Graph[string, *Config]
		mainConfig *Config
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &configGraph{
				graph:      tt.fields.graph,
				mainConfig: tt.fields.mainConfig,
			}
			if got := c.ListConfigs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configGraph_MainConfig(t *testing.T) {
	type fields struct {
		graph      graph.Graph[string, *Config]
		mainConfig *Config
	}
	tests := []struct {
		name   string
		fields fields
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &configGraph{
				graph:      tt.fields.graph,
				mainConfig: tt.fields.mainConfig,
			}
			if got := c.MainConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MainConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configGraph_RemoveConfig(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
					).
					Insert(
						NewContext(context_type.TypeInclude, "test.conf"),
						2,
					),
				1,
			),
		0,
	)
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "a.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "2exist.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "inedge.conf").Insert(NewContext(context_type.TypeInclude, "test.conf"), 0).(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "test.conf").Insert(NewContext(context_type.TypeInclude, "outedge.conf"), 0).(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}

	nonexistentConfig := NewContext(context_type.TypeConfig, "C:\\test\\nonexistent_config.conf").(*Config)
	nonexistentConfig.ConfigPath, err = newConfigPath(testMain, nonexistentConfig.ContextValue)
	if err != nil {
		t.Fatal(err)
	}

	testConfig, err := testMain.GetConfig("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		wantErr bool
	}{
		{
			name:    "remove config does not exist",
			fields:  testMain.graph(),
			args:    args{config: nonexistentConfig},
			wantErr: true,
		},
		{
			name:    "remove main config",
			fields:  testMain.graph(),
			args:    args{config: testMain.MainConfig()},
			wantErr: true,
		},
		{
			name:    "normal test",
			fields:  testMain.graph(),
			args:    args{config: testConfig},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.RemoveConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("RemoveConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_RemoveEdge(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
			),
		0,
	)
	err = testMain.addVertex(NewContext(context_type.TypeConfig, "a.conf").(*Config))
	if err != nil {
		t.Fatal(err)
	}
	a, err := testMain.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.addVertex(NewContext(context_type.TypeConfig, "b.conf").(*Config))
	if err != nil {
		t.Fatal(err)
	}
	b, err := testMain.GetConfig("C:\\test\\b.conf")
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddEdge(a, b)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.addVertex(NewContext(context_type.TypeConfig, "c.conf").(*Config))
	if err != nil {
		t.Fatal(err)
	}
	c, err := testMain.GetConfig("C:\\test\\c.conf")
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddEdge(b, c)
	if err != nil {
		t.Fatal(err)
	}
	notInGraphConfig := NewContext(context_type.TypeConfig, "notingraph.conf").(*Config)
	type args struct {
		src     *Config
		dst     *Config
		keepDst bool
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		wantErr bool
	}{
		{
			name:   "removed edge not found",
			fields: testMain.graph(),
			args: args{
				src: b,
				dst: notInGraphConfig,
			},
			wantErr: true,
		},
		{
			name:   "destination has edge",
			fields: testMain.graph(),
			args: args{
				src: a,
				dst: b,
			},
			wantErr: false,
		},
		{
			name:   "remove edge and destination",
			fields: testMain.graph(),
			args: args{
				src: b,
				dst: c,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.RemoveEdge(tt.args.src, tt.args.dst, tt.args.keepDst); (err != nil) != tt.wantErr {
				t.Errorf("RemoveEdge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_RenameConfig(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
					).
					Insert(
						NewContext(context_type.TypeInclude, "test.conf"),
						2,
					),
				1,
			),
		0,
	)
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "a.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "2exist.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "inedge.conf").Insert(NewContext(context_type.TypeInclude, "test.conf"), 0).(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "test.conf").Insert(NewContext(context_type.TypeInclude, "outedge.conf"), 0).(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	renew2existConfig, err := testMain.GetConfig("C:\\test\\2exist.conf")
	if err != nil {
		t.Fatal(err)
	}
	renew2existConfig.ConfigPath = aConfig.ConfigPath

	//inEdgeConfig, err := testMain.GetConfig("C:\\test\\inedge.conf")
	//if err != nil {
	//	t.Fatal(err)
	//}

	testConfig, err := testMain.GetConfig("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}
	// testMain.AddEdge(inEdgeConfig, testConfig)
	// testMain.AddEdge(testConfig, NewContext(context_type.TypeConfig, "outedge.conf").(*Config))
	test2Main := testMain.Clone().(MainContext)
	test3Main := testMain.Clone().(MainContext)

	type args struct {
		oldFullPath string
		newPath     string
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		wantErr bool
	}{
		{
			name:   "not exist config",
			fields: testMain.graph(),
			args: args{
				oldFullPath: "notexist.conf",
				newPath:     "test.conf",
			},
			wantErr: true,
		},
		{
			name:   "need not renew config",
			fields: testMain.graph(),
			args: args{
				oldFullPath: configHash(testMain.MainConfig()),
				newPath:     configHash(testMain.MainConfig()),
			},
			wantErr: false,
		},
		{
			name:   "renew to exist config",
			fields: testMain.graph(),
			args: args{
				oldFullPath: configHash(aConfig),
				newPath:     "2exist.conf",
			},
			wantErr: true,
		},
		{
			name:   "normal test",
			fields: testMain.graph(),
			args: args{
				oldFullPath: configHash(testConfig),
				newPath:     "modified.conf",
			},
			wantErr: false,
		},
		{
			name:   "rename main config(relative path)",
			fields: testMain.graph(),
			args: args{
				oldFullPath: configHash(testMain.MainConfig()),
				newPath:     "modified_nginx.conf",
			},
			wantErr: false,
		},
		{
			name:   "rename main config(absolut path)",
			fields: test2Main.graph(),
			args: args{
				oldFullPath: configHash(test2Main.MainConfig()),
				newPath:     "C:\\test\\another_modified_nginx.conf",
			},
			wantErr: false,
		},
		{
			name:   "rename main config to another dir(relative path)",
			fields: test3Main.graph(),
			args: args{
				oldFullPath: configHash(test3Main.MainConfig()),
				newPath:     "../another_dir/modified_nginx.conf",
			},
			wantErr: true,
		},
		{
			name:   "rename main config to another dir(absolut path)",
			fields: test3Main.graph(),
			args: args{
				oldFullPath: configHash(test3Main.MainConfig()),
				newPath:     "D:\\another_dir\\modified_nginx.conf",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.RenameConfig(tt.args.oldFullPath, tt.args.newPath); (err != nil) != tt.wantErr {
				t.Errorf("RenameConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_Topology(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
					).
					Insert(
						NewContext(context_type.TypeInclude, "a.conf"),
						2,
					),
				1,
			),
		0,
	)
	testMain.AddConfig(NewContext(context_type.TypeConfig, "a.conf").Insert(NewContext(context_type.TypeInclude, "b.conf"), 0).(*Config))
	testMain.AddConfig(NewContext(context_type.TypeConfig, "b.conf").Insert(NewContext(context_type.TypeInclude, "c.conf"), 0).(*Config))
	testMain.AddConfig(NewContext(context_type.TypeConfig, "c.conf").Insert(NewContext(context_type.TypeInclude, "d.conf"), 0).(*Config))
	testMain.AddConfig(NewContext(context_type.TypeConfig, "d.conf").(*Config))
	testMain.AddConfig(
		NewContext(context_type.TypeConfig, "e.conf").
			Insert(
				NewContext(context_type.TypeInclude, "f.conf"),
				0,
			).
			Insert(
				NewContext(context_type.TypeInclude, "g.conf"),
				1,
			).(*Config),
	)
	testMain.AddConfig(NewContext(context_type.TypeConfig, "f.conf").(*Config))
	testMain.AddConfig(NewContext(context_type.TypeConfig, "g.conf").(*Config))
	a, _ := testMain.GetConfig("C:\\test\\a.conf")
	b, _ := testMain.GetConfig("C:\\test\\b.conf")
	c, _ := testMain.GetConfig("C:\\test\\c.conf")
	d, _ := testMain.GetConfig("C:\\test\\d.conf")

	tests := []struct {
		name   string
		fields ConfigGraph
		want   []*Config
	}{
		{
			name:   "generate only one tree starting from the main config",
			fields: testMain.graph(),
			want:   []*Config{testMain.MainConfig(), a, b, c, d},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if got := c.Topology(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Topology() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configGraph_addVertex(t *testing.T) {
	type fields struct {
		graph      graph.Graph[string, *Config]
		mainConfig *Config
	}
	type args struct {
		v *Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &configGraph{
				graph:      tt.fields.graph,
				mainConfig: tt.fields.mainConfig,
			}
			if err := c.addVertex(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("addVertex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_checkOperatedVertex(t *testing.T) {
	type fields struct {
		graph      graph.Graph[string, *Config]
		mainConfig *Config
	}
	type args struct {
		v *Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &configGraph{
				graph:      tt.fields.graph,
				mainConfig: tt.fields.mainConfig,
			}
			if err := c.checkOperatedVertex(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("checkOperatedVertex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_cleanupGraph(t *testing.T) {
	type fields struct {
		graph      graph.Graph[string, *Config]
		mainConfig *Config
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
			c := &configGraph{
				graph:      tt.fields.graph,
				mainConfig: tt.fields.mainConfig,
			}
			if err := c.cleanupGraph(); (err != nil) != tt.wantErr {
				t.Errorf("cleanupGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_removeVertex(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
					).
					Insert(
						NewContext(context_type.TypeInclude, "a.conf"),
						2,
					),
				1,
			),
		0,
	)
	testMain.AddConfig(NewContext(context_type.TypeConfig, "a.conf").Insert(NewContext(context_type.TypeInclude, "b.conf"), 0).(*Config))
	testMain.AddConfig(NewContext(context_type.TypeConfig, "b.conf").Insert(NewContext(context_type.TypeInclude, "c.conf"), 0).(*Config))
	testMain.AddConfig(NewContext(context_type.TypeConfig, "c.conf").Insert(NewContext(context_type.TypeInclude, "d.conf"), 0).(*Config))
	testMain.AddConfig(NewContext(context_type.TypeConfig, "d.conf").(*Config))
	testMain.AddConfig(
		NewContext(context_type.TypeConfig, "e.conf").
			Insert(
				NewContext(context_type.TypeInclude, "a.conf"),
				0,
			).
			Insert(
				NewContext(context_type.TypeInclude, "f.conf"),
				1,
			).
			Insert(
				NewContext(context_type.TypeInclude, "g.conf"),
				2,
			).(*Config),
	)
	testMain.AddConfig(NewContext(context_type.TypeConfig, "f.conf").(*Config))
	testMain.AddConfig(NewContext(context_type.TypeConfig, "g.conf").(*Config))
	a, _ := testMain.GetConfig("C:\\test\\a.conf")
	c, _ := testMain.GetConfig("C:\\test\\c.conf")
	d, _ := testMain.GetConfig("C:\\test\\d.conf")
	e, _ := testMain.GetConfig("C:\\test\\e.conf")
	e.Child(0).(*Include).load()
	err = testMain.graph().(*configGraph).graph.RemoveEdge(configHash(c), configHash(d))
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		fields  *configGraph
		args    args
		wantErr bool
	}{
		{
			name:    "normal test",
			fields:  testMain.graph().(*configGraph),
			args:    args{d},
			wantErr: false,
		},
		// Sub objects that are not in the topology will no longer be loaded, and this use case will not exist
		//{
		//	name:    "config has no in edge but out edges",
		//	fields:  testMain.graph().(*configGraph),
		//	args:    args{e},
		//	wantErr: true,
		//},
		{
			name:    "config has edges",
			fields:  testMain.graph().(*configGraph),
			args:    args{a},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.removeVertex(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("removeVertex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_renderGraph(t *testing.T) {
	type fields struct {
		graph      graph.Graph[string, *Config]
		mainConfig *Config
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
			c := &configGraph{
				graph:      tt.fields.graph,
				mainConfig: tt.fields.mainConfig,
			}
			if err := c.renderGraph(); (err != nil) != tt.wantErr {
				t.Errorf("renderGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_rerenderGraph(t *testing.T) {
	type fields struct {
		graph      graph.Graph[string, *Config]
		mainConfig *Config
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
			c := &configGraph{
				graph:      tt.fields.graph,
				mainConfig: tt.fields.mainConfig,
			}
			if err := c.rerenderGraph(); (err != nil) != tt.wantErr {
				t.Errorf("rerenderGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_setFatherFor(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
			),
		0,
	)
	testMain.AddConfig(NewContext(context_type.TypeConfig, "a.conf").(*Config))
	a, _ := testMain.GetConfig("C:\\test\\a.conf")
	absolutPathConfig := NewContext(context_type.TypeConfig, "D:\\absolut_path\\a.conf").(*Config)
	absolutPathConfig.ConfigPath, _ = context.NewAbsConfigPath("D:\\absolut_path\\a.conf")
	diffMain, _ := NewMain("C:\\test2\\nginx.conf")
	diffGraphConfig := NewContext(context_type.TypeConfig, "different_graph.conf").(*Config)
	err = diffMain.AddConfig(diffGraphConfig)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		fields  *configGraph
		args    args
		wantErr bool
	}{
		{
			name:    "nil config",
			fields:  testMain.graph().(*configGraph),
			wantErr: true,
		},
		{
			name:    "config in the other graph",
			fields:  testMain.graph().(*configGraph),
			args:    args{diffGraphConfig},
			wantErr: true,
		},
		{
			name:    "config clone",
			fields:  testMain.graph().(*configGraph),
			args:    args{a.Clone().(*Config)},
			wantErr: false,
		},
		{
			name:    "same graph config",
			fields:  testMain.graph().(*configGraph),
			args:    args{a},
			wantErr: false,
		},
		{
			name:    "new config",
			fields:  testMain.graph().(*configGraph),
			args:    args{NewContext(context_type.TypeConfig, "new.conf").(*Config)},
			wantErr: false,
		},
		{
			name:    "absolut path config",
			fields:  testMain.graph().(*configGraph),
			args:    args{absolutPathConfig},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.setFatherFor(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("setFatherFor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configHash(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
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
			),
		0,
	)
	type args struct {
		t *Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil config",
			want: "",
		},
		{
			name: "new config with no config path",
			args: args{t: NewContext(context_type.TypeConfig, "new.conf").(*Config)},
			want: "",
		},
		{
			name: "normal test",
			args: args{t: testMain.MainConfig()},
			want: "C:\\test\\nginx.conf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := configHash(tt.args.t); got != tt.want {
				t.Errorf("configHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newConfigGraph(t *testing.T) {
	type args struct {
		mainConfig *Config
	}
	tests := []struct {
		name    string
		args    args
		want    ConfigGraph
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfigGraph(tt.args.mainConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("newConfigGraph() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConfigGraph() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newConfigPath(t *testing.T) {
	// main config
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	relPathMain, _ := NewMain("C:\\test\\nginx.conf")
	relPathMain.MainConfig().ConfigPath = &context.RelConfigPath{}
	absCP, _ := context.NewAbsConfigPath("D:\\test\\test.conf")
	relCP, _ := context.NewRelConfigPath(testMain.MainConfig().ConfigPath.BaseDir(), "relative.conf")
	type args struct {
		configgraph   ConfigGraph
		newconfigpath string
	}
	tests := []struct {
		name    string
		args    args
		want    context.ConfigPath
		wantErr bool
	}{
		{
			name: "nil graph",
			args: args{
				configgraph:   nil,
				newconfigpath: "test",
			},
			wantErr: true,
		},
		{
			name: "null new config path",
			args: args{
				configgraph:   relPathMain.graph(),
				newconfigpath: "",
			},
			wantErr: true,
		},
		{
			name: "config paths are relative path in graph's main config and input",
			args: args{
				configgraph:   relPathMain.graph(),
				newconfigpath: "test.config",
			},
			wantErr: true,
		},
		{
			name: "absolut path",
			args: args{
				configgraph:   testMain.graph(),
				newconfigpath: "D:\\test\\test.conf",
			},
			want: absCP,
		},
		{
			name: "relative path",
			args: args{
				configgraph:   testMain.graph(),
				newconfigpath: "relative.conf",
			},
			want: relCP,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfigPath(tt.args.configgraph, tt.args.newconfigpath)
			if (err != nil) != tt.wantErr {
				t.Errorf("newConfigPath() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConfigPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_registerConfigBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerConfigBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerConfigBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
