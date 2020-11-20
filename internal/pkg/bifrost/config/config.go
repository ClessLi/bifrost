package config

import (
	kitlog "github.com/go-kit/kit/log"
	"log"
	"os"
	"time"
)

var Logger *log.Logger
var KitLogger kitlog.Logger

func init() {
	Logger = log.New(os.Stderr, "", log.LstdFlags)

	KitLogger = kitlog.NewLogfmtLogger(os.Stderr)
	KitLogger = kitlog.With(KitLogger, "ts", kitlog.TimestampFormat(func() time.Time { return time.Now().Local() }, time.RFC3339Nano))
	KitLogger = kitlog.With(KitLogger, "caller", kitlog.DefaultCaller)

}
