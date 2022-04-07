package config

import (
	"log"
	"os"
	"time"

	kitlog "github.com/go-kit/kit/log"
)

var (
	Logger    *log.Logger
	KitLogger kitlog.Logger
)

func init() {
	Logger = log.New(os.Stderr, "", log.LstdFlags)

	KitLogger = kitlog.NewLogfmtLogger(os.Stderr)
	KitLogger = kitlog.With(
		KitLogger,
		"ts",
		kitlog.TimestampFormat(func() time.Time { return time.Now().Local() }, time.RFC3339Nano),
	)
	KitLogger = kitlog.With(KitLogger, "caller", kitlog.DefaultCaller)
}
