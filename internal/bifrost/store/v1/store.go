package v1

var client StoreFactory

type StoreFactory interface {
	WebServerConfig() WebServerConfigStore
	WebServerStatistics() WebServerStatisticsStore
	WebServerStatus() WebServerStatusStore
	WebServerLogWatcher() WebServerLogWatcherStore
	WebServerBinCMD() WebServerBinCMDStore
	Close() error
}

func Client() StoreFactory {
	return client
}

func SetClient(factory StoreFactory) {
	client = factory
}
