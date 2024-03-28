package local

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
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
