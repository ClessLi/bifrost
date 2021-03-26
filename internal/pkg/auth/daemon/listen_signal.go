package daemon

import (
	"os"
	"os/signal"
	"syscall"
)

func ListenSignal(sigChan chan<- int) {
	procSigs := make(chan os.Signal, 1)
	signal.Notify(procSigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case s := <-procSigs:
		Log(NOTICE, "Get system signal: %s", s.String())
		sigChan <- 9
		Log(DEBUG, "Stop listen system signal")
	}
}
