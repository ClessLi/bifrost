package bifrost

import (
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"os"
	"os/signal"
	"syscall"
)

func ListenSignal(sigChan chan<- int) {
	procSigs := make(chan os.Signal, 1)
	signal.Notify(procSigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case s := <-procSigs:
		utils.Logger.NoticeF("Get system signal: %s", s.String())
		sigChan <- 9
		utils.Logger.Debug("Stop listen system signal")
	}
}
