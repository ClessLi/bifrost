package nginx

import (
	"os"
	"os/exec"
)

func SaveWithCheck(config *Config, ng string) (caches Caches, err error) {
	caches, err = config.Save()
	if err != nil {
		return
	}
	return caches, check(ng, config.Value)
}

func check(ng string, s string) error {
	cmd := exec.Command(ng, "-tc", s)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
