package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"reflect"
	"testing"
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
	testIncludes := NewContext(context_type.TypeInclude, "conf.d\\include*conf").(*Include)
	emptyMain, err := NewMain("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}
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
					).
					Insert(testIncludes, 2),
				1,
			),
		0,
	)
	location1conf := NewContext(context_type.TypeConfig, "conf.d\\include.location1.conf").
		Insert(NewContext(context_type.TypeLocation, "~ /test1"), 0).(*Config)
	location1conf.ConfigPath, err = newConfigPath(testMain.graph(), location1conf.Value())
	if err != nil {
		t.Fatal(err)
	}
	location2conf := NewContext(context_type.TypeConfig, "conf.d\\include.location2.conf").
		Insert(NewContext(context_type.TypeLocation, "^~ /test2"), 0).(*Config)
	location2conf.ConfigPath, err = newConfigPath(testMain.graph(), location1conf.Value())
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(location1conf)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(location2conf)
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
			want:    []byte(`{"main-config":"C:\\test\\test.conf","configs":{"C:\\test\\test.conf":{"context-type":"config","value":"C:\\test\\test.conf"}}}`),
			wantErr: false,
		},
		{
			name:   "normal test",
			fields: fields{ConfigGraph: testMain.graph()},
			want: []byte(
				`{"main-config":"C:\\test\\test.conf","configs":{"C:\\test\\test.conf":{"context-type":"config","value":"C:\\test\\test.conf","params":[{"context-type":"http","params":[{"context-type":"inline_comment","value":"test comment"},{"context-type":"server","params":[{"context-type":"directive","value":"server_name testserver"},{"context-type":"location","value":"~ /test"},{"context-type":"include","value":"conf.d\\include*conf","params":["conf.d\\include.location1.conf","conf.d\\include.location2.conf"]}]}]}]},"conf.d\\include.location1.conf":{"context-type":"config","value":"conf.d\\include.location1.conf","params":[{"context-type":"location","value":"~ /test1"}]},"conf.d\\include.location2.conf":{"context-type":"config","value":"conf.d\\include.location2.conf","params":[{"context-type":"location","value":"^~ /test2"}]}}}`,
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

func TestMain_Modify(t *testing.T) {
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

func TestMain_QueryAllByKeyWords(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if got := m.QueryAllByKeyWords(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryAllByKeyWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain_QueryByKeyWords(t *testing.T) {
	type fields struct {
		ConfigGraph ConfigGraph
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Main{
				ConfigGraph: tt.fields.ConfigGraph,
			}
			if got := m.QueryByKeyWords(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryByKeyWords() = %v, want %v", got, tt.want)
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
						!reflect.DeepEqual(got.MainConfig().BasicContext.ContextValue, tt.want.MainConfig().BasicContext.ContextValue) ||
						!reflect.DeepEqual(got.MainConfig().BasicContext.ContextType, tt.want.MainConfig().BasicContext.ContextType) ||
						!reflect.DeepEqual(got.MainConfig().BasicContext.Children, tt.want.MainConfig().BasicContext.Children) ||
						!isSameFather(got.MainConfig(), tt.want.MainConfig()))) {
				t.Errorf("NewMain() got = %v, want %v", got, tt.want)
			}
		})
	}
}
