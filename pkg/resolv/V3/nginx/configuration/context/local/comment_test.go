package local

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
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
		want   context.PosSet
	}{
		{
			name: "has no children",
			want: context.NewPosSet(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.ChildrenPosSet().QueryAll(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
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
			want: context.NotFoundPos(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Comment{
				Comments:      tt.fields.Comments,
				Inline:        tt.fields.Inline,
				fatherContext: tt.fields.fatherContext,
			}
			if got := c.ChildrenPosSet().QueryOne(tt.args.kw); !reflect.DeepEqual(got, tt.want) {
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

func Test_commentsToContextConverter_Convert(t *testing.T) {
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
	lines, err := testMain.ConfigLines(false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(strings.Join(lines, "\n"))

	toBeConvertingMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	toBeConvertingMain.Insert(
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
			NewContext(context_type.TypeComment, "server {    # disabled server"),
			2,
		).Insert(
			NewContext(context_type.TypeComment, "    listen 8080;"),
			3,
		).Insert(
			NewContext(context_type.TypeComment, "    server_name example.com;"),
			4,
		).Insert(
			NewContext(context_type.TypeComment, "    # location ~ /disabled-location {"),
			5,
		).Insert(
			NewContext(context_type.TypeComment, "    #     proxy_pass http://disabled-url;"),
			6,
		).Insert(
			NewContext(context_type.TypeComment, "    # }"),
			7,
		).Insert(
			NewContext(context_type.TypeComment, "    include conf.d/disabled_location.conf;"),
			8,
		).Insert(
			NewContext(context_type.TypeComment, "    include conf.d/strange_location.conf;"),
			9,
		).Insert(
			NewContext(context_type.TypeComment, "}"),
			10,
		),
		0,
	)

	err = toBeConvertingMain.AddConfig(
		NewContext(context_type.TypeConfig, "conf.d/disabled_location.conf").Insert(
			NewContext(context_type.TypeComment, "# disabled config"),
			0,
		).Insert(
			NewContext(context_type.TypeComment, "location ~ /test {"),
			1,
		).Insert(
			NewContext(context_type.TypeComment, "return 404;"),
			2,
		).Insert(
			NewContext(context_type.TypeComment, "}"),
			3,
		).Insert(
			NewContext(context_type.TypeComment, "location ~ /has-disabled-ctx {"),
			4,
		).Insert(
			NewContext(context_type.TypeComment, "# disabled if ctx"),
			5,
		).Insert(
			NewContext(context_type.TypeComment, "# if ($is_enabled ~* false) {"),
			6,
		).Insert(
			NewContext(context_type.TypeComment, "# # set $is_enabled true;"),
			7,
		).Insert(
			NewContext(context_type.TypeComment, "# return 404;"),
			8,
		).Insert(
			NewContext(context_type.TypeComment, "# }"),
			9,
		).Insert(
			NewContext(context_type.TypeComment, "}"),
			10,
		).Insert(
			NewContext(context_type.TypeComment, "# }"),
			11,
		).Insert(
			NewContext(context_type.TypeComment, "}"),
			12,
		).(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = toBeConvertingMain.AddConfig(
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
	lines2, err := toBeConvertingMain.ConfigLines(false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(strings.Join(lines2, "\n"))
	j, _ := json.Marshal(testMain)
	fmt.Println(string(j))
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "normal test",
			args: args{ctx: toBeConvertingMain},
			want: testMain,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := commentsToContextConverter{}
			got, gotMarshalErr := json.Marshal(c.Convert(tt.args.ctx))
			if gotMarshalErr != nil {
				t.Fatal(gotMarshalErr)
			}
			want, wantMarshalErr := json.Marshal(tt.want)
			if wantMarshalErr != nil {
				t.Fatal(wantMarshalErr)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Convert() ==> jsonMarshal = %v, want %v", string(got), string(want))
			}
		})
	}
}

func Test_commentsToContextConverter_sliceContinuousIndexes(t *testing.T) {
	type args struct {
		indexes []int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "[ 1 ]",
			args: args{indexes: []int{1}},
			want: [][]int{{1}},
		},
		{
			name: "[ 1, 2 ]",
			args: args{indexes: []int{1, 2}},
			want: [][]int{{1, 2}},
		},
		{
			name: "[1, 3, 4, 6, 7 ,8, 11 ]",
			args: args{indexes: []int{1, 3, 4, 6, 7, 8, 11}},
			want: [][]int{{1}, {3, 4}, {6, 7, 8}, {11}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := commentsToContextConverter{}
			if got := c.sliceContinuousIndexes(tt.args.indexes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sliceContinuousIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}
