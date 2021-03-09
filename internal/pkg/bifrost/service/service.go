package service

type Service interface {
	Viewer() Viewer
	Updater() Updater
	Watcher() Watcher
	// TODO: WatchLog暂用锁机制，一个日志文件仅允许一个终端访问
}

type service struct {
	viewer  Viewer
	updater Updater
	watcher Watcher
}

func NewService(viewer Viewer, updater Updater, watcher Watcher) Service {
	if viewer == nil {
		panic("viewer is nil")
	}
	if updater == nil {
		panic("updater is nil")
	}
	if watcher == nil {
		panic("watcher is nil")
	}
	return &service{
		viewer:  viewer,
		updater: updater,
		watcher: watcher,
	}
}

func (s service) Viewer() Viewer {
	if s.viewer == nil {
		panic("viewer is nil")
	}
	return s.viewer
}

func (s service) Updater() Updater {
	if s.updater == nil {
		panic("updater is nil")
	}
	return s.updater
}

func (s service) Watcher() Watcher {
	if s.watcher == nil {
		panic("watcher is nil")
	}
	return s.watcher
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service

type ViewerMiddleware func(Viewer) Viewer
type UpdaterMiddleware func(Updater) Updater
type WatcherMiddleware func(Watcher) Watcher
