package local

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
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

type ProxyPass interface {
	context.Context
	PosVerify() error
	ReparseParams() error
	ConnectivityCheck() ProxyPass
}

type ProxiedAddress struct {
	DomainName string     `json:"domain-name"`
	Port       uint16     `json:"port"`
	Sockets    []*Socket  `json:"sockets"`
	ResolveErr *JSONError `json:"resolve-err"`
}

type Socket struct {
	IPv4            string                 `json:"ipv4"`
	Port            uint16                 `json:"port"`
	TCPConnectivity v1.NetworkConnectivity `json:"tcp-connectivity"`
	UDPConnectivity v1.NetworkConnectivity `json:"udp-connectivity"`
}

func (s *Socket) valid() error {
	if s.Port == 0 || fmt.Sprint(net.ParseIP(s.IPv4).To4()) == "<nil>" {
		return errors.New("invalid address")
	}

	return nil
}

func (s *Socket) tcpCheck() error {
	if err := s.valid(); err != nil {
		return err
	}
	s.TCPConnectivity = utilsV3.SocketConnectivityCheck(utilsV3.TCP, s.IPv4, strconv.Itoa(int(s.Port)), time.Millisecond*100)

	return nil
}

func (s *Socket) udpCheck() error {
	if err := s.valid(); err != nil {
		return err
	}
	s.UDPConnectivity = utilsV3.SocketConnectivityCheck(utilsV3.UDP, s.IPv4, strconv.Itoa(int(s.Port)), time.Millisecond*100)

	return nil
}

func (s *Socket) allCheck() error {
	err := s.tcpCheck()
	if err != nil {
		return err
	}

	return s.udpCheck()
}

func newSocket(ipv4 string, port uint16) *Socket {
	return &Socket{
		IPv4:            ipv4,
		Port:            port,
		TCPConnectivity: v1.NetUnknown,
		UDPConnectivity: v1.NetUnknown,
	}
}

func newSockets(ipv4s []string, port uint16) []*Socket {
	sockets := make([]*Socket, len(ipv4s))
	for i, ipv4 := range ipv4s {
		sockets[i] = newSocket(ipv4, port)
	}

	return sockets
}

type HTTPProxyPass struct {
	Directive
	OriginalURL string
	Protocol    string
	Addresses   []ProxiedAddress
	URI         string
}

func (h *HTTPProxyPass) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Enabled     bool                     `json:"enabled,omitempty"`
		ContextType context_type.ContextType `json:"context-type"`
		Value       string                   `json:"value"`
		ProxyPass   struct {
			OriginalURL string           `json:"original-url"`
			Protocol    string           `json:"protocol"`
			Addresses   []ProxiedAddress `json:"addresses"`
			URI         string           `json:"uri"`
		} `json:"proxy-pass,omitempty"`
	}{
		Enabled:     h.IsEnabled(),
		ContextType: context_type.TypeDirHTTPProxyPass,
		Value:       h.Params,
		ProxyPass: struct {
			OriginalURL string           `json:"original-url"`
			Protocol    string           `json:"protocol"`
			Addresses   []ProxiedAddress `json:"addresses"`
			URI         string           `json:"uri"`
		}{OriginalURL: h.OriginalURL, Protocol: h.Protocol, Addresses: h.Addresses, URI: h.URI},
	})
}

func (h *HTTPProxyPass) UnmarshalJSON(bytes []byte) error {
	tmp := struct {
		Enabled     bool                     `json:"enabled,omitempty"`
		ContextType context_type.ContextType `json:"context-type"`
		Value       string                   `json:"value"`
		ProxyPass   struct {
			OriginalURL string           `json:"original-url"`
			Protocol    string           `json:"protocol"`
			Addresses   []ProxiedAddress `json:"addresses"`
			URI         string           `json:"uri"`
		} `json:"proxy-pass,omitempty"`
	}{}
	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}
	if tmp.ContextType != context_type.TypeDirHTTPProxyPass {
		return errors.WithCode(code.ErrV3InvalidContext, "invalid context-type: %s", tmp.ContextType)
	}
	h.enabled = tmp.Enabled
	h.Params = tmp.Value
	h.OriginalURL = tmp.ProxyPass.OriginalURL
	h.Protocol = tmp.ProxyPass.Protocol
	h.Addresses = tmp.ProxyPass.Addresses
	h.URI = tmp.ProxyPass.URI

	return nil
}

func (h *HTTPProxyPass) Clone() context.Context {
	cloneProxiedAddress := make([]ProxiedAddress, len(h.Addresses))
	for i, address := range h.Addresses {
		cloneProxiedAddress[i] = ProxiedAddress{
			DomainName: address.DomainName,
			Port:       address.Port,
			Sockets:    address.Sockets,
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
	err = h.Directive.SetValue("proxy_pass " + v)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			h.Name = "proxy_pass"
			h.Params = p
		}
	}()

	return h.ReparseParams()
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

func (h *HTTPProxyPass) ReparseParams() (err error) {
	err = h.PosVerify()
	if err != nil {
		return err
	}
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
	switch protocol {
	case "http":
		defaultProxyPassPort = 80
	case "https":
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
			Sockets:    newSockets(ipv4s, port),
		})

		return nil
	}

	// parse upstream servers
	httpCtx := h.FatherPosSet().
		QueryOne(context.NewKeyWordsByType(context_type.TypeHttp).
			SetIsToLeafQuery(false)).
		Target()
	if httpCtx.Error() != nil {
		return errors.WithCode(code.ErrV3InvalidContext, "cannot query father HTTP Context. error: %++v", httpCtx.Error())
	}
	upstreamPos := httpCtx.ChildrenPosSet().QueryOne(
		context.NewKeyWordsByType(context_type.TypeUpstream).
			SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc).
			AppendMatchingFilter(func(targetCtx context.Context) bool {
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
			Sockets:    newSockets(ipv4s, port),
			ResolveErr: ToJSONError(err),
		})

		return nil
	}

	return upstreamPos.QueryAll(context.NewKeyWordsByType(context_type.TypeDirective).
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
						Sockets:    newSockets(srvIPv4s, srvPort),
						ResolveErr: ToJSONError(err),
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

func (h *HTTPProxyPass) PosVerify() error {
	// the father Context must in: Location, If in Location, Limit_except
	notMatchedError := errors.WithCode(code.ErrV3InvalidOperation, "Father context type not supported."+
		" Only the following types are supported: `Location`, `If` in `Location`, `Limit_except`")
	fatherPoses := FatherPosSetWithoutInclude(h)
	if err := fatherPoses.Error(); err != nil {
		if errors.IsCode(err, code.ErrV3CannotQueryFatherPosSetFromMainContext) {
			return notMatchedError
		}

		return err
	}

	return fatherPoses.Map(func(pos context.Pos) (context.Pos, error) {
		switch pos.Target().Type() {
		case context_type.TypeLocation:
			return pos, nil
		case context_type.TypeLimitExcept:
			return pos, nil
		case context_type.TypeIf:
			subFatherPoses := FatherPosSetWithoutInclude(pos.Target())
			if err := subFatherPoses.Error(); err != nil {
				if errors.IsCode(err, code.ErrV3CannotQueryFatherPosSetFromMainContext) {
					return pos, notMatchedError
				}

				return pos, err
			}

			return pos, subFatherPoses.Map(func(pos context.Pos) (context.Pos, error) {
				if pos.Target().Type() != context_type.TypeLocation {
					return pos, notMatchedError
				}

				return pos, nil
			}).Error()
		default:
			return pos, notMatchedError
		}
	}).Error()
}

func (h *HTTPProxyPass) ConnectivityCheck() ProxyPass {
	if err := h.ReparseParams(); err != nil {
		return newErrProxyPass(err)
	}

	var errs []error
	for _, address := range h.Addresses {
		for _, socket := range address.Sockets {
			errs = append(errs, socket.tcpCheck())
		}
	}
	err := errors.NewAggregate(errs)
	if err != nil {
		return newErrProxyPass(err)
	}

	return h
}

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

func (s *StreamProxyPass) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Enabled     bool                     `json:"enabled,omitempty"`
		ContextType context_type.ContextType `json:"context-type"`
		Value       string                   `json:"value"`
		ProxyPass   struct {
			OriginalAddress string           `json:"original-address"`
			Addresses       []ProxiedAddress `json:"addresses"`
		} `json:"proxy-pass,omitempty"`
	}{
		Enabled:     s.IsEnabled(),
		ContextType: context_type.TypeDirStreamProxyPass,
		Value:       s.Params,
		ProxyPass: struct {
			OriginalAddress string           `json:"original-address"`
			Addresses       []ProxiedAddress `json:"addresses"`
		}{OriginalAddress: s.OriginalAddress, Addresses: s.Addresses},
	})
}

func (s *StreamProxyPass) UnmarshalJSON(bytes []byte) error {
	tmp := struct {
		Enabled     bool                     `json:"enabled,omitempty"`
		ContextType context_type.ContextType `json:"context-type"`
		Value       string                   `json:"value"`
		ProxyPass   struct {
			OriginalAddress string           `json:"original-address"`
			Addresses       []ProxiedAddress `json:"addresses"`
		} `json:"proxy-pass,omitempty"`
	}{}
	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}
	if tmp.ContextType != context_type.TypeDirStreamProxyPass {
		return errors.WithCode(code.ErrV3InvalidContext, "invalid context-type: %s", tmp.ContextType)
	}
	s.enabled = tmp.Enabled
	s.Params = tmp.Value
	s.OriginalAddress = tmp.ProxyPass.OriginalAddress
	s.Addresses = tmp.ProxyPass.Addresses

	return nil
}

func (s *StreamProxyPass) Clone() context.Context {
	cloneProxiedAddress := make([]ProxiedAddress, len(s.Addresses))
	for i, address := range s.Addresses {
		cloneProxiedAddress[i] = ProxiedAddress{
			DomainName: address.DomainName,
			Port:       address.Port,
			Sockets:    address.Sockets,
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
	err = s.Directive.SetValue("proxy_pass " + v)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			s.Name = "proxy_pass"
			s.Params = p
		}
	}()

	return s.ReparseParams()
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

func (s *StreamProxyPass) ReparseParams() (err error) {
	err = s.PosVerify()
	if err != nil {
		return err
	}
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
			Sockets:    newSockets(ipv4s, port),
			ResolveErr: ToJSONError(err),
		})

		return nil
	}

	// parse upstream server
	if len(strings.TrimSpace(upstreamServer)) == 0 {
		return errors.WithCode(code.ErrV3InvalidOperation, "the upstream server must not be empty")
	}
	streamCtx := s.FatherPosSet().
		QueryOne(context.NewKeyWordsByType(context_type.TypeStream).
			SetIsToLeafQuery(false)).
		Target()
	if streamCtx.Error() != nil {
		return errors.WithCode(code.ErrV3InvalidContext, "cannot find father Stream Context. error: %++v", streamCtx.Error())
	}

	upstreamPos := streamCtx.ChildrenPosSet().QueryOne(
		context.NewKeyWordsByType(context_type.TypeUpstream).
			SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc).
			AppendMatchingFilter(func(targetCtx context.Context) bool {
				return targetCtx.Value() == upstreamServer
			}),
	)
	if err = upstreamPos.Target().Error(); err != nil {
		return errors.WithCode(code.ErrV3InvalidOperation, "unknown upstream server: '%s'. failed to query Upstream(Stream) Context, error: %v", upstreamServer, err)
	}

	return upstreamPos.QueryAll(context.NewKeyWordsByType(context_type.TypeDirective).
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
							Sockets:    newSockets(ipv4s, srvPort),
							ResolveErr: ToJSONError(err),
						})
						existAddrMap[addr] = true
					}
				}

				return dirPos, nil
			},
		).Error()
}

func (s *StreamProxyPass) PosVerify() error {
	// return nil, if itself is not enabled
	if !s.enabled {
		return nil
	}
	// the set father Context must be Server in Stream
	notMatchedError := errors.WithCode(code.ErrV3InvalidOperation, "Father context type not supported."+
		" Only the following types are supported: `Server` which is in `Stream`")
	fatherPoses := FatherPosSetWithoutInclude(s)
	if err := fatherPoses.Error(); err != nil {
		if errors.IsCode(err, code.ErrV3CannotQueryFatherPosSetFromMainContext) {
			return notMatchedError
		}

		return err
	}

	return fatherPoses.Map(func(pos context.Pos) (context.Pos, error) {
		if pos.Target().Type() != context_type.TypeServer {
			return pos, notMatchedError
		}

		return pos, FatherPosSetWithoutInclude(pos.Target()).Map(func(pos context.Pos) (context.Pos, error) {
			if pos.Target().Type() != context_type.TypeStream {
				return pos, notMatchedError
			}

			return pos, nil
		}).Error()
	}).Error()
}

func (s *StreamProxyPass) ConnectivityCheck() ProxyPass {
	if err := s.ReparseParams(); err != nil {
		return newErrProxyPass(err)
	}

	var errs []error
	for _, address := range s.Addresses {
		for _, socket := range address.Sockets {
			errs = append(errs, socket.allCheck())
		}
	}
	err := errors.NewAggregate(errs)
	if err != nil {
		return newErrProxyPass(err)
	}

	return s
}

func registerStreamProxyPassBuilder() error {
	builderMap[context_type.TypeDirStreamProxyPass] = func(value string) context.Context {
		if !dirStreamProxyPassMustCompile.MatchString("proxy_pass " + value) { //nolint:nestif
			return context.ErrContext(errors.Errorf("invalid value `%s`", value))
		}

		return &StreamProxyPass{Directive: *NewContext(context_type.TypeDirective, "proxy_pass "+value).(*Directive)}
	}

	return nil
}

type tmpProxyPass struct {
	Directive
	isInitialized bool
	rwLock        *sync.RWMutex
}

func (t *tmpProxyPass) MarshalJSON() ([]byte, error) {
	// initialize before marshaling
	p := t.verifyAndInitTo()
	if err := p.Error(); err != nil {
		return nil, err
	}
	// result marshaled data of Directive, when it has not been initialized.
	if !t.isInitialized {
		return t.Directive.MarshalJSON()
	}

	return json.Marshal(p)
}

func (t *tmpProxyPass) verifyAndInitTo() context.Context {
	t.rwLock.Lock()
	defer t.rwLock.Unlock()

	if t.isInitialized {
		return context.ErrContext(errors.New("ProxyPass is already initialized"))
	}

	// if tmpProxyPass has not been included into the MainContext, skip initializing
	if t.fatherContext == nil || !errors.IsCode(context.SetPos(t.fatherContext, 0).
		QueryAll(
			context.NewKeyWordsByType(context_type.TypeMain).
				SetIsToLeafQuery(false),
		).Error(), code.ErrV3CannotQueryFatherPosSetFromMainContext) {
		return t
	}

	// verify ProxyPass
	var p ProxyPass

	var verifyErr error
	if dirHTTPProxyPassMustCompile.MatchString(t.Value()) { //nolint:nestif
		hp := &HTTPProxyPass{Directive: t.Directive}
		verifyErr = hp.PosVerify()
		if verifyErr == nil {
			p = hp
		}
	}
	if p == nil && dirStreamProxyPassMustCompile.MatchString(t.Value()) {
		sp := &StreamProxyPass{Directive: t.Directive}
		verifyErr = sp.PosVerify()
		if verifyErr == nil {
			p = sp
		}
	}

	if verifyErr != nil {
		return context.ErrContext(verifyErr)
	} else if p == nil {
		return context.ErrContext(errors.Errorf("invalid value `%s`", t.Value()))
	}

	// initializing ProxyPass instance
	for i, pos := range t.fatherContext.ChildrenPosSet().List() {
		if pos.Target() == t {
			c := t.fatherContext.Modify(p, i)
			if c.Error() != nil {
				return c
			}
			t.isInitialized = true

			return p
		}
	}

	return context.ErrContext(errors.New("failed to initialize proxy pass"))
}

func (t *tmpProxyPass) PosVerify() error {
	return t.verifyAndInitTo().Error()
}

func (t *tmpProxyPass) ReparseParams() error {
	p := t.verifyAndInitTo()
	if err := p.Error(); err != nil {
		return err
	}
	if p != t {
		return p.(ProxyPass).ReparseParams()
	}

	return nil
}

func (t *tmpProxyPass) ConnectivityCheck() ProxyPass {
	p := t.verifyAndInitTo()
	if err := p.Error(); err != nil {
		return newErrProxyPass(err)
	}
	if p != t {
		return p.(ProxyPass).ConnectivityCheck()
	}

	return t
}

func (t *tmpProxyPass) SetFather(ctx context.Context) error {
	return t.Directive.SetFather(ctx)
}

func (t *tmpProxyPass) Type() context_type.ContextType {
	return context_type.TypeDirUninitializedProxyPass
}

func registerTmpProxyPassBuilder() error {
	builderMap[context_type.TypeDirUninitializedProxyPass] = func(value string) context.Context {
		return context.ErrContext(errors.New("cannot build a Directive of UninitializedProxyPass"))
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
	dirBuilder, ok := inStackParseFuncMap[context_type.TypeDirective]
	if !ok {
		return errors.New("directive parse function is not found in stack parse function map")
	}
	inStackParseFuncMap[context_type.TypeDirective] = func(data []byte, idx *int) context.Context {
		dir := dirBuilder(data, idx)
		if dir == context.NullContext() {
			return dir
		}
		if dir.Type() == context_type.TypeDirective {
			if dirHTTPProxyPassMustCompile.MatchString(dir.Value()) || dirStreamProxyPassMustCompile.MatchString(dir.Value()) { //nolint:nestif
				return &tmpProxyPass{
					Directive: *dir.(*Directive),
					rwLock:    new(sync.RWMutex),
				}
			}
		}

		return dir
	}

	return nil
}

type errProxyPass struct {
	*context.ErrorContext
}

func (e errProxyPass) PosVerify() error {
	return e.Error()
}

func (e errProxyPass) ReparseParams() error {
	return e.Error()
}

func (e errProxyPass) ConnectivityCheck() ProxyPass {
	return e
}

func newErrProxyPass(err error) ProxyPass {
	return &errProxyPass{context.ErrContext(err).(*context.ErrorContext)}
}
