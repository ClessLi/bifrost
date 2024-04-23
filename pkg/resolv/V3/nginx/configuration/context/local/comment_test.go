package local

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"reflect"
	"testing"
)

func TestComment_Child(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
		fatherContext context.Context
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErrCtx bool
	}{
		{
			name:       "has no children",
			wantErrCtx: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.Child(tt.args.idx); (got.Error() != nil) != tt.wantErrCtx {
				t.Errorf("Child() return context's error = %v, wantErrCtx %v", got.Error(), tt.wantErrCtx)
			}
		})
	}
}

func TestComment_Clone(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
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
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_ConfigLines(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
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
			name:   "null comment",
			fields: fields{Comments: ""},
			args:   args{isDumping: true},
			want: []string{
				"#",
			},
			wantErr: false,
		},
		{
			name:   "only space comment",
			fields: fields{Comments: "    \t"},
			args:   args{isDumping: true},
			want: []string{
				"#",
			},
			wantErr: false,
		},
		{
			name:   "normal comment",
			fields: fields{Comments: "test comment"},
			args:   args{isDumping: false},
			want: []string{
				"# test comment",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			got, err := c.ConfigLines(tt.args.isDumping)
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

func TestComment_Error(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
		fatherContext context.Context
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
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if err := c.Error(); (err != nil) != tt.wantErr {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComment_Father(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
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
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.Father(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Father() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_HasChild(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "has no children",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.HasChild(); got != tt.want {
				t.Errorf("HasChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_Insert(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
		fatherContext context.Context
	}
	type args struct {
		ctx context.Context
		idx int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErrCtx bool
	}{
		{
			name:       "return error",
			wantErrCtx: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.Insert(tt.args.ctx, tt.args.idx); (got.Error() != nil) != tt.wantErrCtx {
				t.Errorf("Insert() return context's error = %v, wantErrCtx %v", got.Error(), tt.wantErrCtx)
			}
		})
	}
}

func TestComment_Len(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "has no children",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_Modify(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
		fatherContext context.Context
	}
	type args struct {
		ctx context.Context
		idx int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErrCtx bool
	}{
		{
			name:       "return error",
			wantErrCtx: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.Modify(tt.args.ctx, tt.args.idx); (got.Error() != nil) != tt.wantErrCtx {
				t.Errorf("Modify() return context's error = %v, wantErrCtx %v", got.Error(), tt.wantErrCtx)
			}
		})
	}
}

func TestComment_QueryAllByKeyWords(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
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
			name: "has no children",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.QueryAllByKeyWords(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryAllByKeyWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_QueryByKeyWords(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
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
			name: "has no children",
			want: context.NullPos(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.QueryByKeyWords(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryByKeyWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_Remove(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
		fatherContext context.Context
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErrCtx bool
	}{
		{
			name:       "return error",
			wantErrCtx: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.Remove(tt.args.idx); (got.Error() != nil) != tt.wantErrCtx {
				t.Errorf("Remove() return context's error = %v, wantErrCtx %v", got.Error(), tt.wantErrCtx)
			}
		})
	}
}

func TestComment_SetFather(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
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
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if err := c.SetFather(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SetFather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComment_SetValue(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
		fatherContext context.Context
	}
	type args struct {
		v string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantGetValue string
	}{
		{
			name:         "set null comment",
			fields:       fields{Comments: "test comment"},
			wantErr:      false,
			wantGetValue: "",
		},
		{
			name:         "set some only space comment",
			fields:       fields{Comments: "test comment"},
			args:         args{v: "    \t "},
			wantErr:      false,
			wantGetValue: "    \t ",
		},
		{
			name:         "set normal comment",
			args:         args{v: "test comment"},
			wantErr:      false,
			wantGetValue: "test comment",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if err := c.SetValue(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if getvalue := c.Value(); getvalue != tt.wantGetValue {
				t.Errorf("get value = `%s`, want `%s`", getvalue, tt.wantGetValue)
			}
		})
	}
}

func TestComment_Type(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   context_type.ContextType
	}{
		{
			name: "only comment",
			want: context_type.TypeComment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_Value(t *testing.T) {
	type fields struct {
		Comments      string
		Inline        bool
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "null comment",
			want: "",
		},
		{
			name:   "some only space comment",
			fields: fields{Comments: "   \t "},
			want:   "   \t ",
		},
		{
			name:   "normal comment",
			fields: fields{Comments: " test comment"},
			want:   " test comment",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_registerCommentParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerCommentParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerCommentParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
