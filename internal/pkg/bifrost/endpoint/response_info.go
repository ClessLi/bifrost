package endpoint

import (
	"bytes"
	"errors"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"time"
)

var ErrWatchResponseInfoCloseTimeout = errors.New("watch response info close timeout")

// BytesResponseInfo 字节数据反馈信息接口对象
type BytesResponseInfo interface {
	Respond() []byte
	Error() string
}

// ErrorResponseInfo 错误反馈信息接口对象
type ErrorResponseInfo interface {
	Error() string
}

// WatchResponseInfo 数据监看反馈信息接口对象
type WatchResponseInfo interface {
	Respond() <-chan BytesResponseInfo
	Close() error
}

type bytesResponseInfo struct {
	Result *bytes.Buffer `json:"result"`
	Err    error         `json:"error"`
}

func (br bytesResponseInfo) Respond() []byte {
	return br.Result.Bytes()
}

func (br bytesResponseInfo) Error() string {
	if br.Err != nil {
		return br.Err.Error()
	}
	return ""
}

func newViewResponseInfo(svcResponseInfo service.ViewResponseInfo) BytesResponseInfo {
	return &bytesResponseInfo{
		Result: bytes.NewBuffer(svcResponseInfo.Bytes()),
		Err:    svcResponseInfo.Error(),
	}
}

type errorResponseInfo struct {
	Err error `json:"error"`
}

func (er errorResponseInfo) Error() string {
	if er.Err != nil {
		return er.Err.Error()
	}
	return ""
}

func newUpdateResponseInfo(svcResponseInfo service.UpdateResponseInfo) ErrorResponseInfo {
	return &errorResponseInfo{Err: svcResponseInfo.Error()}
}

type watchResponseInfo struct {
	Result           chan BytesResponseInfo `json:"result"`
	signalChan       chan int
	serviceCloseFunc func() error
}

func (wr watchResponseInfo) Respond() <-chan BytesResponseInfo {
	return wr.Result
}

func (wr watchResponseInfo) Close() error {
	select {
	case wr.signalChan <- 9:
		return wr.serviceCloseFunc()
	case <-time.After(time.Second * 10):
		return ErrWatchResponseInfoCloseTimeout
	}

}

func newWatchResponseInfo(svcResponseInfo service.WatchResponseInfo, signalChan chan int, bytesResponseInfoChan chan BytesResponseInfo) WatchResponseInfo {
	go func() {
		for {
			select {
			case b := <-svcResponseInfo.BytesChan():
				bytesResponseInfoChan <- &bytesResponseInfo{
					Result: bytes.NewBuffer(b),
					Err:    nil,
				}
			case err := <-svcResponseInfo.TransferErrorChan():
				bytesResponseInfoChan <- &bytesResponseInfo{
					Result: bytes.NewBuffer([]byte("")),
					Err:    err,
				}
			case sig := <-signalChan:
				if sig == 9 {
					return
				}
			}
		}
	}()
	return &watchResponseInfo{
		Result:           bytesResponseInfoChan,
		serviceCloseFunc: svcResponseInfo.Close,
		signalChan:       signalChan,
	}
}
