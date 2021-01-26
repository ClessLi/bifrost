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

func NewService(offstageService *BifrostService) Service {
	return &service{
		viewer:  NewViewer(offstageService),
		updater: NewUpdater(offstageService),
		watcher: NewWatcher(offstageService),
	}
}

func (s service) Viewer() Viewer {
	return s.viewer
}

func (s service) Updater() Updater {
	return s.updater
}

func (s service) Watcher() Watcher {
	return s.watcher
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service

type ViewerMiddleware func(Viewer) Viewer
type UpdaterMiddleware func(Updater) Updater
type WatcherMiddleware func(Watcher) Watcher
