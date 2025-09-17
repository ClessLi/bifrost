package main

import (
	"math/rand"
	"os"
	"runtime"

	"github.com/ClessLi/bifrost/internal/bifrost"

	"github.com/marmotedu/component-base/pkg/time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	bifrost.NewApp("bifrost").Run()
}
