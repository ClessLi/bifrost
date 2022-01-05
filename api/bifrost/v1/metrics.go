package v1

const ( // State
	UnknownState State = iota
	Disabled
	Initializing
	Abnormal
	Normal
)

type Metrics struct {
	OS             string           `json:"system"`
	Time           string           `json:"time"`
	Cpu            string           `json:"cpu"`
	Mem            string           `json:"mem"`
	Disk           string           `json:"disk"`
	StatusList     []*WebServerInfo `json:"status-list"`
	BifrostVersion string           `json:"bifrost-version"`
}

type WebServerInfo struct {
	Name    string `json:"name"`
	Status  State  `json:"status"`
	Version string `json:"version"`
}

type State int
