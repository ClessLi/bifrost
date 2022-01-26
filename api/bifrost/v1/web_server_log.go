package v1

type WebServerLog struct {
	Lines <-chan []byte `json:"lines"`
}

type WebServerLogWatchRequest struct {
	ServerName          *ServerName `json:"server-name"`
	LogName             string      `json:"log-path"`
	FilteringRegexpRule string      `json:"filtering-regexp-rule"`
}
