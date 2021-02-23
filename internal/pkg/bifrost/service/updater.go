package service

type Updater interface {
	Update(UpdateRequestInfo) UpdateResponseInfo
}

type updater struct {
	offstage offstageUpdater
}

func NewUpdater(offstage offstageUpdater) Updater {
	return &updater{offstage: offstage}
}

func (u *updater) Update(req UpdateRequestInfo) UpdateResponseInfo {
	serverName := req.GetServerName()
	var err error
	switch req.GetRequestType() {
	case UpdateConfig:
		err = u.offstage.UpdateConfig(serverName, req.GetData())
	default:
		err = UnknownRequestType
	}
	return NewUpdateResponseInfo(serverName, err)
}
