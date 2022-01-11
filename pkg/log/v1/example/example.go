package main

import (
	"context"
	"flag"

	log "github.com/ClessLi/bifrost/pkg/log/v1"
)

var (
	h bool

	level  int
	format string
)

func main() {
	flag.BoolVar(&h, "h", false, "Print this help.")
	flag.IntVar(&level, "l", 0, "Log level.")
	flag.StringVar(&format, "f", "console", "log output format.")

	flag.Parse()

	if h {
		flag.Usage()

		return
	}

	// logger配置
	opts := &log.Options{
		Level:            "debug",
		Format:           "console",
		EnableColor:      true,
		DisableCaller:    true,
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{"error.log", "stderr"},
	}
	// 初始化全局logger
	log.Init(opts)
	defer log.Flush()

	// Debug、Info(with field)、Warnf、Errorw使用
	log.Debug("This is a debug message")
	log.Info("This is a info message", log.Int32("int_key", 10))
	log.Warnf("This is a formatted %s message", "warn")
	log.Errorw("Message printed with Errorw", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")

	// WithValues使用
	lv := log.WithValues("X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")
	lv.Infow("Info message printed with [WithValues] logger")
	lv.Infow("Debug message printed with [WithValues] logger")

	// Context使用
	ctx := lv.WithContext(context.Background())
	lc := log.FromContext(ctx)
	lc.Info("Message printed with [WithContext] logger")

	ln := lv.WithName("test")
	ln.Info("Message printed with [WithName] logger")

	// V level使用
	log.V(1).Info("This is a V level message")
	log.V(1).Infow("This is a V level message with fields", "X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")
}
