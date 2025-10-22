package v1

type WebServerConfig struct {
	ServerName           *ServerName `json:"server-name"`
	JsonData             []byte      `json:"data"`
	OriginalFingerprints []byte      `json:"original-fingerprints"`
}

type ServerNames []ServerName

type ServerName struct {
	Name string `json:"name"`
}

type WebServerConfigContextPos struct {
	ServerName           *ServerName `json:"server-name"`
	ContextPos           *ContextPos `json:"context-pos"`
	OriginalFingerprints []byte      `json:"original-fingerprints"`
}

type ContextPos struct {
	ConfigPath string  `json:"config-path"`
	PosIndex   []int32 `json:"pos-index"`
}

type ContextData struct {
	JsonData []byte `json:"data"`
}
