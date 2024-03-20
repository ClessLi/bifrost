package local

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileLoader(t *testing.T) {
	type args struct {
		configpath string
	}
	tests := []struct {
		name string
		args args
		want Loader
	}{
		{
			name: "normal test",
			args: args{configpath: "C:\\test\\test.conf"},
			want: &fileLoader{
				mainConfigAbsPath: filepath.Clean("C:\\test\\test.conf"),
				contextStack:      newContextStack(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileLoader(tt.args.configpath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileLoader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonLoader(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want Loader
	}{
		{
			name: "normal test",
			args: args{data: []byte("{\"testdata\": \"test\"}")},
			want: &jsonLoader{
				unmarshaler: &mainUnmarshaler{unmarshalContext: new(jsonUnmarshalMain)},
				jsonBytes:   []byte("{\"testdata\": \"test\"}"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JsonLoader(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonLoader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contextStack_current(t *testing.T) {
	testCurrentCtx := NewContext(context_type.TypeServer, "")
	type fields struct {
		contexts []context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    context.Context
		wantErr bool
	}{
		{
			name:    "empty stack",
			want:    context.NullContext(),
			wantErr: true,
		},
		{
			name: "normal test",
			fields: fields{contexts: []context.Context{
				NewContext(context_type.TypeHttp, ""),
				testCurrentCtx,
			}},
			want:    testCurrentCtx,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &contextStack{
				contexts: tt.fields.contexts,
			}
			got, err := s.current()
			if (err != nil) != tt.wantErr {
				t.Errorf("current() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("current() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contextStack_pop(t *testing.T) {
	testPopCtx := NewContext(context_type.TypeLocation, "~ /test")
	testCurrentCtx := NewContext(context_type.TypeServer, "")
	type fields struct {
		contexts []context.Context
	}
	tests := []struct {
		name        string
		fields      fields
		wantPop     context.Context
		wantCurrent context.Context
		wantErr     bool
	}{
		{
			name:    "empty stack",
			wantPop: context.NullContext(),
			wantErr: true,
		},
		{
			name: "after pop, the stack becomes empty",
			fields: fields{contexts: []context.Context{
				testPopCtx,
			}},
			wantPop: testPopCtx,
			wantErr: false,
		},
		{
			name: "normal test",
			fields: fields{contexts: []context.Context{
				NewContext(context_type.TypeHttp, ""),
				testCurrentCtx,
				testPopCtx,
			}},
			wantPop:     testPopCtx,
			wantCurrent: testCurrentCtx,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &contextStack{
				contexts: tt.fields.contexts,
			}
			got, err := s.pop()
			if (err != nil) != tt.wantErr {
				t.Errorf("pop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantPop) {
				t.Errorf("pop() got = %v, want %v", got, tt.wantPop)
				return
			}
			if s != nil && len(s.contexts) > 0 {
				if currentCtx := s.contexts[len(s.contexts)-1]; !reflect.DeepEqual(currentCtx, tt.wantCurrent) {
					t.Errorf("after pop(), currentCtx = %v, wantCurrent %v", currentCtx, tt.wantCurrent)
				}
			} else {
				if tt.wantCurrent != nil {
					t.Errorf("after pop(), wantCurrent %v != nil", tt.wantCurrent)
				}
			}
		})
	}
}

func Test_contextStack_push(t *testing.T) {
	type fields struct {
		contexts []context.Context
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
			name:    "push nil",
			wantErr: true,
		},
		{
			name:    "push an error context",
			args:    args{ctx: context.NullContext()},
			wantErr: true,
		},
		{
			name:    "normal test",
			fields:  fields{contexts: make([]context.Context, 0)},
			args:    args{ctx: NewContext(context_type.TypeHttp, "")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &contextStack{
				contexts: tt.fields.contexts,
			}
			if err := s.push(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fileLoader_Load(t *testing.T) {
	simpleMain := NewContext(context_type.TypeMain, filepath.Join(os.Getenv("GOPATH"), "src/bifrost", "test/config_test/simple_nginx.conf")).(*Main)
	simpleMain.
		Insert(NewComment("user  nobody;", false), 0).
		Insert(NewDirective("worker_processes", "1"), 1).
		Insert(
			NewContext(context_type.TypeEvents, "").
				Insert(NewDirective("worker_connections", "1024"), 0),
			2,
		)
	type fields struct {
		mainConfigAbsPath string
		configGraph       ConfigGraph
		contextStack      *contextStack
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Main
		wantErr bool
	}{
		{
			name: "main config path is not an absolute one",
			fields: fields{
				mainConfigAbsPath: "test.conf",
				configGraph:       nil,
				contextStack:      newContextStack(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "normal test",
			fields: fields{
				mainConfigAbsPath: simpleMain.MainConfig().FullPath(),
				configGraph:       nil,
				contextStack:      newContextStack(),
			},
			want:    simpleMain,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileLoader{
				mainConfigAbsPath: tt.fields.mainConfigAbsPath,
				configGraph:       tt.fields.configGraph,
				contextStack:      tt.fields.contextStack,
			}
			got, err := f.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.want && got == nil {
				return
			}
			gotLines, err := got.ConfigLines(false)
			if err != nil {
				t.Fatal(err)
			}
			wantLines, err := tt.want.ConfigLines(false)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotLines, wantLines) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileLoader_load(t *testing.T) {
	simpleConfig := NewContext(context_type.TypeConfig, filepath.Join(os.Getenv("GOPATH"), "src/bifrost", "test/config_test/simple_nginx.conf")).(*Config)
	var err error
	simpleConfig.ConfigPath, err = context.NewAbsConfigPath(simpleConfig.Value())
	if err != nil {
		t.Fatal(err)
	}
	invalidDirectiveConfig := NewContext(context_type.TypeConfig, filepath.Join(os.Getenv("GOPATH"), "src/bifrost", "test/config_test/invalid_directive_nginx.conf")).(*Config)
	invalidDirectiveConfig.ConfigPath, err = context.NewAbsConfigPath(invalidDirectiveConfig.Value())
	if err != nil {
		t.Fatal(err)
	}
	missEndBraceConfig := NewContext(context_type.TypeConfig, filepath.Join(os.Getenv("GOPATH"), "src/bifrost", "test/config_test/miss_end_brace_nginx.conf")).(*Config)
	missEndBraceConfig.ConfigPath, err = context.NewAbsConfigPath(missEndBraceConfig.Value())
	if err != nil {
		t.Fatal(err)
	}
	includeMain := NewContext(context_type.TypeMain, filepath.Join(os.Getenv("GOPATH"), "src/bifrost", "test/config_test/include_nginx.conf")).(*Main)
	type fields struct {
		mainConfigAbsPath string
		configGraph       ConfigGraph
		contextStack      *contextStack
	}
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "config is not found",
			args:    args{config: NewContext(context_type.TypeConfig, "").(*Config)},
			wantErr: true,
		},
		{
			name:    "normal test",
			fields:  fields{contextStack: newContextStack()},
			args:    args{config: simpleConfig},
			wantErr: false,
		},
		{
			name:    "invalid directive config",
			fields:  fields{contextStack: newContextStack()},
			args:    args{config: invalidDirectiveConfig},
			wantErr: true,
		},
		{
			name:    "missing end brace config",
			fields:  fields{contextStack: newContextStack()},
			args:    args{config: missEndBraceConfig},
			wantErr: true,
		},
		{
			name: "include configs",
			fields: fields{
				contextStack:      newContextStack(),
				mainConfigAbsPath: includeMain.Value(),
				configGraph:       includeMain,
			},
			args:    args{config: includeMain.MainConfig()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileLoader{
				mainConfigAbsPath: tt.fields.mainConfigAbsPath,
				configGraph:       tt.fields.configGraph,
				contextStack:      tt.fields.contextStack,
			}
			if err := f.load(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fileLoader_loadInclude(t *testing.T) {
	simpleMain := NewContext(context_type.TypeMain, filepath.Join(os.Getenv("GOPATH"), "src/bifrost", "test/config_test/simple_nginx.conf")).(*Main)
	type fields struct {
		mainConfigAbsPath string
		configGraph       ConfigGraph
		contextStack      *contextStack
	}
	absPath := filepath.Join(os.Getenv("GOPATH"), "src/bifrost", "test/config_test/conf.d/test*conf")
	relPath := "conf.d/test*.conf"
	noMatchAbsPath := filepath.Join(os.Getenv("GOPATH"), "src/bifrost", "test/config_test/conf.d/test.aa*conf")
	noMatchRelPath := "conf.d/test.aa*.conf"

	absPathInclude := NewContext(context_type.TypeInclude, absPath).(*Include)
	absPathInclude.fatherContext = simpleMain.MainConfig()
	relPathInclude := NewContext(context_type.TypeInclude, relPath).(*Include)
	relPathInclude.fatherContext = simpleMain.MainConfig()
	cycleInclude := NewContext(context_type.TypeInclude, "conf.d/cycle1.conf").(*Include)
	cycleInclude.fatherContext = simpleMain.MainConfig()
	type args struct {
		include *Include
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "absolute path include",
			fields: fields{
				mainConfigAbsPath: simpleMain.Value(),
				configGraph:       simpleMain.ConfigGraph,
				contextStack:      newContextStack(),
			},
			args:    args{include: absPathInclude},
			wantErr: false,
		},
		{
			name: "relative path include",
			fields: fields{
				mainConfigAbsPath: simpleMain.Value(),
				configGraph:       simpleMain.ConfigGraph,
				contextStack:      newContextStack(),
			},
			args:    args{include: relPathInclude},
			wantErr: false,
		},
		{
			name: "no match absolute path include",
			fields: fields{
				mainConfigAbsPath: simpleMain.Value(),
				configGraph:       simpleMain.ConfigGraph,
				contextStack:      newContextStack(),
			},
			args:    args{include: NewContext(context_type.TypeInclude, noMatchAbsPath).(*Include)},
			wantErr: false,
		},
		{
			name: "no match relative path include",
			fields: fields{
				mainConfigAbsPath: simpleMain.Value(),
				configGraph:       simpleMain.ConfigGraph,
				contextStack:      newContextStack(),
			},
			args:    args{include: NewContext(context_type.TypeInclude, noMatchRelPath).(*Include)},
			wantErr: false,
		},
		{
			name: "cycle include",
			fields: fields{
				mainConfigAbsPath: simpleMain.Value(),
				configGraph:       simpleMain.ConfigGraph,
				contextStack:      newContextStack(),
			},
			args:    args{include: cycleInclude},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileLoader{
				mainConfigAbsPath: tt.fields.mainConfigAbsPath,
				configGraph:       tt.fields.configGraph,
				contextStack:      tt.fields.contextStack,
			}
			if err := f.loadInclude(tt.args.include); (err != nil) != tt.wantErr {
				t.Errorf("loadInclude() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_jsonLoader_Load(t *testing.T) {
	type fields struct {
		unmarshaler *mainUnmarshaler
		jsonBytes   []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Main
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonLoader{
				unmarshaler: tt.fields.unmarshaler,
				jsonBytes:   tt.fields.jsonBytes,
			}
			got, err := j.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newContextStack(t *testing.T) {
	tests := []struct {
		name string
		want *contextStack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newContextStack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newContextStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBlankLine(t *testing.T) {
	type args struct {
		data []byte
		idx  *int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBlankLine(tt.args.data, tt.args.idx); got != tt.want {
				t.Errorf("parseBlankLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBraceEnd(t *testing.T) {
	type args struct {
		data []byte
		idx  *int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBraceEnd(tt.args.data, tt.args.idx); got != tt.want {
				t.Errorf("parseBraceEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseErrLine(t *testing.T) {
	type args struct {
		data   []byte
		idx    *int
		config *Config
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
			if err := parseErrLine(tt.args.data, tt.args.idx, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("parseErrLine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
