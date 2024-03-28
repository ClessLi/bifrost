package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"reflect"
	"testing"
)

func TestBasicContext_Child(t *testing.T) {
	testChildCtx := NewContext(context_type.TypeHttp, "")
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
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
			name:   "negative index",
			fields: fields{Children: []context.Context{context.NullContext()}},
			args:   args{idx: -1},
			want:   context.ErrContext(errors.WithCode(code.ErrV3ContextIndexOutOfRange, "index(%d) out of range", -1)),
		},
		{
			name:   "index larger than children's length",
			fields: fields{Children: []context.Context{context.NullContext()}},
			args:   args{idx: 1},
			want:   context.ErrContext(errors.WithCode(code.ErrV3ContextIndexOutOfRange, "index(%d) out of range", 1)),
		},
		{
			name:   "null context child",
			fields: fields{Children: []context.Context{NewContext(context_type.TypeLocation, "~ /test"), context.NullContext()}},
			args:   args{idx: 1},
			want:   context.NullContext(),
		},
		{
			name:   "normal child",
			fields: fields{Children: []context.Context{context.NullContext(), context.NullContext(), testChildCtx, NewContext(context_type.TypeStream, "")}},
			args:   args{idx: 2},
			want:   testChildCtx,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if got := b.Child(tt.args.idx); !reflect.DeepEqual(got, tt.want) {
				if got.Type() != context_type.TypeErrContext || got.Error().Error() != tt.want.Error().Error() {
					t.Errorf("Child() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestBasicContext_Clone(t *testing.T) {
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
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
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if got := b.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicContext_ConfigLines(t *testing.T) {
	_2levelctx := NewContext(context_type.TypeHttp, "")
	_2levelctx.Insert(NewComment("test comment", false), _2levelctx.Len())

	hasComments2levelctx := _2levelctx.Clone().
		Insert(NewComment("the first inline comment", true), 0).
		Insert(NewComment("inline comment after inline comment", true), 1).
		Insert(NewComment("inline comment after comment", true), 3).
		Insert(NewDirective("test_directive", "aaaa bbbb\n cccc"), 4).
		Insert(NewComment("inline comment after other context", true), 5)

	_3levelctx := NewContext(context_type.TypeHttp, "").Insert(NewContext(context_type.TypeServer, "").Insert(NewDirective("server_name", "testserver"), 0), 0)

	includeCtx := NewContext(context_type.TypeInclude, "conf.d/*.conf").(*Include)
	withIncludeCtx, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	withIncludeCtx.Insert(NewContext(context_type.TypeHttp, "").
		Insert(includeCtx, 0),
		0)
	configPath, err := context.NewRelConfigPath("C:\\test", "conf.d\\server.conf")
	if err != nil {
		t.Fatal(err)
	}
	includeConfig := &Config{
		BasicContext: newBasicContext(context_type.TypeConfig, nullHeadString, nullTailString),
		ConfigPath:   configPath,
	}
	includeConfig.self = includeConfig
	includeConfig.ContextValue = "conf.d\\server.conf"
	includeConfig.Insert(NewContext(context_type.TypeServer, "").
		Insert(NewDirective("server_name", "testserver"), 0).
		Insert(NewContext(context_type.TypeLocation, "~ /test"), 1),
		0)
	err = includeCtx.InsertConfig(includeConfig)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
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
			name: "1 level context, not for dumping",
			fields: fields{
				ContextType:    context_type.TypeStream,
				ContextValue:   "",
				Children:       make([]context.Context, 0),
				father:         context.NullContext(),
				self:           context.NullContext(),
				headStringFunc: nonValueBraceHeadString,
				tailStringFunc: braceTailString,
			},
			args: args{isDumping: false},
			want: []string{
				"stream {",
				"}",
			},
			wantErr: false,
		},
		{
			name: "2 level context, not for dumping",
			fields: fields{
				ContextType:    _2levelctx.Type(),
				ContextValue:   _2levelctx.Value(),
				Children:       _2levelctx.(*Http).Children,
				father:         _2levelctx.Father(),
				self:           _2levelctx,
				headStringFunc: _2levelctx.(*Http).headStringFunc,
				tailStringFunc: _2levelctx.(*Http).tailStringFunc,
			},
			args: args{isDumping: false},
			want: []string{
				"http {",
				"    # test comment",
				"}",
			},
			wantErr: false,
		},
		{
			name: "comments children context, not for dumping",
			fields: fields{
				ContextType:    hasComments2levelctx.Type(),
				ContextValue:   hasComments2levelctx.Value(),
				Children:       hasComments2levelctx.(*Http).Children,
				father:         hasComments2levelctx.Father(),
				self:           hasComments2levelctx,
				headStringFunc: hasComments2levelctx.(*Http).headStringFunc,
				tailStringFunc: hasComments2levelctx.(*Http).tailStringFunc,
			},
			args: args{isDumping: false},
			want: []string{
				"http {    # the first inline comment",
				"    # inline comment after inline comment",
				"    # test comment",
				"    # inline comment after comment",
				"    test_directive aaaa bbbb\n cccc;    # inline comment after other context",
				"}",
			},
		},
		{
			name: "3 level context, not for dumping",
			fields: fields{
				ContextType:    _3levelctx.Type(),
				ContextValue:   _3levelctx.Value(),
				Children:       _3levelctx.(*Http).Children,
				father:         _3levelctx.Father(),
				self:           _3levelctx,
				headStringFunc: _3levelctx.(*Http).headStringFunc,
				tailStringFunc: _3levelctx.(*Http).tailStringFunc,
			},
			args: args{isDumping: false},
			want: []string{
				"http {",
				"    server {",
				"        server_name testserver;",
				"    }",
				"}",
			},
		},
		{
			name: "with include context, not for dumping",
			fields: fields{
				ContextType:    withIncludeCtx.Child(0).Type(),
				ContextValue:   withIncludeCtx.Child(0).Value(),
				Children:       withIncludeCtx.Child(0).(*Http).Children,
				father:         withIncludeCtx.Child(0).Father(),
				self:           withIncludeCtx.Child(0),
				headStringFunc: withIncludeCtx.Child(0).(*Http).headStringFunc,
				tailStringFunc: withIncludeCtx.Child(0).(*Http).tailStringFunc,
			},
			args: args{isDumping: false},
			want: []string{
				"http {",
				"    # include <== conf.d/*.conf",
				"    server {",
				"        server_name testserver;",
				"        location ~ /test {",
				"        }",
				"    }",
				"}",
			},
			wantErr: false,
		},
		{
			name: "child ConfigLines() return error, not for dumping",
			fields: fields{
				ContextType:    context_type.TypeHttp,
				ContextValue:   "",
				Children:       []context.Context{context.NullContext()},
				father:         context.NullContext(),
				self:           context.NullContext(),
				headStringFunc: nonValueBraceHeadString,
				tailStringFunc: braceTailString,
			},
			args:    args{isDumping: false},
			want:    nil,
			wantErr: true,
		},
		{
			name: "child is nil, not for dumping",
			fields: fields{
				ContextType:    context_type.TypeHttp,
				ContextValue:   "",
				Children:       []context.Context{nil},
				father:         context.NullContext(),
				self:           context.NullContext(),
				headStringFunc: nonValueBraceHeadString,
				tailStringFunc: braceTailString,
			},
			args:    args{isDumping: false},
			want:    nil,
			wantErr: true,
		},
		{
			name: "1 level context, for dumping",
			fields: fields{
				ContextType:    context_type.TypeStream,
				ContextValue:   "",
				Children:       make([]context.Context, 0),
				father:         context.NullContext(),
				self:           context.NullContext(),
				headStringFunc: nonValueBraceHeadString,
				tailStringFunc: braceTailString,
			},
			args: args{isDumping: true},
			want: []string{
				"stream {",
				"}",
			},
			wantErr: false,
		},
		{
			name: "2 level context, for dumping",
			fields: fields{
				ContextType:    _2levelctx.Type(),
				ContextValue:   _2levelctx.Value(),
				Children:       _2levelctx.(*Http).Children,
				father:         _2levelctx.Father(),
				self:           _2levelctx,
				headStringFunc: _2levelctx.(*Http).headStringFunc,
				tailStringFunc: _2levelctx.(*Http).tailStringFunc,
			},
			args: args{isDumping: true},
			want: []string{
				"http {",
				"    # test comment",
				"}",
			},
			wantErr: false,
		},
		{
			name: "comments children context, for dumping",
			fields: fields{
				ContextType:    hasComments2levelctx.Type(),
				ContextValue:   hasComments2levelctx.Value(),
				Children:       hasComments2levelctx.(*Http).Children,
				father:         hasComments2levelctx.Father(),
				self:           hasComments2levelctx,
				headStringFunc: hasComments2levelctx.(*Http).headStringFunc,
				tailStringFunc: hasComments2levelctx.(*Http).tailStringFunc,
			},
			args: args{isDumping: true},
			want: []string{
				"http {    # the first inline comment",
				"    # inline comment after inline comment",
				"    # test comment",
				"    # inline comment after comment",
				"    test_directive aaaa bbbb\n cccc;    # inline comment after other context",
				"}",
			},
		},
		{
			name: "3 level context, for dumping",
			fields: fields{
				ContextType:    _3levelctx.Type(),
				ContextValue:   _3levelctx.Value(),
				Children:       _3levelctx.(*Http).Children,
				father:         _3levelctx.Father(),
				self:           _3levelctx,
				headStringFunc: _3levelctx.(*Http).headStringFunc,
				tailStringFunc: _3levelctx.(*Http).tailStringFunc,
			},
			args: args{isDumping: true},
			want: []string{
				"http {",
				"    server {",
				"        server_name testserver;",
				"    }",
				"}",
			},
		},
		{
			name: "with include context, for dumping",
			fields: fields{
				ContextType:    withIncludeCtx.Child(0).Type(),
				ContextValue:   withIncludeCtx.Child(0).Value(),
				Children:       withIncludeCtx.Child(0).(*Http).Children,
				father:         withIncludeCtx.Child(0).Father(),
				self:           withIncludeCtx.Child(0),
				headStringFunc: withIncludeCtx.Child(0).(*Http).headStringFunc,
				tailStringFunc: withIncludeCtx.Child(0).(*Http).tailStringFunc,
			},
			args: args{isDumping: true},
			want: []string{
				"http {",
				"    include conf.d/*.conf;",
				"}",
			},
			wantErr: false,
		},
		{
			name: "child ConfigLines() return error, for dumping",
			fields: fields{
				ContextType:    context_type.TypeHttp,
				ContextValue:   "",
				Children:       []context.Context{context.NullContext()},
				father:         context.NullContext(),
				self:           context.NullContext(),
				headStringFunc: nonValueBraceHeadString,
				tailStringFunc: braceTailString,
			},
			args:    args{isDumping: true},
			want:    nil,
			wantErr: true,
		},
		{
			name: "child is nil, for dumping",
			fields: fields{
				ContextType:    context_type.TypeHttp,
				ContextValue:   "",
				Children:       []context.Context{nil},
				father:         context.NullContext(),
				self:           context.NullContext(),
				headStringFunc: nonValueBraceHeadString,
				tailStringFunc: braceTailString,
			},
			args:    args{isDumping: true},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			got, err := b.ConfigLines(tt.args.isDumping)
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

func TestBasicContext_Error(t *testing.T) {
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "nil error",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if err := b.Error(); (err != nil) != tt.wantErr {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasicContext_Father(t *testing.T) {
	testFatherLocation := NewContext(context_type.TypeLocation, "~ /test")
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name:   "null context",
			fields: fields{father: context.NullContext()},
			want:   context.NullContext(),
		},
		{
			name:   "test father location",
			fields: fields{father: testFatherLocation},
			want:   testFatherLocation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if got := b.Father(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Father() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicContext_HasChild(t *testing.T) {
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "nil children",
			fields: fields{Children: nil},
			want:   false,
		},
		{
			name:   "empty children",
			fields: fields{Children: make([]context.Context, 0)},
			want:   false,
		},
		{
			name:   "has children",
			fields: fields{Children: []context.Context{context.NullContext()}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if got := b.HasChild(); got != tt.want {
				t.Errorf("HasChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicContext_Insert(t *testing.T) {
	testCtx := NewContext(context_type.TypeServer, "").
		Insert(NewDirective("server_name", "testserver"), 0).
		Insert(NewContext(context_type.TypeLocation, "~ /test"), 1)
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
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
			name:   "insert into negative index",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: nil,
				idx: -1,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3ContextIndexOutOfRange, "index(%d) out of range", -1)),
		},
		{
			name:   "insert nil",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: nil,
				idx: 0,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert nil")),
		},
		{
			name:   "insert null context",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: context.NullContext(),
				idx: 0,
			},
			want: context.NullContext().(*context.ErrorContext).AppendError(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert error context")),
		},
		{
			name:   "insert error context",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: context.ErrContext(errors.New("test error")),
				idx: 0,
			},
			want: context.ErrContext(errors.New("test error")).(*context.ErrorContext).AppendError(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert error context")),
		},
		{
			name:   "insert invalid error context",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: &Location{BasicContext{ContextType: context_type.TypeErrContext}},
				idx: 0,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert invalid context")),
		},
		{
			name:   "insert config context",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: NewContext(context_type.TypeConfig, "test.conf"),
				idx: 0,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert config context")),
		},
		{
			name:   "insert invalid config context",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: &Http{BasicContext{ContextType: context_type.TypeConfig}},
				idx: 0,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert invalid context")),
		},
		{
			name: "insert normal context into index within range",
			fields: fields{
				ContextType:    testCtx.Type(),
				ContextValue:   testCtx.Value(),
				Children:       testCtx.(*Server).Children,
				father:         testCtx.Father(),
				headStringFunc: testCtx.(*Server).headStringFunc,
				tailStringFunc: testCtx.(*Server).tailStringFunc,
			},
			args: args{
				ctx: NewContext(context_type.TypeLocation, "~ /test2"),
				idx: 0,
			},
			want: testCtx,
		},
		{
			name: "insert normal context into index beyond range",
			fields: fields{
				ContextType:    testCtx.Type(),
				ContextValue:   testCtx.Value(),
				Children:       testCtx.(*Server).Children,
				father:         testCtx.Father(),
				headStringFunc: testCtx.(*Server).headStringFunc,
				tailStringFunc: testCtx.(*Server).tailStringFunc,
			},
			args: args{
				ctx: NewContext(context_type.TypeLocation, "~ /test3"),
				idx: testCtx.Len() + 1,
			},
			want: testCtx,
		},
		{
			name: "inserted context cannot set father",
			fields: fields{
				ContextType:    testCtx.Type(),
				ContextValue:   testCtx.Value(),
				Children:       testCtx.(*Server).Children,
				father:         testCtx.Father(),
				headStringFunc: testCtx.(*Server).headStringFunc,
				tailStringFunc: testCtx.(*Server).tailStringFunc,
			},
			args: args{
				ctx: &Main{},
				idx: 0,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3SetFatherContextFailed, "cannot set father for MainContext")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &struct {
				BasicContext
			}{BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}}
			c.self = c
			actualIdx := tt.args.idx
			if tt.args.idx > c.Len() {
				actualIdx = c.Len()
			}
			got := c.Insert(tt.args.ctx, tt.args.idx)
			if (got.Type() == context_type.TypeErrContext && got.Error().Error() != tt.want.Error().Error()) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
			if got.Type() != context_type.TypeErrContext && !reflect.DeepEqual(got.Child(actualIdx), tt.args.ctx) {
				t.Errorf("Insert() context into corresponding index(%d) is %v, want %v", actualIdx, got.Child(actualIdx), tt.args.ctx)
			}
		})
	}
}

func TestBasicContext_Len(t *testing.T) {
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "has no children",
			fields: fields{Children: make([]context.Context, 0)},
			want:   0,
		},
		{
			name: "nil children",
			want: 0,
		},
		{
			name:   "has one child",
			fields: fields{Children: []context.Context{context.NullContext()}},
			want:   1,
		},
		{
			name:   "has some children",
			fields: fields{Children: []context.Context{context.NullContext(), context.NullContext(), context.NullContext()}},
			want:   3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if got := b.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicContext_Modify(t *testing.T) {
	testCtx := NewContext(context_type.TypeServer, "").
		Insert(NewDirective("server_name", "testserver"), 0).
		Insert(NewContext(context_type.TypeLocation, "~ /test"), 1)
	test2Ctx := testCtx.Clone()
	hasErrChildCtx := testCtx.Clone()
	hasErrChildCtx.(*Server).Children = append(hasErrChildCtx.(*Server).Children, context.NullContext())
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
	}
	type args struct {
		ctx context.Context
		idx int
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         context.Context
		wantModified bool
	}{
		{
			name:   "modify negative index",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: NewDirective("test", ""),
				idx: -1,
			},
			want:         context.ErrContext(errors.WithCode(code.ErrV3ContextIndexOutOfRange, "index(%d) out of range", -1)).(*context.ErrorContext).AppendError(context.ErrInsertIntoErrorContext),
			wantModified: false,
		},
		{
			name:   "modify to nil",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: nil,
				idx: 0,
			},
			want:         context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to modify to nil")),
			wantModified: false,
		},
		{
			name:   "modify to a null context",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: context.NullContext(),
				idx: 0,
			},
			want:         context.NullContext().(*context.ErrorContext).AppendError(errors.WithCode(code.ErrV3InvalidOperation, "refuse to modify to error context")),
			wantModified: false,
		},
		{
			name:   "modify to an error context",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: context.ErrContext(errors.New("test error")),
				idx: 0,
			},
			want: context.ErrContext(errors.New("test error")).(*context.ErrorContext).
				AppendError(errors.WithCode(code.ErrV3InvalidOperation, "refuse to modify to error context")),
			wantModified: false,
		},
		{
			name:   "modify to an invalid error context",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				ctx: &Location{BasicContext{ContextType: context_type.TypeErrContext}},
				idx: 0,
			},
			want:         context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to modify to invalid context")),
			wantModified: false,
		},
		{
			name: "modify normal context into index within range",
			fields: fields{
				ContextType:    testCtx.Type(),
				ContextValue:   testCtx.Value(),
				Children:       testCtx.(*Server).Children,
				father:         testCtx.Father(),
				headStringFunc: testCtx.(*Server).headStringFunc,
				tailStringFunc: testCtx.(*Server).tailStringFunc,
			},
			args: args{
				ctx: NewContext(context_type.TypeLocation, "~ /test2"),
				idx: 0,
			},
			want:         testCtx,
			wantModified: true,
		},
		{
			name: "modify normal context into index beyond range",
			fields: fields{
				ContextType:    testCtx.Type(),
				ContextValue:   testCtx.Value(),
				Children:       testCtx.(*Server).Children,
				father:         testCtx.Father(),
				headStringFunc: testCtx.(*Server).headStringFunc,
				tailStringFunc: testCtx.(*Server).tailStringFunc,
			},
			args: args{
				ctx: NewContext(context_type.TypeLocation, "~ /test3"),
				idx: testCtx.Len() + 1,
			},
			want:         testCtx,
			wantModified: true,
		},
		{
			name: "modified context release father context error",
			fields: fields{
				ContextType:    hasErrChildCtx.Type(),
				ContextValue:   hasErrChildCtx.Value(),
				Children:       hasErrChildCtx.(*Server).Children,
				father:         hasErrChildCtx.Father(),
				headStringFunc: hasErrChildCtx.(*Server).headStringFunc,
				tailStringFunc: hasErrChildCtx.(*Server).tailStringFunc,
			},
			args: args{
				ctx: NewContext(context_type.TypeLocation, "~ /test4"),
				idx: hasErrChildCtx.Len() - 1,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3SetFatherContextFailed, context.NullContext().(*context.ErrorContext).
				AppendError(context.ErrSetFatherToErrorContext).Error().Error())).(*context.ErrorContext).
				AppendError(context.ErrInsertIntoErrorContext),
			wantModified: false,
		},
		{
			name: "modify to modified context, itself",
			fields: fields{
				ContextType:    test2Ctx.Type(),
				ContextValue:   test2Ctx.Value(),
				Children:       test2Ctx.(*Server).Children,
				father:         test2Ctx.Father(),
				headStringFunc: test2Ctx.(*Server).headStringFunc,
				tailStringFunc: test2Ctx.(*Server).tailStringFunc,
			},
			args: args{
				ctx: test2Ctx.Child(0),
				idx: 0,
			},
			want:         test2Ctx,
			wantModified: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &struct {
				BasicContext
			}{BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}}
			c.self = c
			actualIdx := tt.args.idx
			if tt.args.idx > c.Len() {
				actualIdx = c.Len()
			}
			var modifiedCtx context.Context
			if actualIdx >= 0 && c.Len() > 0 {
				modifiedCtx = c.Child(actualIdx)
			}
			got := c.Modify(tt.args.ctx, tt.args.idx)
			if (got.Type() == context_type.TypeErrContext && got.Error().Error() != tt.want.Error().Error()) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Modify() = %v, want %v", got, tt.want)
			}
			if tt.wantModified {
				if !reflect.DeepEqual(got.Child(actualIdx), tt.args.ctx) {
					t.Errorf("Modify() context into corresponding index(%d) is %v, want %v", actualIdx, got.Child(actualIdx), tt.args.ctx)
				}
				if modifiedCtx == c.Child(actualIdx) {
					t.Errorf("Modify() context = %v, want to modify to %v", modifiedCtx, tt.args.ctx)
				}
			} else {
				if modifiedCtx != nil && modifiedCtx != c.Child(actualIdx) {
					t.Errorf("Modify() context = %v does not appear to be modified, but is modified to %v", modifiedCtx, got.Child(actualIdx))
				}
			}
		})
	}
}

func TestBasicContext_QueryAllByKeyWords(t *testing.T) {
	testFather := NewContext(context_type.TypeServer, "").
		Insert(NewContext(context_type.TypeLocation, "~ /test"), 0).
		Insert(NewContext(context_type.TypeLocation, "/text"), 1).
		Insert(NewContext(context_type.TypeLocation, "~ /test2"), 2)
	testContext := NewContext(context_type.TypeHttp, "").Insert(testFather, 0).(*Http)
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
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
				ContextType:    testContext.ContextType,
				ContextValue:   testContext.ContextValue,
				Children:       testContext.Children,
				father:         testContext.father,
				self:           testContext.self,
				headStringFunc: testContext.headStringFunc,
				tailStringFunc: testContext.tailStringFunc,
			},
			args: args{kw: context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue("test")},
			want: []context.Pos{
				context.SetPos(testFather, 0),
				context.SetPos(testFather, 2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if got := b.QueryAllByKeyWords(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryAllByKeyWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicContext_QueryByKeyWords(t *testing.T) {
	testFather := NewContext(context_type.TypeServer, "").
		Insert(NewContext(context_type.TypeLocation, "~ /test"), 0).
		Insert(NewContext(context_type.TypeLocation, "/text"), 1).
		Insert(NewContext(context_type.TypeLocation, "~ /test2"), 2)
	testContext := NewContext(context_type.TypeHttp, "").Insert(testFather, 0).(*Http)
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
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
				ContextType:    testContext.ContextType,
				ContextValue:   testContext.ContextValue,
				Children:       testContext.Children,
				father:         testContext.father,
				self:           testContext.self,
				headStringFunc: testContext.headStringFunc,
				tailStringFunc: testContext.tailStringFunc,
			},
			args: args{kw: context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue("test")},
			want: context.SetPos(testFather, 0),
		},
		{
			name: "has not been matched context",
			fields: fields{
				ContextType:    testContext.ContextType,
				ContextValue:   testContext.ContextValue,
				Children:       testContext.Children,
				father:         testContext.father,
				self:           testContext.self,
				headStringFunc: testContext.headStringFunc,
				tailStringFunc: testContext.tailStringFunc,
			},
			args: args{kw: context.NewKeyWords(context_type.TypeComment)},
			want: context.NotFoundPos(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if got := b.QueryByKeyWords(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryByKeyWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicContext_Remove(t *testing.T) {
	testCtx := NewContext(context_type.TypeServer, "").
		Insert(NewDirective("server_name", "testserver"), 0).
		Insert(NewContext(context_type.TypeLocation, "~ /test"), 1)
	hasErrChildCtx := testCtx.Clone()
	hasErrChildCtx.(*Server).Children = append(hasErrChildCtx.(*Server).Children, context.NullContext())
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           context.Context
		wantRemovedCtx context.Context
	}{
		{
			name:   "remove negative index",
			fields: fields{Children: make([]context.Context, 0)},
			args: args{
				idx: -1,
			},
			want:           context.ErrContext(errors.WithCode(code.ErrV3ContextIndexOutOfRange, "index(%d) out of range", -1)),
			wantRemovedCtx: nil,
		},
		{
			name: "remove the normal context, whose index is within range",
			fields: fields{
				ContextType:    testCtx.Type(),
				ContextValue:   testCtx.Value(),
				Children:       testCtx.(*Server).Children,
				father:         testCtx.Father(),
				headStringFunc: testCtx.(*Server).headStringFunc,
				tailStringFunc: testCtx.(*Server).tailStringFunc,
			},
			args:           args{idx: 0},
			want:           testCtx,
			wantRemovedCtx: testCtx.Child(0),
		},
		{
			name: "remove the normal context, whose index is beyond range",
			fields: fields{
				ContextType:    testCtx.Type(),
				ContextValue:   testCtx.Value(),
				Children:       testCtx.(*Server).Children,
				father:         testCtx.Father(),
				headStringFunc: testCtx.(*Server).headStringFunc,
				tailStringFunc: testCtx.(*Server).tailStringFunc,
			},
			args:           args{idx: testCtx.Len()},
			want:           testCtx,
			wantRemovedCtx: nil,
		},
		{
			name: "removed context release father context error",
			fields: fields{
				ContextType:    hasErrChildCtx.Type(),
				ContextValue:   hasErrChildCtx.Value(),
				Children:       hasErrChildCtx.(*Server).Children,
				father:         hasErrChildCtx.Father(),
				headStringFunc: hasErrChildCtx.(*Server).headStringFunc,
				tailStringFunc: hasErrChildCtx.(*Server).tailStringFunc,
			},
			args: args{
				idx: hasErrChildCtx.Len() - 1,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3SetFatherContextFailed, context.NullContext().(*context.ErrorContext).
				AppendError(context.ErrSetFatherToErrorContext).Error().Error())),
			wantRemovedCtx: hasErrChildCtx.Child(hasErrChildCtx.Len() - 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &struct {
				BasicContext
			}{BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}}
			c.self = c
			var removedCtx context.Context
			if tt.args.idx >= 0 && tt.args.idx < c.Len() {
				removedCtx = c.Child(tt.args.idx)
			}
			got := c.Remove(tt.args.idx)
			if (got.Type() == context_type.TypeErrContext && got.Error().Error() != tt.want.Error().Error()) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(removedCtx, tt.wantRemovedCtx) {
				t.Errorf("Remove() the corresponding index(%d) context = %v, want %v", tt.args.idx, removedCtx, tt.wantRemovedCtx)
			}
		})
	}
}

func TestBasicContext_SetFather(t *testing.T) {
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
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
			name:    "nil error",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if err := b.SetFather(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SetFather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasicContext_SetValue(t *testing.T) {
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
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
			name:    "nil error",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if err := b.SetValue(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasicContext_Type(t *testing.T) {
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
	}
	tests := []struct {
		name   string
		fields fields
		want   context_type.ContextType
	}{
		{
			name:   "normal test",
			fields: fields{ContextType: context_type.TypeLocation},
			want:   context_type.TypeLocation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if got := b.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicContext_Value(t *testing.T) {
	type fields struct {
		ContextType    context_type.ContextType
		ContextValue   string
		Children       []context.Context
		father         context.Context
		self           context.Context
		headStringFunc func(ctxType context_type.ContextType, value string) string
		tailStringFunc func() string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "normal test",
			fields: fields{ContextValue: "test_value"},
			want:   "test_value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicContext{
				ContextType:    tt.fields.ContextType,
				ContextValue:   tt.fields.ContextValue,
				Children:       tt.fields.Children,
				father:         tt.fields.father,
				self:           tt.fields.self,
				headStringFunc: tt.fields.headStringFunc,
				tailStringFunc: tt.fields.tailStringFunc,
			}
			if got := b.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newBasicContext(t *testing.T) {
	type args struct {
		ctxType context_type.ContextType
		head    func(context_type.ContextType, string) string
		tail    func() string
	}
	tests := []struct {
		name string
		args args
		want BasicContext
	}{
		{
			name: "normal test",
			args: args{
				ctxType: context_type.TypeLocation,
				head:    hasValueBraceHeadString,
				tail:    braceTailString,
			},
			want: BasicContext{
				ContextType:    context_type.TypeLocation,
				Children:       make([]context.Context, 0),
				father:         context.NullContext(),
				headStringFunc: hasValueBraceHeadString,
				tailStringFunc: braceTailString,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newBasicContext(tt.args.ctxType, tt.args.head, tt.args.tail); !reflect.DeepEqual(got.father, tt.want.father) ||
				got.ContextType != tt.want.ContextType ||
				got.ContextValue != tt.want.ContextValue ||
				!reflect.DeepEqual(got.Children, tt.want.Children) {
				t.Errorf("newBasicContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
