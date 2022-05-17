package main

import (
	"math/rand"
	"os"
	"runtime"

	"github.com/marmotedu/component-base/pkg/time"

	"github.com/yongPhone/bifrost/internal/bifrost"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	bifrost.NewApp("bifrost").Run()
}
