package endpoint

// ViewRequestInfo 信息展示请求信息
type ViewRequestInfo struct {
	ViewType   string `json:"view_type"`
	ServerName string `json:"server_name"`
	Token      string `json:"token"`
}

// UpdateRequestInfo 数据更新请求信息
type UpdateRequestInfo struct {
	UpdateType string `json:"update_type"`
	ServerName string `json:"server_name"`
	Token      string `json:"token"`
	Data       []byte `json:"data"`
}

// WatchRequestInfo 数据监看请求信息
type WatchRequestInfo struct {
	WatchType   string `json:"watch_type"`
	ServerName  string `json:"server_name"`
	Token       string `json:"token"`
	WatchObject string `json:"watch_object"`
}
