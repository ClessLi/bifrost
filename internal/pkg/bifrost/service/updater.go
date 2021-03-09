package service

type Updater interface {
	Update(UpdateRequestInfo) UpdateResponseInfo
}

type updater struct {
	offstage offstageUpdater
}

func NewUpdater(offstage offstageUpdater) Updater {

	if offstage == nil {
		panic("offstage is nil")
	}

	return &updater{offstage: offstage}
}

func (u *updater) Update(req UpdateRequestInfo) UpdateResponseInfo {
	var err error

	if req == nil {
		err = ErrNilRequestInfo
		return NewUpdateResponseInfo("", err)
	}

	serverName := req.GetServerName()

	switch req.GetRequestType() {
	case UpdateConfig:
		err = u.offstage.UpdateConfig(serverName, req.GetData())
	default:
		err = ErrUnknownRequestType
	}

	return NewUpdateResponseInfo(serverName, err)
}
