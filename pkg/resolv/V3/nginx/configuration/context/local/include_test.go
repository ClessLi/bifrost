package local

import (
	"reflect"
	"sync"
	"testing"

	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/marmotedu/errors"
)

func TestInclude_Child(t *testing.T) {
	type fields struct {
		ContextValue  string
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
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Child(tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Child() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestInclude_ChildConfig(t *testing.T) {
//	testMain, err := NewMain("C:\\test\\nginx.conf")
//	if err != nil {
//		t.Fatal(err)
//	}
//	type fields struct {
//		enabled       bool
//		ContextValue  string
//		fatherContext context.Context
//		loadLocker    *sync.RWMutex
//	}
//	type args struct {
//		fullpath string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *Config
//		wantErr bool
//	}{
//		{
//			name: "normal test",
//			fields: fields{
//				enabled:       true,
//				ContextValue:  testMain.Value(),
//				fatherContext: testMain.MainConfig(),
//				loadLocker:    new(sync.RWMutex),
//			},
//			args:    args{fullpath: testMain.Value()},
//			want:    testMain.MainConfig(),
//			wantErr: false,
//		},
//		{
//			name: "config not found",
//			fields: fields{
//				enabled:       true,
//				ContextValue:  testMain.Value(),
//				fatherContext: context.NullContext(),
//				loadLocker:    new(sync.RWMutex),
//			},
//			args:    args{fullpath: testMain.Value()},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			i := &Include{
//				enabled:       tt.fields.enabled,
//				ContextValue:  tt.fields.ContextValue,
//				fatherContext: tt.fields.fatherContext,
//				loadLocker:    tt.fields.loadLocker,
//			}
//			got, err := i.ChildConfig(tt.args.fullpath)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ChildConfig() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ChildConfig() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestInclude_Clone(t *testing.T) {
	type fields struct {
		enabled       bool
		ContextValue  string
		fatherContext context.Context
		loadLocker    *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name:   "normal test",
			fields: fields{ContextValue: "/test/*.conf", loadLocker: new(sync.RWMutex)},
			want:   NewContext(context_type.TypeInclude, "/test/*.conf").Disable(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				fatherContext: tt.fields.fatherContext,
				loadLocker:    tt.fields.loadLocker,
			}
			if got := i.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_ConfigLines(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testConfig := NewContext(context_type.TypeConfig, "conf.d/test.conf").
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
		).(*Config)
	err = testMain.AddConfig(testConfig)
	if err != nil {
		t.Fatal(err)
	}
	errTestMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	errConfigPath, err := newConfigPath(errTestMain.graph(), "conf.d/err.conf")
	if err != nil {
		t.Fatal(err)
	}
	err = errTestMain.graph().(*configGraph).graph.AddVertex(&Config{
		BasicContext: BasicContext{
			ContextType:    context_type.TypeConfig,
			ContextValue:   "/test/err.conf",
			Children:       []context.Context{context.NullContext()},
			father:         context.NullContext(),
			self:           context.NullContext(),
			headStringFunc: nullHeadString,
			tailStringFunc: nullTailString,
		},
		ConfigPath: errConfigPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		enabled       bool
		ContextValue  string
		fatherContext context.Context
		loadLocker    *sync.RWMutex
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
			fields: fields{enabled: true, ContextValue: "C:\\test\\conf.d\\*.conf", loadLocker: new(sync.RWMutex)},
			args:   args{isDumping: true},
			want: []string{
				"include C:\\test\\conf.d\\*.conf;",
			},
			wantErr: false,
		},
		{
			name:   "for dumping with disabled include context",
			fields: fields{enabled: false, ContextValue: "C:\\test\\conf.d\\*.conf", loadLocker: new(sync.RWMutex)},
			args:   args{isDumping: true},
			want: []string{
				"# include C:\\test\\conf.d\\*.conf;",
			},
			wantErr: false,
		},
		{
			name: "error test, not for dumping",
			fields: fields{
				enabled:       true,
				ContextValue:  "C:\\test\\conf.d\\*.conf",
				fatherContext: errTestMain.MainConfig(),
				loadLocker:    new(sync.RWMutex),
			},
			args:    args{isDumping: false},
			want:    nil,
			wantErr: true,
		},
		{
			name: "normal test, not for dumping",
			fields: fields{
				enabled:       true,
				ContextValue:  "C:\\test\\conf.d\\*.conf",
				fatherContext: testMain.MainConfig(),
				loadLocker:    new(sync.RWMutex),
			},
			args: args{isDumping: false},
			want: []string{
				"# include <== C:\\test\\conf.d\\*.conf",
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
		{
			name: "normal test, not for dumping with disabled include context",
			fields: fields{
				enabled:       false,
				ContextValue:  "C:\\test\\conf.d\\*.conf",
				fatherContext: testMain.MainConfig(),
				loadLocker:    new(sync.RWMutex),
			},
			args: args{isDumping: false},
			want: []string{
				"# # include <== C:\\test\\conf.d\\*.conf",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				fatherContext: tt.fields.fatherContext,
				loadLocker:    tt.fields.loadLocker,
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
	testInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			),
		0,
	)
	type fields struct {
		enabled       bool
		ContextValue  string
		fatherContext context.Context
		loadLocker    *sync.RWMutex
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
				enabled:       testInclude.enabled,
				ContextValue:  testInclude.ContextValue,
				fatherContext: testInclude.fatherContext,
				loadLocker:    testInclude.loadLocker,
			},
			want:    testMain.MainConfig(),
			wantErr: false,
		},
		{
			name:    "father config not found",
			fields:  fields{fatherContext: context.NullContext(), loadLocker: testInclude.loadLocker},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "found a father main context",
			fields:  fields{fatherContext: testMain, loadLocker: testInclude.loadLocker},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				fatherContext: tt.fields.fatherContext,
				loadLocker:    tt.fields.loadLocker,
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
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testConfig := NewContext(context_type.TypeConfig, "conf.d/test.conf").
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
		).(*Config)
	err = testMain.AddConfig(testConfig)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		enabled       bool
		ContextValue  string
		fatherContext context.Context
		loadLocker    *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "has no child",
			fields: fields{
				enabled:       true,
				ContextValue:  "conf.d/*.conf",
				fatherContext: context.NullContext(),
				loadLocker:    new(sync.RWMutex),
			},
			want: false,
		},
		{
			name: "has children",
			fields: fields{
				enabled:       true,
				ContextValue:  "conf.d/*.conf",
				fatherContext: testMain.MainConfig(),
				loadLocker:    new(sync.RWMutex),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				fatherContext: tt.fields.fatherContext,
				loadLocker:    tt.fields.loadLocker,
			}
			if got := i.HasChild(); got != tt.want {
				t.Errorf("HasChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Insert(t *testing.T) {
	type fields struct {
		enabled       bool
		ContextValue  string
		fatherContext context.Context
		loadLocker    *sync.RWMutex
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
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				fatherContext: tt.fields.fatherContext,
				loadLocker:    tt.fields.loadLocker,
			}
			if got := i.Insert(tt.args.ctx, tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestInclude_InsertConfig(t *testing.T) {
//	testMain, err := NewMain("C:\\test\\nginx.conf")
//	if err != nil {
//		t.Fatal(err)
//	}
//	testInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").(*Include)
//	disabledInclude := NewContext(context_type.TypeInclude, "*.conf").Disable().(*Include)
//	err = testMain.Insert(
//		NewContext(context_type.TypeHttp, "").
//			Insert(
//				testInclude,
//				0,
//			).
//			Insert(
//				disabledInclude,
//				1,
//			),
//		0,
//	).Error()
//	if err != nil {
//		t.Fatal(err)
//	}
//	relConfig := NewContext(context_type.TypeConfig, "conf.d/relative.conf").(*Config)
//	relConfig.ConfigPath, _ = newConfigPath(testMain, relConfig.ContextValue)
//	absConfig := NewContext(context_type.TypeConfig, "C:\\test\\conf.d\\absolut.conf").(*Config)
//	absConfig.ConfigPath, _ = newConfigPath(testMain, absConfig.ContextValue)
//	notMatchConfig := NewContext(context_type.TypeConfig, "test\\conf.d\\test.conf").(*Config)
//	notMatchConfig.ConfigPath, _ = newConfigPath(testMain, notMatchConfig.ContextValue)
//
//	unboundedConfig := NewContext(context_type.TypeConfig, "unbound.conf")
//	unboundedInclude := NewContext(context_type.TypeInclude, "unbound.*.conf").(*Include)
//	unboundedConfig.Insert(unboundedInclude, 0)
//	type fields struct {
//		enabled       bool
//		ContextValue  string
//		fatherContext context.Context
//		loadLocker    *sync.RWMutex
//	}
//	type args struct {
//		configs []*Config
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "nil configs",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			wantErr: true,
//		},
//		{
//			name: "include a nil config",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			args:    args{configs: []*Config{nil}},
//			wantErr: true,
//		},
//		{
//			name:    "Include context has no father config",
//			fields:  fields{fatherContext: context.NullContext(), loadLocker: new(sync.RWMutex)},
//			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
//			wantErr: true,
//		},
//		{
//			name: "Include context has no father main",
//			fields: fields{
//				enabled:       unboundedInclude.enabled,
//				ContextValue:  unboundedInclude.ContextValue,
//				fatherContext: unboundedInclude.fatherContext,
//				loadLocker:    unboundedInclude.loadLocker,
//			},
//			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
//			wantErr: true,
//		},
//		{
//			name: "the included config has no config path",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "conf.d\\test.conf").(*Config)}},
//			wantErr: false,
//		},
//		{
//			name: "the included config has no config path, and not matched up",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "conf.d\\test\\test.conf").(*Config)}},
//			wantErr: true,
//		},
//		{
//			name: "has not matched config",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			args: args{configs: []*Config{
//				relConfig,
//				absConfig,
//				notMatchConfig,
//			}},
//			wantErr: true,
//		},
//		{
//			name: "normal test",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			args: args{configs: []*Config{
//				relConfig,
//				absConfig,
//			}},
//			wantErr: false,
//		},
//		{
//			name: "test disabled include insert configs",
//			fields: fields{
//				enabled:       disabledInclude.enabled,
//				ContextValue:  disabledInclude.ContextValue,
//				fatherContext: disabledInclude.fatherContext,
//				loadLocker:    disabledInclude.loadLocker,
//			},
//			args: args{configs: []*Config{
//				unboundedConfig.(*Config),
//			}},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			i := &Include{
//				enabled:       tt.fields.enabled,
//				ContextValue:  tt.fields.ContextValue,
//				fatherContext: tt.fields.fatherContext,
//				loadLocker:    tt.fields.loadLocker,
//			}
//			err := i.InsertConfig(tt.args.configs...)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("InsertConfig() error = %+v, wantErr %+v", err, tt.wantErr)
//			}
//			if err != nil {
//				t.Log(fmt.Sprintf("%+v", err))
//			}
//		})
//	}
//}

func TestInclude_Len(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").(*Include)
	disabledInclude := NewContext(context_type.TypeInclude, "*.conf").Disable().(*Include)
	noConfigInclude := NewContext(context_type.TypeInclude, "null.conf").(*Include)
	err = testMain.AddConfig(NewContext(context_type.TypeConfig, "conf.d/1.conf").(*Config))
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(NewContext(context_type.TypeConfig, "conf.d/2.conf").(*Config))
	err = testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			).
			Insert(
				disabledInclude,
				1,
			).
			Insert(
				noConfigInclude,
				2,
			),
		0,
	).Error()
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		enabled       bool
		ContextValue  string
		fatherContext context.Context
		loadLocker    *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "disabled include",
			fields: fields{
				enabled:       disabledInclude.enabled,
				ContextValue:  disabledInclude.ContextValue,
				fatherContext: disabledInclude.fatherContext,
				loadLocker:    disabledInclude.loadLocker,
			},
			want: 1,
		},
		{
			name: "null configs",
			fields: fields{
				enabled:       noConfigInclude.enabled,
				ContextValue:  noConfigInclude.ContextValue,
				fatherContext: noConfigInclude.fatherContext,
				loadLocker:    noConfigInclude.loadLocker,
			},
			want: 0,
		},
		{
			name: "has some configs",
			fields: fields{
				enabled:       testInclude.enabled,
				ContextValue:  testInclude.ContextValue,
				fatherContext: testInclude.fatherContext,
				loadLocker:    testInclude.loadLocker,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				fatherContext: tt.fields.fatherContext,
				loadLocker:    tt.fields.loadLocker,
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
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Modify(tt.args.ctx, tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Modify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_QueryAllByKeyWords(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").(*Include)
	disabledInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").Disable().(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			).
			Insert(
				disabledInclude,
				1,
			),
		0,
	)
	aConfig := NewContext(context_type.TypeConfig, "conf.d/a.conf").(*Config)
	aConfig.ConfigPath, _ = newConfigPath(testMain, aConfig.ContextValue)
	bConfig := NewContext(context_type.TypeConfig, "C:\\test\\conf.d\\b.conf").(*Config)
	bConfig.ConfigPath, _ = newConfigPath(testMain, bConfig.ContextValue)
	err = testMain.AddConfig(aConfig)
	if err != nil {
		t.Fatal(err)
	}

	err = testMain.AddConfig(bConfig)
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
		enabled       bool
		ContextValue  string
		fatherContext context.Context
		loadLocker    *sync.RWMutex
	}
	type args struct {
		kw context.KeyWords
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   context.PosSet
	}{
		{
			name: "normal test",
			fields: fields{
				enabled:       testInclude.enabled,
				ContextValue:  testInclude.ContextValue,
				fatherContext: testInclude.fatherContext,
				loadLocker:    testInclude.loadLocker,
			},
			args: args{kw: context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue("test")},
			want: context.NewPosSet().Append(
				context.SetPos(aFather, 0),
				context.SetPos(bFather, 1),
			),
		},
		{
			name: "test disabled include context",
			fields: fields{
				enabled:       disabledInclude.enabled,
				ContextValue:  disabledInclude.ContextValue,
				fatherContext: disabledInclude.fatherContext,
				loadLocker:    disabledInclude.loadLocker,
			},
			args: args{kw: context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue("test")},
			want: context.NewPosSet(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				fatherContext: tt.fields.fatherContext,
				loadLocker:    tt.fields.loadLocker,
			}
			if got := i.ChildrenPosSet().QueryAll(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
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
	testInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").(*Include)
	disabledInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").Disable().(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			).
			Insert(
				disabledInclude,
				1,
			),
		0,
	)
	aConfig := NewContext(context_type.TypeConfig, "conf.d/a.conf").(*Config)
	aConfig.ConfigPath, _ = newConfigPath(testMain, aConfig.ContextValue)
	bConfig := NewContext(context_type.TypeConfig, "C:\\test\\conf.d\\b.conf").(*Config)
	bConfig.ConfigPath, _ = newConfigPath(testMain, bConfig.ContextValue)
	err = testMain.AddConfig(aConfig)
	if err != nil {
		t.Fatal(err)
	}

	err = testMain.AddConfig(bConfig)
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
		enabled       bool
		ContextValue  string
		fatherContext context.Context
		loadLocker    *sync.RWMutex
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
				enabled:       testInclude.enabled,
				ContextValue:  testInclude.ContextValue,
				fatherContext: testInclude.fatherContext,
				loadLocker:    testInclude.loadLocker,
			},
			args: args{kw: context.NewKeyWords(context_type.TypeLocation).SetStringMatchingValue("test")},
			want: context.SetPos(aFather, 0),
		},
		{
			name: "test disabled include context",
			fields: fields{
				enabled:       disabledInclude.enabled,
				ContextValue:  disabledInclude.ContextValue,
				fatherContext: disabledInclude.fatherContext,
				loadLocker:    disabledInclude.loadLocker,
			},
			args: args{kw: context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue("test")},
			want: context.NotFoundPos(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				fatherContext: tt.fields.fatherContext,
				loadLocker:    tt.fields.loadLocker,
			}
			if got := i.ChildrenPosSet().QueryOne(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryByKeyWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Remove(t *testing.T) {
	type fields struct {
		ContextValue  string
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
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Remove(tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestInclude_RemoveConfig(t *testing.T) {
//	testMain, err := NewMain("C:\\test\\nginx.conf")
//	if err != nil {
//		t.Fatal(err)
//	}
//	testInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").(*Include)
//	disabledInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").Disable().(*Include)
//	testMain.Insert(
//		NewContext(context_type.TypeHttp, "").
//			Insert(
//				testInclude,
//				0,
//			).
//			Insert(
//				disabledInclude,
//				1,
//			),
//		0,
//	)
//	relConfig := NewContext(context_type.TypeConfig, "conf.d/relative.conf").(*Config)
//	relConfig.ConfigPath, _ = newConfigPath(testMain, relConfig.ContextValue)
//	absConfig := NewContext(context_type.TypeConfig, "C:\\test\\conf.d\\absolut.conf").(*Config)
//	absConfig.ConfigPath, _ = newConfigPath(testMain, absConfig.ContextValue)
//	notMatchConfig := NewContext(context_type.TypeConfig, "conf.d\\test\\test.conf").(*Config)
//	notMatchConfig.ConfigPath, _ = newConfigPath(testMain, notMatchConfig.ContextValue)
//	err = testMain.(*Main).ConfigGraph.(*configGraph).addVertex(notMatchConfig)
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = testMain.AddConfig(relConfig)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = testMain.AddConfig(absConfig)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	onlyAddedIntoGraphConfig := NewContext(context_type.TypeConfig, "C:\\test\\conf.d\\unbound.conf").(*Config)
//	onlyAddedIntoGraphConfig.ConfigPath, _ = newConfigPath(testMain, onlyAddedIntoGraphConfig.ContextValue)
//	err = testMain.(*Main).ConfigGraph.(*configGraph).addVertex(onlyAddedIntoGraphConfig)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	type fields struct {
//		enabled       bool
//		ContextValue  string
//		fatherContext context.Context
//		loadLocker    *sync.RWMutex
//	}
//	type args struct {
//		configs []*Config
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			name:    "nil configs",
//			wantErr: true,
//		},
//		{
//			name: "removed a nil config",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			args:    args{configs: []*Config{nil}},
//			wantErr: false,
//		},
//		{
//			name:    "Include context has no father config",
//			fields:  fields{fatherContext: context.NullContext(), loadLocker: new(sync.RWMutex)},
//			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
//			wantErr: true,
//		},
//		{
//			name: "removed config has no config path",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "test.conf").(*Config)}},
//			wantErr: false,
//		},
//		{
//			name: "removed config is not included in",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			args:    args{configs: []*Config{notMatchConfig}},
//			wantErr: false,
//		},
//		{
//			name: "removed config is not included into father config",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			args:    args{configs: []*Config{onlyAddedIntoGraphConfig}},
//			wantErr: true,
//		},
//		{
//			name: "normal test",
//			fields: fields{
//				enabled:       testInclude.enabled,
//				ContextValue:  testInclude.ContextValue,
//				fatherContext: testInclude.fatherContext,
//				loadLocker:    testInclude.loadLocker,
//			},
//			args: args{configs: []*Config{
//				relConfig,
//				absConfig,
//			}},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			i := &Include{
//				enabled:       tt.fields.enabled,
//				ContextValue:  tt.fields.ContextValue,
//				fatherContext: tt.fields.fatherContext,
//				loadLocker:    tt.fields.loadLocker,
//			}
//			if err := i.RemoveConfig(tt.args.configs...); (err != nil) != tt.wantErr {
//				t.Errorf("RemoveConfig() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestInclude_SetFather(t *testing.T) {
	type fields struct {
		enabled       bool
		ContextValue  string
		loadLocker    *sync.RWMutex
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
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				loadLocker:    tt.fields.loadLocker,
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
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
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

//func Test_registerIncludeParseFunc(t *testing.T) {
//	tests := []struct {
//		name    string
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := registerIncludeParseFunc(); (err != nil) != tt.wantErr {
//				t.Errorf("registerIncludeParseFunc() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestInclude_MarshalJSON(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").(*Include)
	disabledInclude := NewContext(context_type.TypeInclude, "conf.d/*.conf").Disable().(*Include)
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				testInclude,
				0,
			).
			Insert(
				disabledInclude,
				1,
			),
		0,
	)
	relConfig := NewContext(context_type.TypeConfig, "conf.d/relative.conf").(*Config)
	relConfig.ConfigPath, _ = newConfigPath(testMain, relConfig.ContextValue)
	absConfig := NewContext(context_type.TypeConfig, "C:\\test\\conf.d\\absolut.conf").(*Config)
	absConfig.ConfigPath, _ = newConfigPath(testMain, absConfig.ContextValue)
	notMatchConfig := NewContext(context_type.TypeConfig, "conf.d\\test\\test.conf").(*Config)
	notMatchConfig.ConfigPath, _ = newConfigPath(testMain, notMatchConfig.ContextValue)
	err = testMain.AddConfig(relConfig)
	if err != nil {
		t.Fatal(err)
	}

	err = testMain.AddConfig(absConfig)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		enabled       bool
		ContextValue  string
		fatherContext context.Context
		loadLocker    *sync.RWMutex
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
				enabled:       testInclude.enabled,
				ContextValue:  testInclude.ContextValue,
				fatherContext: testInclude.fatherContext,
				loadLocker:    testInclude.loadLocker,
			},
			want:    []byte(`{"enabled":true,"context-type":"include","value":"conf.d/*.conf"}`),
			wantErr: false,
		},
		{
			name: "test disabled include context",
			fields: fields{
				enabled:       disabledInclude.enabled,
				ContextValue:  disabledInclude.ContextValue,
				fatherContext: disabledInclude.fatherContext,
				loadLocker:    disabledInclude.loadLocker,
			},
			want:    []byte(`{"context-type":"include","value":"conf.d/*.conf"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				fatherContext: tt.fields.fatherContext,
				loadLocker:    tt.fields.loadLocker,
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

func TestInclude_IsEnabled(t *testing.T) {
	type fields struct {
		enabled       bool
		ContextValue  string
		loadLocker    *sync.RWMutex
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Include{
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				loadLocker:    tt.fields.loadLocker,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.IsEnabled(); got != tt.want {
				t.Errorf("IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Enable(t *testing.T) {
	type fields struct {
		enabled       bool
		ContextValue  string
		loadLocker    *sync.RWMutex
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
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				loadLocker:    tt.fields.loadLocker,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Enable(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Enable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Disable(t *testing.T) {
	type fields struct {
		enabled       bool
		ContextValue  string
		loadLocker    *sync.RWMutex
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
				enabled:       tt.fields.enabled,
				ContextValue:  tt.fields.ContextValue,
				loadLocker:    tt.fields.loadLocker,
				fatherContext: tt.fields.fatherContext,
			}
			if got := i.Disable(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Disable() = %v, want %v", got, tt.want)
			}
		})
	}
}
