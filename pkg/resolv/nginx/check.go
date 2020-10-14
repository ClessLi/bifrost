package nginx

import (
	"os"
	"os/exec"
)

func Check(config *Config, ng string) (err error) {
	err = config.Save()
	if err != nil {
		return
	}
	return check(ng, config.Value)
}

func check(ng string, s string) error {
	cmd := exec.Command(ng, "-tc", s)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
