package service

type Monitor interface {
	DisplayStatus() ([]byte, error)
}
