package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
)

func init() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf("Usage %s <config file name>...\n", os.Args[0])
}

func main() {
	var configs []*nginx.Config
	cacheSlice := make([]nginx.Caches, 0)
	for _, s := range os.Args[1:] {
		path, caches, loadErr := nginx.Load(s)
		if loadErr != nil {
			fmt.Println(loadErr)
			usage()
			os.Exit(2)
		}
		cacheSlice = append(cacheSlice, caches)
		config, confErr := caches.GetConfig(path)
		if confErr != nil {
			fmt.Println(confErr)
			os.Exit(3)
		}
		configs = append(configs, config)
	}

	for i := range configs {
		bakPath, bakErr := nginx.Backup(configs[i], "nginx.conf", 7, 1, filepath.Dir(configs[i].Value))
		if bakErr != nil && bakErr != nginx.NoBackupRequired {
			fmt.Printf("backup to %s failed, cased by %s", bakPath, bakErr)
			os.Exit(4)
		}
	}

	for _, config := range configs {
		_, saveErr := nginx.Save(config)
		if saveErr != nil {
			fmt.Println(saveErr)
			os.Exit(5)
		}
	}
}
