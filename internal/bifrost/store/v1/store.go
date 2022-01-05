package v1

var client StoreFactory

type StoreFactory interface {
	WebServerConfig() WebServerConfigStore
	Close() error
}

func Client() StoreFactory {
	return client
}

func SetClient(factory StoreFactory) {
	client = factory
}
