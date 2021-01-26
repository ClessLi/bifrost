package service

type RequestType int

const (
	Unknown RequestType = iota
	DisplayConfig
	GetConfig
	ShowStatistics
	DisplayStatus

	UpdateConfig

	WatchLog
)
