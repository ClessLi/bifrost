package local

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/dominikbraun/graph"
	"reflect"
	"testing"
)

func TestConfig_Clone(t *testing.T) {
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
						),
					1,
				),
			0,
		).(*Main)
	testCloneChildren := make([]context.Context, 0)
	for _, child := range testMain.Children {
		testCloneChildren = append(testCloneChildren, child.Clone())
	}
	type fields struct {
		BasicContext BasicContext
		Graph        ConfigGraph
		ConfigPath   context.ConfigPath
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name: "clone without graph and ConfigPath",
			fields: fields{
				BasicContext: testMain.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   testMain.ConfigPath,
			},
			want: &Config{
				BasicContext: BasicContext{
					ContextType:    context_type.TypeConfig,
					ContextValue:   testMain.ContextValue,
					Children:       testCloneChildren,
					father:         testMain.father,
					self:           context.NullContext(),
					headStringFunc: testMain.headStringFunc,
					tailStringFunc: testMain.tailStringFunc,
				},
				Graph:      nil,
				ConfigPath: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				Graph:        tt.fields.Graph,
				ConfigPath:   tt.fields.ConfigPath,
			}
			got := c.Clone()
			gotlines, err := got.ConfigLines(false)
			if err != nil {
				t.Errorf("got.ConfigLines() return error: %v", err)
				return
			}
			wantlines, err := tt.want.ConfigLines(false)
			if err != nil {
				t.Errorf("want.ConfigLines() return error: %v", err)
				return
			}
			if reflect.DeepEqual(got.(*Config).Children, tt.want.(*Config).Children) ||
				!reflect.DeepEqual(gotlines, wantlines) ||
				!reflect.DeepEqual(got.Father(), tt.want.Father()) ||
				got.Type() != tt.want.Type() ||
				got.Value() != tt.want.Value() ||
				!reflect.DeepEqual(got.(*Config).Graph, tt.want.(*Config).Graph) ||
				!reflect.DeepEqual(got.(*Config).ConfigPath, tt.want.(*Config).ConfigPath) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_ConfigLines(t *testing.T) {
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
						),
					1,
				),
			0,
		).(*Main)
	type fields struct {
		BasicContext BasicContext
		Graph        ConfigGraph
		ConfigPath   context.ConfigPath
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
			name: "normal test",
			fields: fields{
				BasicContext: testMain.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   testMain.ConfigPath,
			},
			args: args{isDumping: false},
			want: []string{
				"http {    # test comment",
				"    server {",
				"        server_name testserver;",
				"        location ~ /test {",
				"        }",
				"    }",
				"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				Graph:        tt.fields.Graph,
				ConfigPath:   tt.fields.ConfigPath,
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

func TestConfig_Father(t *testing.T) {
	type fields struct {
		BasicContext BasicContext
		Graph        ConfigGraph
		ConfigPath   context.ConfigPath
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name: "null context father",
			want: context.NullContext(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				Graph:        tt.fields.Graph,
				ConfigPath:   tt.fields.ConfigPath,
			}
			if got := c.Father(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Father() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IncludeConfig(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	_, err := testMain.IncludeConfig(
		NewContext(context_type.TypeConfig, "C:\\test\\existing.conf").(*Config),
		NewContext(context_type.TypeConfig, "a.conf").(*Config),
		NewContext(context_type.TypeConfig, "b.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	existingConfig, err := testMain.Graph.GetConfig("C:\\test\\existing.conf")
	if err != nil {
		t.Fatal(err)
	}
	nullPathConfig := NewContext(context_type.TypeConfig, "").(*Config)
	nullPathConfig.ConfigPath = &context.AbsConfigPath{}
	nullPathConfig.Graph = testMain.Graph
	err = nullPathConfig.Graph.(*configGraph).graph.AddVertex(nullPathConfig)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.Graph.(*configGraph).graph.AddEdge("", existingConfig.FullPath())
	if err != nil {
		t.Fatal(err)
	}
	// for cycles include
	aConfig, err := testMain.Graph.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	bConfig, err := testMain.Graph.GetConfig("C:\\test\\b.conf")
	if err != nil {
		t.Fatal(err)
	}
	cConfig := NewContext(context_type.TypeConfig, "c.conf").(*Config)
	_, err = aConfig.IncludeConfig(bConfig)
	if err != nil {
		t.Fatal(err)
	}

	// different main config
	diffTestMain := NewContext(context_type.TypeMain, "C:\\test2\\nginx.conf").(*Main)
	// main with invalid value
	invalidTestMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").(*Main)
	invalidTestMainCP := invalidTestMain.ConfigPath
	invalidTestMain.ConfigPath = nil
	newcpFailedConfig := NewContext(context_type.TypeConfig, "test\\test2.conf").(*Config)
	// different graph config
	diffGraphConfig := NewContext(context_type.TypeConfig, "C:\\test\\test2.conf").(*Config)
	diffGraphConfig.Graph = diffTestMain.Graph
	diffGraphConfigPath, _ := newConfigPath(testMain.Graph, diffGraphConfig.Value())

	// different base dir config
	diffBaseDirConfig := NewContext(context_type.TypeConfig, "C:\\test2\\test.conf").(*Config)
	diffBaseDirConfig.ConfigPath, _ = context.NewRelConfigPath("C:\\test2", "test.conf")
	err = testMain.Graph.AddConfig(diffBaseDirConfig)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		BasicContext BasicContext
		Graph        ConfigGraph
		ConfigPath   context.ConfigPath
	}
	type args struct {
		configs []*Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Config
		wantErr bool
	}{
		{
			name: "has not been added to a graph",
			fields: fields{
				BasicContext: BasicContext{},
				Graph:        nil,
				ConfigPath:   testMain.ConfigPath,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "has not been banded with a config path",
			fields: fields{
				BasicContext: BasicContext{},
				Graph:        testMain.Graph,
				ConfigPath:   nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "insert nil",
			fields: fields{
				BasicContext: testMain.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   testMain.ConfigPath,
			},
			args:    args{nil},
			want:    nil,
			wantErr: true,
		},
		{
			name: "insert empty config list",
			fields: fields{
				BasicContext: testMain.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   testMain.ConfigPath,
			},
			args:    args{make([]*Config, 0)},
			want:    make([]*Config, 0),
			wantErr: false,
		},
		{
			name: "insert nil config",
			fields: fields{
				BasicContext: testMain.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   testMain.ConfigPath,
			},
			args:    args{configs: []*Config{nil}},
			want:    make([]*Config, 0),
			wantErr: true,
		},
		{
			name: "insert another graph config",
			fields: fields{
				BasicContext: testMain.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   testMain.ConfigPath,
			},
			args: args{configs: []*Config{diffGraphConfig}},
			want: []*Config{{
				BasicContext: diffGraphConfig.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   diffGraphConfigPath,
			}},
			wantErr: false,
		},
		{
			name: "failed to build config path for included config",
			fields: fields{
				BasicContext: invalidTestMain.BasicContext,
				Graph:        invalidTestMain.Graph,
				ConfigPath:   invalidTestMainCP,
			},
			args:    args{configs: []*Config{newcpFailedConfig}},
			want:    make([]*Config, 0),
			wantErr: true,
		},
		{
			name: "different base dir config path",
			fields: fields{
				BasicContext: testMain.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   testMain.ConfigPath,
			},
			args:    args{configs: []*Config{diffBaseDirConfig}},
			want:    make([]*Config, 0),
			wantErr: true,
		},
		{
			name: "include an existing config, but not the same config pointer",
			fields: fields{
				BasicContext: testMain.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   testMain.ConfigPath,
			},
			args:    args{configs: []*Config{NewContext(context_type.TypeConfig, "existing.conf").(*Config)}},
			want:    make([]*Config, 0),
			wantErr: true,
		},
		{
			name: "add config edge error",
			fields: fields{
				BasicContext: nullPathConfig.BasicContext,
				Graph:        nullPathConfig.Graph,
				ConfigPath:   nullPathConfig.ConfigPath,
			},
			args:    args{configs: []*Config{existingConfig}},
			want:    make([]*Config, 0),
			wantErr: true,
		},
		{
			name: "cycles include",
			fields: fields{
				BasicContext: bConfig.BasicContext,
				Graph:        bConfig.Graph,
				ConfigPath:   bConfig.ConfigPath,
			},
			args:    args{configs: []*Config{aConfig}},
			want:    make([]*Config, 0),
			wantErr: true,
		},
		{
			name: "include valid and invalid configs",
			fields: fields{
				BasicContext: bConfig.BasicContext,
				Graph:        bConfig.Graph,
				ConfigPath:   bConfig.ConfigPath,
			},
			args:    args{configs: []*Config{aConfig, existingConfig}},
			want:    []*Config{existingConfig},
			wantErr: true,
		},
		{
			name: "include valid configs",
			fields: fields{
				BasicContext: testMain.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   testMain.ConfigPath,
			},
			args:    args{configs: []*Config{aConfig, bConfig, cConfig}},
			want:    []*Config{aConfig, bConfig, cConfig},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				Graph:        tt.fields.Graph,
				ConfigPath:   tt.fields.ConfigPath,
			}
			isSame := func(got, want []*Config) bool {
				if len(got) != len(want) {
					return false
				}
				for i := range got {
					if wantCache, _ := c.Graph.GetConfig(want[i].FullPath()); got[i] != want[i] && got[i] != wantCache {
						return false
					}
				}
				return true
			}
			got, err := c.IncludeConfig(tt.args.configs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncludeConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !isSame(got, tt.want) {
				t.Errorf("IncludeConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetFather(t *testing.T) {
	type fields struct {
		BasicContext BasicContext
		Graph        ConfigGraph
		ConfigPath   context.ConfigPath
	}
	type args struct {
		in0 context.Context
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
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				Graph:        tt.fields.Graph,
				ConfigPath:   tt.fields.ConfigPath,
			}
			if err := c.SetFather(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("SetFather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_SetValue(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	_, err := testMain.IncludeConfig(NewContext(context_type.TypeConfig, "a.conf").(*Config))
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.Graph.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		BasicContext BasicContext
		Graph        ConfigGraph
		ConfigPath   context.ConfigPath
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
			name: "failed to set value",
			fields: fields{
				BasicContext: testMain.BasicContext,
				Graph:        testMain.Graph,
				ConfigPath:   testMain.ConfigPath,
			},
			args:    args{v: "C:\\test\\a.conf"},
			wantErr: true,
		},
		{
			name: "set value",
			fields: fields{
				BasicContext: aConfig.BasicContext,
				Graph:        aConfig.Graph,
				ConfigPath:   aConfig.ConfigPath,
			},
			args:    args{v: "b.conf"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				Graph:        tt.fields.Graph,
				ConfigPath:   tt.fields.ConfigPath,
			}
			cache, err := c.Graph.GetConfig(c.FullPath())
			if err != nil {
				t.Fatal(err)
			}
			if err := cache.SetValue(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_isInGraph(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	nilGraphConfig := NewContext(context_type.TypeConfig, "nilgraph").(*Config)
	nilGraphConfig.ConfigPath, _ = newConfigPath(testMain.Graph, nilGraphConfig.Value())

	nilConfigPathConfig := NewContext(context_type.TypeConfig, "nilpath").(*Config)
	nilConfigPathConfig.Graph = testMain.Graph

	notInGraphConfig := NewContext(context_type.TypeConfig, "notingraph").(*Config)
	notInGraphConfig.ConfigPath, _ = newConfigPath(testMain.Graph, notInGraphConfig.Value())
	notInGraphConfig.Graph = testMain.Graph

	inGraphConfig := NewContext(context_type.TypeConfig, "ingraph").(*Config)
	inGraphConfig.ConfigPath, _ = newConfigPath(testMain.Graph, inGraphConfig.Value())
	inGraphConfig.Graph = testMain.Graph
	err := inGraphConfig.Graph.AddConfig(inGraphConfig)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		BasicContext BasicContext
		Graph        ConfigGraph
		ConfigPath   context.ConfigPath
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "config has nil graph",
			fields: fields{
				BasicContext: nilGraphConfig.BasicContext,
				Graph:        nilGraphConfig.Graph,
				ConfigPath:   nilGraphConfig.ConfigPath,
			},
			want: false,
		},
		{
			name: "config has nil config path",
			fields: fields{
				BasicContext: nilConfigPathConfig.BasicContext,
				Graph:        nilConfigPathConfig.Graph,
				ConfigPath:   nilConfigPathConfig.ConfigPath,
			},
			want: false,
		},
		{
			name: "config has not been added into a graph",
			fields: fields{
				BasicContext: notInGraphConfig.BasicContext,
				Graph:        notInGraphConfig.Graph,
				ConfigPath:   notInGraphConfig.ConfigPath,
			},
			want: false,
		},
		{
			name: "config has been added into a graph",
			fields: fields{
				BasicContext: inGraphConfig.BasicContext,
				Graph:        inGraphConfig.Graph,
				ConfigPath:   inGraphConfig.ConfigPath,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BasicContext: tt.fields.BasicContext,
				Graph:        tt.fields.Graph,
				ConfigPath:   tt.fields.ConfigPath,
			}
			if got := c.isInGraph(); got != tt.want {
				t.Errorf("isInGraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_modifyPathInGraph(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	_, err := testMain.IncludeConfig(
		NewContext(context_type.TypeConfig, "a.conf").(*Config),
		NewContext(context_type.TypeConfig, "b.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.Graph.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	bConfig, err := testMain.Graph.GetConfig("C:\\test\\b.conf")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  *Config
		args    args
		wantErr bool
	}{
		{
			name:    "config has not been added into a graph",
			fields:  &Config{BasicContext: aConfig.BasicContext},
			args:    args{path: "b.conf"},
			wantErr: false,
		},
		{
			name:    "modify to null string path",
			fields:  aConfig,
			args:    args{path: ""},
			wantErr: true,
		},
		{
			name:    "modify to same path",
			fields:  aConfig,
			args:    args{path: aConfig.Value()},
			wantErr: false,
		},
		{
			name:    "modify to another path already exist in graph",
			fields:  aConfig,
			args:    args{path: bConfig.Value()},
			wantErr: true,
		},
		{
			name: "modify same config, but different from cache in the graph",
			fields: &Config{
				BasicContext: aConfig.BasicContext,
				Graph:        aConfig.Graph,
				ConfigPath:   aConfig.ConfigPath,
			},
			args:    args{path: "c.conf"},
			wantErr: false,
		},
		{
			name:    "modify config path",
			fields:  aConfig,
			args:    args{path: "c.conf"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.modifyPathInGraph(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("modifyPathInGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_AddConfig(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	_, err := testMain.IncludeConfig(NewContext(context_type.TypeConfig, "a.conf").(*Config))
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.Graph.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	bConfig := NewContext(context_type.TypeConfig, "b.conf").(*Config)
	bConfig.Graph = testMain.Graph
	bConfig.ConfigPath, _ = newConfigPath(testMain.Graph, bConfig.Value())
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		wantErr bool
	}{
		{
			name:    "add nil config",
			fields:  testMain.Graph,
			args:    args{config: nil},
			wantErr: true,
		},
		{
			name:    "add an already exist config",
			fields:  testMain.Graph,
			args:    args{config: aConfig},
			wantErr: true,
		},
		{
			name:    "normal test",
			fields:  testMain.Graph,
			args:    args{config: bConfig},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.AddConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("AddConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_AddEdge(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	_, err := testMain.IncludeConfig(
		NewContext(context_type.TypeConfig, "a.conf").(*Config),
		NewContext(context_type.TypeConfig, "b.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.Graph.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	bConfig, err := testMain.Graph.GetConfig("C:\\test\\b.conf")
	if err != nil {
		t.Fatal(err)
	}

	otherMain := NewContext(context_type.TypeMain, "C:\\test1\\nginx.conf").(*Main)
	_, err = otherMain.IncludeConfig(
		NewContext(context_type.TypeConfig, "other.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	inOtherGraphConfig, err := otherMain.Graph.GetConfig("C:\\test1\\other.conf")
	if err != nil {
		t.Fatal(err)
	}

	nullpathConfig := NewContext(context_type.TypeConfig, "").(*Config)
	nullpathConfig.ConfigPath = &context.AbsConfigPath{}

	excludeConfig := NewContext(context_type.TypeConfig, "exclude.conf").(*Config)
	excludeConfig.ConfigPath, _ = newConfigPath(testMain.Graph, excludeConfig.Value())

	type args struct {
		src *Config
		dst *Config
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		wantErr bool
	}{
		{
			name:   "nil source config",
			fields: testMain.Graph,
			args: args{
				src: nil,
				dst: aConfig,
			},
			wantErr: true,
		},
		{
			name:   "nil destination config",
			fields: testMain.Graph,
			args: args{
				src: aConfig,
				dst: nil,
			},
			wantErr: true,
		},
		{
			name:   "source config with null config path",
			fields: testMain.Graph,
			args: args{
				src: nullpathConfig,
				dst: aConfig,
			},
			wantErr: true,
		},
		{
			name:   "destination config with null config path",
			fields: testMain.Graph,
			args: args{
				src: aConfig,
				dst: nullpathConfig,
			},
			wantErr: true,
		},
		{
			name:   "source config is exclude from the graph",
			fields: testMain.Graph,
			args: args{
				src: excludeConfig,
				dst: aConfig,
			},
			wantErr: true,
		},
		{
			name:   "destination config is exclude from the graph",
			fields: testMain.Graph,
			args: args{
				src: aConfig,
				dst: excludeConfig,
			},
			wantErr: true,
		},
		{
			name:   "source config is in the other graph",
			fields: testMain.Graph,
			args: args{
				src: inOtherGraphConfig,
				dst: aConfig,
			},
			wantErr: true,
		},
		{
			name:   "destination config is in the other graph",
			fields: testMain.Graph,
			args: args{
				src: aConfig,
				dst: inOtherGraphConfig,
			},
			wantErr: true,
		},
		{
			name:   "normal test",
			fields: testMain.Graph,
			args: args{
				src: aConfig,
				dst: bConfig,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.AddEdge(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("AddEdge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_GetConfig(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	_, err := testMain.IncludeConfig(
		NewContext(context_type.TypeConfig, "a.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.Graph.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		fullpath string
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name:    "wrong config path",
			fields:  testMain.Graph,
			args:    args{fullpath: "wrong/config/path.conf"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "normal test",
			fields:  testMain.Graph,
			args:    args{fullpath: "C:\\test\\a.conf"},
			want:    aConfig,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			got, err := c.GetConfig(tt.args.fullpath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configGraph_MainConfig(t *testing.T) {
	type fields struct {
		graph      graph.Graph[string, *Config]
		mainConfig *Config
	}
	tests := []struct {
		name   string
		fields fields
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &configGraph{
				graph:      tt.fields.graph,
				mainConfig: tt.fields.mainConfig,
			}
			if got := c.MainConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MainConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configGraph_RemoveEdge(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	testMain.IncludeConfig(NewContext(context_type.TypeConfig, "a.conf").(*Config))
	a, _ := testMain.Graph.GetConfig("C:\\test\\a.conf")
	a.IncludeConfig(NewContext(context_type.TypeConfig, "b.conf").(*Config))
	b, _ := testMain.Graph.GetConfig("C:\\test\\b.conf")
	b.IncludeConfig(NewContext(context_type.TypeConfig, "c.conf").(*Config))
	c, _ := testMain.Graph.GetConfig("C:\\test\\c.conf")
	notInGraphConfig := NewContext(context_type.TypeConfig, "notingraph.conf").(*Config)
	type args struct {
		src *Config
		dst *Config
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		wantErr bool
	}{
		{
			name:   "removed edge not found",
			fields: testMain.Graph,
			args: args{
				src: b,
				dst: notInGraphConfig,
			},
			wantErr: true,
		},
		{
			name:   "destination has edge",
			fields: testMain.Graph,
			args: args{
				src: a,
				dst: b,
			},
			wantErr: false,
		},
		{
			name:   "remove edge and destination",
			fields: testMain.Graph,
			args: args{
				src: b,
				dst: c,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.RemoveEdge(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("RemoveEdge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_RenewConfigPath(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	_, err := testMain.IncludeConfig(
		NewContext(context_type.TypeConfig, "a.conf").(*Config),
		NewContext(context_type.TypeConfig, "2exist.conf").(*Config),
		NewContext(context_type.TypeConfig, "inedge.conf").(*Config),
		NewContext(context_type.TypeConfig, "test.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	aConfig, err := testMain.Graph.GetConfig("C:\\test\\a.conf")
	if err != nil {
		t.Fatal(err)
	}
	renew2existConfig, err := testMain.Graph.GetConfig("C:\\test\\2exist.conf")
	if err != nil {
		t.Fatal(err)
	}
	renew2existConfig.ConfigPath = aConfig.ConfigPath

	inEdgeConfig, err := testMain.Graph.GetConfig("C:\\test\\inedge.conf")
	if err != nil {
		t.Fatal(err)
	}

	testConfig, err := testMain.Graph.GetConfig("C:\\test\\test.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = inEdgeConfig.IncludeConfig(testConfig)
	_, _ = testConfig.IncludeConfig(NewContext(context_type.TypeConfig, "outedge.conf").(*Config))
	testConfig.ConfigPath, err = newConfigPath(testConfig.Graph, "modified.conf")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		fullpath string
	}
	tests := []struct {
		name    string
		fields  ConfigGraph
		args    args
		wantErr bool
	}{
		{
			name:    "not exist config",
			fields:  testMain.Graph,
			args:    args{fullpath: "notexist.conf"},
			wantErr: true,
		},
		{
			name:    "need not renew config",
			fields:  testMain.Graph,
			args:    args{fullpath: testMain.FullPath()},
			wantErr: false,
		},
		{
			name:    "renew to exist config",
			fields:  testMain.Graph,
			args:    args{fullpath: "C:\\test\\2exist.conf"},
			wantErr: true,
		},
		{
			name:    "normal test",
			fields:  testMain.Graph,
			args:    args{fullpath: "C:\\test\\test.conf"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.RenewConfigPath(tt.args.fullpath); (err != nil) != tt.wantErr {
				t.Errorf("RenewConfigPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_Topology(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	testMain.IncludeConfig(NewContext(context_type.TypeConfig, "a.conf").(*Config))
	a, _ := testMain.Graph.GetConfig("C:\\test\\a.conf")
	a.IncludeConfig(NewContext(context_type.TypeConfig, "b.conf").(*Config))
	b, _ := testMain.Graph.GetConfig("C:\\test\\b.conf")
	b.IncludeConfig(NewContext(context_type.TypeConfig, "c.conf").(*Config))
	c, _ := testMain.Graph.GetConfig("C:\\test\\c.conf")
	c.IncludeConfig(NewContext(context_type.TypeConfig, "d.conf").(*Config))
	d, _ := testMain.Graph.GetConfig("C:\\test\\d.conf")
	e := NewContext(context_type.TypeConfig, "e.conf").(*Config)
	e.ConfigPath, _ = newConfigPath(testMain.Graph, e.Value())
	err := testMain.Graph.AddConfig(e)
	if err != nil {
		t.Fatal(err)
	}
	_, err = e.IncludeConfig(
		a,
		NewContext(context_type.TypeConfig, "f.conf").(*Config),
		NewContext(context_type.TypeConfig, "g.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fields ConfigGraph
		want   []*Config
	}{
		{
			name:   "generate only one tree starting from the main config",
			fields: testMain.Graph,
			want:   []*Config{testMain.Config, a, b, c, d},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if got := c.Topology(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Topology() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configGraph_removeConfig(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	testMain.IncludeConfig(NewContext(context_type.TypeConfig, "a.conf").(*Config))
	a, _ := testMain.Graph.GetConfig("C:\\test\\a.conf")
	a.IncludeConfig(NewContext(context_type.TypeConfig, "b.conf").(*Config))
	b, _ := testMain.Graph.GetConfig("C:\\test\\b.conf")
	b.IncludeConfig(NewContext(context_type.TypeConfig, "c.conf").(*Config))
	c, _ := testMain.Graph.GetConfig("C:\\test\\c.conf")
	c.IncludeConfig(NewContext(context_type.TypeConfig, "d.conf").(*Config))
	d, _ := testMain.Graph.GetConfig("C:\\test\\d.conf")
	err := testMain.Graph.(*configGraph).graph.RemoveEdge(configHash(c), configHash(d))
	if err != nil {
		t.Fatal(err)
	}
	e := NewContext(context_type.TypeConfig, "e.conf").(*Config)
	e.ConfigPath, _ = newConfigPath(testMain.Graph, e.Value())
	err = testMain.Graph.AddConfig(e)
	if err != nil {
		t.Fatal(err)
	}
	_, err = e.IncludeConfig(
		a,
		NewContext(context_type.TypeConfig, "f.conf").(*Config),
		NewContext(context_type.TypeConfig, "g.conf").(*Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		fields  *configGraph
		args    args
		wantErr bool
	}{
		{
			name:    "normal test",
			fields:  testMain.Graph.(*configGraph),
			args:    args{d},
			wantErr: false,
		},
		{
			name:    "config has no in edge but out edges",
			fields:  testMain.Graph.(*configGraph),
			args:    args{e},
			wantErr: true,
		},
		{
			name:    "config has edges",
			fields:  testMain.Graph.(*configGraph),
			args:    args{a},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.removeConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("removeConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configGraph_setGraphFor(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	testMain.IncludeConfig(NewContext(context_type.TypeConfig, "a.conf").(*Config))
	a, _ := testMain.Graph.GetConfig("C:\\test\\a.conf")
	diffGraphConfig := NewContext(context_type.TypeConfig, "different_graph.conf").(*Config)
	diffGraphConfig.Graph = &configGraph{}
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		fields  *configGraph
		args    args
		wantErr bool
	}{
		{
			name:    "nil config",
			fields:  testMain.Graph.(*configGraph),
			wantErr: true,
		},
		{
			name:    "config in the other graph",
			fields:  testMain.Graph.(*configGraph),
			args:    args{diffGraphConfig},
			wantErr: true,
		},
		{
			name:    "config clone",
			fields:  testMain.Graph.(*configGraph),
			args:    args{a.Clone().(*Config)},
			wantErr: false,
		},
		{
			name:    "same graph config",
			fields:  testMain.Graph.(*configGraph),
			args:    args{a},
			wantErr: false,
		},
		{
			name:    "new config",
			fields:  testMain.Graph.(*configGraph),
			args:    args{NewContext(context_type.TypeConfig, "new.conf").(*Config)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.setGraphFor(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("setGraphFor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configHash(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf   ").
		Insert(
			NewContext(context_type.TypeHttp, "").
				Insert(NewComment("test comment", true), 0).
				Insert(
					NewContext(context_type.TypeServer, "").
						Insert(NewDirective("server_name", "testserver"), 0).
						Insert(
							NewContext(context_type.TypeLocation, "~ /test"),
							1,
						),
					1,
				),
			0,
		).(*Main)
	type args struct {
		t *Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil config",
			want: "",
		},
		{
			name: "new config with no config path",
			args: args{t: NewContext(context_type.TypeConfig, "new.conf").(*Config)},
			want: "",
		},
		{
			name: "normal test",
			args: args{t: testMain.Config},
			want: "C:\\test\\nginx.conf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := configHash(tt.args.t); got != tt.want {
				t.Errorf("configHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newConfigPath(t *testing.T) {
	// main config
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").(*Main)
	relPathMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").(*Main)
	relPathMain.ConfigPath = &context.RelConfigPath{}
	absCP, _ := context.NewAbsConfigPath("D:\\test\\test.conf")
	relCP, _ := context.NewRelConfigPath(testMain.ConfigPath.BaseDir(), "relative.conf")
	type args struct {
		configgraph   ConfigGraph
		newconfigpath string
	}
	tests := []struct {
		name    string
		args    args
		want    context.ConfigPath
		wantErr bool
	}{
		{
			name: "nil graph",
			args: args{
				configgraph:   nil,
				newconfigpath: "test",
			},
			wantErr: true,
		},
		{
			name: "null new config path",
			args: args{
				configgraph:   relPathMain.Graph,
				newconfigpath: "",
			},
			wantErr: true,
		},
		{
			name: "config paths are relative path in graph's main config and input",
			args: args{
				configgraph:   relPathMain.Graph,
				newconfigpath: "test.config",
			},
			wantErr: true,
		},
		{
			name: "absolut path",
			args: args{
				configgraph:   testMain.Graph,
				newconfigpath: "D:\\test\\test.conf",
			},
			want: absCP,
		},
		{
			name: "relative path",
			args: args{
				configgraph:   testMain.Graph,
				newconfigpath: "relative.conf",
			},
			want: relCP,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfigPath(tt.args.configgraph, tt.args.newconfigpath)
			if (err != nil) != tt.wantErr {
				t.Errorf("newConfigPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConfigPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newMainConfig(t *testing.T) {
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
	err = setGraphForMainConfig(testMainConfig)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		abspath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
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
			want:    testMainConfig,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newMainConfig(tt.args.abspath)
			if (err != nil) != tt.wantErr {
				t.Errorf("newMainConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != (tt.want == nil) ||
				(got != nil &&
					(!reflect.DeepEqual(got.ConfigPath, tt.want.ConfigPath) ||
						!reflect.DeepEqual(got.BasicContext.ContextValue, tt.want.BasicContext.ContextValue) ||
						!reflect.DeepEqual(got.BasicContext.ContextType, tt.want.BasicContext.ContextType) ||
						!reflect.DeepEqual(got.BasicContext.Children, tt.want.BasicContext.Children) ||
						!reflect.DeepEqual(got.BasicContext.father, tt.want.BasicContext.father))) {
				t.Errorf("newMainConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_registerConfigBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerConfigBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerConfigBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setGraphForMainConfig(t *testing.T) {
	testMain := NewContext(context_type.TypeMain, "C:\\test\\nginx.conf").(*Main)
	cloneMain := testMain.Clone().(*Main)
	cloneMain.ConfigPath = nil
	cloneMain.Graph = &configGraph{mainConfig: cloneMain.Config}
	newMainConf := NewContext(context_type.TypeConfig, "C:\\test\\newmain.conf").(*Config)
	type args struct {
		mainConfig *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "set graph for main config already exist",
			args:    args{mainConfig: testMain.Config},
			wantErr: true,
		},
		{
			name:    "main config is in another graph",
			args:    args{cloneMain.Config},
			wantErr: true,
		},
		{
			name:    "normal test",
			args:    args{mainConfig: newMainConf},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setGraphForMainConfig(tt.args.mainConfig); (err != nil) != tt.wantErr {
				t.Errorf("setGraphForMainConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
