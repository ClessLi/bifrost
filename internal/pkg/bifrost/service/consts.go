package service

type RequestType int

type State int

const ( // RequestType
	UnknownReqType RequestType = iota
	DisplayConfig
	GetConfig
	ShowStatistics
	DisplayServersStatus

	UpdateConfig

	WatchLog
)

const ( // State
	UnknownState State = iota
	Disabled
	Initializing
	Abnormal
	Normal
)
