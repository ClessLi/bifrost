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
