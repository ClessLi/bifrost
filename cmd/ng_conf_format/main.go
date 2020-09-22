package main

import (
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"os"
	"path/filepath"
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
	for _, s := range os.Args[1:] {
		config, loadErr := nginx.Load(s)
		if loadErr != nil {
			fmt.Println(loadErr)
			usage()
			os.Exit(2)
		}
		configs = append(configs, config)
	}

	for _, config := range configs {
		bakPath, bakErr := nginx.Backup(config, "nginx.conf", 7, 1, filepath.Dir(config.Value))
		if bakErr != nil && bakErr != nginx.NoBackupRequired {
			fmt.Printf("backup to %s failed, cased by %s", bakPath, bakErr)
			os.Exit(3)
		}
	}

	for _, config := range configs {
		saveErr := nginx.Save(config)
		if saveErr != nil {
			fmt.Println(saveErr)
			os.Exit(4)
		}
	}
}
