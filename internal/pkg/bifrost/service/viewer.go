package service

type Viewer interface {
	View(ViewRequestInfo) ViewResponseInfo
}

type viewer struct {
	offstage offstageViewer
}

func NewViewer(offstage offstageViewer) Viewer {
	return &viewer{offstage: offstage}
}

func (v viewer) View(req ViewRequestInfo) ViewResponseInfo {
	serverName := req.GetServerName()
	data := make([]byte, 0)
	var err error
	switch req.GetRequestType() {
	case DisplayConfig:
		data, err = v.offstage.DisplayConfig(serverName)
	case GetConfig:
		data, err = v.offstage.GetConfig(serverName)
	case ShowStatistics:
		data, err = v.offstage.ShowStatistics(serverName)
	case DisplayStatus:
		data, err = v.offstage.DisplayStatus()
		serverName = ""
	default:
		err = UnknownRequestType
	}
	return NewViewResponseInfo(serverName, data, err)
}
