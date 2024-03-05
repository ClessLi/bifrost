package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"reflect"
	"testing"
)

func TestDirective_Child(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
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
			name: "has no child",
			want: context.ErrContext(errors.WithCode(code.V3ErrInvalidOperation, "directive has no child context")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.Child(tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Child() = %v, want %v", got.Error(), tt.want.Error())
			}
		})
	}
}

func TestDirective_ConfigLines(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
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
			name: "view directive without params, not for dumping",
			fields: fields{
				Name:          "   test_directive",
				Params:        "",
				fatherContext: context.NullContext(),
			},
			args:    args{isDumping: false},
			want:    []string{"test_directive;"},
			wantErr: false,
		},
		{
			name: "view directive with params, not for dumping",
			fields: fields{
				Name:          "  test_directive  ",
				Params:        "'  test_param1\n   test_param2'    ",
				fatherContext: context.NullContext(),
			},
			args:    args{isDumping: false},
			want:    []string{"test_directive '  test_param1\n   test_param2';"},
			wantErr: false,
		},
		{
			name: "view directive without params, for dumping",
			fields: fields{
				Name:          "test_directive",
				Params:        "",
				fatherContext: context.NullContext(),
			},
			args:    args{isDumping: true},
			want:    []string{"test_directive;"},
			wantErr: false,
		},
		{
			name: "view directive with params, for dumping",
			fields: fields{
				Name:          "  test_directive",
				Params:        "test_param1\n   test_param2  ",
				fatherContext: context.NullContext(),
			},
			args:    args{isDumping: true},
			want:    []string{"test_directive test_param1\n   test_param2;"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			got, err := d.ConfigLines(tt.args.isDumping)
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

func TestDirective_Error(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
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
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if err := d.Error(); (err != nil) != tt.wantErr {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirective_Father(t *testing.T) {
	testFatherLocation := NewContext(context_type.TypeLocation, "~ /test")
	type fields struct {
		Name          string
		Params        string
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name:   "null context",
			fields: fields{fatherContext: context.NullContext()},
			want:   context.NullContext(),
		},
		{
			name: "test father location",
			fields: fields{
				Name:          "test_directive",
				Params:        "param1",
				fatherContext: testFatherLocation,
			},
			want: testFatherLocation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.Father(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Father() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirective_HasChild(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "no child",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.HasChild(); got != tt.want {
				t.Errorf("HasChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirective_Insert(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
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
			name: "cannot insert child",
			want: context.ErrContext(errors.WithCode(code.V3ErrInvalidOperation, "directive cannot insert context")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.Insert(tt.args.ctx, tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirective_Len(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "length 0",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirective_Modify(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
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
			name: "cannot modify child",
			want: context.ErrContext(errors.WithCode(code.V3ErrInvalidOperation, "directive cannot modify context")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.Modify(tt.args.ctx, tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Modify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirective_QueryAllByKeyWords(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.QueryAllByKeyWords(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryAllByKeyWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirective_QueryByKeyWords(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.QueryByKeyWords(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryByKeyWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirective_Remove(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
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
			name: "cannot remove child",
			want: context.ErrContext(errors.WithCode(code.V3ErrInvalidOperation, "directive cannot remove context")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.Remove(tt.args.idx); got.Error().Error() != tt.want.Error().Error() {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirective_SetFather(t *testing.T) {
	testFatherLocation := NewContext(context_type.TypeLocation, "~ /test")
	type fields struct {
		Name          string
		Params        string
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
		{
			name: "set null context father",
			fields: fields{
				Name:          "test_directive",
				Params:        "",
				fatherContext: testFatherLocation,
			},
			args:    args{ctx: context.NullContext()},
			wantErr: false,
		},
		{
			name: "set test father location",
			fields: fields{
				Name:          "test_directive",
				Params:        "",
				fatherContext: context.NullContext(),
			},
			args:    args{ctx: testFatherLocation},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if err := d.SetFather(tt.args.ctx); reflect.DeepEqual(d.fatherContext, tt.args.ctx) && (err != nil) != tt.wantErr {
				t.Errorf("SetFather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirective_SetValue(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
		fatherContext context.Context
	}
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    fields
		wantErr bool
	}{
		{
			name: "null value",
			fields: fields{
				Name:          "test_directive",
				Params:        "param1",
				fatherContext: context.NullContext(),
			},
			args: args{v: "    "},
			want: fields{
				Name:          "test_directive",
				Params:        "param1",
				fatherContext: context.NullContext(),
			},
			wantErr: true,
		},
		{
			name: "has param",
			fields: fields{
				Name:          "test_directive",
				Params:        "",
				fatherContext: context.NullContext(),
			},
			args: args{v: "  test_directive_2    aaa   bbb\n ccc  "},
			want: fields{
				Name:          "test_directive_2",
				Params:        "aaa   bbb\n ccc",
				fatherContext: context.NullContext(),
			},
			wantErr: false,
		},
		{
			name: "has not param",
			fields: fields{
				Name:          "test_directive",
				Params:        "aaa",
				fatherContext: context.NullContext(),
			},
			args: args{v: "test_directive_3   "},
			want: fields{
				Name:          "test_directive_3",
				Params:        "",
				fatherContext: context.NullContext(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if err := d.SetValue(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := fields{
				Name:          d.Name,
				Params:        d.Params,
				fatherContext: d.fatherContext,
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetValue() got fields = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirective_Type(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   context_type.ContextType
	}{
		{
			name: "only type directive",
			want: context_type.TypeDirective,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirective_Value(t *testing.T) {
	type fields struct {
		Name          string
		Params        string
		fatherContext context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "without params",
			fields: fields{
				Name:          "   test_directive",
				Params:        "",
				fatherContext: context.NullContext(),
			},
			want: "test_directive",
		},
		{
			name: "with params",
			fields: fields{
				Name:          "  test_directive  ",
				Params:        "  test_param1\n   test_param2    ",
				fatherContext: context.NullContext(),
			},
			want: "test_directive test_param1\n   test_param2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Directive{
				Name:          tt.fields.Name,
				Params:        tt.fields.Params,
				fatherContext: tt.fields.fatherContext,
			}
			if got := d.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDirective(t *testing.T) {
	type args struct {
		name   string
		params string
	}
	tests := []struct {
		name string
		args args
		want *Directive
	}{
		{
			name: "new directive without params",
			args: args{
				name:   " test_directive  ",
				params: "",
			},
			want: &Directive{
				Name:          "test_directive",
				Params:        "",
				fatherContext: context.NullContext(),
			},
		},
		{
			name: "new directive with params",
			args: args{
				name:   "  test_directive ",
				params: " param1  param2\n param3   ",
			},
			want: &Directive{
				Name:          "test_directive",
				Params:        "param1  param2\n param3",
				fatherContext: context.NullContext(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDirective(tt.args.name, tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDirective() = %v, want %v", got, tt.want)
			}
		})
	}
}
