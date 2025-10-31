package local

import (
	"encoding/json"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"

	"github.com/marmotedu/errors"
)

func fakeHostsResolver() utilsV3.DomainNameResolver {
	return utilsV3.NewIPv4Hosts(map[string][]net.IP{
		"baidu.com":   {net.IPv4(10, 1, 11, 111), net.IPv4(10, 1, 12, 122)},
		"example.com": {net.IPv4(10, 1, 12, 122), net.IPv4(10, 1, 13, 133)},
		"test.cn":     {net.IPv4(10, 1, 12, 122), net.IPv4(10, 1, 13, 133)},
		"abc.com":     {net.IPv4(10, 2, 12, 122)},
	})
}

//nolint:funlen
func fakeProxyPassTestMainCtx() (MainContext, error) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		return testMain, err
	}
	testMain.Insert(
		NewContext(context_type.TypeHttp, "").
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeLocation, "~ /baidu").
							Insert(
								NewContext(context_type.TypeDirHTTPProxyPass, "https://baidu.com"),
								0,
							).
							Insert(
								NewContext(context_type.TypeLimitExcept, "POST").
									Insert(
										NewContext(context_type.TypeDirHTTPProxyPass, "http://abc.com"),
										0,
									),
								1,
							),
						0,
					).
					Insert(
						NewContext(context_type.TypeLocation, "~ /example").
							Insert(
								NewContext(context_type.TypeIf, "( $http_gray ~* test )").
									Insert(
										NewContext(context_type.TypeDirHTTPProxyPass, "http://test.cn"),
										0,
									),
								0,
							).
							Insert(
								NewContext(context_type.TypeDirHTTPProxyPass, "https://example.com"),
								1,
							),
						1,
					).
					Insert(
						NewContext(context_type.TypeLocation, "~ /upstream").
							Insert(
								NewContext(context_type.TypeDirHTTPProxyPass, "http://upstream_server/proxy/to/upstreamserver"),
								0,
							),
						1,
					).
					Insert(
						NewContext(context_type.TypeLocation, "~ /upstream2").
							Insert(
								NewContext(context_type.TypeDirHTTPProxyPass, "http://upstream_server:8080/proxy/to/upstreamserver8080"),
								0,
							),
						2,
					).
					Insert(
						NewContext(context_type.TypeLocation, "~ /ipv4").
							Insert(
								NewContext(context_type.TypeDirHTTPProxyPass, "http://127.0.0.1/localiptest"),
								0,
							),
						3,
					).
					Insert(
						NewContext(context_type.TypeLocation, "~ /unknown-upstream").
							Insert(
								NewContext(context_type.TypeDirHTTPProxyPass, "http://unknown.upstream/localiptest"),
								0,
							),
						4,
					).
					Insert(
						NewContext(context_type.TypeLocation, "~ /has-unknown-server").
							Insert(
								NewContext(context_type.TypeDirHTTPProxyPass, "http://has_unknown_server/proxy/to/upstreamserver"),
								0,
							),
						5,
					).
					Insert(
						NewContext(context_type.TypeLocation, "~ /has-error-server").
							Insert(
								NewContext(context_type.TypeDirHTTPProxyPass, "http://has_error_server/proxy/to/errorupstreamserver"),
								0,
							),
						6,
					).
					Insert(
						NewContext(context_type.TypeLocation, "~ /wrong-pos-stream-proxy-pass").
							Insert(
								NewContext(context_type.TypeDirStreamProxyPass, "error_pos_stream_proxy_pass_2"),
								0,
							),
						7,
					).
					Insert(
						NewContext(context_type.TypeDirHTTPProxyPass, "http://error.pos.http.proxy.pass2"),
						8,
					).
					Insert(
						NewContext(context_type.TypeLocation, "~ /disabled-unknown-upstream").
							Insert(
								NewContext(context_type.TypeDirHTTPProxyPass, "https://unknown.upstream/localiptest"),
								0,
							).Disable(),
						9,
					),
				0,
			).
			Insert(
				NewContext(context_type.TypeUpstream, "upstream_server").
					Insert(
						NewContext(context_type.TypeDirective, "server test.cn:443"),
						0,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server example.com:8443 backup"),
						1,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server 127.0.0.1:8080"),
						2,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server 127.0.0.1"),
						3,
					),
				0,
			).
			Insert(
				NewContext(context_type.TypeUpstream, "has_unknown_server").
					Insert(
						NewContext(context_type.TypeDirective, "server test.cn:443"),
						0,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server example.com:8443 backup"),
						1,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server 127.0.0.1:8080"),
						2,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server unknown.domain:8080"),
						4,
					),
				1,
			).
			Insert(
				NewContext(context_type.TypeUpstream, "has_error_server").
					Insert(
						NewContext(context_type.TypeDirective, "server test.cn:443"),
						0,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server "),
						1,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server 127.0.0.1:8080"),
						2,
					),
				2,
			),
		0,
	).Insert(
		NewContext(context_type.TypeStream, "").
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeDirStreamProxyPass, "abc.com:22"),
						0,
					),
				0,
			).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeDirStreamProxyPass, "upstream_server"),
						0,
					),
				1,
			).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeDirStreamProxyPass, "127.0.0.1:9876"),
						0,
					),
				2,
			).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeDirStreamProxyPass, "unknown.domain:9876"),
						0,
					),
				3,
			).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeDirStreamProxyPass, "has_unknown_server"),
						0,
					),
				5,
			).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeDirStreamProxyPass, "has_error_server"),
						0,
					),
				6,
			).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeDirStreamProxyPass, "unknown_upstream"),
						0,
					),
				7,
			).
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeComment, "disabled stream server"),
						0,
					).
					Insert(
						NewContext(context_type.TypeDirStreamProxyPass, "unknown.domain:9876"),
						1,
					).Disable(),
				8,
			).
			Insert(
				NewContext(context_type.TypeUpstream, "upstream_server").
					Insert(
						NewContext(context_type.TypeDirective, "server test.cn:22 backup"),
						0,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server example.com:22"),
						1,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server 127.0.0.1:123"),
						2,
					),
				0,
			).
			Insert(
				NewContext(context_type.TypeUpstream, "has_unknown_server").
					Insert(
						NewContext(context_type.TypeDirective, "server test.cn:22 backup"),
						0,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server example.com:22"),
						1,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server 127.0.0.1:123"),
						2,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server unknown.domain:8080"),
						3,
					),
				1,
			).
			Insert(
				NewContext(context_type.TypeUpstream, "has_error_server").
					Insert(
						NewContext(context_type.TypeDirective, "server test.cn:22 backup"),
						0,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server example.com:22"),
						1,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server wrong.server"),
						2,
					).
					Insert(
						NewContext(context_type.TypeDirective, "server 127.0.0.1:123"),
						3,
					),
				2,
			),
		1,
	).Insert(
		NewContext(context_type.TypeDirHTTPProxyPass, "http://error.pos.http.proxy.pass"),
		2,
	).Insert(
		NewContext(context_type.TypeDirStreamProxyPass, "error_pos_stream_proxy_pass"),
		3,
	)

	return testMain, nil
}

func TestHTTPProxyPass_Clone(t *testing.T) {
	type fields struct {
		Directive   Directive
		OriginalURL string
		Protocol    string
		Addresses   []ProxiedAddress
		URI         string
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
			h := &HTTPProxyPass{
				Directive:   tt.fields.Directive,
				OriginalURL: tt.fields.OriginalURL,
				Protocol:    tt.fields.Protocol,
				Addresses:   tt.fields.Addresses,
				URI:         tt.fields.URI,
			}
			if got := h.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPProxyPass_SetFather(t *testing.T) {
	type fields struct {
		Directive   Directive
		OriginalURL string
		Protocol    string
		Addresses   []ProxiedAddress
		URI         string
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
			h := &HTTPProxyPass{
				Directive:   tt.fields.Directive,
				OriginalURL: tt.fields.OriginalURL,
				Protocol:    tt.fields.Protocol,
				Addresses:   tt.fields.Addresses,
				URI:         tt.fields.URI,
			}
			if err := h.SetFather(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SetFather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPProxyPass_SetValue(t *testing.T) {
	type fields struct {
		Directive   Directive
		OriginalURL string
		Protocol    string
		Addresses   []ProxiedAddress
		URI         string
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
			h := &HTTPProxyPass{
				Directive:   tt.fields.Directive,
				OriginalURL: tt.fields.OriginalURL,
				Protocol:    tt.fields.Protocol,
				Addresses:   tt.fields.Addresses,
				URI:         tt.fields.URI,
			}
			if err := h.SetValue(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPProxyPass_Type(t *testing.T) {
	type fields struct {
		Directive   Directive
		OriginalURL string
		Protocol    string
		Addresses   []ProxiedAddress
		URI         string
	}
	tests := []struct {
		name   string
		fields fields
		want   context_type.ContextType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPProxyPass{
				Directive:   tt.fields.Directive,
				OriginalURL: tt.fields.OriginalURL,
				Protocol:    tt.fields.Protocol,
				Addresses:   tt.fields.Addresses,
				URI:         tt.fields.URI,
			}
			if got := h.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPProxyPass_ReparseParams(t *testing.T) {
	utilsV3.SetDomainNameResolver(fakeHostsResolver())
	testMain, err := fakeProxyPassTestMainCtx()
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		proxyPass *HTTPProxyPass
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "http location proxy pass to an address",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`baidu`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass in a limit expect",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`baidu`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass).SetStringMatchingValue("abc.com"),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass in an If context",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`example`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass).SetStringMatchingValue("test.cn"),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass with an ipv4 address",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`ipv4`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass to upstream servers, with default port",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`/upstream$`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass to upstream servers, with a specify port",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`upstream2`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass to unknown domain name",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`unknown-upstream`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass to upstream servers, which has some servers with unknown domain name",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`has-unknown-server`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass with an error position",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass).SetStringMatchingValue("error.pos.http.proxy.pass"),
					).Target().(*HTTPProxyPass),
			},
			wantErr: true,
		},
		// {
		// 	name: "http location proxy pass to upstream servers, which has some wrong upstream servers",
		// 	fields: fields{
		// 		proxyPass: testMain.MainConfig().ChildrenPosSet().
		// 			QueryOne(
		// 				context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`has-error-server`),
		// 			).
		// 			QueryOne(
		// 				context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
		// 			).Target().(*HTTPProxyPass),
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.fields.proxyPass
			err := h.ReparseParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("reparseParams() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.IsCode(err, code.ErrV3DomainNameResolutionFailed) {
				t.Logf("Err: %++v", err)
			} else {
				aj, err := json.Marshal(h.Addresses)
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("OriginalURL: '%s'\nProtocol: '%s'\nAddresses: '%s'\nURI: '%s'", h.OriginalURL, h.Protocol, aj, h.URI)
			}
		})
	}
}

func TestHTTPProxyPass_PosVerify(t *testing.T) {
	utilsV3.SetDomainNameResolver(fakeHostsResolver())
	testMain, err := fakeProxyPassTestMainCtx()
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		proxyPass *HTTPProxyPass
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		wantErrCode int
	}{
		{
			name: "http location proxy pass to an address",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`baidu`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass in a limit expect",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`baidu`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass).SetStringMatchingValue("abc.com"),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass in an If context",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`example`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass).SetStringMatchingValue("test.cn"),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass with an error position",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass).AppendMatchingFilter(func(targetCtx context.Context) bool {
							return targetCtx.Value() == "proxy_pass http://error.pos.http.proxy.pass"
						}),
					).Target().(*HTTPProxyPass),
			},
			wantErr:     true,
			wantErrCode: code.ErrV3InvalidOperation,
		},
		{
			name: "http location proxy pass with an error position#2",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass).AppendMatchingFilter(func(targetCtx context.Context) bool {
							return targetCtx.Value() == "proxy_pass http://error.pos.http.proxy.pass2"
						}),
					).Target().(*HTTPProxyPass),
			},
			wantErr:     true,
			wantErrCode: code.ErrV3InvalidOperation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.fields.proxyPass
			if err := h.PosVerify(); (err != nil) != tt.wantErr {
				t.Errorf("PosVerify() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr && !errors.IsCode(err, tt.wantErrCode) {
				t.Errorf("PosVerify() error = %+v, wantErrCode %v", err, tt.wantErrCode)
			}
		})
	}
}

//func TestHTTPProxyPass_setFatherVerify(t *testing.T) {
//	type fields struct {
//		Directive   Directive
//		OriginalURL string
//		Protocol    string
//		Addresses   []ProxiedAddress
//		URI         string
//	}
//	type args struct {
//		ctx context.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := &HTTPProxyPass{
//				Directive:   tt.fields.Directive,
//				OriginalURL: tt.fields.OriginalURL,
//				Protocol:    tt.fields.Protocol,
//				Addresses:   tt.fields.Addresses,
//				URI:         tt.fields.URI,
//			}
//			if err := h.setFatherVerify(tt.args.ctx); (err != nil) != tt.wantErr {
//				t.Errorf("setFatherVerify() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestStreamProxyPass_Clone(t *testing.T) {
	type fields struct {
		Directive       Directive
		OriginalAddress string
		Addresses       []ProxiedAddress
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
			s := &StreamProxyPass{
				Directive:       tt.fields.Directive,
				OriginalAddress: tt.fields.OriginalAddress,
				Addresses:       tt.fields.Addresses,
			}
			if got := s.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamProxyPass_SetFather(t *testing.T) {
	type fields struct {
		Directive       Directive
		OriginalAddress string
		Addresses       []ProxiedAddress
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
			s := &StreamProxyPass{
				Directive:       tt.fields.Directive,
				OriginalAddress: tt.fields.OriginalAddress,
				Addresses:       tt.fields.Addresses,
			}
			if err := s.SetFather(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SetFather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStreamProxyPass_SetValue(t *testing.T) {
	type fields struct {
		Directive       Directive
		OriginalAddress string
		Addresses       []ProxiedAddress
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
			s := &StreamProxyPass{
				Directive:       tt.fields.Directive,
				OriginalAddress: tt.fields.OriginalAddress,
				Addresses:       tt.fields.Addresses,
			}
			if err := s.SetValue(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStreamProxyPass_Type(t *testing.T) {
	type fields struct {
		Directive       Directive
		OriginalAddress string
		Addresses       []ProxiedAddress
	}
	tests := []struct {
		name   string
		fields fields
		want   context_type.ContextType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamProxyPass{
				Directive:       tt.fields.Directive,
				OriginalAddress: tt.fields.OriginalAddress,
				Addresses:       tt.fields.Addresses,
			}
			if got := s.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamProxyPass_ReparseParams(t *testing.T) {
	utilsV3.SetDomainNameResolver(fakeHostsResolver())
	testMain, err := fakeProxyPassTestMainCtx()
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		proxyPass *StreamProxyPass
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "stream proxy pass to an address",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("abc.com:22"),
					).Target().(*StreamProxyPass),
			},
		},
		{
			name: "stream proxy pass to upstream servers",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("upstream_server"),
					).Target().(*StreamProxyPass),
			},
		},
		{
			name: "stream proxy pass to unknown server",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("unknown.domain"),
					).Target().(*StreamProxyPass),
			},
		},
		{
			name: "stream proxy pass to upstream servers, which has some servers with unknown domain name",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("has_unknown_server"),
					).Target().(*StreamProxyPass),
			},
		},
		{
			name: "stream proxy pass with an error position",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("error_pos_stream_proxy_pass"),
					).Target().(*StreamProxyPass),
			},
			wantErr: true,
		},
		{
			name: "stream proxy pass to upstream servers, which has some wrong servers",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("has_error_server"),
					).Target().(*StreamProxyPass),
			},
			wantErr: true,
		},
		{
			name: "stream proxy pass to an unknown upstream context",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("unknown_upstream"),
					).Target().(*StreamProxyPass),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.proxyPass
			err := s.ReparseParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("reparseParams() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.IsCode(err, code.ErrV3DomainNameResolutionFailed) {
				t.Logf("Err: %++v", err)
			} else {
				aj, err := json.Marshal(s.Addresses)
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("OriginalAddress: '%s'\nAddresses: '%s'", s.OriginalAddress, aj)
			}
		})
	}
}

func TestStreamProxyPass_PosVerify(t *testing.T) {
	utilsV3.SetDomainNameResolver(fakeHostsResolver())
	testMain, err := fakeProxyPassTestMainCtx()
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		proxyPass *StreamProxyPass
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		wantErrCode int
	}{
		{
			name: "stream proxy pass to an address",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("abc.com:22"),
					).Target().(*StreamProxyPass),
			},
		},
		{
			name: "stream proxy pass with an error position",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).AppendMatchingFilter(func(targetCtx context.Context) bool {
							return targetCtx.Value() == "proxy_pass error_pos_stream_proxy_pass"
						}),
					).Target().(*StreamProxyPass),
			},
			wantErr:     true,
			wantErrCode: code.ErrV3InvalidOperation,
		},
		{
			name: "stream proxy pass with an error position #2",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).AppendMatchingFilter(func(targetCtx context.Context) bool {
							return targetCtx.Value() == "proxy_pass error_pos_stream_proxy_pass_2"
						}),
					).Target().(*StreamProxyPass),
			},
			wantErr:     true,
			wantErrCode: code.ErrV3InvalidOperation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.proxyPass
			if err := s.PosVerify(); (err != nil) != tt.wantErr {
				t.Errorf("PosVerify() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr && !errors.IsCode(err, tt.wantErrCode) {
				t.Errorf("PosVerify() error = %+v, wantErrCode %v", err, tt.wantErrCode)
			}
		})
	}
}

func Test_parseHTTPURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name         string
		args         args
		wantProtocol string
		wantHost     string
		wantUri      string
		wantPort     uint16
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProtocol, gotHost, gotUri, gotPort, err := parseHTTPURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHTTPURL() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if gotProtocol != tt.wantProtocol {
				t.Errorf("parseHTTPURL() gotProtocol = %v, want %v", gotProtocol, tt.wantProtocol)
			}
			if gotHost != tt.wantHost {
				t.Errorf("parseHTTPURL() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if gotUri != tt.wantUri {
				t.Errorf("parseHTTPURL() gotUri = %v, want %v", gotUri, tt.wantUri)
			}
			if gotPort != tt.wantPort {
				t.Errorf("parseHTTPURL() gotPort = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}

func Test_registerHTTPProxyPassBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerHTTPProxyPassBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerHTTPProxyPassBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerProxyPassParseFunc(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerProxyPassParseFunc(); (err != nil) != tt.wantErr {
				t.Errorf("registerProxyPassParseFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerStreamProxyPassBuilder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerStreamProxyPassBuilder(); (err != nil) != tt.wantErr {
				t.Errorf("registerStreamProxyPassBuilder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPProxyPass_MarshalJSON(t *testing.T) {
	utilsV3.SetDomainNameResolver(fakeHostsResolver())
	testMain, err := fakeProxyPassTestMainCtx()
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		proxyPass *HTTPProxyPass
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "http location proxy pass to an address",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetRegexpMatchingValue(`baidu`),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
			want: []byte(`{"enabled":true,"context-type":"dir_http_proxy_pass","value":"https://baidu.com","proxy-pass":{"original-url":"https://baidu.com","protocol":"https","addresses":[{"domain-name":"baidu.com","port":443,"sockets":[{"ipv4":"10.1.11.111","port":443,"tcp-connectivity":0,"udp-connectivity":0},{"ipv4":"10.1.12.122","port":443,"tcp-connectivity":0,"udp-connectivity":0}],"resolve-err":null}],"uri":""}}`), //nolint:lll
		},
		{
			name: "http location proxy pass with some errors",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).
							SetStringMatchingValue("unknown-upstream").
							SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
			want: []byte("{\"enabled\":true,\"context-type\":\"dir_http_proxy_pass\",\"value\":\"http://unknown.upstream/localiptest\",\"proxy-pass\":{\"original-url\":\"http://unknown.upstream/localiptest\",\"protocol\":\"http\",\"addresses\":[{\"domain-name\":\"unknown.upstream\",\"port\":80,\"sockets\":[],\"resolve-err\":{\"message\":\"Domain name resolution failed\",\"error\":\"the domain name resolution record for `unknown.upstream` does not exist\",\"code\":110020}}],\"uri\":\"/localiptest\"}}"), //nolint:lll
		},
		{
			name: "http location proxy pass to upstream servers, which have some errors",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeLocation).SetStringMatchingValue("has-unknown-server"),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
			want: []byte("{\"enabled\":true,\"context-type\":\"dir_http_proxy_pass\",\"value\":\"http://has_unknown_server/proxy/to/upstreamserver\",\"proxy-pass\":{\"original-url\":\"http://has_unknown_server/proxy/to/upstreamserver\",\"protocol\":\"http\",\"addresses\":[{\"domain-name\":\"test.cn\",\"port\":443,\"sockets\":[{\"ipv4\":\"10.1.12.122\",\"port\":443,\"tcp-connectivity\":0,\"udp-connectivity\":0},{\"ipv4\":\"10.1.13.133\",\"port\":443,\"tcp-connectivity\":0,\"udp-connectivity\":0}],\"resolve-err\":null},{\"domain-name\":\"example.com\",\"port\":8443,\"sockets\":[{\"ipv4\":\"10.1.12.122\",\"port\":8443,\"tcp-connectivity\":0,\"udp-connectivity\":0},{\"ipv4\":\"10.1.13.133\",\"port\":8443,\"tcp-connectivity\":0,\"udp-connectivity\":0}],\"resolve-err\":null},{\"domain-name\":\"127.0.0.1\",\"port\":8080,\"sockets\":[{\"ipv4\":\"127.0.0.1\",\"port\":8080,\"tcp-connectivity\":0,\"udp-connectivity\":0}],\"resolve-err\":null},{\"domain-name\":\"unknown.domain\",\"port\":8080,\"sockets\":[],\"resolve-err\":{\"message\":\"Domain name resolution failed\",\"error\":\"the domain name resolution record for `unknown.domain` does not exist\",\"code\":110020}}],\"uri\":\"/proxy/to/upstreamserver\"}}"), //nolint:lll
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.fields.proxyPass
			err = h.ReparseParams()
			if err != nil {
				t.Fatal(err)
			}
			got, err := h.MarshalJSON()
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

func TestStreamProxyPass_MarshalJSON(t *testing.T) {
	utilsV3.SetDomainNameResolver(fakeHostsResolver())
	testMain, err := fakeProxyPassTestMainCtx()
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		proxyPass *StreamProxyPass
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "stream proxy pass to an address",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("abc.com:22"),
					).Target().(*StreamProxyPass),
			},
			want: []byte(`{"enabled":true,"context-type":"dir_stream_proxy_pass","value":"abc.com:22","proxy-pass":{"original-address":"abc.com:22","addresses":[{"domain-name":"abc.com","port":22,"sockets":[{"ipv4":"10.2.12.122","port":22,"tcp-connectivity":0,"udp-connectivity":0}],"resolve-err":null}]}}`), //nolint:lll
		},
		{
			name: "stream proxy pass to upstream servers, which have unknown server",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("has_unknown_server"),
					).Target().(*StreamProxyPass),
			},
			want: []byte("{\"enabled\":true,\"context-type\":\"dir_stream_proxy_pass\",\"value\":\"has_unknown_server\",\"proxy-pass\":{\"original-address\":\"has_unknown_server\",\"addresses\":[{\"domain-name\":\"test.cn\",\"port\":22,\"sockets\":[{\"ipv4\":\"10.1.12.122\",\"port\":22,\"tcp-connectivity\":0,\"udp-connectivity\":0},{\"ipv4\":\"10.1.13.133\",\"port\":22,\"tcp-connectivity\":0,\"udp-connectivity\":0}],\"resolve-err\":null},{\"domain-name\":\"example.com\",\"port\":22,\"sockets\":[{\"ipv4\":\"10.1.12.122\",\"port\":22,\"tcp-connectivity\":0,\"udp-connectivity\":0},{\"ipv4\":\"10.1.13.133\",\"port\":22,\"tcp-connectivity\":0,\"udp-connectivity\":0}],\"resolve-err\":null},{\"domain-name\":\"127.0.0.1\",\"port\":123,\"sockets\":[{\"ipv4\":\"127.0.0.1\",\"port\":123,\"tcp-connectivity\":0,\"udp-connectivity\":0}],\"resolve-err\":null},{\"domain-name\":\"unknown.domain\",\"port\":8080,\"sockets\":[],\"resolve-err\":{\"message\":\"Domain name resolution failed\",\"error\":\"the domain name resolution record for `unknown.domain` does not exist\",\"code\":110020}}]}}"), //nolint:lll
		},
		{
			name: "stream proxy pass to an address, which with an unknown domain name",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).
							SetStringMatchingValue("unknown.domain").
							SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc),
					).Target().(*StreamProxyPass),
			},
			want: []byte("{\"enabled\":true,\"context-type\":\"dir_stream_proxy_pass\",\"value\":\"unknown.domain:9876\",\"proxy-pass\":{\"original-address\":\"unknown.domain:9876\",\"addresses\":[{\"domain-name\":\"unknown.domain\",\"port\":9876,\"sockets\":[],\"resolve-err\":{\"message\":\"Domain name resolution failed\",\"error\":\"the domain name resolution record for `unknown.domain` does not exist\",\"code\":110020}}]}}"), //nolint:lll
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.proxyPass
			err = s.ReparseParams()
			if err != nil {
				t.Fatal(err)
			}
			got, err := s.MarshalJSON()
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

//nolint:funlen
func Test_tmpProxyPass_verifyAndInitTo(t1 *testing.T) {
	utilsV3.SetDomainNameResolver(fakeHostsResolver())
	testMain, err := fakeProxyPassTestMainCtx()
	if err != nil {
		t1.Fatal(err)
	}
	initializedTmp := &tmpProxyPass{
		isInitialized: true,
		rwLock:        new(sync.RWMutex),
	}
	excludedTmp := &tmpProxyPass{
		isInitialized: false,
		rwLock:        new(sync.RWMutex),
	}
	httpTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "https://example.com",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	httpWrongPosTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "https://wrong.pos",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	streamTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "example.com",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	streamWrongPosTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "wrong.pos",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	wrongTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "unknownProtocol://example.com",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	invalidOperatedTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "https://example.com",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "excluded.conf").
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeLocation, "/test").
							Insert(excludedTmp, 0),
						0,
					),
				0,
			).(*Config),
	)
	if err != nil {
		t1.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "http_proxy_pass.conf").
			Insert(httpTmp, 0).(*Config),
	)
	if err != nil {
		t1.Fatal(err)
	}
	testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeHttp)).
		QueryOne(context.NewKeyWordsByType(context_type.TypeServer)).Target().
		Insert(
			NewContext(context_type.TypeLocation, "~ /included-proxy-pass").
				Insert(
					NewContext(context_type.TypeIf, "( $request_uri =~ included-proxy-pass )").
						Insert(
							NewContext(context_type.TypeInclude, "http_proxy_pass.conf"),
							0,
						),
					0,
				),
			0,
		)

	testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeHttp)).Target().
		Insert(httpWrongPosTmp, 0)

	testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeStream)).Target().
		Insert(
			NewContext(context_type.TypeServer, "").
				Insert(streamTmp, 0),
			0,
		)

	testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeStream)).Target().
		Insert(streamWrongPosTmp, 0)

	loc := NewContext(context_type.TypeLocation, "~ /invalid-operated-proxy-pass").
		Insert(NewContext(context_type.TypeComment, "stack"), 0)
	testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeHttp)).
		QueryOne(context.NewKeyWordsByType(context_type.TypeServer)).Target().
		Insert(
			NewContext(context_type.TypeLocation, "/wrong-url").
				Insert(wrongTmp, 0),
			0,
		).
		Insert(loc, 1)
	invalidOperatedTmp.fatherContext = loc

	type fields struct {
		proxyPass *tmpProxyPass
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name: "initialized tmpProxyPass",
			fields: fields{
				proxyPass: initializedTmp,
			},
			want: context.ErrContext(errors.New("ProxyPass is already initialized")),
		},
		{
			name: "excluded to main context",
			fields: fields{
				proxyPass: excludedTmp,
			},
			want: excludedTmp,
		},
		{
			name: "http proxy pass",
			fields: fields{
				proxyPass: httpTmp,
			},
			want: &HTTPProxyPass{Directive: httpTmp.Directive},
		},
		{
			name: "wrong pos http proxy pass",
			fields: fields{
				proxyPass: httpWrongPosTmp,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "Father context type not supported."+
				" Only the following types are supported: `Location`, `If` in `Location`, `Limit_except`")),
		},
		{
			name: "stream proxy pass",
			fields: fields{
				proxyPass: streamTmp,
			},
			want: &StreamProxyPass{Directive: streamTmp.Directive},
		},
		{
			name: "wrong pos stream proxy pass",
			fields: fields{
				proxyPass: streamWrongPosTmp,
			},
			want: context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "Father context type not supported."+
				" Only the following types are supported: `Server` which is in `Stream`")),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields.proxyPass
			if got := t.verifyAndInitTo(); reflect.DeepEqual(got, tt.want) {
			} else if got.Error() == nil && tt.want.Error() == nil && got.Value() == tt.want.Value() {
			} else if got.Error() != nil && tt.want.Error() != nil {
				if got.Error().Error() != tt.want.Error().Error() {
					t1.Errorf("verifyAndInitTo().Error() = %#v, want ErrContext %#v", got.Error(), tt.want.Error())
				}
			} else {
				t1.Errorf("verifyAndInitTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:funlen
func Test_tmpProxyPass_MarshalJSON(t1 *testing.T) {
	utilsV3.SetDomainNameResolver(fakeHostsResolver())
	testMain, err := fakeProxyPassTestMainCtx()
	if err != nil {
		t1.Fatal(err)
	}
	initializedTmp := &tmpProxyPass{
		isInitialized: true,
		rwLock:        new(sync.RWMutex),
	}
	excludedTmp := &tmpProxyPass{
		isInitialized: false,
		Directive: Directive{
			enabled: true,
			Name:    "proxy_pass",
			Params:  "https://example.com",
		},
		rwLock: new(sync.RWMutex),
	}
	httpTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "https://example.com",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	httpWrongPosTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "https://wrong.pos",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	streamTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "example.com",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	streamWrongPosTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "wrong.pos",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	wrongTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "unknownProtocol://example.com",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	invalidOperatedTmp := &tmpProxyPass{
		Directive: Directive{
			enabled:       true,
			Name:          "proxy_pass",
			Params:        "https://example.com",
			fatherContext: context.NullContext(),
		},
		rwLock: new(sync.RWMutex),
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "excluded.conf").
			Insert(
				NewContext(context_type.TypeServer, "").
					Insert(
						NewContext(context_type.TypeLocation, "/test").
							Insert(excludedTmp, 0),
						0,
					),
				0,
			).(*Config),
	)
	if err != nil {
		t1.Fatal(err)
	}
	err = testMain.AddConfig(
		NewContext(context_type.TypeConfig, "http_proxy_pass.conf").
			Insert(httpTmp, 0).(*Config),
	)
	if err != nil {
		t1.Fatal(err)
	}
	testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeHttp)).
		QueryOne(context.NewKeyWordsByType(context_type.TypeServer)).Target().
		Insert(
			NewContext(context_type.TypeLocation, "~ /included-proxy-pass").
				Insert(
					NewContext(context_type.TypeIf, "( $request_uri =~ included-proxy-pass )").
						Insert(
							NewContext(context_type.TypeInclude, "http_proxy_pass.conf"),
							0,
						),
					0,
				),
			0,
		)

	testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeHttp)).Target().
		Insert(httpWrongPosTmp, 0)

	testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeStream)).Target().
		Insert(
			NewContext(context_type.TypeServer, "").
				Insert(streamTmp, 0),
			0,
		)

	testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeStream)).Target().
		Insert(streamWrongPosTmp, 0)

	loc := NewContext(context_type.TypeLocation, "~ /invalid-operated-proxy-pass").
		Insert(NewContext(context_type.TypeComment, "stack"), 0)
	testMain.MainConfig().ChildrenPosSet().QueryOne(context.NewKeyWordsByType(context_type.TypeHttp)).
		QueryOne(context.NewKeyWordsByType(context_type.TypeServer)).Target().
		Insert(
			NewContext(context_type.TypeLocation, "/wrong-url").
				Insert(wrongTmp, 0),
			0,
		).
		Insert(loc, 1)
	invalidOperatedTmp.fatherContext = loc

	type fields struct {
		proxyPass *tmpProxyPass
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "initialized tmpProxyPass",
			fields: fields{
				proxyPass: initializedTmp,
			},
			wantErr: true,
		},
		{
			name: "excluded to main context",
			fields: fields{
				proxyPass: excludedTmp,
			},
			want: []byte("{\"enabled\":true,\"context-type\":\"directive\",\"value\":\"proxy_pass https://example.com\"}"),
		},
		{
			name: "http proxy pass",
			fields: fields{
				proxyPass: httpTmp,
			},
			want: []byte("{\"enabled\":true,\"context-type\":\"dir_http_proxy_pass\",\"value\":\"https://example.com\",\"proxy-pass\":{\"original-url\":\"\",\"protocol\":\"\",\"addresses\":null,\"uri\":\"\"}}"),
		},
		{
			name: "wrong pos http proxy pass",
			fields: fields{
				proxyPass: httpWrongPosTmp,
			},
			wantErr: true,
		},
		{
			name: "stream proxy pass",
			fields: fields{
				proxyPass: streamTmp,
			},
			want: []byte("{\"enabled\":true,\"context-type\":\"dir_stream_proxy_pass\",\"value\":\"example.com\",\"proxy-pass\":{\"original-address\":\"\",\"addresses\":null}}"),
		},
		{
			name: "wrong pos stream proxy pass",
			fields: fields{
				proxyPass: streamWrongPosTmp,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields.proxyPass
			got, err := t.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t1.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}
