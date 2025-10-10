package local

import (
	"reflect"
	"testing"

	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/marmotedu/errors"
)

func TestMain_Child(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if got := m.Child(tt.args.idx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Child() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain_Clone(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if got := m.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain_ConfigLines(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			got, err := m.ConfigLines(tt.args.isDumping)
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

func TestMain_Error(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if err := m.Error(); (err != nil) != tt.wantErr {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMain_Father(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if got := m.Father(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Father() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain_HasChild(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if got := m.HasChild(); got != tt.want {
				t.Errorf("HasChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain_Insert(t *testing.T) {
	testMain, err := NewMain("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		idx int
	}
	tests := []struct {
		name   string
		fields MainContext
		args   args
		want   context.Context
	}{
		{
			name:   "error context result",
			fields: testMain,
			args: args{
				ctx: nil,
				idx: 0,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert nil")),
		},
		{
			name:   "return main context itself",
			fields: testMain,
			args: args{
				ctx: NewContext(context_type.TypeDirective, "test"),
				idx: 0,
			},
			want: testMain,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.fields
			if got := m.Insert(tt.args.ctx, tt.args.idx); !reflect.DeepEqual(got, tt.want) {
				if got.Type() != tt.want.Type() || (got.Type() == context_type.TypeErrContext && got.Error().Error() != tt.want.Error().Error()) {
					t.Errorf("Insert() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMain_Len(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if got := m.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain_MarshalJSON(t *testing.T) {
	emptyMain, err := NewMain("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}

	// enabled/disabled contexts test data
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeInlineComment, "enabled server with enabled children configs"),
						0,
					).
					Insert(
						NewContext(context_type.TypeInclude, "conf.d/enabled.conf"),
						1,
					),
				0,
			).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeInlineComment, "enabled server with disabled children configs"),
						0,
					).
					Insert(
						NewContext(context_type.TypeInclude, "conf.d/disabled.conf"),
						1,
					),
				1,
			).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeInlineComment, "enabled server with disabled include context"),
						0,
					).
					Insert(
						NewContext(context_type.TypeInclude, "conf.d/enabled.conf").Disable(),
						1,
					),
				2,
			).
			Insert(
				NewContext(context_type.TypeServer, "").Disable().
					Insert(
						NewContext(context_type.TypeInlineComment, "disabled server with enabled children configs"),
						0,
					).
					Insert(
						NewContext(context_type.TypeInclude, "conf.d/enabled.conf"),
						1,
					),
				3,
			).
			Insert(
				NewContext(context_type.TypeServer, "").Disable().
					Insert(
						NewContext(context_type.TypeInlineComment, "disabled server with disabled children configs"),
						0,
					).
					Insert(
						NewContext(context_type.TypeInclude, "conf.d/disabled.conf"),
						1,
					),
				4,
			).
			Insert(
				NewContext(context_type.TypeServer, "").Disable().
					Insert(
						NewContext(context_type.TypeInlineComment, "disabled server with disabled include context"),
						0,
					).
					Insert(
						NewContext(context_type.TypeInclude, "conf.d/enabled.conf").Disable(),
						1,
					),
				5,
			),
		0,
	)
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "conf.d/enabled.conf").
			Insert(
				NewContext(context_type.TypeLocation, "~ /test").
					Insert(NewContext(context_type.TypeDirective, "return 200 'test'"), 0),
				0,
			).(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "conf.d/disabled.conf").Disable().
			Insert(
				NewContext(context_type.TypeComment, "disabled config"),
				0,
			).
			Insert(
				NewContext(context_type.TypeLocation, "~ /test").
					Insert(NewContext(context_type.TypeDirective, "return 404"), 0),
				1,
			).(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		ConfigGraph ConfigGraph
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty main",
			fields:  fields{ConfigGraph: emptyMain.graph()},
			want:    []byte(`{"main-config":"C:\\test\\test.conf","configs":{"C:\\test\\test.conf":{"enabled":true,"context-type":"config","value":"C:\\test\\test.conf"}}}`),
			wantErr: false,
		},
		{
			name:   "normal test",
			fields: fields{ConfigGraph: testMain.graph()},
			want: []byte(
				`{"main-config":"C:\\test\\nginx.conf","configs":{"C:\\test\\nginx.conf":{"enabled":true,"context-type":"config","value":"C:\\test\\nginx.conf","params":[{"enabled":true,"context-type":"http","params":[{"enabled":true,"context-type":"server","params":[{"context-type":"inline_comment","value":"enabled server with enabled children configs"},{"enabled":true,"context-type":"include","value":"conf.d/enabled.conf"}]},{"enabled":true,"context-type":"server","params":[{"context-type":"inline_comment","value":"enabled server with disabled children configs"},{"enabled":true,"context-type":"include","value":"conf.d/disabled.conf"}]},{"enabled":true,"context-type":"server","params":[{"context-type":"inline_comment","value":"enabled server with disabled include context"},{"context-type":"include","value":"conf.d/enabled.conf"}]},{"context-type":"server","params":[{"context-type":"inline_comment","value":"disabled server with enabled children configs"},{"enabled":true,"context-type":"include","value":"conf.d/enabled.conf"}]},{"context-type":"server","params":[{"context-type":"inline_comment","value":"disabled server with disabled children configs"},{"enabled":true,"context-type":"include","value":"conf.d/disabled.conf"}]},{"context-type":"server","params":[{"context-type":"inline_comment","value":"disabled server with disabled include context"},{"context-type":"include","value":"conf.d/enabled.conf"}]}]}]},"conf.d/disabled.conf":{"context-type":"config","value":"conf.d/disabled.conf","params":[{"context-type":"comment","value":"disabled config"},{"enabled":true,"context-type":"location","value":"~ /test","params":[{"enabled":true,"context-type":"directive","value":"return 404"}]}]},"conf.d/enabled.conf":{"enabled":true,"context-type":"config","value":"conf.d/enabled.conf","params":[{"enabled":true,"context-type":"location","value":"~ /test","params":[{"enabled":true,"context-type":"directive","value":"return 200 'test'"}]}]}}}`, //nolint:lll
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			got, err := m.MarshalJSON()
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

func TestMain_Modify(t *testing.T) { //nolint:dupl
	testMain, err := NewMain("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		idx int
	}
	tests := []struct {
		name   string
		fields MainContext
		args   args
		want   context.Context
	}{
		{
			name:   "error context result",
			fields: testMain,
			args: args{
				ctx: nil,
				idx: 0,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert nil")),
		},
		{
			name:   "return main context itself",
			fields: testMain,
			args: args{
				ctx: NewContext(context_type.TypeDirective, "test"),
				idx: 0,
			},
			want: testMain,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.fields
			if got := m.Modify(tt.args.ctx, tt.args.idx); !reflect.DeepEqual(got, tt.want) {
				if got.Type() != tt.want.Type() || (got.Type() == context_type.TypeErrContext && got.Error().Error() != tt.want.Error().Error()) {
					t.Errorf("Modify() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMain_Remove(t *testing.T) {
	testMain, err := NewMain("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name   string
		fields MainContext
		args   args
		want   context.Context
	}{
		{
			name:   "error context result",
			fields: testMain,
			args:   args{idx: -1},
			want:   context.ErrContext(errors.WithCode(code.ErrV3ContextIndexOutOfRange, "index(%d) out of range", -1)),
		},
		{
			name:   "return main itself",
			fields: testMain,
			args:   args{idx: 0},
			want:   testMain,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.fields
			if got := m.Remove(tt.args.idx); !reflect.DeepEqual(got, tt.want) {
				if got.Type() != tt.want.Type() || (got.Type() == context_type.TypeErrContext && got.Error().Error() != tt.want.Error().Error()) {
					t.Errorf("Remove() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMain_SetFather(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if err := m.SetFather(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SetFather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMain_SetValue(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if err := m.SetValue(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMain_Type(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
	}
	tests := []struct {
		name   string
		fields fields
		want   context_type.ContextType
	}{
		{
			name: "main context type",
			want: context_type.TypeMain,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if got := m.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain_Value(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if got := m.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NewMain(t *testing.T) {
	confpath, err := context.NewAbsConfigPath("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	testMainConfig := &Config{
		BasicContext: newBasicContext(context_type.TypeConfig, nullHeadString, nullTailString),
		ConfigPath:   confpath,
	}
	testMainConfig.self = testMainConfig
	testMainConfig.ContextValue = "C:\\test\\nginx.conf"
	mg, err := newConfigGraph(testMainConfig)
	if err != nil {
		t.Fatal(err)
	}
	testMain := &Main{mg}
	testMain.MainConfig().father = testMain
	type args struct {
		abspath string
	}
	tests := []struct {
		name    string
		args    args
		want    MainContext
		wantErr bool
	}{
		{
			name:    "not a absolute path",
			args:    args{abspath: "not/absolute/path"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "normal test",
			args:    args{abspath: "C:\\test\\nginx.conf"},
			want:    testMain,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMain(tt.args.abspath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMain() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			isSameFather := func(got, want *Config) bool {
				return got.father.Type() == want.father.Type() && got.father.Value() == want.father.Value()
			}
			if (got == nil) != (tt.want == nil) ||
				(got != nil &&
					(!reflect.DeepEqual(got.MainConfig().ConfigPath, tt.want.MainConfig().ConfigPath) ||
						!reflect.DeepEqual(got.MainConfig().ContextValue, tt.want.MainConfig().ContextValue) ||
						!reflect.DeepEqual(got.MainConfig().ContextType, tt.want.MainConfig().ContextType) ||
						!reflect.DeepEqual(got.MainConfig().Children, tt.want.MainConfig().Children) ||
						!isSameFather(got.MainConfig(), tt.want.MainConfig()))) {
				t.Errorf("NewMain() got = %v, want %v", got, tt.want)
			}
		})
	}
}
