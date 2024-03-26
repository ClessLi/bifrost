package configuration

import (
	"bytes"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"reflect"
	"sync"
	"testing"
)

func TestNewNginxConfigFromJsonBytes(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    NginxConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNginxConfigFromJsonBytes(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNginxConfigFromJsonBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNginxConfigFromJsonBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNginxConfigFromPath(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		args    args
		want    NginxConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNginxConfigFromFS(tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNginxConfigFromFS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNginxConfigFromFS() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dumpMainContext(t *testing.T) {
	type args struct {
		m *local.Main
	}
	tests := []struct {
		name string
		args args
		want map[string]*bytes.Buffer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dumpMainContext(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dumpMainContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newNginxConfig(t *testing.T) {
	type args struct {
		maincontext *local.Main
	}
	tests := []struct {
		name    string
		args    args
		want    NginxConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newNginxConfig(tt.args.maincontext)
			if (err != nil) != tt.wantErr {
				t.Errorf("newNginxConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newNginxConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nginxConfig_Dump(t *testing.T) {
	type fields struct {
		mainContext *local.Main
		rwLocker    *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*bytes.Buffer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &nginxConfig{
				mainContext: tt.fields.mainContext,
				rwLocker:    tt.fields.rwLocker,
			}
			if got := n.Dump(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nginxConfig_Json(t *testing.T) {
	type fields struct {
		mainContext *local.Main
		rwLocker    *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &nginxConfig{
				mainContext: tt.fields.mainContext,
				rwLocker:    tt.fields.rwLocker,
			}
			if got := n.Json(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nginxConfig_Main(t *testing.T) {
	type fields struct {
		mainContext *local.Main
		rwLocker    *sync.RWMutex
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
			n := &nginxConfig{
				mainContext: tt.fields.mainContext,
				rwLocker:    tt.fields.rwLocker,
			}
			if got := n.Main(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Main() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nginxConfig_TextLines(t *testing.T) {
	type fields struct {
		mainContext *local.Main
		rwLocker    *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &nginxConfig{
				mainContext: tt.fields.mainContext,
				rwLocker:    tt.fields.rwLocker,
			}
			if got := n.TextLines(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TextLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nginxConfig_UpdateFromJsonBytes(t *testing.T) {
	type fields struct {
		mainContext *local.Main
		rwLocker    *sync.RWMutex
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal test",
			fields: fields{
				mainContext: local.NewContext(context_type.TypeMain, "C:\\test\\test.conf").(*local.Main),
				rwLocker:    new(sync.RWMutex),
			},
			args: args{data: []byte(
				`{
    "main-config": "C:\\test\\test.conf",
    "configs":
    {
        "C:\\test\\test.conf":
        [
            {
                "http":
                {
                    "params":
                    [
                        {
                            "inline": true, "comments": "test comment"
                        },
                        {
                            "server":
                            {
                                "params":
                                [
                                    {
                                        "directive": "server_name", "params": "testserver"
                                    },
                                    {
                                        "location": {"value": "~ /test"}
                                    },
                                    {
                                        "include":
                                        {
                                            "value": "conf.d\\include*conf",
                                            "params": ["conf.d\\include.location1.conf", "conf.d\\include.location2.conf"]
                                        }
                                    }
                                ]
                            }
                        }
                    ]
                }
            }
        ],
        "conf.d\\include.location1.conf":
        [
            {
                "location": {"value": "~ /test1"}
            }
        ],
        "conf.d\\include.location2.conf":
        [
            {
                "location": {"value": "^~ /test2"}
            }
        ]
    }
}`,
			)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &nginxConfig{
				mainContext: tt.fields.mainContext,
				rwLocker:    tt.fields.rwLocker,
			}
			if err := n.UpdateFromJsonBytes(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFromJsonBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_nginxConfig_renewMainContext(t *testing.T) {
	type fields struct {
		mainContext *local.Main
		rwLocker    *sync.RWMutex
	}
	type args struct {
		m *local.Main
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
			n := &nginxConfig{
				mainContext: tt.fields.mainContext,
				rwLocker:    tt.fields.rwLocker,
			}
			if err := n.renewMainContext(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("renewMainContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
