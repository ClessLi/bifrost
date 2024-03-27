package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"reflect"
	"testing"
)

func TestBuildBasicContextConfig_BasicContext(t *testing.T) {
	type fields struct {
		ContextType    context_type.ContextType
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
	}
	tests := []struct {
		name   string
		fields fields
		want   BasicContext
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BuildBasicContextConfig{
				ContextType:    tt.fields.ContextType,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if got := b.BasicContext(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BasicContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
	testMain := NewContext(context_type.TypeMain, "C:\\test\\test.conf").(*Main)
	type args struct {
		ctx context.Context
		idx int
	}
	tests := []struct {
		name   string
		fields *Main
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
				ctx: NewDirective("test", ""),
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
	testMain := NewContext(context_type.TypeMain, "C:\\test\\test.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						).
						Insert(testIncludes, 2),
					1,
				),
			0,
		).(*Main)
	var err error
	location1conf := NewContext(context_type.TypeConfig, "conf.d\\include.location1.conf").
		Insert(NewContext(context_type.TypeLocation, "~ /test1"), 0).(*Config)
	location1conf.ConfigPath, err = newConfigPath(testMain.ConfigGraph, location1conf.Value())
	if err != nil {
		t.Fatal(err)
	}
	location2conf := NewContext(context_type.TypeConfig, "conf.d\\include.location2.conf").
		Insert(NewContext(context_type.TypeLocation, "^~ /test2"), 0).(*Config)
	location2conf.ConfigPath, err = newConfigPath(testMain.ConfigGraph, location1conf.Value())
	if err != nil {
		t.Fatal(err)
	}
	err = testIncludes.InsertConfig(location1conf, location2conf)
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
			fields:  fields{ConfigGraph: NewContext(context_type.TypeMain, "C:\\test\\test.conf").(*Main).ConfigGraph},
			want:    []byte(`{"main-config":"C:\\test\\test.conf","configs":{"C:\\test\\test.conf":[]}}`),
			wantErr: false,
		},
		{
			name:   "normal test",
			fields: fields{ConfigGraph: testMain.ConfigGraph},
			want: []byte(
				`{"main-config":"C:\\test\\test.conf",` +
					`"configs":{"C:\\test\\test.conf":[{` +
					`"http":{"params":[{"comments":"test comment","inline":true},` +
					`{"server":{"params":[{"directive":"server_name","params":"testserver"},` +
					`{"location":{"value":"~ /test"}},` +
					`{"include":{"value":"conf.d\\include*conf","params":["conf.d\\include.location1.conf","conf.d\\include.location2.conf"]}}` +
					`]}}` +
					`]}` +
					`}],` +
					`"conf.d\\include.location1.conf":[{"location":{"value":"~ /test1"}}],` +
					`"conf.d\\include.location2.conf":[{"location":{"value":"^~ /test2"}}]` +
					`}` +
					`}`,
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
	testMain := NewContext(context_type.TypeMain, "C:\\test\\test.conf").(*Main)
	type args struct {
		ctx context.Context
		idx int
	}
	tests := []struct {
		name   string
		fields *Main
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
				ctx: NewDirective("test", ""),
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
	testMain := NewContext(context_type.TypeMain, "C:\\test\\test.conf").(*Main)
	type args struct {
		idx int
	}
	tests := []struct {
		name   string
		fields *Main
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

func TestNewContext(t *testing.T) {
	type args struct {
		contextType context_type.ContextType
		value       string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContext(tt.args.contextType, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptsApplyTo(t *testing.T) {
	type args struct {
		opts context.BuildOptions
	}
	tests := []struct {
		name    string
		args    args
		want    BuildBasicContextConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OptsApplyTo(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("OptsApplyTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OptsApplyTo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterBuilder(t *testing.T) {
	type args struct {
		opts      context.BuildOptions
		registrar ContextBuilderRegistrar
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
			if err := RegisterBuilder(tt.args.opts, tt.args.registrar); (err != nil) != tt.wantErr {
				t.Errorf("RegisterBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterParseFunc(t *testing.T) {
	type args struct {
		opts       parseFuncBuildOptions
		parserFunc map[context_type.ContextType]parseFunc
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
			if err := RegisterParseFunc(tt.args.opts, tt.args.parserFunc); (err != nil) != tt.wantErr {
				t.Errorf("RegisterParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_braceTailString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := braceTailString(); got != tt.want {
				t.Errorf("braceTailString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_directiveTailString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := directiveTailString(); got != tt.want {
				t.Errorf("directiveTailString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasValueBraceHeadString(t *testing.T) {
	type args struct {
		ctxType context_type.ContextType
		value   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasValueBraceHeadString(tt.args.ctxType, tt.args.value); got != tt.want {
				t.Errorf("hasValueBraceHeadString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nonValueBraceHeadString(t *testing.T) {
	type args struct {
		ctxType context_type.ContextType
		in1     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nonValueBraceHeadString(tt.args.ctxType, tt.args.in1); got != tt.want {
				t.Errorf("nonValueBraceHeadString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nullHeadString(t *testing.T) {
	type args struct {
		in0 context_type.ContextType
		in1 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nullHeadString(tt.args.in0, tt.args.in1); got != tt.want {
				t.Errorf("nullHeadString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nullTailString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nullTailString(); got != tt.want {
				t.Errorf("nullTailString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_registerContextBuilders(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerContextBuilders(); (err != nil) != tt.wantErr {
				t.Errorf("registerContextBuilders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerContextParseFuncs(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerContextParseFuncs(); (err != nil) != tt.wantErr {
				t.Errorf("registerContextParseFuncs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerEventsBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerEventsBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerEventsBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerEventsParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerEventsParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerEventsParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerGeoBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerGeoBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerGeoBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerGeoParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerGeoParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerGeoParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerHttpBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerHttpBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerHttpBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerHttpParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerHttpParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerHttpParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerIfBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerIfBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerIfBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerIfParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerIfParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerIfParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerLimitExceptBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerLimitExceptBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerLimitExceptBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerLimitExceptParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerLimitExceptParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerLimitExceptParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerLocationBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerLocationBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerLocationBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerLocationParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerLocationParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerLocationParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerMainBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerMainBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerMainBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerMapBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerMapBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerMapBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerMapParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerMapParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerMapParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerServerBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerServerBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerServerBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerServerParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerServerParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerServerParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerStreamBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerStreamBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerStreamBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerStreamParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerStreamParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerStreamParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerTypesBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerTypesBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerTypesBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerTypesParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerTypesParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerTypesParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerUpstreamBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerUpstreamBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerUpstreamBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerUpstreamParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerUpstreamParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerUpstreamParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
