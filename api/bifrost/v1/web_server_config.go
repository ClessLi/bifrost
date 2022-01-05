package v1

type WebServerConfig struct {
	ServerName *ServerName `json:"server-name"`
	JsonData   []byte      `json:"data"`
}

type ServerName struct {
	Name string `json:"name"`
}
