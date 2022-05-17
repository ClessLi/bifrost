package fake

import (
	pbv1 "github.com/yongPhone/bifrost/api/protobuf-spec/bifrostpb/v1"
	log "github.com/yongPhone/bifrost/pkg/log/v1"
)

type webServerLogWatcher struct{}

func (w webServerLogWatcher) Watch(request *pbv1.LogWatchRequest, stream pbv1.WebServerLogWatcher_WatchServer) error {
	log.Infof("watch web server log '%s'", request.ServerName)

	return nil
}
