package local

import (
	"reflect"
	"testing"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
)

func TestPosBasedOnConfig(t *testing.T) {
	testMain, err := fakeProxyPassTestMainCtx()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.ContextPos
		wantErr bool
	}{
		{
			name: "test a first level context",
			args: args{
				ctx: testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeHttp)).Target(),
			},
			want: &v1.ContextPos{
				ConfigPath: testMain.MainConfig().FullPath(),
				PosIndex:   []int32{0},
			},
		},
		{
			name: "test a deeper level context",
			args: args{
				ctx: testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("error_pos_stream_proxy_pass_2")).Target(),
			},
			want: &v1.ContextPos{
				ConfigPath: testMain.MainConfig().FullPath(),
				PosIndex:   []int32{0, 3, 7, 0},
			},
		},
		{
			name: "test with main config",
			args: args{
				ctx: testMain,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PosBasedOnConfig(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PosBasedOnConfig() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PosBasedOnConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
