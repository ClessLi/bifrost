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
	//os.Stderr = utils.Logf
	Logger = log.New(os.Stderr, "", log.LstdFlags)
	//Logger = log.New(utils.Stdoutf, "", log.LstdFlags)

	KitLogger = kitlog.NewLogfmtLogger(os.Stderr)
	//KitLogger = kitlog.NewLogfmtLogger(utils.Logf)
	KitLogger = kitlog.With(KitLogger, "ts", kitlog.TimestampFormat(func() time.Time { return time.Now().Local() }, time.RFC3339Nano))
	KitLogger = kitlog.With(KitLogger, "caller", kitlog.DefaultCaller)

}
