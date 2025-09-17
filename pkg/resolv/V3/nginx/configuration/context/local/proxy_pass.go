package local

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"

	"github.com/marmotedu/errors"
)

var (
	dirHTTPProxyPassMustCompile   = regexp.MustCompile(`^\s*proxy_pass\s+("\s*https?://[^";\s#]+"|'\s*https?://[^';\s#]+'|\s*https?://[^;\s#]+)`)
	dirStreamProxyPassMustCompile = regexp.MustCompile(`^\s*proxy_pass\s+(` + S1 + `)$`)
	httpURLMustCompile            = regexp.MustCompile(`^(http|https)://([^/]+\S*)$`)
	addrAndURIMustCompile         = regexp.MustCompile(`^([^/]+)(/\S*)?$`)
	streamAddrWithPortMustCompile = regexp.MustCompile(`^([^/:\s]+):(\d+)$`)
	streamUpstreamMustCompile     = regexp.MustCompile(`^([^/\s]+)$`)
	streamUnixAddrMustCompile     = regexp.MustCompile(`^unix:(/[^/]\s*)$`)
	addrWithPortMustCompile       = regexp.MustCompile(`^([^/:\s]+):(\d+)`)
	upstreamSrvDirMustCompile     = regexp.MustCompile(`^server\s+(\S+)`)
)

type ProxiedAddress struct {
	DomainName string
	Port       uint16
	IPv4s      []string
	ResolveErr error
}

type HTTPProxyPass struct {
	Directive
	OriginalURL string
	Protocol    string
	Addresses   []ProxiedAddress
	URI         string
}

func (h *HTTPProxyPass) Clone() context.Context {
	cloneProxiedAddress := make([]ProxiedAddress, len(h.Addresses))
	for i, address := range h.Addresses {
		cloneProxiedAddress[i] = ProxiedAddress{
			DomainName: address.DomainName,
			Port:       address.Port,
			IPv4s:      address.IPv4s,
		}
	}

	return &HTTPProxyPass{
		Directive:   *h.Directive.Clone().(*Directive),
		OriginalURL: h.OriginalURL,
		Protocol:    h.Protocol,
		Addresses:   cloneProxiedAddress,
		URI:         h.URI,
	}
}

func (h *HTTPProxyPass) SetValue(v string) (err error) {
	p := h.Params
	err = h.Directive.SetValue(v)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			h.Name = "proxy_pass"
			h.Params = p
		}
	}()
	if subMatch := dirHTTPProxyPassMustCompile.FindSubmatch([]byte(v)); len(subMatch) != 2 {
		return errors.WithCode(code.ErrV3InvalidOperation, "the set value does not conform to the syntax rules, the syntax is `proxy_pass <URL>`")
	}

	return h.reparseParams()
}

func (h *HTTPProxyPass) SetFather(ctx context.Context) error {
	//err := h.setFatherVerify(ctx)
	//if err != nil {
	//	return err
	//}

	return h.Directive.SetFather(ctx)
}

func (h *HTTPProxyPass) Type() context_type.ContextType {
	return context_type.TypeDirHTTPProxyPass
}

func (h *HTTPProxyPass) reparseParams() (err error) {
	var (
		defaultProxyPassPort uint16
		proxiedAddrs         []ProxiedAddress
		existAddrMap         = make(map[string]bool)
	)
	url := strings.TrimSpace(strings.Trim(h.Params, `'"`))

	// parse url
	protocol, host, uri, port, err := parseHTTPURL(url)
	if err != nil {
		return err
	}

	// protocol must be "http" or "https"
	if protocol == "http" {
		defaultProxyPassPort = 80
	} else if protocol == "https" {
		defaultProxyPassPort = 443
	}

	defer func() {
		if err != nil {
			return
		}
		h.OriginalURL = url
		h.Protocol = protocol
		h.Addresses = proxiedAddrs
		h.URI = uri
	}()

	// parse IPv4s
	ipv4s, err := utilsV3.ResolveDomainNameToIPv4(host)
	if err == nil {
		if port == 0 {
			port = defaultProxyPassPort
		}
		proxiedAddrs = append(proxiedAddrs, ProxiedAddress{
			DomainName: host,
			Port:       port,
			IPv4s:      ipv4s,
		})

		return nil
	}

	// parse upstream servers
	httpCtx := getFatherContextByType(h, context_type.TypeHttp)
	if httpCtx.Error() != nil {
		return errors.WithCode(code.ErrV3InvalidContext, "cannot query father HTTP Context. error: %++v", httpCtx.Error())
	}
	upstreamPos := httpCtx.ChildrenPosSet().QueryOne(
		context.NewKeyWords(context_type.TypeUpstream).
			SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc).
			SetMatchingFilter(func(targetCtx context.Context) bool {
				return targetCtx.Value() == host
			}),
	)

	if posErr := upstreamPos.Target().Error(); posErr != nil {
		if port == 0 {
			port = defaultProxyPassPort
		}
		proxiedAddrs = append(proxiedAddrs, ProxiedAddress{
			DomainName: host,
			Port:       port,
			IPv4s:      ipv4s,
			ResolveErr: err,
		})

		return nil
	}

	return upstreamPos.QueryAll(context.NewKeyWords(context_type.TypeDirective).
		SetRegexpMatchingValue(`^server\s+`).
		SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
		Filter(
			func(dirPos context.Pos) bool {
				return upstreamSrvDirMustCompile.MatchString(dirPos.Target().Value())
			},
		).
		Map(
			func(dirPos context.Pos) (context.Pos, error) {
				srvSubMatch := upstreamSrvDirMustCompile.FindStringSubmatch(dirPos.Target().Value())
				if len(srvSubMatch) != 2 {
					return dirPos, errors.WithCode(code.ErrV3InvalidContext, "the directive `server` must be `server <address> [parameters]`,"+
						" the current is: '%s'", dirPos.Target().Value())
				}
				srvAddr := srvSubMatch[1]
				var (
					srvHost  string
					srvIPv4s []string
					srvPort  uint16
				)
				if subMatch := addrWithPortMustCompile.FindStringSubmatch(srvAddr); len(subMatch) == 3 {
					srvHost = subMatch[1]
					sp, err := strconv.Atoi(subMatch[2])
					if err != nil {
						return dirPos, errors.WithCode(code.ErrV3InvalidOperation, "the port must be an integer, error: %s", err.Error())
					}
					srvPort = uint16(sp)
				} else { // TODO: parse error <address> string
					srvHost = srvAddr
				}
				if port != 0 {
					srvPort = port
				}
				if srvPort == 0 {
					srvPort = defaultProxyPassPort
				}
				addr := srvHost + ":" + strconv.Itoa(int(srvPort))
				if !existAddrMap[addr] {
					srvIPv4s, err = utilsV3.ResolveDomainNameToIPv4(srvHost)
					proxiedAddrs = append(proxiedAddrs, ProxiedAddress{
						DomainName: srvHost,
						Port:       srvPort,
						IPv4s:      srvIPv4s,
						ResolveErr: err,
					})
					existAddrMap[addr] = true
				}

				return dirPos, nil
			},
		).Error()
}

func parseHTTPURL(url string) (protocol, host, uri string, port uint16, err error) {
	var addressAndURI, address string

	// parse protocol
	if urlMatch := httpURLMustCompile.FindStringSubmatch(url); len(urlMatch) == 3 {
		protocol = urlMatch[1]
		addressAndURI = urlMatch[2]
	} else {
		err = errors.WithCode(code.ErrV3InvalidOperation, "the URL protocol must be `http` or `https`")

		return protocol, host, uri, port, err
	}

	// parse address and uri
	if subMatch := addrAndURIMustCompile.FindStringSubmatch(addressAndURI); len(subMatch) >= 2 {
		address = subMatch[1]
		if len(subMatch) == 3 {
			uri = subMatch[2]
		}
	} else {
		err = errors.WithCode(code.ErrV3InvalidOperation, "the address and URI must match the format `<host>[:<port>]/<URI>`")

		return protocol, host, uri, port, err
	}

	// parse host and port
	if subMatch := addrWithPortMustCompile.FindStringSubmatch(address); len(subMatch) == 3 {
		host = subMatch[1]
		p, err := strconv.Atoi(subMatch[2])
		if err != nil {
			err = errors.WithCode(code.ErrV3InvalidOperation, "the port must be an integer, error: %s", err.Error())

			return protocol, host, uri, port, err
		}
		port = uint16(p)
	} else {
		host = address
	}

	return protocol, host, uri, port, err
}

//func (h *HTTPProxyPass) setFatherVerify(ctx context.Context) error {
//	// TODO
//	// the father Context must in: Location, If in Location, Limit_except
//	if fatherT := ctx.Type(); fatherT != context_type.TypeLocation &&
//		fatherT != context_type.TypeLimitExcept &&
//		(fatherT != context_type.TypeIf || ctx.Father().Type() != context_type.TypeLocation) {
//		return errors.WithCode(code.ErrV3InvalidOperation, "Father context type not supported."+
//			" Only the following types are supported: `Location`, `If` in `Location`, `Limit_except`")
//	}
//
//	return nil
//}

func registerHTTPProxyPassBuilder() error {
	builderMap[context_type.TypeDirHTTPProxyPass] = func(value string) context.Context {
		if !dirHTTPProxyPassMustCompile.MatchString("proxy_pass " + value) { //nolint:nestif
			return context.ErrContext(errors.Errorf("invalid value `%s`", value))
		}

		return &HTTPProxyPass{Directive: *NewContext(context_type.TypeDirective, "proxy_pass "+value).(*Directive)}
	}

	return nil
}

type StreamProxyPass struct {
	Directive
	OriginalAddress string
	Addresses       []ProxiedAddress
}

func (s *StreamProxyPass) Clone() context.Context {
	cloneProxiedAddress := make([]ProxiedAddress, len(s.Addresses))
	for i, address := range s.Addresses {
		cloneProxiedAddress[i] = ProxiedAddress{
			DomainName: address.DomainName,
			Port:       address.Port,
			IPv4s:      address.IPv4s,
		}
	}

	return &StreamProxyPass{
		Directive:       *s.Directive.Clone().(*Directive),
		OriginalAddress: s.OriginalAddress,
		Addresses:       cloneProxiedAddress,
	}
}

func (s *StreamProxyPass) SetValue(v string) (err error) {
	p := s.Params
	err = s.Directive.SetValue(v)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			s.Name = "proxy_pass"
			s.Params = p
		}
	}()
	if subMatch := dirStreamProxyPassMustCompile.FindSubmatch([]byte(v)); len(subMatch) != 2 {
		return errors.WithCode(code.ErrV3InvalidOperation, "the set value does not conform to the syntax rules, the syntax is `proxy_pass <Address>`")
	}

	return s.reparseParams()
}

func (s *StreamProxyPass) SetFather(ctx context.Context) error {
	//err := s.setFatherVerify(ctx)
	//if err != nil {
	//	return err
	//}

	return s.Directive.SetFather(ctx)
}

func (s *StreamProxyPass) Type() context_type.ContextType {
	return context_type.TypeDirStreamProxyPass
}

func (s *StreamProxyPass) reparseParams() (err error) {
	var (
		host, upstreamServer string
		port                 uint16
		proxiedAddrs         []ProxiedAddress
		existAddrMap         = make(map[string]bool)
	)
	address := strings.TrimSpace(strings.Trim(s.Params, `'"`))

	// parse protocol
	switch {
	case streamAddrWithPortMustCompile.MatchString(address): // pass proxied address, not upstream server
		addressMatch := streamAddrWithPortMustCompile.FindStringSubmatch(address)
		host = addressMatch[1]
		p, err := strconv.Atoi(addressMatch[2])
		if err != nil {
			return errors.WithCode(code.ErrV3InvalidOperation, "the port must be an integer, error: %s", err.Error())
		}
		port = uint16(p)
	case streamUpstreamMustCompile.MatchString(address): // pass upstream server, not proxied address
		upstreamServer = streamUpstreamMustCompile.FindStringSubmatch(address)[1]
	case streamUnixAddrMustCompile.MatchString(address):
		// TODO: parse UNIX-domain socket
		return nil
	default:
		return errors.WithCode(code.ErrV3InvalidOperation, "the address can be specified as a domain name or IP address, and a port")
	}

	defer func() {
		if err != nil {
			return
		}
		s.OriginalAddress = address
		s.Addresses = proxiedAddrs
	}()

	// parse proxied address
	if len(strings.TrimSpace(host)) > 0 {
		if port == 0 {
			return errors.WithCode(code.ErrV3InvalidOperation, "the port must be set")
		}
		ipv4s, err := utilsV3.ResolveDomainNameToIPv4(host)
		proxiedAddrs = append(proxiedAddrs, ProxiedAddress{
			DomainName: host,
			Port:       port,
			IPv4s:      ipv4s,
			ResolveErr: err,
		})

		return nil
	}

	// parse upstream server
	if len(strings.TrimSpace(upstreamServer)) == 0 {
		return errors.WithCode(code.ErrV3InvalidOperation, "the upstream server must not be empty")
	}
	streamCtx := getFatherContextByType(s, context_type.TypeStream)
	if streamCtx.Error() != nil {
		return errors.WithCode(code.ErrV3InvalidContext, "cannot find father Stream Context. error: %++v", streamCtx.Error())
	}

	upstreamPos := streamCtx.ChildrenPosSet().QueryOne(
		context.NewKeyWords(context_type.TypeUpstream).
			SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc).
			SetMatchingFilter(func(targetCtx context.Context) bool {
				return targetCtx.Value() == upstreamServer
			}),
	)
	if err = upstreamPos.Target().Error(); err != nil {
		return errors.WithCode(code.ErrV3InvalidOperation, "unknown upstream server: '%s'. failed to query Upstream(Stream) Context, error: %v", upstreamServer, err)
	}

	return upstreamPos.QueryAll(context.NewKeyWords(context_type.TypeDirective).
		SetRegexpMatchingValue(`^server\s+`).
		SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
		Map(
			func(dirPos context.Pos) (context.Pos, error) {
				if srvSubMatch := upstreamSrvDirMustCompile.FindStringSubmatch(dirPos.Target().Value()); len(srvSubMatch) == 2 {
					srvAddr := srvSubMatch[1]
					var (
						srvHost string
						srvPort uint16
					)
					switch {
					case addrWithPortMustCompile.MatchString(srvAddr):
						subAddrMatch := addrWithPortMustCompile.FindStringSubmatch(srvAddr)
						srvHost = subAddrMatch[1]
						sp, err := strconv.Atoi(subAddrMatch[2])
						if err != nil {
							return dirPos, errors.WithCode(code.ErrV3InvalidOperation, "the port must be an integer, error: %s", err.Error())
						}
						srvPort = uint16(sp)
					case streamUnixAddrMustCompile.MatchString(srvAddr):
						// TODO: parse UNIX-domain socket
						return dirPos, nil
					default:
						return dirPos, errors.WithCode(code.ErrV3InvalidContext, "this address in the server directive under the upstream context of the stream context is not standardized")
					}
					addr := srvHost + ":" + strconv.Itoa(int(srvPort))
					if !existAddrMap[addr] {
						ipv4s, err := utilsV3.ResolveDomainNameToIPv4(srvHost)
						proxiedAddrs = append(proxiedAddrs, ProxiedAddress{
							DomainName: srvHost,
							Port:       srvPort,
							IPv4s:      ipv4s,
							ResolveErr: err,
						})
						existAddrMap[addr] = true
					}
				}

				return dirPos, nil
			},
		).Error()
}

//func (s *StreamProxyPass) setFatherVerify(ctx context.Context) error {
//	// the set father Context must be Server in Stream
//  // TODO
//	if ctx.Type() != context_type.TypeServer && ctx.Father().Type() != context_type.TypeStream {
//		return errors.WithCode(code.ErrV3InvalidOperation, "Father context type not supported. Only the following types are supported: `Server` which is in `Stream`")
//	}
//
//	return nil
//}

func registerStreamProxyPassBuilder() error {
	builderMap[context_type.TypeDirStreamProxyPass] = func(value string) context.Context {
		if !dirStreamProxyPassMustCompile.MatchString("proxy_pass " + value) { //nolint:nestif
			return context.ErrContext(errors.Errorf("invalid value `%s`", value))
		}

		return &StreamProxyPass{Directive: *NewContext(context_type.TypeDirective, "proxy_pass "+value).(*Directive)}
	}

	return nil
}

func registerProxyPassParseFunc() error {
	if inStackParseFuncMap[context_type.TypeDirective] == nil {
		err := registerDirectiveParseFunc()
		if err != nil {
			return err
		}
	}
	dirBuilder := inStackParseFuncMap[context_type.TypeDirective]
	inStackParseFuncMap[context_type.TypeDirective] = func(data []byte, idx *int) context.Context {
		dir := dirBuilder(data, idx)
		if dir == context.NullContext() {
			return dir
		}
		if dir.Type() == context_type.TypeDirective {
			if dirHTTPProxyPassMustCompile.MatchString(dir.Value()) { //nolint:nestif
				return &StreamProxyPass{Directive: *dir.(*Directive)}
			} else if dirStreamProxyPassMustCompile.MatchString(dir.Value()) { //nolint:nestif
				return &HTTPProxyPass{Directive: *dir.(*Directive)}
			}
		}

		return dir
	}

	return nil
}
