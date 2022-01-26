package main

import (
	"github.com/ClessLi/bifrost/internal/bifrost"
	"github.com/marmotedu/component-base/pkg/time"
	"math/rand"
	"os"
	"runtime"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	bifrost.NewApp("bifrost").Run()
}
