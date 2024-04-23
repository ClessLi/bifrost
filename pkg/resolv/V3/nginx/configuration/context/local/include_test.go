package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"reflect"
	"testing"
)

func TestInclude_Child(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   context.Context
	}{
		{
			name: "error context",
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "include cannot get child config by index")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Child(tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Child() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_ChildConfig(t *testing.T) {
	testConfig := NewContext(context_type.TypeConfig, "C:\\test\\test.conf").(*Config)
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		fullpath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name:    "normal test",
			fields:  fields{Configs: map[string]*Config{testConfig.Value(): testConfig}},
			args:    args{fullpath: testConfig.Value()},
			want:    testConfig,
			wantErr: false,
		},
		{
			name:    "config not found",
			fields:  fields{Configs: make(map[string]*Config)},
			args:    args{fullpath: testConfig.Value()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			got, err := i.ChildConfig(tt.args.fullpath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChildConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChildConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Clone(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name:   "normal test",
			fields: fields{ContextValue: "/test/*.conf"},
			want:   NewContext(context_type.TypeInclude, "/test/*.conf"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_ConfigLines(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
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
			name:   "for dumping",
			fields: fields{ContextValue: "/test/*.conf"},
			args:   args{isDumping: true},
			want: []string{
				"include /test/*.conf;",
			},
			wantErr: false,
		},
		{
			name: "error test, not for dumping",
			fields: fields{
				ContextValue: "/test/*.conf",
				Configs: map[string]*Config{
					"/test/err.conf": {BasicContext: BasicContext{
						ContextType:    "",
						ContextValue:   "",
						Children:       []context.Context{context.NullContext()},
						father:         context.NullContext(),
						self:           context.NullContext(),
						headStringFunc: nullHeadString,
						tailStringFunc: nullTailString,
					}},
				},
			},
			args:    args{isDumping: false},
			want:    nil,
			wantErr: true,
		},
		{
			name: "normal test, not for dumping",
			fields: fields{
				ContextValue: "/test/*.conf",
				Configs: map[string]*Config{
					"/test/test.conf": NewContext(context_type.TypeConfig, "/test/test.conf").
						Insert(
							NewContext(context_type.TypeHttp, "").
								Insert(
									NewContext(context_type.TypeServer, "").
										Insert(
											NewContext(context_type.TypeLocation, "~ /test").
												Insert(
													NewContext(context_type.TypeDirective, "proxy_pass https://www.baidu.com"),
													0,
												),
											0,
										).
										Insert(
											NewContext(context_type.TypeDirective, "server_name testserver.com"),
											0,
										),
									0,
								),
							0,
						).(*Config),
				},
			},
			args: args{isDumping: false},
			want: []string{
				"# include <== /test/*.conf",
				"http {",
				"    server {",
				"        server_name testserver.com;",
				"        location ~ /test {",
				"            proxy_pass https://www.baidu.com;",
				"        }",
				"    }",
				"}",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			got, err := i.ConfigLines(tt.args.isDumping)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigLines() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Error(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "return nil",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if err := i.Error(); (err != nil) != tt.wantErr {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInclude_Father(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Father(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Father() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_FatherConfig(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testInclude := NewContext(context_type.TypeInclude, "*.conf").(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			),
		0,
	)
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Config
		wantErr bool
	}{
		{
			name: "normal test",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			want:    testMain.MainConfig(),
			wantErr: false,
		},
		{
			name:    "father config not found",
			fields:  fields{fatherContext: context.NullContext()},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "found a father main context",
			fields:  fields{fatherContext: testMain},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			got, err := i.FatherConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("FatherConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FatherConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_HasChild(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "has no child",
			want: false,
		},
		{
			name:   "has children",
			fields: fields{Configs: map[string]*Config{"/test.conf": NewContext(context_type.TypeConfig, "/test.conf").(*Config)}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.HasChild(); got != tt.want {
				t.Errorf("HasChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Insert(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		ctx context.Context
		idx int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   context.Context
	}{
		{
			name: "return error context",
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "include cannot insert by index")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Insert(tt.args.ctx, tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_InsertConfig(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testInclude := NewContext(context_type.TypeInclude, "*.conf").(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			),
		0,
	)
	relConfig := NewContext(context_type.TypeConfig, "relative.conf").(*Config)
	relConfig.ConfigPath, _ = newConfigPath(testMain, relConfig.ContextValue)
	absConfig := NewContext(context_type.TypeConfig, "C:\\test\\absolut.conf").(*Config)
	absConfig.ConfigPath, _ = newConfigPath(testMain, absConfig.ContextValue)
	notMatchConfig := NewContext(context_type.TypeConfig, "test\\test.conf").(*Config)
	notMatchConfig.ConfigPath, _ = newConfigPath(testMain, notMatchConfig.ContextValue)

	unboundedConfig := NewContext(context_type.TypeConfig, "unbound.conf")
	unboundedInclude := NewContext(context_type.TypeInclude, "unbound.*.conf").(*Include)
	unboundedConfig.Insert(unboundedInclude, 0)
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		configs []*Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "nil configs",
			wantErr: true,
		},
		{
			name: "include a nil config",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args:    args{configs: []*Config{nil}},
			wantErr: true,
		},
		{
			name:    "Include context has no father config",
			fields:  fields{fatherContext: context.NullContext()},
			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
			wantErr: true,
		},
		{
			name: "Include context has no father main",
			fields: fields{
				ContextValue:  unboundedInclude.ContextValue,
				Configs:       unboundedInclude.Configs,
				fatherContext: unboundedInclude.fatherContext,
			},
			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
			wantErr: true,
		},
		{
			name: "the included config has no config path",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
			wantErr: false,
		},
		{
			name: "the included config has no config path, and not matched up",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test\\test.conf").(*Config)}},
			wantErr: true,
		},
		{
			name: "has not matched config",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args: args{configs: []*Config{
				relConfig,
				absConfig,
				notMatchConfig,
			}},
			wantErr: true,
		},
		{
			name: "normal test",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args: args{configs: []*Config{
				relConfig,
				absConfig,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if err := i.InsertConfig(tt.args.configs...); (err != nil) != tt.wantErr {
				t.Errorf("InsertConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInclude_Len(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "nil configs map",
			want: 0,
		},
		{
			name:   "null configs",
			fields: fields{Configs: make(map[string]*Config)},
			want:   0,
		},
		{
			name: "has some configs",
			fields: fields{Configs: map[string]*Config{
				"1.conf": NewContext(context_type.TypeConfig, "1.conf").(*Config),
				"2.conf": NewContext(context_type.TypeConfig, "2.conf").(*Config),
			}},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Modify(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		ctx context.Context
		idx int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   context.Context
	}{
		{
			name: "return error context",
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "include cannot modify by index")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Modify(tt.args.ctx, tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Modify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_ModifyConfig(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testInclude := NewContext(context_type.TypeInclude, "*.conf").(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			),
		0,
	)
	relConfig := NewContext(context_type.TypeConfig, "relative.conf").(*Config)
	relConfig.ConfigPath, _ = newConfigPath(testMain, relConfig.ContextValue)
	absConfig := NewContext(context_type.TypeConfig, "C:\\test\\absolut.conf").(*Config)
	absConfig.ConfigPath, _ = newConfigPath(testMain, absConfig.ContextValue)
	notMatchConfig := NewContext(context_type.TypeConfig, "test\\test.conf").(*Config)
	notMatchConfig.ConfigPath, _ = newConfigPath(testMain, notMatchConfig.ContextValue)
	err = testInclude.InsertConfig(relConfig, absConfig)
	if err != nil {
		t.Fatal(err)
	}
	modifiedRelConfig := relConfig.Clone().(*Config)
	modifiedRelConfig.ConfigPath, _ = newConfigPath(testMain, modifiedRelConfig.Value())
	modifiedRelConfig.Insert(NewContext(context_type.TypeDirective, "keepalive_timeout 300s"), 0)
	modifiedAbsConfig := absConfig.Clone().(*Config)
	modifiedAbsConfig.ConfigPath, _ = newConfigPath(testMain, modifiedAbsConfig.Value())
	modifiedAbsConfig.Insert(NewContext(context_type.TypeDirective, "proxy_http_version 1.1"), 0)

	unboundedConfig := NewContext(context_type.TypeConfig, "unbound.conf")
	unboundedInclude := NewContext(context_type.TypeInclude, "unbound.*.conf").(*Include)
	unboundedConfig.Insert(unboundedInclude, 0)
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		configs []*Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "nil configs",
			wantErr: true,
		},
		{
			name: "a nil in configs",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args:    args{configs: []*Config{nil}},
			wantErr: true,
		},
		{
			name: "Include context has no father main",
			fields: fields{
				ContextValue:  unboundedInclude.ContextValue,
				Configs:       unboundedInclude.Configs,
				fatherContext: unboundedInclude.fatherContext,
			},
			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
			wantErr: true,
		},
		{
			name: "config has no config path",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
			wantErr: false,
		},
		{
			name: "the modified config has no config path, and not matched up",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test\\test.conf").(*Config)}},
			wantErr: true,
		},
		{
			name: "has no matched config",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args: args{configs: []*Config{
				relConfig,
				absConfig,
				notMatchConfig,
			}},
			wantErr: true,
		},
		{
			name: "normal test",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args: args{configs: []*Config{
				modifiedRelConfig,
				modifiedAbsConfig,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if err := i.ModifyConfig(tt.args.configs...); (err != nil) != tt.wantErr {
				t.Errorf("ModifyConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInclude_QueryAllByKeyWords(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testInclude := NewContext(context_type.TypeInclude, "*.conf").(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			),
		0,
	)
	aConfig := NewContext(context_type.TypeConfig, "a.conf").(*Config)
	aConfig.ConfigPath, _ = newConfigPath(testMain, aConfig.ContextValue)
	bConfig := NewContext(context_type.TypeConfig, "C:\\test\\b.conf").(*Config)
	bConfig.ConfigPath, _ = newConfigPath(testMain, bConfig.ContextValue)
	err = testInclude.InsertConfig(aConfig, bConfig)
	if err != nil {
		t.Fatal(err)
	}

	aFather := NewContext(context_type.TypeServer, "")
	aConfig.Insert(
		aFather.
			Insert(NewContext(context_type.TypeLocation, "~ /test"), 0),
		0,
	)
	bFather := NewContext(context_type.TypeServer, "")
	bConfig.Insert(
		bFather.
			Insert(NewContext(context_type.TypeLocation, "/text"), 0).
			Insert(NewContext(context_type.TypeLocation, "/test1"), 1),
		0,
	)
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		kw context.KeyWords
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []context.Pos
	}{
		{
			name: "normal test",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args: args{kw: context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue("test")},
			want: []context.Pos{
				context.SetPos(aFather, 0),
				context.SetPos(bFather, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.QueryAllByKeyWords(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryAllByKeyWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_QueryByKeyWords(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testInclude := NewContext(context_type.TypeInclude, "*.conf").(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			),
		0,
	)
	aConfig := NewContext(context_type.TypeConfig, "a.conf").(*Config)
	aConfig.ConfigPath, _ = newConfigPath(testMain, aConfig.ContextValue)
	bConfig := NewContext(context_type.TypeConfig, "C:\\test\\b.conf").(*Config)
	bConfig.ConfigPath, _ = newConfigPath(testMain, bConfig.ContextValue)
	err = testInclude.InsertConfig(aConfig, bConfig)
	if err != nil {
		t.Fatal(err)
	}

	aFather := NewContext(context_type.TypeServer, "")
	aConfig.Insert(
		aFather.
			Insert(NewContext(context_type.TypeLocation, "~ /test"), 0),
		0,
	)
	bFather := NewContext(context_type.TypeServer, "")
	bConfig.Insert(
		bFather.
			Insert(NewContext(context_type.TypeLocation, "/text"), 0).
			Insert(NewContext(context_type.TypeLocation, "/test1"), 1),
		0,
	)
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		kw context.KeyWords
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   context.Pos
	}{
		{
			name: "normal test",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args: args{kw: context.NewKeyWords(context_type.TypeLocation).SetStringMatchingValue("test")},
			want: context.SetPos(aFather, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.QueryByKeyWords(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryByKeyWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Remove(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   context.Context
	}{
		{
			name: "return error context",
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "include cannot remove by index")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Remove(tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_RemoveConfig(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testInclude := NewContext(context_type.TypeInclude, "*.conf").(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			),
		0,
	)
	relConfig := NewContext(context_type.TypeConfig, "relative.conf").(*Config)
	relConfig.ConfigPath, _ = newConfigPath(testMain, relConfig.ContextValue)
	absConfig := NewContext(context_type.TypeConfig, "C:\\test\\absolut.conf").(*Config)
	absConfig.ConfigPath, _ = newConfigPath(testMain, absConfig.ContextValue)
	notMatchConfig := NewContext(context_type.TypeConfig, "test\\test.conf").(*Config)
	notMatchConfig.ConfigPath, _ = newConfigPath(testMain, notMatchConfig.ContextValue)
	err = testInclude.InsertConfig(relConfig, absConfig)
	if err != nil {
		t.Fatal(err)
	}

	onlyAddedIntoGraphConfig := NewContext(context_type.TypeConfig, "C:\\test\\unbound.conf").(*Config)
	onlyAddedIntoGraphConfig.ConfigPath, _ = newConfigPath(testMain, onlyAddedIntoGraphConfig.ContextValue)
	err = testMain.AddConfig(onlyAddedIntoGraphConfig)
	if err != nil {
		t.Fatal(err)
	}
	testInclude.Configs[configHash(onlyAddedIntoGraphConfig)] = onlyAddedIntoGraphConfig

	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		configs []*Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "nil configs",
			wantErr: true,
		},
		{
			name: "removed a nil config",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args:    args{configs: []*Config{nil}},
			wantErr: false,
		},
		{
			name:    "Include context has no father config",
			fields:  fields{fatherContext: context.NullContext()},
			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
			wantErr: true,
		},
		{
			name: "removed config has no config path",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
			wantErr: false,
		},
		{
			name: "removed config is not included in",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args:    args{configs: []*Config{notMatchConfig}},
			wantErr: false,
		},
		{
			name: "removed config is not included into father config",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args:    args{configs: []*Config{onlyAddedIntoGraphConfig}},
			wantErr: true,
		},
		{
			name: "normal test",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			args: args{configs: []*Config{
				relConfig,
				absConfig,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if err := i.RemoveConfig(tt.args.configs...); (err != nil) != tt.wantErr {
				t.Errorf("RemoveConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInclude_SetFather(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if err := i.SetFather(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SetFather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInclude_SetValue(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
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
			name:    "return error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if err := i.SetValue(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInclude_Type(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   context_type.ContextType
	}{
		{
			name: "return include",
			want: context_type.TypeInclude,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Value(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
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
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_matchConfigPath(t *testing.T) {
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "malformed include patten value",
			fields:  fields{ContextValue: "][a"},
			args:    args{path: "test.conf"},
			wantErr: true,
		},
		{
			name:    "path is not matched",
			fields:  fields{ContextValue: "C:\\aaa\\*.conf"},
			args:    args{path: "test.conf"},
			wantErr: true,
		},
		{
			name:    "match relative path",
			fields:  fields{ContextValue: "..\\aaa\\*.conf"},
			args:    args{path: "..\\aaa\\test.conf"},
			wantErr: false,
		},
		{
			name:    "match absolut path",
			fields:  fields{ContextValue: "C:\\aaa\\*.conf"},
			args:    args{path: "C:\\aaa\\test.conf"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			if err := i.matchConfigPath(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("matchConfigPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerIncludeBuild(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerIncludeBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerIncludeBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerIncludeParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerIncludeParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerIncludeParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInclude_MarshalJSON(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testInclude := NewContext(context_type.TypeInclude, "*.conf").(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			),
		0,
	)
	relConfig := NewContext(context_type.TypeConfig, "relative.conf").(*Config)
	relConfig.ConfigPath, _ = newConfigPath(testMain, relConfig.ContextValue)
	absConfig := NewContext(context_type.TypeConfig, "C:\\test\\absolut.conf").(*Config)
	absConfig.ConfigPath, _ = newConfigPath(testMain, absConfig.ContextValue)
	notMatchConfig := NewContext(context_type.TypeConfig, "test\\test.conf").(*Config)
	notMatchConfig.ConfigPath, _ = newConfigPath(testMain, notMatchConfig.ContextValue)
	err = testInclude.InsertConfig(relConfig, absConfig)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		ContextValue  string
		Configs       map[string]*Config
		fatherContext context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "normal test",
			fields: fields{
				ContextValue:  testInclude.ContextValue,
				Configs:       testInclude.Configs,
				fatherContext: testInclude.fatherContext,
			},
			want:    []byte(`{"context-type":"include","value":"*.conf","params":["relative.conf","C:\\test\\absolut.conf"]}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				ContextValue:  tt.fields.ContextValue,
				Configs:       tt.fields.Configs,
				fatherContext: tt.fields.fatherContext,
			}
			got, err := i.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}
