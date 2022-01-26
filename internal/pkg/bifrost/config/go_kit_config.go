package config

import (
	kitlog "github.com/go-kit/kit/log"
	"io"
	"time"
)

//var Logger *log.Logger
//var KitLogger kitlog.Logger

//func init() {

//}

func KitLogger(out io.Writer) kitlog.Logger {
	//os.Stderr = utils.Logf
	//Logger = log.New(out, "", log.LstdFlags)
	//Logger = log.New(utils.Stdoutf, "", log.LstdFlags)

	kitLogger := kitlog.NewLogfmtLogger(out)
	//KitLogger = kitlog.NewLogfmtLogger(utils.Logf)
	kitLogger = kitlog.With(kitLogger, "ts", kitlog.TimestampFormat(func() time.Time { return time.Now().Local() }, time.RFC3339Nano))
	kitLogger = kitlog.With(kitLogger, "caller", kitlog.DefaultCaller)
	return kitLogger
}
