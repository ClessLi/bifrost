package nginx

import (
	"context"
	"encoding/json"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"

	"github.com/marmotedu/errors"
)

type webServerConfigStore struct {
	configs map[string]configuration.NginxConfig
}

func (w *webServerConfigStore) GetServerNames(ctx context.Context) (*v1.ServerNames, error) {
	serverNames := make(v1.ServerNames, 0)
	for name := range w.configs {
		serverNames = append(serverNames, v1.ServerName{Name: name})
	}

	return &serverNames, nil
}

func (w *webServerConfigStore) Get(ctx context.Context, servername *v1.ServerName) (*v1.WebServerConfig, error) {
	if conf, has := w.configs[servername.Name]; has {
		jdata := conf.Json()
		if len(jdata) == 0 {
			return nil, errors.WithCode(code.ErrInvalidConfig, "the config belonging to '%s' web server is null", servername.Name)
		}

		fp := utilsV3.NewConfigFingerprinter(conf.Dump()).Fingerprints()
		fpdata, err := json.Marshal(fp)
		if err != nil {
			return nil, errors.WithCode(code.ErrInvalidConfig, "failed to marshal fingerprints of the web server config, cased by: %v", err)
		}

		return &v1.WebServerConfig{
			ServerName:           servername,
			JsonData:             jdata,
			OriginalFingerprints: fpdata,
		}, nil
	}

	return nil, errors.WithCode(code.ErrConfigurationNotFound, "the web server '%s' not found", servername.Name)
}

func (w *webServerConfigStore) ConnectivityCheckOfProxiedServers(ctx context.Context, pos *v1.WebServerConfigContextPos) (*v1.ContextData, error) {
	// check original fingerprints
	var ofp utilsV3.ConfigFingerprints
	err := json.Unmarshal(pos.OriginalFingerprints, &ofp)
	if err != nil {
		return nil, errors.WithCode(code.ErrV3InvalidConfigFingerprints, "failed to unmarshal fingerprints of the web server config, cased by: %v", err)
	}

	// find web server config
	conf, has := w.configs[pos.ServerName.Name]
	if !has {
		return nil, errors.WithCode(code.ErrConfigurationNotFound, "the web server '%s' not found", pos.ServerName.Name)
	}

	// match fingerprints
	if utilsV3.NewConfigFingerprinter(w.configs[pos.ServerName.Name].Dump()).Diff(ofp) {
		return nil, errors.WithCode(code.ErrV3InvalidConfigFingerprints, "the original fingerprints(%v) submitted for the config update request does not match the web server config fingerprints", ofp)
	}

	// verify the pos indexes
	if len(pos.ContextPos.PosIndex) < 1 {
		return nil, errors.WithCode(code.ErrV3InvalidOperation, "the length of posIndex is less than 1")
	}

	// find target config
	targetConfig, err := conf.Main().GetConfig(pos.ContextPos.ConfigPath)
	if err != nil {
		return nil, errors.WithCode(code.ErrConfigurationNotFound, "the config file(%s) belonging to '%s' web server not found, cased by: %v", pos.ContextPos.ConfigPath, pos.ServerName.Name, err)
	}

	// find target proxy_pass
	target := targetConfig.Child(int(pos.ContextPos.PosIndex[0]))
	for _, index := range pos.ContextPos.PosIndex[1:] {
		target = target.Child(int(index))
	}
	if err = target.Error(); err != nil {
		return nil, err
	}

	// verify type of the proxy_pass
	proxyPass, ok := target.(local.ProxyPass)
	if !ok {
		return nil, errors.WithCode(code.ErrV3InvalidContext, "the context to be checked is not a `ProxyPass` context")
	}

	// check proxied servers' connectivity
	proxyPass = proxyPass.ConnectivityCheck()
	if err = proxyPass.Error(); err != nil {
		return nil, err
	}

	// marshal
	jsondata, err := json.Marshal(proxyPass)
	if err != nil {
		return nil, err
	}

	return &v1.ContextData{JsonData: jsondata}, nil
}

func (w *webServerConfigStore) Update(ctx context.Context, config *v1.WebServerConfig) error {
	// check original fingerprints
	var ofp utilsV3.ConfigFingerprints
	err := json.Unmarshal(config.OriginalFingerprints, &ofp)
	if err != nil {
		return errors.WithCode(code.ErrInvalidConfig, "failed to unmarshal fingerprints of the web server config, cased by: %v", err)
	}

	if conf, has := w.configs[config.ServerName.Name]; has {
		if utilsV3.NewConfigFingerprinter(w.configs[config.ServerName.Name].Dump()).Diff(ofp) {
			return errors.WithCode(code.ErrV3InvalidConfigFingerprints, "the original fingerprints(%v) submitted for the config update request does not match the web server config fingerprints", ofp)
		}

		return conf.UpdateFromJsonBytes(config.JsonData)
	}

	return errors.WithCode(code.ErrConfigurationNotFound, "the web server '%s' not found", config.ServerName.Name)
}

var _ storev1.WebServerConfigStore = &webServerConfigStore{}

func newNginxConfigStore(store *webServerStore) storev1.WebServerConfigStore {
	return &webServerConfigStore{configs: store.configsManger.GetConfigs()}
}
