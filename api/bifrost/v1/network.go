package v1

const (
	NetUnknown NetworkConnectivity = iota
	NetReachable
	NetUnreachable
)

type NetworkConnectivity int
