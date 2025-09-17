package local

import (
	"encoding/json"
	"net"
	"reflect"
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

func TestHTTPProxyPass_reparseParams(t *testing.T) {
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
						context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`baidu`),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass in a limit expect",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`baidu`),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirHTTPProxyPass).SetStringMatchingValue("abc.com"),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass in an If context",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`example`),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirHTTPProxyPass).SetStringMatchingValue("test.cn"),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass with an ipv4 address",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`ipv4`),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass to upstream servers, with default port",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`/upstream$`),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass to upstream servers, with a specify port",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`upstream2`),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass to unknown domain name",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`unknown-upstream`),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass to upstream servers, which has some servers with unknown domain name",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`has-unknown-server`),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirHTTPProxyPass),
					).Target().(*HTTPProxyPass),
			},
		},
		{
			name: "http location proxy pass with an error position",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeDirHTTPProxyPass).SetStringMatchingValue("error.pos.http.proxy.pass"),
					).Target().(*HTTPProxyPass),
			},
			wantErr: true,
		},
		// {
		// 	name: "http location proxy pass to upstream servers, which has some wrong upstream servers",
		// 	fields: fields{
		// 		proxyPass: testMain.MainConfig().ChildrenPosSet().
		// 			QueryOne(
		// 				context.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`has-error-server`),
		// 			).
		// 			QueryOne(
		// 				context.NewKeyWords(context_type.TypeDirHTTPProxyPass),
		// 			).Target().(*HTTPProxyPass),
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.fields.proxyPass
			err := h.reparseParams()
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

func TestStreamProxyPass_reparseParams(t *testing.T) {
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
						context.NewKeyWords(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("abc.com:22"),
					).Target().(*StreamProxyPass),
			},
		},
		{
			name: "stream proxy pass to upstream servers",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("upstream_server"),
					).Target().(*StreamProxyPass),
			},
		},
		{
			name: "stream proxy pass to unknown server",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("unknown.domain"),
					).Target().(*StreamProxyPass),
			},
		},
		{
			name: "stream proxy pass to upstream servers, which has some servers with unknown domain name",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("has_unknown_server"),
					).Target().(*StreamProxyPass),
			},
		},
		{
			name: "http location proxy pass with an error position",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("error_pos_stream_proxy_pass"),
					).Target().(*StreamProxyPass),
			},
			wantErr: true,
		},
		{
			name: "stream proxy pass to upstream servers, which has some wrong servers",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("has_error_server"),
					).Target().(*StreamProxyPass),
			},
			wantErr: true,
		},
		{
			name: "stream proxy pass to an unknown upstream context",
			fields: fields{
				proxyPass: testMain.MainConfig().ChildrenPosSet().
					QueryOne(
						context.NewKeyWords(context_type.TypeStream),
					).
					QueryOne(
						context.NewKeyWords(context_type.TypeDirStreamProxyPass).SetStringMatchingValue("unknown_upstream"),
					).Target().(*StreamProxyPass),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.proxyPass
			err := s.reparseParams()
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

//func TestStreamProxyPass_setFatherVerify(t *testing.T) {
//	type fields struct {
//		Directive       Directive
//		OriginalAddress string
//		Addresses       []ProxiedAddress
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
//			s := &StreamProxyPass{
//				Directive:       tt.fields.Directive,
//				OriginalAddress: tt.fields.OriginalAddress,
//				Addresses:       tt.fields.Addresses,
//			}
//			if err := s.setFatherVerify(tt.args.ctx); (err != nil) != tt.wantErr {
//				t.Errorf("setFatherVerify() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

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
