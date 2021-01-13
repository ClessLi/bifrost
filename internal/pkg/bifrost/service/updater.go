package service

type Updater interface {
	Update(UpdateRequester) UpdateResponder
}

type updater struct {
	offstage offstageUpdater
}

func NewUpdater(offstage offstageUpdater) Updater {
	return &updater{offstage: offstage}
}

func (u *updater) Update(req UpdateRequester) UpdateResponder {
	serverName := req.GetServerName()
	var err error
	switch req.GetRequestType() {
	case UpdateConfig:
		err = u.offstage.UpdateConfig(serverName, req.GetData())
	default:
		err = UnknownRequestType
	}
	return NewUpdateResponder(serverName, err)
}
