package service

type Viewer interface {
	View(ViewRequestInfo) ViewResponseInfo
}

type viewer struct {
	offstage offstageViewer
}

func NewViewer(offstage offstageViewer) Viewer {
	if offstage == nil {
		panic("offstage is nil")
	}

	return &viewer{offstage: offstage}
}

func (v viewer) View(req ViewRequestInfo) ViewResponseInfo {
	var err error
	var data []byte

	if req == nil {
		err = ErrNilRequestInfo
		return NewViewResponseInfo("", data, err)
	}

	serverName := req.GetServerName()

	switch req.GetRequestType() {
	case DisplayConfig:
		data, err = v.offstage.DisplayConfig(serverName)
	case GetConfig:
		data, err = v.offstage.GetConfig(serverName)
	case ShowStatistics:
		data, err = v.offstage.ShowStatistics(serverName)
	case DisplayServersStatus:
		data, err = v.offstage.DisplayServersStatus()
		serverName = ""
	default:
		err = ErrUnknownRequestType
	}

	return NewViewResponseInfo(serverName, data, err)
}
